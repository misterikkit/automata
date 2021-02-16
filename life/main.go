package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/misterikkit/automata/life/game"
	"github.com/misterikkit/automata/life/gene"
	"github.com/misterikkit/automata/life/tui"
)

func main() {
	defer fmt.Println("bye!")
	seed := time.Now().Unix()
	defer fmt.Printf("Seed: %v\n", seed)
	rand.Seed(seed)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var gn gene.Gene
	if len(os.Args) > 1 {
		gn = gene.FromString(os.Args[1])
	} else {
		gn = gene.Random()
	}
	defer fmt.Printf("Gene: %+v\n", gn)
	rule := gn.AsRule()

	g := game.New(20, 30)
	g = g.Next(game.Random)

	t, err := tui.New(fmt.Sprintf("Gene %v\tEsc to exit", gn), func(e tui.Event) {
		switch e {
		case tui.Escape:
			cancel()
		case tui.Enter:
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
