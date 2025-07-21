package main

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"time"
)

// ================================================= Data we have to move

// Data is the structure of the data we are copying
type Data struct {
	Line string
}

// ================================================= Discovered abstractions

// Puller declares behavior for pulling data.
//
// Writes to d
type Puller interface {
	Pull(d *Data) error
}

// Storer declares behavior for storing data.
//
// Reads from d
type Storer interface {
	Store(d *Data) error
}

/* removed because it is not used
// PullStorer declares behavior for both pulling and storing.
type PullStorer interface {
	Puller
	Storer
}
*/
// ================================================= Stateful primitives

// Xenia is a system we need to pull data from.
type Xenia struct {
	Host    string
	Timeout time.Duration
}

// Pull knows how to pull data out of Xenia.
func (x *Xenia) Pull(d *Data) error {
	switch rand.Intn(10) {
	case 1, 9:
		return io.EOF

	case 5:
		return errors.New("error reading data from Xenia")

	default:
		d.Line = "Data"
		fmt.Println("[x.Pull] In:", d.Line)
		return nil
	}
}

// Pillar is the system we need to store data into.
type Pillar struct {
	Host    string
	Timeout time.Duration
}

// Store knows how to store data into Pillar.
func (p *Pillar) Store(d *Data) error {
	fmt.Println("[p.Store] Out:", d.Line)
	return nil
}

// ================================================= Combined structs
/* removed because it is not used
// System wraps Xenia and Pillar together into a single system.
type System struct {
	Puller
	Storer
}
*/

// =================================================

// pull knows how to pull bulks of data from any Puller.
func pull(p Puller, data []Data) (int, error) {
	for i := range data {
		if err := p.Pull(&data[i]); err != nil {
			return i, err
		}
	}

	return len(data), nil
}

// store knows how to store bulks of data into any Storer.
func store(s Storer, data []Data) (int, error) {
	for i := range data {
		if err := s.Store(&data[i]); err != nil {
			return i, err
		}
	}

	return len(data), nil
}

// Copy knows how to pull and store data from the System
func Copy(p Puller, s Storer, batch int) error {
	data := make([]Data, batch)

	for {
		i, err := pull(p, data)
		if i > 0 {
			if _, err := store(s, data[:i]); err != nil {
				return err
			}
		}

		if err != nil {
			return err
		}
	}
}

// =================================================

func main() {
	x := &Xenia{
		Host:    "localhost:8000",
		Timeout: time.Second,
	}
	p := &Pillar{
		Host:    "localhost:9000",
		Timeout: time.Second,
	}

	if err := Copy(x, p, 5); err != io.EOF {
		fmt.Println(err)
	}
}
