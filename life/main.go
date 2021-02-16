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
	genCount := flag.Int("gens", 0, "if >0, run a fixed number of generations and exit")
	life := flag.Bool("life", false, "Use game of life rules")
	density := flag.Float64("density", 0.01, "Probability of initial cell state being alive")
	init := flag.String("initialize", "random", "Type of initial state to use. One of (random, smallrandom)")

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
	if *life {
		gn = gene.FromString("001100000000100000") // B23/S3
	}

	// Initialize Game
	g := game.New(40, 100)
	switch *init {
	case "random":
		g = g.Next(game.RandomSparse(float32(*density)))
	case "smallrandom":
		newG := g.Next(game.Random)
		for i := 0; i < 10; i++ {
			for j := 0; j < 10; j++ {
				g[i][j] = newG[i][j]
			}
		}
	}

	if *genCount > 0 {
		runAuto(&g, gn, *genCount)
	} else {
		runInteractive(context.Background(), &g, gn)
	}
	fmt.Printf("Final state:\n%v", tui.Fmt(g))
	fmt.Printf("Gene: %+v\n", gn)
	fmt.Printf("Seed: %v\n", seed)
}

func runAuto(g *game.Game, gn gene.Gene, n int) {
	rule := gn.AsRule()
	for i := 0; i < n; i++ {
		*g = g.Next(rule)
	}
}

func runInteractive(ctx context.Context, g *game.Game, gn gene.Gene) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Initialize ui, and wire events
	pauseCh := make(chan struct{}, 1)
	stepCh := make(chan struct{}, 1)
	t, err := tui.New(fmt.Sprintf("Gene %v\tEnter:play/pause\t➡:step\tEsc:exit", gn), func(e tui.Event) {
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
	paused := true
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
