package unify

import (
	"bytes"
	"fmt"
)

// A Term is the subject of unification. There are two types, Var and Apply.
type Term interface {
	term()
	matches(Term) bool
}

// Var represents a metavariable in a logical statement.
type Var struct {
	Name string
}

// Apply represents the application of a function to some arguments.
type Apply struct {
	Fn   string
	Args []Term
}

func (v Var) term()   {}
func (a Apply) term() {}

func (v Var) String() string {
	return fmt.Sprintf("%s", v.Name)
}

func (a Apply) String() string {
	var b bytes.Buffer
	b.WriteString(a.Fn)
	b.WriteRune('(')
	for i, arg := range a.Args {
		if i > 0 {
			b.WriteString(", ")
		}
		fmt.Fprint(&b, arg)
	}
	b.WriteRune(')')
	return b.String()
}

func (v Var) matches(t Term) bool {
	if t, ok := t.(Var); ok {
		return v == t
	}
	return false
}

func (a Apply) matches(t Term) bool {
	if t, ok := t.(Apply); ok {
		if a.Fn != t.Fn {
			return false
		}
		if len(a.Args) != len(t.Args) {
			return false
		}
		for i, x := range a.Args {
			if !t.Args[i].matches(x) {
				return false
			}
		}
		return true
	}
	return false
}
