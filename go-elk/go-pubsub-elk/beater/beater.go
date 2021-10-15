package beater

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"runtime"
	"sync"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/elastic/beats/v7/libbeat/beat"
	"github.com/elastic/beats/v7/libbeat/common"
	"github.com/elastic/beats/v7/libbeat/logp"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Pubsubbeat struct {
	done         chan struct{}
	config       *config.Config
	client       beat.Client
	pubsubClient *pubsub.Client
	subscription *pubsub.Subscription
	logger       *logp.Logger
	zippers      *sync.Pool
}

func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	config, err := config.GetAndValidateConfig(cfg)
	if err != nil {
		return nil, err
	}

	logger := logp.NewLogger(fmt.Sprintf("PubSub: %s/%s/%s", config.Project, config.Topic, config.Subscription.Name))
	logger.Infof("config retrieved: %+v", config)

	client, err := createPubsubClient(config)
	if err != nil {
		return nil, err
	}

	subscription, err := getOrCreateSubscription(client, config)
	if err != nil {
		return nil, err
	}

	connectionPoolSize := config.Subscription.ConnectionPoolSize
	subscription.ReceiveSettings.Synchronous = false // explicit
	subscription.ReceiveSettings.NumGoroutines = connectionPoolSize

	if connectionPoolSize == 1 {
		logger.Warnf("Pub/Sub streaming pull has a per-subscriber throughput limit, https://cloud.google.com/pubsub/quotas")
		logger.Warnf("Use `subscription.connection_pool_size` to increase the number of subscribers.")
	}

	bt := &Pubsubbeat{
		done:         make(chan struct{}),
		config:       config,
		pubsubClient: client,
		subscription: subscription,
		logger:       logger,
		zippers:      &sync.Pool{New: func() interface{} { return new(gzip.Reader) }},
	}
	return bt, nil
}

var keysToIgnore = map[string]bool{
	"password": true,
}

var keysToFlatten = map[string]bool{
	"tlv":        true,
	"connection": true,
}

func (bt *Pubsubbeat) Run(b *beat.Beat) error {
	bt.logger.Info("pubsubbeat is running! Hit CTRL-C to stop it.")

	var err error
	bt.client, err = b.Publisher.Connect()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-bt.done
		// The beat is stopping...
		bt.logger.Info("cancelling PubSub receive context...")
		cancel()
		bt.logger.Info("closing PubSub client...")
		bt.pubsubClient.Close()
	}()

	err = bt.subscription.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
		// This callback is invoked concurrently by multiple goroutines
		var datetime time.Time

		if m.Attributes["pubsubbeat.compression"] == "gzip" {
			err = bt.decompress(m)
			if err != nil {
				bt.logger.Warnf("failed to decompress gzip: %s", err)
				m.Nack()
				return
			}
		}

		var rawRecords [][]byte
		if m.Attributes["pubsubbeat.batch_ndjson"] == "true" {
			bt.logger.Infof("Incomfing messsage: %s", m.Data)
			rawRecords = bytes.Split(m.Data, []byte("\n"))
		} else {
			rawRecords = [][]byte{m.Data}
		}

		var batch []beat.Event

		for _, rawRecord := range rawRecords {
			if len(rawRecord) == 0 {
				continue
			}

			eventMap := common.MapStr{
				"publish_time": m.PublishTime,
			}

			var unmarshalErr error
			if bt.config.Json.FieldsUnderRoot {
				unmarshalErr = json.Unmarshal(rawRecord, &eventMap)
				if unmarshalErr == nil && bt.config.Json.FieldsUseTimestamp {
					var timeErr error
					timestamp := eventMap[bt.config.Json.FieldsTimestampName]
					delete(eventMap, bt.config.Json.FieldsTimestampName)
					datetime, timeErr = time.Parse(bt.config.Json.FieldsTimestampFormat, timestamp.(string))
					if timeErr != nil {
						bt.logger.Errorf("Failed to format timestamp string as time. Using time.Now(): %s", timeErr)
					}
				}
			} else {
				var jsonData map[string]interface{}
				unmarshalErr = json.Unmarshal(rawRecord, &jsonData)
				if unmarshalErr == nil {
					for key, val := range jsonData {
						if _, found := keysToIgnore[key]; found {
							continue
						}

						if _, found := keysToFlatten[key]; found {
							hash, ok := val.(map[string]interface{})
							if val == nil || !ok {
								log.Printf("Value %v is not map[string]interface{} type", val)
								continue
							}
							for k, v := range hash {
								if _, found := keysToIgnore[k]; found {
									continue
								}

								eventMap["app_"+fmt.Sprintf("%s_%s", key, k)] = v
							}

							continue
						}

						eventMap["app_"+key] = val
					}
				}
			}

			if unmarshalErr != nil {
				bt.logger.Warnf("failed to decode json message: %s", unmarshalErr)
				bt.logger.Warnf("Original message: %s", rawRecord)
				if bt.config.Json.AddErrorKey {
					eventMap["error"] = common.MapStr{
						"key":              "json",
						"message":          fmt.Sprintf("failed to decode json message: %s", unmarshalErr),
						"original_message": string(rawRecord),
					}
				}
			}

			if datetime.IsZero() {
				datetime = time.Now()
			}
			batch = append(batch, beat.Event{
				Timestamp: datetime,
				Fields:    eventMap,
			})
		}

		bt.client.PublishAll(batch)

		// TODO: Evaluate using AckHandler.
		m.Ack()
	})

	if err != nil {
		return fmt.Errorf("fail to receive message from subscription %q: %v", bt.subscription.String(), err)
	}

	return nil
}

