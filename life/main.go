package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/misterikkit/automata/life/game"
	"github.com/misterikkit/automata/life/gene"
	"github.com/misterikkit/automata/life/tui"
)

func main() {
	defer fmt.Println("bye!")
	rand.Seed(time.Now().Unix())
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	g := game.New(20, 30)
	gene := gene.Random()
	rule := gene.AsRule()
	defer fmt.Printf("Gene: %+v\n", gene)
	// g = g.Next(game.Random)

	t, err := tui.New(fmt.Sprintf("Gene %v\tEsc to exit", gene), func(e tui.Event) {
		if e == tui.Escape {
			cancel()
		}
	})
	if err != nil {
		log.Fatal(err)
	}
	defer t.Close()
	t.DrawGame(g)

	tick := time.NewTicker(100 * time.Millisecond)
	defer tick.Stop()

loop:
	for {
		select {
		case <-tick.C:
			g = g.Next(rule)
			t.DrawGame(g)
		case <-ctx.Done():
			break loop
		}
	}

}
