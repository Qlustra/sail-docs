package main

import (
	"fmt"
	"reflect"

	"github.com/qlustra/sail"
)

// Driver retains the populated declaration and registered frame types, never
// the values of a frame passed to New.
type Driver struct {
	nav    *sail.Navigator[Grammar]
	model  *Grammar
	frames map[string]reflect.Type
}

func New(chart any, registry map[string]any) (*Driver, error) {
	nav, err := sail.Compile[Grammar]()
	if err != nil {
		return nil, err
	}
	model, err := nav.Populate(chart)
	if err != nil {
		return nil, err
	}
	frames := make(map[string]reflect.Type, len(registry))
	for name, value := range registry {
		t := reflect.TypeOf(value)
		if name == "" || t == nil || t.Kind() != reflect.Struct {
			return nil, fmt.Errorf("invalid frame registration %q", name)
		}
		frames[name] = t
	}
	return &Driver{nav: nav, model: model, frames: frames}, nil
}

// step asks only for an immediate child, preserving the recursive route.
func (d *Driver) step(previous *sail.SpotSelection[*Command], match func(*Command) bool) (sail.SpotSelection[*Command], *Command, error) {
	selector := sail.Select[*Command](match)
	var selection sail.SpotSelection[*Command]
	var command *Command
	var err error
	if previous == nil {
		selection, command, err = d.nav.Next(d.model, selector)
	} else {
		selection, command, err = previous.Next(selector)
	}
	if err != nil {
		return selection, nil, err
	}
	if selection.Ambiguous() {
		return selection, nil, fmt.Errorf("ambiguous command")
	}
	if !selection.Matched() {
		return selection, nil, fmt.Errorf("command not found")
	}
	return selection, command, nil
}

// Describe returns declaration-derived usage lines for an explicit command path.
func (d *Driver) Describe(path ...string) ([]string, error) {
	if len(path) == 0 {
		return nil, fmt.Errorf("empty path")
	}
	var previous *sail.SpotSelection[*Command]
	lines := make([]string, 0, len(path))
	for _, name := range path {
		selection, command, err := d.step(previous, func(c *Command) bool { return c.Name.Get() == name })
		if err != nil {
			return nil, err
		}
		line := command.Name.Get()
		for _, option := range command.Body.Options {
			switch option.Kind.Get() {
			case "flag", "value":
				line += " " + option.Name.Get()
				if option.Kind.Get() == "value" {
					line += " <value>"
				}
			case "arg":
				line += " <" + option.Name.Get() + ">"
			default:
				return nil, fmt.Errorf("unknown option kind %q", option.Kind.Get())
			}
		}
		lines = append(lines, line)
		previous = &selection
	}
	return lines, nil
}

// Render emits literal argv elements from separate frame values. It does not
// parse a command line, quote shell text, or execute a process.
func (d *Driver) Render(frames ...any) ([]string, error) {
	if len(frames) == 0 {
		return nil, fmt.Errorf("empty frames")
	}
	var previous *sail.SpotSelection[*Command]
	var args []string
	for _, frame := range frames {
		t := reflect.TypeOf(frame)
		if t == nil || t.Kind() != reflect.Struct {
			return nil, fmt.Errorf("invalid frame %T", frame)
		}
		selection, command, err := d.step(previous, func(c *Command) bool {
			registered, ok := d.frames[c.Frame.Get()]
			return ok && registered == t
		})
		if err != nil {
			return nil, err
		}
		if command.Name.Get() == "" {
			return nil, fmt.Errorf("empty command name")
		}
		args = append(args, command.Name.Get())
		value := reflect.ValueOf(frame)
		for _, option := range command.Body.Options {
			fieldName, optionName := option.Field.Get(), option.Name.Get()
			field := value.FieldByName(fieldName)
			if !field.IsValid() || !field.CanInterface() || optionName == "" {
				return nil, fmt.Errorf("invalid option field %q", fieldName)
			}
			switch option.Kind.Get() {
			case "flag":
				if field.Kind() != reflect.Bool {
					return nil, fmt.Errorf("flag %q requires bool", fieldName)
				}
				if field.Bool() {
					args = append(args, optionName)
				}
			case "value", "arg":
				if field.Kind() != reflect.String {
					return nil, fmt.Errorf("option %q requires string", fieldName)
				}
				if option.Kind.Get() == "value" {
					args = append(args, optionName)
				}
				args = append(args, field.String())
			default:
				return nil, fmt.Errorf("unknown option kind %q", option.Kind.Get())
			}
		}
		previous = &selection
	}
	return args, nil
}
