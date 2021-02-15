package main

import (
	"fmt"
	"log"

	tcell "github.com/gdamore/tcell/v2"
)

func box(s tcell.Screen, y, x, h, w int) {
	for i := 0; i < w; i++ {
		s.SetContent(x+i, y, tcell.RuneHLine, nil, tcell.StyleDefault)
		s.SetContent(x+i, y+h, tcell.RuneHLine, nil, tcell.StyleDefault)
	}
	for j := 0; j < h; j++ {
		s.SetContent(x, y+j, tcell.RuneVLine, nil, tcell.StyleDefault)
		s.SetContent(x+w, y+j, tcell.RuneVLine, nil, tcell.StyleDefault)
	}
	s.SetContent(x, y, tcell.RuneULCorner, nil, tcell.StyleDefault)
	s.SetContent(x, y+h, tcell.RuneLLCorner, nil, tcell.StyleDefault)
	s.SetContent(x+w, y, tcell.RuneURCorner, nil, tcell.StyleDefault)
	s.SetContent(x+w, y+h, tcell.RuneLRCorner, nil, tcell.StyleDefault)
	// for xx := x; xx < x+w; xx++ {
	// 	for yy := y; yy < y+h; yy++ {
	// 		s.SetContent(xx, yy, tcell.RuneHLine, nil, tcell.StyleDefault)
	// 	}
	// }
	s.Sync()
}

func main() {
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatal(err)
	}
	if err := s.Init(); err != nil {
		log.Fatal(err)
	}
	defer s.Fini()

	box(s, 0, 0, 25, 40)
	for {
		e := s.PollEvent()
		if ek, ok := e.(*tcell.EventKey); ok {
			fmt.Printf("%+v\n", ek)
			if ek.Key() == tcell.KeyEsc {
				break
			}
		}
	}
}
