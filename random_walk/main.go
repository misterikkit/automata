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

	// TODO: switch to wall.NewMaze
	m := NewMaze(*h, *w)

	log.Printf("%-20v %-14q -> %-20v (%v)", "sender", "event", "recipient", "param")

	// start the engine
	start := time.Now()
	m.Run(context.Background())
	end := time.Now()

	fmt.Printf("Generated %dx%d maze in %v\n", *h, *w, end.Sub(start))
	fmt.Println(m)
}
