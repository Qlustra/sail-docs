package main

import (
	"bytes"
	"slices"
	"testing"

	"github.com/qlustra/sail"
)

func TestPopulationExample(t *testing.T) {
	var out bytes.Buffer
	if err := showPopulation(&out); err != nil {
		t.Fatal(err)
	}
	want := "Service: reporter\n  scalar: port\n  enum: format\n    text (declared default: true)\n    json (declared default: false)\nBuild:    Render reports (port on declaration receiver: 0)\nPopulate: Render reports (port on declaration receiver: 0)\n"
	if got := out.String(); got != want {
		t.Fatalf("output:\n%s\nwant:\n%s", got, want)
	}
}

type ExtraSettings struct {
	Promoted string `service:"setting" name:"promoted"`
}

func TestPromotedSettingsAndReuseAfterFailure(t *testing.T) {
	nav, err := sail.Compile[serviceLegend]()
	if err != nil {
		t.Fatal(err)
	}
	type validChart struct {
		Server struct {
			First         string `service:"setting" name:"first"`
			ExtraSettings `service:"setting" name:"ignored embedding"`
			Unknown       string `service:"other,unrecognized"`
			Last          string `service:"setting" name:"last"`
		} `service:"service" name:"reporter"`
	}
	check := func(m *serviceLegend) {
		t.Helper()
		var names []string
		for _, setting := range m.Entries.Service.Settings.Scalars {
			names = append(names, setting.Name.Get())
		}
		if want := []string{"first", "promoted", "last"}; !slices.Equal(names, want) {
			t.Fatalf("settings: got %v, want %v", names, want)
		}
	}
	built, err := nav.Build[validChart]()
	if err != nil {
		t.Fatal(err)
	}
	check(built)

	// A selected property cannot be reached through a pointer embedding.
	type badChart struct {
		Server struct{ *ExtraSettings } `service:"service"`
	}
	if m, err := nav.Build[badChart](); err == nil || m != nil {
		t.Fatalf("pointer path: model=%v, err=%v", m, err)
	}
	var configured validChart
	configured.Server.Promoted = "runtime value"
	populated, err := nav.Populate(configured)
	if err != nil {
		t.Fatal(err)
	}
	check(populated)
	if populated == built {
		t.Fatal("population reused the first model")
	}
}

func TestEnumAcceptanceChecksDeclaredType(t *testing.T) {
	type definedBool bool
	type chart struct {
		Server struct {
			Format struct {
				Text definedBool `service:"option,default" name:"text"`
			} `service:"setting,enum" name:"format"`
		} `service:"service" name:"reporter"`
	}
	nav, err := sail.Compile[serviceLegend]()
	if err != nil {
		t.Fatal(err)
	}
	if m, err := nav.Build[chart](); err == nil || m != nil {
		t.Fatalf("distinct defined bool: model=%v, err=%v", m, err)
	}
}
