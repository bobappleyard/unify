package unify

import (
	"bytes"
	"fmt"

	"github.com/pkg/errors"
)

// Errors returned by Unify. These will have context using errors.Wrap, so if you need to
// perform logic on them use error.Unwrap.
var (
	ErrOccursIn     = errors.New("var occurs in term")
	ErrMismatchTerm = errors.New("terms are mismatched")
)

// Subs is a set of substitutions.
type Subs map[Var]Term

func (s Subs) String() string {
	var b bytes.Buffer
	b.WriteRune('{')
	first := true
	for v, t := range s {
		if !first {
			b.WriteString(", ")
		}
		first = false
		fmt.Fprintf(&b, "%s -> %s", v, t)
	}
	b.WriteRune('}')
	return b.String()
}

// Unify takes a pair of terms and produces the set of substitutions so as to make source
// the same as target. If no such set exists, an error is returned.
func Unify(source, target Term, s Subs) (Subs, error) {
	if source.matches(target) {
		return nil, nil
	}
	switch source := source.(type) {
	case Var:
		if occursIn(source, target) {
			return nil, wrapErr(ErrOccursIn, source, target)
		}
		return addVar(source, target, s)

	case Apply:
		switch target := target.(type) {
		case Var:
			return Unify(target, source, s)

		case Apply:
			return unifyApplications(source, target, s)
		}
	}
	panic("unreachable")
}

// Eval applies a set of substitutions to a term.
func Eval(t Term, s Subs) Term {
	switch t := t.(type) {
	case Var:
		if v, ok := s[t]; ok {
			return v
		}
		return t

	case Apply:
		args := make([]Term, len(t.Args))
		for i, t := range t.Args {
			args[i] = Eval(t, s)
		}
		return Apply{t.Fn, args}
	}
	panic("unreachable")
}

func unifyApplications(source, target Apply, s Subs) (Subs, error) {
	if source.Fn != target.Fn {
		return nil, wrapErr(ErrMismatchTerm, source, target)
	}
	if len(source.Args) != len(target.Args) {
		return nil, wrapErr(ErrMismatchTerm, source, target)
	}
	merged := s
	for i := range source.Args {
		subs, err := Unify(source.Args[i], target.Args[i], merged)
		if err != nil {
			return nil, wrapErr(err, source, target)
		}
		merged = subs
	}
	return merged, nil
}

func wrapErr(err error, source, target Term) error {
	return errors.Wrapf(err, "while unifying %v with %v", source, target)
}

func occursIn(x Var, t Term) bool {
	switch t := t.(type) {
	case Var:
		return x == t

	case Apply:
		for _, t := range t.Args {
			if occursIn(x, t) {
				return true
			}
		}
		return false
	}
	panic("unreachable")
}

func addVar(v Var, t Term, s Subs) (Subs, error) {
	// fmt.Printf("\tadding %s -> %s to %s\n", v, t, s)
	if sub, ok := s[v]; ok {
		next, err := Unify(t, sub, s)
		if err != nil {
			return nil, err
		}
		return next, nil
	}
	res := Subs{}
	res[v] = Eval(t, s)
	for w, t := range s {
		if occursIn(v, t) {
			t = Eval(t, res)
		}
		res[w] = t
	}
	return res, nil
}

func (s Subs) matches(t Subs) bool {
	if len(s) != len(t) {
		return false
	}
	for v, u := range s {
		other, ok := t[v]
		if !ok {
			return false
		}
		if !other.matches(u) {
			return false
		}
	}
	return true
}
