package cmd

import (
	"os"

	"github.com/Funfun/go-snippets/go-elk/go-pubsub-elk/beater"
	"github.com/elastic/beats/v7/libbeat/cfgfile"
	"github.com/elastic/beats/v7/libbeat/cmd"
	"github.com/elastic/beats/v7/libbeat/cmd/instance"
	"github.com/elastic/beats/v7/libbeat/common"
)

// Name of this beat
var Name = "pubsubbeat"
var configOverrides = common.MustNewConfigFrom(map[string]interface{}{
	"pubsubbeat": map[string]interface{}{
		"project_id":                         "set-me",
		"topic":                              "my-indexs",
		"subscription.name":                  "my-indexs-pubsubbeat",
		"subscription.create":                false,
		"subscription.retain_acked_messages": false,
		"subscription.retention_duration":    "168h",
		"subscription.connection_pool_size":  2,
		"json.enabled":                       true,
		"json.add_error_key":                 true,
	},

	"logging.level": "info",
	"logging.selectors": []interface{}{
		"*",
	},
	"http.enabled": true,
	"http.host":    "0.0.0.0",
	"http.port":    5066,

	"processors": []interface{}{
		map[string]interface{}{
			"drop_fields": map[string]interface{}{
				"fields": []interface{}{
					"ecs",
					"host.name",
					"agent.ephemeral_id",
					"agent.hostname",
					"agent.type",
					"agent.id",
				},
				"ignore_missing": true,
			},
		},
	},

	"output.elasticsearch": map[string]interface{}{
		"compression_level": 9,
		"index":             "my-indexs-%{[agent.version]}-%{+yyyy.MM.dd}",
		"worker":            2,
		"max_retries":       5,
		"bulk_max_size":     3200,
		"timeout":           30,
		"hosts":             []string{"localhost:9200"},
		"username":          "elastic",
		"password":          "changeme",
	},

	"setup.ilm.enabled":        true,
	"setup.ilm.rollover_alias": "my-indexs",
	"setup.ilm.policy_name":    "my-indexs_policy",

	"setup.template.enabled":   true,
	"setup.template.type":      "index",
	"setup.template.name":      "my-index",
	"setup.template.pattern":   "my-index-%{[agent.version]}-*",
	"setup.template.fields":    "${path.config}/fields.yml",
	"setup.template.overwrite": true,
	"setup.template.settings": map[string]interface{}{
		"index.number_of_shards":                         1,
		"index.number_of_replicas":                       1,
		"index.refresh_interval":                         "250s",
		"index.routing.allocation.total_shards_per_node": 8, // the correct to go 20 shards to total 8 GB
		"index.codec":                                    "best_compression",
		"_source.enabled":                                true,
	},
})

var elasticCredentialsOverride = common.MustNewConfigFrom(map[string]interface{}{
	"cloud.id":   os.Getenv("ELASTIC_CLOUD_ID"),
	"cloud.auth": os.Getenv("ELASTIC_CLOUD_AUTH"),
})

var pubsubbeatCfgOverride = common.MustNewConfigFrom(map[string]interface{}{
	"pubsubbeat.project_id":        os.Getenv("PUBSUBBEAT_PROJECT_ID"),
	"pubsubbeat.topic":             os.Getenv("PUBSUBBEAT_TOPIC"),
	"pubsubbeat.subscription.name": os.Getenv("PUBSUBBEAT_SUBSCRIPTION"),
})

var loglevelCfgOverride = common.MustNewConfigFrom(map[string]interface{}{
	"logging.level": os.Getenv("LOGGING_LEVEL"),
})

var prodPathHome = common.MustNewConfigFrom(map[string]interface{}{
	"path.home": os.Getenv("PATH_HOME"),
})

var settings = instance.Settings{
	Name: Name,
	ConfigOverrides: []cfgfile.ConditionalOverride{
		{
			Check:  always,
			Config: configOverrides,
		},
		{
			Check:  isElasticCloud,
			Config: elasticCredentialsOverride,
		},
		{
			Check:  isOverridePubsubCfg,
			Config: pubsubbeatCfgOverride,
		},
		{
			Check:  isOverrideLoglevel,
			Config: loglevelCfgOverride,
		},
		{
			Check:  isOverridePath,
			Config: prodPathHome,
		},
	},
}

// RootCmd to handle beats cli
var RootCmd = cmd.GenRootCmdWithSettings(beater.New, settings)

var always = func(_ *common.Config) bool {
	return true
}

var isElasticCloud = func(_ *common.Config) bool {
	return os.Getenv("ELASTIC_CLOUD_ID") != "" && os.Getenv("ELASTIC_CLOUD_AUTH") != ""
}

var isOverridePubsubCfg = func(_ *common.Config) bool {
	return os.Getenv("PUBSUBBEAT_PROJECT_ID") != "" && os.Getenv("PUBSUBBEAT_TOPIC") != "" && os.Getenv("PUBSUBBEAT_SUBSCRIPTION") != ""
}

var isOverrideLoglevel = func(_ *common.Config) bool {
	return os.Getenv("LOGGING_LEVEL") != ""
}

var isOverridePath = func(_ *common.Config) bool {
	return os.Getenv("PATH_HOME") != ""
}
