package main

import "context"

type Event struct {
	Name string
	// TODO args
}

type Object struct {
	listeners []*Object
	events    chan Event
	script    func(Event)
}

func New(script func(Event)) *Object { return &Object{events: make(chan Event, 10), script: script} }

func (o *Object) Run(ctx context.Context) {
	for {
		select {
		case e := <-o.events:
			o.script(e) // TODO: per-object lock?
			for _, l := range o.listeners {
				l.events <- e
			}
		case <-ctx.Done():
			return
		}
	}
}

func (o *Object) Listen(target *Object) { o.listeners = append(o.listeners, target) }

func main() {

}
