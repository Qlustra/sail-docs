// Command validation demonstrates three separate Sail validation boundaries.
package main

import (
	"fmt"
	"os"

	"github.com/qlustra/sail"
	"github.com/qlustra/sail/legend"
)

type option struct {
	legend.Landmark `sail:"landmark" family:"option"`
}

type options struct {
	legend.Lane
	Option option `sail:"spot" spot:"single,required"`
}

type language struct {
	legend.Legend `sail:"legend" mark:"demo"`
	Options       options `sail:"lane" lane:"options"`
}

type duplicateLanes struct {
	legend.Legend `sail:"legend" mark:"demo"`
	First         options `sail:"lane" lane:"options" exploration:"eager"`
	Second        options `sail:"lane" lane:"options"`
}

type validChart struct {
	Choice bool `demo:"option"`
}

type missingOption struct{}

func run() error {
	badNavigator, err := sail.Compile[duplicateLanes]()
	if err == nil || badNavigator != nil {
		return fmt.Errorf("Compile did not reject duplicate lane IDs")
	}
	fmt.Println("Compile: rejected duplicate lane IDs")

	navigator, err := sail.Compile[language]()
	if err != nil || navigator == nil {
		return fmt.Errorf("Compile did not accept the valid Legend")
	}
	model, err := navigator.Build[validChart]()
	if err != nil || model == nil || !model.Options.Option.Present() {
		return fmt.Errorf("Build did not accept the Chart with its required option")
	}
	fmt.Println("Build: accepted Chart with required option")

	missing, err := navigator.Build[missingOption]()
	if err == nil || missing != nil {
		return fmt.Errorf("Build did not reject the missing required option")
	}
	fmt.Println("Build: rejected missing required option")

	plot := navigator.PlotCourse()
	plot.Waypoint(sail.Select[*option]())
	plot.Waypoint(sail.Select[*option]())
	course, err := plot.Prepare()
	if err == nil || course != nil {
		return fmt.Errorf("Prepare did not reject continuing past a terminal Landmark")
	}
	fmt.Println("Prepare: rejected waypoint after terminal Landmark")
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
