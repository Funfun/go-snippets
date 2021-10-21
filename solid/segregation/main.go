package main

// bad
type BadStorage interface {
	Query()
	Send()
}

// good
type GoodStorage interface {
	Querable
	Sender
}

type Querable interface {
	Query()
}
type Sender interface {
	Sender()
}
