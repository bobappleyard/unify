package unify

import (
	"fmt"
)

// A Term is the subject of unification. There are two types, Var and Apply.
type Term interface {
	term()
	Matches(Term) bool
}

// Var represents a metavariable in a logical statement.
type Var struct {
	Of interface{}
}

// Apply represents the application of a function to some arguments.
type Apply struct {
	Fn   interface{}
	Args []Term
}

func (v Var) term()   {}
func (a Apply) term() {}

func (v Var) String() string {
	return fmt.Sprintf("%v", v.Of)
}

func (a Apply) Format(s fmt.State, c rune) {
	fmt.Fprintf(s, "%v(", a.Fn)
	for i, arg := range a.Args {
		if i > 0 {
			fmt.Fprint(s, ", ")
		}
		fmt.Fprint(s, arg)
	}
	fmt.Fprint(s, ")")
}

func (v Var) Matches(t Term) bool {
	if t, ok := t.(Var); ok {
		return v == t
	}
	return false
}

func (a Apply) Matches(t Term) bool {
	if t, ok := t.(Apply); ok {
		if a.Fn != t.Fn {
			return false
		}
		if len(a.Args) != len(t.Args) {
			return false
		}
		for i, x := range a.Args {
			if !t.Args[i].Matches(x) {
				return false
			}
		}
		return true
	}
	return false
}
