package main

import (
	"errors"
	"fmt"

	"github.com/qlustra/intel/rules/typerule"
	"github.com/qlustra/sail"
	"github.com/qlustra/sail/legend"
)

type Document interface{ Read() string }
type Descriptive interface{ Summary() string }

type summary struct {
	legend.Qualifier `sail:"qualifier"`
	Text             legend.Barrel[string, Descriptive] `sail:"barrel" populate:"Summary" barrel:"required"`
}

type document struct {
	legend.Landmark `sail:"landmark" family:"document"`
	Name            legend.Feature `sail:"feature" symbol:"name"`
	Summary         summary        `sail:"qualifier" symbol:"summary"`
}

func (document) Accept() typerule.RuleAPI { return typerule.Type().Satisfies[Document]() }

type folder struct {
	legend.Region `sail:"region" family:"folder"`
	Name          legend.Feature `sail:"feature" symbol:"name"`
	Entries       entries        `sail:"lane"`
}

type entries struct {
	legend.Lane
	Folders   []folder   `sail:"spot" spot:"multi"`
	Documents []document `sail:"spot" spot:"multi"`
}

type model struct {
	legend.Legend `sail:"legend" mark:"catalog"`
	Entries       entries `sail:"lane"`
}

// Catalog stores declaration metadata, not the Chart used to create it.
type Catalog struct {
	nav   *sail.Navigator[model]
	model *model
}

// New populates from the Chart's type; it does not read its live values.
func New(chart any) (*Catalog, error) {
	nav, err := sail.Compile[model]()
	if err != nil {
		return nil, err
	}
	m, err := nav.Populate(chart)
	if err != nil {
		return nil, err
	}
	return &Catalog{nav: nav, model: m}, nil
}

// Preview lists documents in declaration order. A summary is populated only
// when the Chart field has the summary qualifier and its type satisfies
// Descriptive.
type PreviewEntry struct {
	Path       string
	Summary    string
	HasSummary bool
}

func (c *Catalog) Preview() ([]PreviewEntry, error) {
	var out []PreviewEntry
	folders := sail.Select[*folder]().Encounter(func(_ sail.Encounter, prefix string, f *folder) string {
		return join(prefix, f.Name.Get())
	})
	docs := sail.Select[*document]().Encounter(func(_ sail.Encounter, prefix string, d *document) string {
		out = append(out, PreviewEntry{Path: join(prefix, d.Name.Get()), Summary: d.Summary.Text.Get(), HasSummary: d.Summary.Text.Present()})
		return prefix
	})
	err := c.nav.Survey(c.model, sail.ChartOrder, "", folders, docs)
	return out, err
}

func join(prefix, name string) string {
	if prefix == "" {
		return name
	}
	return prefix + "/" + name
}

// Read resolves an exact route and reads the matching field from workspace.
// Route selection completes before any live field is accessed.
func (c *Catalog) Read(workspace any, names ...string) (string, error) {
	if len(names) == 0 {
		return "", errors.New("empty document route")
	}
	plot := c.nav.PlotCourse()
	for _, name := range names[:len(names)-1] {
		plot.Waypoint(sail.Select[*folder](func(f *folder) bool { return f.Name.Get() == name }))
	}
	last := names[len(names)-1]
	plot.Waypoint(sail.Select[*document](func(d *document) bool { return d.Name.Get() == last }))
	course, err := plot.Prepare()
	if err != nil {
		return "", err
	}
	route, err := course.RequireOneRoute(c.model, sail.ChartOrder)
	if err != nil {
		return "", err
	}
	d := route[len(route)-1].(*document)
	location, check := d.ChartLocation()
	if !check.OK() {
		return "", fmt.Errorf("locate document: %s", check.Why())
	}
	field, check := location.Resolve(workspace)
	if !check.OK() {
		return "", fmt.Errorf("resolve document: %s", check.Why())
	}
	address, check := field.Address()
	if !check.OK() {
		return "", fmt.Errorf("address document: %s", check.Why())
	}
	reader, ok := address.(Document)
	if !ok {
		return "", fmt.Errorf("document field %T does not implement Document", address)
	}
	return reader.Read(), nil
}