func (bt *Pubsubbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}

func (bt *Pubsubbeat) decompress(m *pubsub.Message) error {
	rc := bt.zippers.Get().(*gzip.Reader)
	if err := rc.Reset(bytes.NewReader(m.Data)); err != nil {
		return fmt.Errorf("rc.Reset: %v", err)
	}
	var data bytes.Buffer
	if _, err := io.Copy(&data, rc); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := rc.Close(); err != nil {
		return fmt.Errorf("gzip.Close: %v", err)
	}
	bt.zippers.Put(rc)
	m.Data = data.Bytes()
	return nil
}

func createPubsubClient(config *config.Config) (*pubsub.Client, error) {
	ctx := context.Background()
	userAgent := fmt.Sprintf(
		"Elastic/Pubsubbeat (%s; %s)", runtime.GOOS, runtime.GOARCH)
	options := []option.ClientOption{option.WithUserAgent(userAgent)}
	if config.CredentialsFile != "" {
		options = append(options, option.WithCredentialsFile(config.CredentialsFile))
	}

	client, err := pubsub.NewClient(ctx, config.Project, options...)
	if err != nil {
		return nil, fmt.Errorf("fail to create pubsub client: %v", err)
	}
	return client, nil
}

func getOrCreateSubscription(client *pubsub.Client, config *config.Config) (*pubsub.Subscription, error) {
	if !config.Subscription.Create {
		subscription := client.Subscription(config.Subscription.Name)
		return subscription, nil
	}

	topic := client.Topic(config.Topic)
	ctx := context.Background()

	subscription, err := client.CreateSubscription(ctx, config.Subscription.Name, pubsub.SubscriptionConfig{
		Topic:               topic,
		RetainAckedMessages: config.Subscription.RetainAckedMessages,
		RetentionDuration:   config.Subscription.RetentionDuration,
		RetryPolicy:         &pubsub.RetryPolicy{MinimumBackoff: 1, MaximumBackoff: 120},
		AckDeadline:         5,
	})

	if st, ok := status.FromError(err); ok && st.Code() == codes.AlreadyExists {
		// The subscription already exists.
		subscription = client.Subscription(config.Subscription.Name)
	} else if ok && st.Code() == codes.NotFound {
		return nil, fmt.Errorf("topic %q does not exists", config.Topic)
	} else if err != nil {
		return nil, fmt.Errorf("fail to create subscription: %v", err)
	}

	return subscription, nil
}
