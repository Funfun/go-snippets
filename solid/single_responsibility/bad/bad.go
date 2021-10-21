package bad

import "database/sql"

/*
	INCORRECT VARIANT WITH "GOD OBJECT" (holds external dependencies and mixed logic)
*/

type Envelope struct {
	message string
	db      *sql.DB
}

func (e *Envelope) Pack(msg string) {
	e.message = msg
}

func (e *Envelope) Unpack() string {
	return e.message
}

func (e *Envelope) SaveToDB() error {

	//SAVING ENVELOPE's MESSAGE TO DATABASE

	return nil
}
