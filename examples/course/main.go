// Course compares immediate selection with complete-route search.
package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/qlustra/sail"
	"github.com/qlustra/sail/legend"
)

// Program declares the grammar and is also the populated model.
type Program struct {
	legend.Legend `sail:"legend" mark:"cli"`
	Entry         struct {
		legend.Lane
		Root Command `sail:"spot" spot:"single,required"`
	} `sail:"lane"`
}

type Command struct {
	legend.Region `sail:"region" family:"command"`
	Name          legend.Feature `sail:"feature" symbol:"name"`
	Help          legend.Feature `sail:"feature" symbol:"help"`
	Children      struct {
		legend.Lane
		Commands []Command `sail:"spot" spot:"multi"`
	} `sail:"lane" lane:"commands"`
}

// Client A is a dead end for version; client B completes the route.
type App struct {
	Root struct {
		A struct {
			Status struct{} `cli:"command" name:"status"`
		} `cli:"command" name:"client" help:"Client A"`
		B struct {
			Version struct{} `cli:"command" name:"version"`
		} `cli:"command" name:"client" help:"Client B"`
	} `cli:"command" name:"tool"`
}

func named(name string) sail.SpotSelector[*Command] {
	return sail.Select(func(command *Command) bool {
		return command.Name.Get() == name
	}).InLane("commands")
}

func describe(route sail.Route) string {
	parts := make([]string, 0, len(route))
	for _, marker := range route {
		command := marker.(*Command)
		name := command.Name.Get()
		if command.Help.Present() {
			name += " (" + command.Help.Get() + ")"
		}
		parts = append(parts, name)
	}
	return strings.Join(parts, " -> ")
}

func run() error {
	navigator, err := sail.Compile[Program]()
	if err != nil {
		return err
	}
	program, err := navigator.Build[App]()
	if err != nil {
		return err
	}

	root, _, err := navigator.Next(program, sail.Select[*Command]())
	if err != nil {
		return err
	}
	if !root.Matched() {
		return fmt.Errorf("expected a unique root command")
	}
	client, _, err := root.Next(named("client"))
	if err != nil {
		return err
	}
	switch {
	case client.Ambiguous():
		fmt.Println("Next client: ambiguous (two immediate matches)")
	case client.Matched():
		fmt.Println("Next client: unique immediate match")
	default:
		fmt.Println("Next client: no immediate match")
	}

	plot := navigator.PlotCourse()
	plot.Waypoint(sail.Select[*Command]())
	plot.Waypoint(named("client"))
	plot.Waypoint(named("version"))
	versionCourse, err := plot.Prepare()
	if err != nil {
		return err
	}
	route, err := versionCourse.RequireOneRoute(program, sail.ChartOrder)
	if err != nil {
		return err
	}
	fmt.Println("RequireOneRoute:", describe(route))

	shortPlot := navigator.PlotCourse()
	shortPlot.Waypoint(sail.Select[*Command]())
	shortPlot.Waypoint(named("client"))
	clientCourse, err := shortPlot.Prepare()
	if err != nil {
		return err
	}
	routes, err := clientCourse.AllRoutes(program, sail.ChartOrder)
	if err != nil {
		return err
	}
	fmt.Printf("AllRoutes: %d\n", len(routes))
	for _, route := range routes {
		fmt.Println(" ", describe(route))
	}
	first, err := clientCourse.FirstRoute(program, sail.ChartOrder)
	if err != nil {
		return err
	}
	fmt.Println("FirstRoute:", describe(first))
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
