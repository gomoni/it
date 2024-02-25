package it_test

import (
	"fmt"

	"github.com/gomoni/it"
)

type pusher struct {
	stack chan string
}

func (y *pusher) push(s string) {
	y.stack <- s
}

func (y pusher) seq() func(func(string) bool) {
	return func(yield func(string) bool) {
		for {
			select {
			case s, open := <-y.stack:
				if !open || !yield(s) {
					return
				}
			}
		}
	}
}

func (y pusher) wait() {
	<-y.stack
}

func Example_break_da_chain() {
	n := []string{"aa", "aaa", "aaaaaaa", "a"}

	// create a method chain
	chain := it.NewChain(it.From(n)).
		Filter(func(s string) bool { return true })

	// break it - with some syntax sugar
	p := pusher{stack: make(chan string)}
	defer p.wait()
	go func() {
		defer close(p.stack)
		for s := range chain.Seq() {
			p.push(s)
		}
	}()

	// continue here
	chain2 := it.NewChain(p.seq()).
		Filter(func(s string) bool { return len(s) > 2 })
	slice := chain2.Slice()
	fmt.Println(slice)
	// Output: [aaa aaaaaaa]
}
