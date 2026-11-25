package main

import (
	"testing"

	"github.com/qlustra/sail"
	"github.com/qlustra/sail/legend"
)

func TestDemonstrate(t *testing.T) {
	if err := demonstrate(); err != nil {
		t.Fatal(err)
	}
}

type unreachableLane struct {
	legend.Lane
	Jobs []job `sail:"spot" spot:"multi,required"`
}

type unreachableLegend struct {
	legend.Legend `sail:"legend" mark:"pipeline"`
	Before        jobsLane        `sail:"lane" exploration:"eager"`
	After         unreachableLane `sail:"lane"`
}

func TestCheckpointBoundary(t *testing.T) {
	nav, err := sail.Compile[checkpointLegend]()
	if err != nil {
		t.Fatal(err)
	}
	if len(nav.Warnings()) != 1 {
		t.Fatalf("advisories = %d; want 1", len(nav.Warnings()))
	}
	model, err := nav.Populate(checkpointChart{})
	if err != nil {
		t.Fatal(err)
	}
	if len(model.Before.Jobs) != 2 || len(model.After.Jobs) != 1 || model.After.Jobs[0].Name.Get() != "C" || !model.After.Checkpoint.Present() {
		t.Fatal("checkpoint did not separate the job occurrences")
	}
	if _, err := nav.Populate(unseparatedChart{}); err == nil {
		t.Fatal("missing boundary accepted")
	}
	if _, err := sail.Compile[unreachableLegend](); err == nil {
		t.Fatal("unreachable required jobs accepted")
	}
}
