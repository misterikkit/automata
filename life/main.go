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
	rand.Seed(*seed)

	// Initialize gene
	var gn gene.Gene
	if len(*geneStr) == 0 {
		gn = gene.Random()
	} else {
		gn = gene.FromString(*geneStr)
	}

	// Initialize Game
	g := game.New(40, 100)
	g = g.Next(game.Random)

	runInteractive(context.Background(), &g, gn)
	fmt.Printf("Final state:\n%v", tui.Fmt(g))
	fmt.Printf("Gene: %+v\n", gn)
	fmt.Printf("Seed: %v\n", seed)
}

func runInteractive(ctx context.Context, g *game.Game, gn gene.Gene) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Initialize ui, and wire events
	pauseCh := make(chan struct{}, 1)
	stepCh := make(chan struct{}, 1)
	t, err := tui.New(fmt.Sprintf("Gene %v\tEsc to exit", gn), func(e tui.Event) {
		switch e {
		case tui.Escape:
			cancel()
		case tui.Enter:
			pauseCh <- struct{}{}
		case tui.Right:
			stepCh <- struct{}{}
		}
	})
	if err != nil {
		log.Fatal(err)
	}
	defer t.Close()
	t.DrawGame(*g) // draw initial state

	rule := gn.AsRule()

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
			*g = g.Next(rule)
			t.DrawGame(*g)

		case <-pauseCh:
			paused = !paused
		case <-stepCh:
			*g = g.Next(rule)
			t.DrawGame(*g)

		case <-ctx.Done():
			break loop
		}
	}
}
