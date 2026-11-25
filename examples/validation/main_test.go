package main

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestRun(t *testing.T) {
	original := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = writer
	runErr := run()
	_ = writer.Close()
	os.Stdout = original
	output, readErr := io.ReadAll(reader)
	_ = reader.Close()
	if runErr != nil {
		t.Fatal(runErr)
	}
	if readErr != nil {
		t.Fatal(readErr)
	}
	want := "Compile: rejected duplicate lane IDs\n" +
		"Build: accepted Chart with required option\n" +
		"Build: rejected missing required option\n" +
		"Prepare: rejected waypoint after terminal Landmark\n"
	if !bytes.Equal(output, []byte(want)) {
		t.Fatalf("unexpected output:\n%s", output)
	}
}
