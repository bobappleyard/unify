package unify

import (
	"testing"

	"github.com/pkg/errors"
)

func BenchmarkUnification(b *testing.B) {
	for i := 0; i < b.N; i++ {
		doUnifyTest(b)
	}
}

func TestUnification(t *testing.T) {
	doUnifyTest(t)
}

func doUnifyTest(t testing.TB) {
	x, y, z := Var{"x"}, Var{"y"}, Var{"z"}
	fn := func(name string) func(ts ...Term) Term { return func(ts ...Term) Term { return Apply{name, ts} } }
	f := fn("f")
	g := fn("g")
	a, b := fn("a")(), fn("b")()

	for i, testCase := range []struct {
		a, b Term
		res  Subs
		err  error
	}{
		{
			a:   a,
			b:   a,
			res: Subs{},
		},
		{
			a:   a,
			b:   b,
			err: ErrMismatchTerm,
		},
		{
			a:   a,
			b:   x,
			res: Subs{x: a},
		},
		{
			a:   x,
			b:   x,
			res: Subs{},
		},
		{
			a:   x,
			b:   y,
			res: Subs{x: y},
		},
		{
			a:   f(a, x),
			b:   f(a, b),
			res: Subs{x: b},
		},
		{
			a:   f(x),
			b:   g(x),
			err: ErrMismatchTerm,
		},
		{
			a:   f(a),
			b:   f(b),
			err: ErrMismatchTerm,
		},
		{
			a:   f(x),
			b:   f(x, y),
			err: ErrMismatchTerm,
		},
		{
			a:   x,
			b:   f(z),
			res: Subs{x: f(z)},
		},
		{
			a:   f(x),
			b:   f(z),
			res: Subs{x: z},
		},
		{
			a:   f(x, x),
			b:   f(y, z),
			res: Subs{x: y, z: y},
		},
		{
			a:   f(g(x)),
			b:   f(y),
			res: Subs{y: g(x)},
		},
		{
			a:   f(x, g(x)),
			b:   f(a, y),
			res: Subs{x: a, y: g(a)},
		},
		{
			a:   f(g(x), x),
			b:   f(y, a),
			res: Subs{x: a, y: g(a)},
		},
		{
			a:   x,
			b:   f(x),
			err: ErrOccursIn,
		},
	} {
		// fmt.Printf("unifying %s with %s\n", testCase.a, testCase.b)
		res, err := Unify(testCase.a, testCase.b, nil)
		if errors.Cause(err) != testCase.err {
			t.Errorf("[%d] expecting error to be '%v'\n got '%v'", i, testCase.err, err)
			continue
		}
		if !res.matches(testCase.res) {
			t.Errorf("[%d] expecting\n\t%v\ngot\n\t%v", i, testCase.res, res)
		}
	}
}
