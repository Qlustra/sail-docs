package main

import (
	"strings"
	"testing"

	"github.com/qlustra/sail"
)

func TestSurveyOrders(t *testing.T) {
	navigator, err := sail.Compile[Workflow]()
	if err != nil {
		t.Fatal(err)
	}
	model, err := navigator.Build[Chart]()
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name  string
		order sail.TraversalOrder
		want  string
	}{
		{"Chart order", sail.ChartOrder, "step /notify/draft\nnote /memo\nstep /notify/deliver\n"},
		{"Legend order", sail.LegendOrder, "step /notify/draft\nstep /notify/deliver\nnote /memo\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var got strings.Builder
			branchPath := sail.Select[*Branch]().Encounter(func(_ sail.Encounter, path string, branch *Branch) string {
				return path + "/" + branch.Name.Get()
			})
			step := sail.Select[*Step]().Encounter(func(_ sail.Encounter, path string, step *Step) string {
				got.WriteString("step " + path + "/" + step.Name.Get() + "\n")
				return path
			})
			note := sail.Select[*Note]().Encounter(func(_ sail.Encounter, path string, note *Note) string {
				got.WriteString("note " + path + "/" + note.Name.Get() + "\n")
				return path
			})
			if err := navigator.Survey(model, tc.order, "", branchPath, step, note); err != nil {
				t.Fatal(err)
			}
			if got.String() != tc.want {
				t.Errorf("survey output:\n%s\nwant:\n%s", got.String(), tc.want)
			}
		})
	}
}
