package good

import "fmt"

// Liskov principle:
// composition struct that satisfies simple interface should do same as parent struct
// usage: use interface for arguments and satisfy interface

type Vehicle interface {
	DriveIn()
}
type Car struct{}

func (*Sedan) DriveIn() {
	fmt.Println("drive in")
}

type Sedan struct {
	Car
}

func NewGarage(vehicle Vehicle) {
	vehicle.DriveIn()
}
