package main

import (
	"reflect"
	"testing"
)

func TestDemonstrate(t *testing.T) {
	got, err := demonstrate()
	if err != nil {
		t.Fatal(err)
	}
	if got.typeName != "main.chart" || got.fieldName != "first" || !reflect.DeepEqual(got.path, []int{1, 0}) {
		t.Errorf("location = type %q, field %q, path %v", got.typeName, got.fieldName, got.path)
	}
	if got.before != "before" || got.after != "after" || got.sibling != "sibling" {
		t.Errorf("live field result = before %q, after %q, sibling %q", got.before, got.after, got.sibling)
	}
	if !got.selectedAddress || got.ignored != "untouched" {
		t.Errorf("resolved address selected field = %t, ignored field = %q", got.selectedAddress, got.ignored)
	}
}
