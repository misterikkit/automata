# Fun with Maze Generating Automata

Inspired by [this paper][1] I wrote a set of experiments using cellular
automata to generate mazes. Before I got as far as writing a fitness function,
I found some [useful "genes"][2] on Wikipedia.

## Try it out

To view an interactive simulation with a random ruleset, run

```
go run ./life
```

and press Enter to begin.

To run a fixed number of generations and print the final state, run

```
go run ./life -gens 50
```

[Maze][3] ruleset

```
go run ./life -gene 0001000000111110000 -gens 100
```

[Mazectric][3] ruleset

```
go run ./life -gene 0001000000111110000 -gens 100
```

I discovered this ruleset which generates a nice style of maze.

```
go run ./life -gene 011100110001000000 -gens 150
```

TODO: add screenshots

[1]: https://scholarworks.unr.edu/bitstream/handle/11714/3433/Adams_unr_0139M_12635.pdf?sequence=1&isAllowed=y
[2]: https://en.wikipedia.org/wiki/Maze_generation_algorithm#Cellular_automaton_algorithms
[3]: https://www.conwaylife.com/wiki/OCA:Maze
