package main

import (
	"fmt"
	"strings"
)

type Attributes struct {
	Str int
	Agi int
	Vig int
	Int int
}

func newAttributes(str, agi, vig, int int) *Attributes {
	return &Attributes{
		Str: str,
		Agi: agi,
		Vig: vig,
		Int: int,
	}
}

func (a Attributes) String() []string {
	text := fmt.
		Sprintf("STR: %d\nAGI: %d\nVIG: %d\nINT: %d", a.Str, a.Agi, a.Vig, a.Int)
	return strings.Split(text, "\n")
}
