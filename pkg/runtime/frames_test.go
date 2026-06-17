package runtime

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCallerNameCanReturnTheNameOfThisFunction(t *testing.T) {
	want := "runtime.TestCallerNameCanReturnTheNameOfThisFunction"
	got := CallerName(0)
	assert.Equal(t, got, want)
}

func level1(f func()) { f() }
func level2(f func()) { level1(f) }
func level3(f func()) { level2(f) }

func TestCallerNameCanSkipSomeFunctions(t *testing.T) {
	cases := []struct {
		fn   func(func())
		skip uint
		want string
	}{
		{fn: level1, skip: 1, want: "runtime.level1"},

		{fn: level2, skip: 1, want: "runtime.level1"},
		{fn: level2, skip: 2, want: "runtime.level2"},

		{fn: level3, skip: 1, want: "runtime.level1"},
		{fn: level3, skip: 2, want: "runtime.level2"},
		{fn: level3, skip: 3, want: "runtime.level3"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(fmt.Sprintf("%s/skip:%d", tc.want, tc.skip), func(t *testing.T) {
			tc.fn(func() {
				got := CallerName(tc.skip)
				assert.Equal(t, got, tc.want)
			})
		})
	}
}

func next(f func()) { f() }

func TestNext(t *testing.T) {
	want := "runtime.TestNext"

	next(func() {
		frame := Frame(1)
		next := Next(&frame)
		got := FrameFunctionName(*next)
		assert.Equal(t, got, want)
	})
}

func globalFuncNameTest() {}

type funcNameTest struct{}

func (fn funcNameTest) ExportedMethod() {}
func (fn funcNameTest) method()         {}

func TestFuncNameSuccess(t *testing.T) {
	var fn funcNameTest
	anonFunc := func() {}

	for name, tc := range map[string]struct {
		fn   interface{}
		want string
	}{
		"global":               {fn: globalFuncNameTest, want: "runtime.globalFuncNameTest"},
		"local":                {fn: anonFunc, want: "runtime.TestFuncNameSuccess.func1"},
		"anonymous":            {fn: func() {}, want: "runtime.TestFuncNameSuccess.func2"},
		"free method":          {fn: funcNameTest.method, want: "runtime.funcNameTest.method"},
		"bind method":          {fn: fn.method, want: "runtime.funcNameTest.method-fm"},
		"free exported method": {fn: funcNameTest.ExportedMethod, want: "runtime.funcNameTest.ExportedMethod"},
		"bind exported method": {fn: fn.ExportedMethod, want: "runtime.funcNameTest.ExportedMethod-fm"},
	} {
		tc := tc
		t.Run(name, func(t *testing.T) {
			got := FuncName(tc.fn)
			assert.Equal(t, got, tc.want)
		})
	}
}

func TestFuncNameFailsWithNonFunction(t *testing.T) {
	defer func() {
		e := recover()
		assert.Equal(t, e == nil, false)
	}()

	FuncName("a string is not a function")

	t.Fatal("unreachable code")
}
