package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())

	m := NewMaze(5, 5)
	// fmt.Print(m)

	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	init := New("init", func(Event, *Object) { cancel() })

	log.Printf("%-20v %-14q -> %-20v (%v)", "sender", "event", "recipient", "param")

	// make sure all _wire events are done
	time.AfterFunc(500*time.Millisecond, func() { m.cells[0][0].cell.Send(init, "visit", init) })
	// start the engine
	m.Run(ctx, init)

	fmt.Printf("Final state:\n%v\n", m)
}
