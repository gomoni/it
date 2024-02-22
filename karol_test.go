package it_test

import (
	// "fmt"
	"testing"

	"github.com/gomoni/it"
)

func TestKarolFrom1(t *testing.T) {
	nn := []int{1, 5, 3, 9, 7, 2, 4, 8, 10, 0}
	fmod2 := func(i int) bool {
		return i%2 == 0
	}

	// Let's have one Filter operation
	{
		s0 := it.From(nn)
		s1 := it.Filter(s0, fmod2)
		slice := it.Slice(s1)
		t.Log(slice)
	}
	{
		// This might be easier to use for just one Filter operation
		slice := it.SimpleFilter(nn, fmod2)
		t.Log(slice)
	}
}

func TestKarolFrom2(t *testing.T) {
	nn := []int{1, 5, 3, 9, 7, 2, 4, 8, 10, 0}
	fNotModulo2 := func(i int) bool {
		return i%2 != 0
	}
	fGreaterThanPrevious := func() func(i int) bool {
		var previous *int
		return func(i int) bool {
			if previous == nil {
				previous = &i
				return true
			}
			if i > *previous {
				previous = &i
				return true
			}
			previous = &i
			return false
		}
	}

	// What about 2+ consecutive Filter operations?
	{
		s0 := it.From(nn)
		s1 := it.Filter(s0, fNotModulo2)
		s2 := it.Filter(s1, fGreaterThanPrevious())
		slice := it.Slice(s2)
		t.Log(slice)
	}
	{
		// Is this better? ¯\_(ツ)_/¯
        // For 2 maybe, for >2 definitely not ... /me thinks
		slice := it.SimpleFilter(
			it.SimpleFilter(nn, fNotModulo2),
			fGreaterThanPrevious(),
		)
		t.Log(slice)
	}
	{
		// Maybe this might be better?
		slice := it.SimpleFilters(nn, fNotModulo2, fGreaterThanPrevious())
		t.Log(slice)
	}

}

