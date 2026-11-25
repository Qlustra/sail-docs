// Command lanes demonstrates required separators and exploration policies.
package main

import (
	"fmt"
	"os"

	"github.com/qlustra/sail"
	"github.com/qlustra/sail/legend"
)

type job struct {
	legend.Landmark `sail:"landmark" family:"job"`
	Name            legend.Feature `sail:"feature" symbol:"name"`
}

type checkpoint struct {
	legend.Landmark `sail:"landmark" family:"checkpoint"`
}

type reusableJobLane struct {
	legend.Lane
	Job job `sail:"spot" spot:"single"`
}

type requiredSeparatorLane struct {
	legend.Lane
	Mark checkpoint `sail:"spot" spot:"single,required"`
}

type optionalSeparatorLane struct {
	legend.Lane
	Mark checkpoint `sail:"spot" spot:"single"`
}

type separatedLegend struct {
	legend.Legend `sail:"legend" mark:"pipeline"`
	Before        reusableJobLane       `sail:"lane"`
	Boundary      requiredSeparatorLane `sail:"lane"`
	After         reusableJobLane       `sail:"lane"`
}

type optionalLegend struct {
	legend.Legend `sail:"legend" mark:"pipeline"`
	Before        reusableJobLane       `sail:"lane"`
	Boundary      optionalSeparatorLane `sail:"lane"`
	After         reusableJobLane       `sail:"lane"`
}

type eagerBeforeLegend struct {
	legend.Legend `sail:"legend" mark:"pipeline"`
	Before        reusableJobLane       `sail:"lane" exploration:"eager"`
	Boundary      optionalSeparatorLane `sail:"lane"`
	After         reusableJobLane       `sail:"lane"`
}

type jobsLane struct {
	legend.Lane
	Jobs []job `sail:"spot" spot:"multi"`
}

type laterLane struct {
	legend.Lane
	Jobs       []job      `sail:"spot" spot:"multi,required"`
	Checkpoint checkpoint `sail:"spot" spot:"single"`
}

type checkpointLegend struct {
	legend.Legend `sail:"legend" mark:"pipeline"`
	Before        jobsLane  `sail:"lane" exploration:"eager"`
	After         laterLane `sail:"lane"`
}

type separatedChart struct {
	A    struct{} `pipeline:"job" name:"A"`
	Mark struct{} `pipeline:"checkpoint"`
	B    struct{} `pipeline:"job" name:"B"`
}

type unseparatedChart struct {
	A struct{} `pipeline:"job" name:"A"`
	B struct{} `pipeline:"job" name:"B"`
}

type checkpointChart struct {
	A    struct{} `pipeline:"job" name:"A"`
	B    struct{} `pipeline:"job" name:"B"`
	Mark struct{} `pipeline:"checkpoint"`
	C    struct{} `pipeline:"job" name:"C"`
}

func demonstrate() error {
	nav, err := sail.Compile[separatedLegend]()
	if err != nil {
		return fmt.Errorf("compile separated lanes: %w", err)
	}
	got, err := nav.Populate(separatedChart{})
	if err != nil {
		return fmt.Errorf("populate separated lanes: %w", err)
	}
	fmt.Printf("required boundary: Before %q, boundary present %v, After %q\n", got.Before.Job.Name.Get(), got.Boundary.Mark.Present(), got.After.Job.Name.Get())
	if _, err := nav.Populate(unseparatedChart{}); err == nil {
		return fmt.Errorf("unseparated jobs unexpectedly populated")
	}
	fmt.Println("missing boundary: population failed")
	if _, err := sail.Compile[optionalLegend](); err == nil {
		return fmt.Errorf("optional boundary unexpectedly compiled with default reuse")
	}
	fmt.Println("optional boundary + default reuse: compilation rejected")
	eager, err := sail.Compile[eagerBeforeLegend]()
	if err != nil {
		return fmt.Errorf("compile eager-before lanes: %w", err)
	}
	if _, err := eager.Populate(unseparatedChart{}); err == nil {
		return fmt.Errorf("eager Before unexpectedly left a job for After")
	}
	fmt.Println("optional boundary + eager Before: compiles; unseparated population fails")

	checkpointNav, err := sail.Compile[checkpointLegend]()
	if err != nil {
		return fmt.Errorf("compile checkpoint lanes: %w", err)
	}
	fmt.Printf("checkpoint dependency advisories: %d\n", len(checkpointNav.Warnings()))
	checkpointGot, err := checkpointNav.Populate(checkpointChart{})
	if err != nil {
		return fmt.Errorf("populate checkpoint lanes: %w", err)
	}
	fmt.Printf("eager checkpoint: Before %d jobs, After %q, checkpoint present %v\n", len(checkpointGot.Before.Jobs), checkpointGot.After.Jobs[0].Name.Get(), checkpointGot.After.Checkpoint.Present())
	type missingCheckpointBoundary struct {
		A     struct{} `pipeline:"job" name:"A"`
		Other string
		C     struct{} `pipeline:"job" name:"C"`
	}
	if _, err := checkpointNav.Populate(missingCheckpointBoundary{}); err == nil {
		return fmt.Errorf("missing later required job unexpectedly populated")
	}
	fmt.Println("missing checkpoint boundary: required later job population failed")
	return nil
}

func main() {
	if err := demonstrate(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
