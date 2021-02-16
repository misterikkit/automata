package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/misterikkit/automata/life/game"
	"github.com/misterikkit/automata/life/gene"
	"github.com/misterikkit/automata/life/tui"
)

func main() {
	seed := flag.Int64("seed", 0, "random seed")
	geneStr := flag.String("gene", "", "automata gene")
	flag.Parse()
	// Set random seed
	if *seed == 0 {
		*seed = time.Now().Unix()
	}
	defer fmt.Printf("Seed: %v\n", seed)
	rand.Seed(*seed)

	// Create game context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize gene
	var gn gene.Gene
	if len(*geneStr) == 0 {
		gn = gene.Random()
	} else {
		gn = gene.FromString(*geneStr)
	}
	defer fmt.Printf("Gene: %+v\n", gn)
	rule := gn.AsRule()

	// Initialize Game
	g := game.New(20, 30)
	g = g.Next(game.Random)

	// Things that print after the TUI closes need to be deferred before it opens
	defer func() {
		fmt.Printf("Final state:\n%v", tui.Fmt(g))
	}()

	// Initialize ui, and wire events
	pauseCh := make(chan struct{}, 1)
	t, err := tui.New(fmt.Sprintf("Gene %v\tEsc to exit", gn), func(e tui.Event) {
		switch e {
		case tui.Escape:
			cancel()
		case tui.Enter:
			pauseCh <- struct{}{}
		}
	})
	if err != nil {
		log.Fatal(err)
	}
	defer t.Close()
	t.DrawGame(g) // draw initial state

	// Run the game loop
	tick := time.NewTicker(100 * time.Millisecond)
	defer tick.Stop()
	paused := false
loop:
	for {
		select {
		case <-tick.C:
			if paused {
				break
			}
			g = g.Next(rule)
			t.DrawGame(g)

		case <-pauseCh:
			paused = !paused

		case <-ctx.Done():
			break loop
		}
	}
}
