package main

import "fmt"

// Mover provides support for moving things.
type Mover interface {
	Move()
	SetState(state string)
}

// Locker provides support for locking and unlocking things.
type Locker interface {
	Lock()
	Unlock()
}

// MoveLocker provides support for moving and locking things.
type MoveLocker interface {
	Mover
	Locker
}

// car represents something you drive.
type car struct{}

// String implements the fmt.Stringer interface.
func (car) String() string {
	return "Vroom!"
}

// cloud represents somewhere you store information.
type cloud struct{}

// String implements the fmt.Stringer interface.
func (cloud) String() string {
	return "Big Data!"
}

// bike represents a concrete type for the example.
type bike struct {
	State string // State of the bike
}

func (b *bike) SetState(state string) {
	b.State = state
	fmt.Println("Setting bike state to:", b.State)
}

// Move can change the position of a bike.
func (bike) Move() {
	fmt.Println("Moving the bike")
}

// Lock prevents a bike from moving.
func (bike) Lock() {
	fmt.Println("Locking the bike")
}

// Unlock allows a bike to be moved.
func (bike) Unlock() {
	fmt.Println("Unlocking the bike")
}

func main() {
	// Declare variables of the MoveLocker and Mover interfaces set to their
	// zero value.
	var ml MoveLocker
	var m Mover

	// Create a value of type bike and assign the value to the MoveLocker
	// interface variable.
	ml = &bike{}

	// An interface value of the type MoveLocker can be implicitly converted into
	// a value of type Mover. They both declare a method named move.
	m = ml

	// Debug/scratch work
	// see what happens to the interface variables when the concrete type is updated
	ml.SetState("new state")

	// prog.go:68: cannot use m (type Mover) as type MoveLocker in assignement:
	// ml = m

	// We can perform a type assertion at runtime to support the assignment

	// Perform a type assertion against the Mover interface to
	if b, ok := m.(*bike); ok {
		ml = b
	}
}
