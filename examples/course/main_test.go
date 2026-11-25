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
	want := "Next client: ambiguous (two immediate matches)\n" +
		"RequireOneRoute: tool -> client (Client B) -> version\n" +
		"AllRoutes: 2\n" +
		"  tool -> client (Client A)\n" +
		"  tool -> client (Client B)\n" +
		"FirstRoute: tool -> client (Client A)\n"
	if !bytes.Equal(output, []byte(want)) {
		t.Fatalf("unexpected output:\n%s", output)
	}
}
