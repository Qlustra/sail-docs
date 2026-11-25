// Command live-field demonstrates declaration metadata and explicit live access.
package main

import (
	"fmt"
	"os"

	"github.com/qlustra/sail"
	"github.com/qlustra/sail/legend"
)

type field struct {
	legend.Landmark `sail:"landmark" family:"field"`
	Name            legend.Feature `sail:"feature" symbol:"name"`
}

type fields struct {
	legend.Lane
	Entries []field `sail:"spot" spot:"multi"`
}

type model struct {
	legend.Legend `sail:"legend" mark:"demo"`
	Fields        fields `sail:"lane"`
}

type NamedFields struct {
	First string `demo:"field" name:"first"`
}

type chart struct {
	Ignored string
	NamedFields
	Second string `demo:"field" name:"second"`
}

// The same layout does not make this the recorded Chart type.
type otherChart chart

type demonstration struct {
	typeName, fieldName string
	path                []int
	before, after       string
	sibling             string
	selectedAddress     bool
	ignored             string
}

func demonstrate() (demonstration, error) {
	navigator, err := sail.Compile[model]()
	if err != nil {
		return demonstration{}, fmt.Errorf("compile: %w", err)
	}
	// Build uses only the Chart type; it neither reads nor retains live values.
	populated, err := navigator.Build[chart]()
	if err != nil {
		return demonstration{}, fmt.Errorf("populate from type: %w", err)
	}
	if len(populated.Fields.Entries) == 0 {
		return demonstration{}, fmt.Errorf("no populated field occurrence")
	}
	occurrence := populated.Fields.Entries[0]
	location, check := occurrence.ChartLocation()
	if !check.OK() || location == nil {
		return demonstration{}, fmt.Errorf("chart location: %v", check)
	}

	live := chart{Ignored: "untouched", NamedFields: NamedFields{First: "before"}, Second: "sibling"}
	resolved, check := location.Resolve(&live)
	if !check.OK() || resolved == nil {
		return demonstration{}, fmt.Errorf("resolve live chart: %v", check)
	}
	address, check := resolved.Address()
	if !check.OK() {
		return demonstration{}, fmt.Errorf("field address: %v", check)
	}
	ptr, ok := address.(*string)
	if !ok || ptr == nil {
		return demonstration{}, fmt.Errorf("resolved field address has type %T, want *string", address)
	}
	before := *ptr
	*ptr = "after"

	other := otherChart{NamedFields: NamedFields{First: "other"}, Second: "safe"}
	rejected, check := location.Resolve(&other)
	if check.OK() || rejected != nil {
		return demonstration{}, fmt.Errorf("resolve different Chart type: expected rejection, got %v", check)
	}
	return demonstration{
		typeName:        location.ChartType().CapturedType().String(),
		fieldName:       occurrence.Name.Get(),
		path:            location.Index(),
		before:          before,
		after:           live.First,
		sibling:         live.Second,
		selectedAddress: ptr == &live.First,
		ignored:         live.Ignored,
	}, nil
}

func main() {
	result, err := demonstrate()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("type evidence: %s, field path %v (name %q)\n", result.typeName, result.path, result.fieldName)
	fmt.Printf("live Resolve: First %q -> %q; Second remains %q\n", result.before, result.after, result.sibling)
	fmt.Println("different Chart type rejected")
}
