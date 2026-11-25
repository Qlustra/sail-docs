package main

import "github.com/qlustra/sail/legend"

// Grammar declares commands at the root and recursively in each command body.
type Grammar struct {
	legend.Legend `sail:"legend" mark:"cli"`
	Body          CommandLane `sail:"lane"`
}

type CommandLane struct {
	legend.Lane
	Commands []Command `sail:"spot" spot:"multi"`
}

type Command struct {
	legend.Region `sail:"region" family:"command"`
	Name          legend.Feature `sail:"feature" symbol:"name"`
	Frame         legend.Feature `sail:"feature" symbol:"frame"`
	Body          BodyLane       `sail:"lane"`
}

type BodyLane struct {
	legend.Lane
	Options  []Option  `sail:"spot" spot:"multi"`
	Commands []Command `sail:"spot" spot:"multi"`
}

type Option struct {
	legend.Landmark `sail:"landmark" family:"option"`
	Name            legend.Feature `sail:"feature" symbol:"name"`
	Field           legend.Feature `sail:"feature" symbol:"field"`
	Kind            legend.Feature `sail:"feature" symbol:"kind"`
}
