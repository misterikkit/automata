package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"time"
)

func main() {
	h := flag.Int("h", 5, "height")
	w := flag.Int("w", 5, "width")
	verbose := flag.Bool("v", false, "enable logging")
	flag.Parse()
	if !*verbose {
		log.SetOutput(io.Discard) // io.Discard is new in go1.16
	}
	rand.Seed(time.Now().Unix())

	m := NewMaze(*h, *w)
	// fmt.Print(m)

	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	init := New("init", func(Event, *Object) { cancel() })

	log.Printf("%-20v %-14q -> %-20v (%v)", "sender", "event", "recipient", "param")

	// make sure all _wire events are done before starting the random walk
	go func() {
		wireWG.Wait()
		m.cells[0][0].cell.Send(init, "visit", init)
	}()
	// start the engine
	start := time.Now()
	count := m.Run(ctx, init)
	end := time.Now()

	fmt.Printf("Generated %dx%d maze in %v\n", *h, *w, end.Sub(start))
	fmt.Println(m)
	fmt.Printf("And it only took %d threads!\n", count)
}
