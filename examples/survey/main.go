package main

import (
	"fmt"
	"log"

	"github.com/qlustra/sail"
	"github.com/qlustra/sail/legend"
)

type Workflow struct {
	legend.Legend `sail:"legend" mark:"workflow"`
	Entries       Entries `sail:"lane" lane:"entries"`
}

type Entries struct {
	legend.Lane
	Branches []Branch `sail:"spot" spot:"multi"`
	Notes    []Note   `sail:"spot" spot:"multi"`
}

type Branch struct {
	legend.Region `sail:"region" family:"branch"`
	Name          legend.Feature `sail:"feature" symbol:"name"`
	Steps         Steps          `sail:"lane" lane:"steps"`
}

type Steps struct {
	legend.Lane
	Items []Step `sail:"spot" spot:"multi"`
}

type Step struct {
	legend.Landmark `sail:"landmark" family:"step"`
	Name            legend.Feature `sail:"feature" symbol:"name"`
}

type Note struct {
	legend.Landmark `sail:"landmark" family:"note"`
	Name            legend.Feature `sail:"feature" symbol:"name"`
}

type Chart struct {
	Draft struct {
		Write struct{} `workflow:"step" name:"draft"`
	} `workflow:"branch" name:"notify"`
	Memo     struct{} `workflow:"note" name:"memo"`
	Delivery struct {
		Send struct{} `workflow:"step" name:"deliver"`
	} `workflow:"branch" name:"notify"`
}

func main() {
	navigator, err := sail.Compile[Workflow]()
	if err != nil {
		log.Fatal(err)
	}
	model, err := navigator.Build[Chart]()
	if err != nil {
		log.Fatal(err)
	}

	branchPath := sail.Select[*Branch]().Encounter(func(_ sail.Encounter, path string, branch *Branch) string {
		return path + "/" + branch.Name.Get()
	})
	step := sail.Select[*Step]().Encounter(func(_ sail.Encounter, path string, step *Step) string {
		fmt.Printf("step %s/%s\n", path, step.Name.Get())
		return path
	})
	note := sail.Select[*Note]().Encounter(func(_ sail.Encounter, path string, note *Note) string {
		fmt.Printf("note %s/%s\n", path, note.Name.Get())
		return path
	})

	for _, order := range []struct {
		name  string
		order sail.TraversalOrder
	}{{"Chart order", sail.ChartOrder}, {"Legend order", sail.LegendOrder}} {
		fmt.Println(order.name)
		if err := navigator.Survey(model, order.order, "", branchPath, step, note); err != nil {
			log.Fatal(err)
		}
	}
}
