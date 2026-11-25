package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/qlustra/intel/rules/typerule"
	"github.com/qlustra/sail"
	"github.com/qlustra/sail/legend"
)

type descriptive interface{ About() string }

// The Legend is both the grammar and the populated model.
type serviceLegend struct {
	legend.Legend `sail:"legend" mark:"service"`
	Entries       struct {
		legend.Lane
		Service serviceRegion `sail:"spot" spot:"single,required"`
	} `sail:"lane"`
}

type serviceRegion struct {
	legend.Region `sail:"region" family:"service"`
	Name          legend.Feature                     `sail:"feature" symbol:"name"`
	Description   legend.Barrel[string, descriptive] `sail:"barrel" populate:"About"`
	Settings      struct {
		legend.Lane
		Scalars []settingLandmark `sail:"spot" spot:"multi"`
		Modes   []modeRegion      `sail:"spot" spot:"multi"`
	} `sail:"lane"`
}

type settingLandmark struct {
	legend.Landmark `sail:"landmark" family:"setting"`
	Name            legend.Feature `sail:"feature" symbol:"name"`
}

type modeRegion struct {
	legend.Region `sail:"region" family:"setting" discriminate:"enum"`
	Name          legend.Feature `sail:"feature" symbol:"name"`
	Choices       struct {
		legend.Lane
		Options []optionLandmark `sail:"spot" spot:"multi,required"`
	} `sail:"lane"`
}

type optionLandmark struct {
	legend.Landmark `sail:"landmark" family:"option"`
	Name            legend.Feature   `sail:"feature" symbol:"name"`
	Default         legend.Qualifier `sail:"qualifier" symbol:"default"`
}

func (optionLandmark) Accept() typerule.RuleAPI { return typerule.Type().Is[bool]() }

// The application Chart is separate from the framework's Legend.
type appChart struct {
	Server server `service:"service" name:"reporter"`
}

type server struct {
	Port   int    `service:"setting" name:"port"`
	Format format `service:"setting,enum" name:"format"`
}

type format struct {
	Text bool `service:"option,default" name:"text"`
	JSON bool `service:"option" name:"json"`
}

func (s server) About() string {
	return fmt.Sprintf("Render reports (port on declaration receiver: %d)", s.Port)
}

func showPopulation(out io.Writer) error {
	navigator, err := sail.Compile[serviceLegend]()
	if err != nil {
		return err
	}
	built, err := navigator.Build[appChart]()
	if err != nil {
		return err
	}
	populated, err := navigator.Populate(appChart{Server: server{Port: 8080, Format: format{JSON: true}}})
	if err != nil {
		return err
	}
	if built == populated {
		return fmt.Errorf("population reused the Legend")
	}
	service := &built.Entries.Service
	fmt.Fprintf(out, "Service: %s\n", service.Name.Get())
	for _, scalar := range service.Settings.Scalars {
		fmt.Fprintf(out, "  scalar: %s\n", scalar.Name.Get())
	}
	for _, mode := range service.Settings.Modes {
		fmt.Fprintf(out, "  enum: %s\n", mode.Name.Get())
		for _, option := range mode.Choices.Options {
			fmt.Fprintf(out, "    %s (declared default: %t)\n", option.Name.Get(), option.Default.Present())
		}
	}
	fmt.Fprintf(out, "Build:    %s\n", service.Description.Get())
	fmt.Fprintf(out, "Populate: %s\n", populated.Entries.Service.Description.Get())
	return nil
}

func main() {
	if err := showPopulation(os.Stdout); err != nil {
		log.Fatal(err)
	}
}
