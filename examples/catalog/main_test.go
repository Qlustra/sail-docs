package main

import (
	"errors"
	"reflect"
	"testing"

	"github.com/qlustra/sail"
)

var readCalls int

type trackedText struct{ Value string }

func (t trackedText) Read() string  { readCalls++; return t.Value }
func (trackedText) Summary() string { return "declaration summary" }

type testWorkspace struct {
	Before trackedText `catalog:"document,summary" name:"before"`
	Folder struct {
		Inside trackedText `catalog:"document" name:"inside"`
		Nested struct {
			Deep trackedText `catalog:"document" name:"deep"`
		} `catalog:"folder" name:"nested"`
	} `catalog:"folder" name:"dir"`
	After trackedText `catalog:"document" name:"after"`
}

func TestPreviewDeclarationOrderAndSelectedInstance(t *testing.T) {
	c, err := New(testWorkspace{})
	if err != nil {
		t.Fatal(err)
	}
	got, err := c.Preview()
	if err != nil {
		t.Fatal(err)
	}
	want := []PreviewEntry{{"before", "declaration summary", true}, {"dir/inside", "", false}, {"dir/nested/deep", "", false}, {"after", "", false}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Preview() = %#v, want %#v", got, want)
	}

	var a, b testWorkspace
	a.Folder.Inside.Value = "first instance"
	b.Folder.Inside.Value = "chosen instance"
	readCalls = 0
	value, err := c.Read(&b, "dir", "inside")
	if err != nil || value != "chosen instance" || readCalls != 1 {
		t.Fatalf("Read(b) = %q, %v; calls %d", value, err, readCalls)
	}
	value, err = c.Read(&a, "dir", "inside")
	if err != nil || value != "first instance" || readCalls != 2 {
		t.Fatalf("Read(a) = %q, %v; calls %d", value, err, readCalls)
	}
}

type ambiguousWorkspace struct {
	A trackedText `catalog:"document" name:"same"`
	B trackedText `catalog:"document" name:"same"`
}

func TestMissingAmbiguousAndWrongWorkspaceHaveNoReadEffects(t *testing.T) {
	c, err := New(ambiguousWorkspace{})
	if err != nil {
		t.Fatal(err)
	}
	readCalls = 0
	if _, err := c.Read(&ambiguousWorkspace{}, "missing"); !errors.Is(err, sail.ErrNoRoute) {
		t.Fatalf("missing route error = %v", err)
	}
	if _, err := c.Read(&ambiguousWorkspace{}, "same"); !errors.Is(err, sail.ErrAmbiguousRoute) {
		t.Fatalf("ambiguous route error = %v", err)
	}
	if readCalls != 0 {
		t.Fatalf("missing/ambiguous route triggered %d reads", readCalls)
	}

	unique, err := New(testWorkspace{})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := unique.Read(&ambiguousWorkspace{}, "before"); err == nil {
		t.Fatal("wrong workspace type accepted")
	}
	if readCalls != 0 {
		t.Fatalf("wrong workspace triggered %d reads", readCalls)
	}
}

func TestEmptyRouteRejected(t *testing.T) {
	c, err := New(testWorkspace{})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := c.Read(&testWorkspace{}); err == nil {
		t.Fatal("empty route accepted")
	}
}
