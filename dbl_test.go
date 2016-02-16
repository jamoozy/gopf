package main

import (
	"testing"
)

func TestDbVarSet(t *testing.T) {
	dv.Set("foo")
	if dv.String() != "foo" {
		t.Fatalf(`Exp: "foo"  Got: "%s"`, dv.String())
	}

	dv.Set("bar")
	if dv.String() != "bar" {
		t.Fatalf(`Exp: "bar"  Got: "%s"`, dv.String())
	}
}

func TestSingleStringRowParser(t *testing.T) {
	t.Fatalf("not impl")
}
