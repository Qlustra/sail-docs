package main

import (
	"reflect"
	"testing"
)

func testDriver(t *testing.T) *Driver {
	t.Helper()
	driver, err := New(Chart{}, map[string]any{
		"dispatch": Dispatch{}, "project": Project{}, "show": Show{},
		"detail": Detail{}, "other": Other{}, "other-show": OtherShow{},
	})
	if err != nil {
		t.Fatal(err)
	}
	return driver
}

func TestDescribeUsesDeclaredRecursiveRoute(t *testing.T) {
	driver := testDriver(t)
	tests := []struct {
		path []string
		want []string
	}{
		{[]string{"dispatch", "project", "show"}, []string{"dispatch --profile <value>", "project --workspace <value>", "show --verbose <target>"}},
		{[]string{"other", "show"}, []string{"other", "show --detail"}},
		{[]string{"dispatch", "project", "show", "detail"}, []string{"dispatch --profile <value>", "project --workspace <value>", "show --verbose <target>", "detail --format <value>"}},
	}
	for _, test := range tests {
		got, err := driver.Describe(test.path...)
		if err != nil || !reflect.DeepEqual(got, test.want) {
			t.Errorf("Describe(%v) = %v, %v; want %v", test.path, got, err, test.want)
		}
	}
	for _, path := range [][]string{nil, {"dispatch", "show"}, {"dispatch", "project", "absent"}} {
		if got, err := driver.Describe(path...); err == nil || got != nil {
			t.Errorf("Describe(%v) = %v, %v; want an error and no result", path, got, err)
		}
	}
}

func TestRenderUsesEachSeparateFrameAndRecursiveParent(t *testing.T) {
	driver := testDriver(t)
	tests := []struct {
		frames []any
		want   []string
	}{
		{[]any{Dispatch{Profile: "prod"}, Project{Workspace: "demo"}, Show{Verbose: true, Target: "api"}}, []string{"dispatch", "--profile", "prod", "project", "--workspace", "demo", "show", "--verbose", "api"}},
		{[]any{Dispatch{Profile: "dev"}, Project{Workspace: "local"}, Show{Target: "web"}}, []string{"dispatch", "--profile", "dev", "project", "--workspace", "local", "show", "web"}},
		{[]any{Dispatch{Profile: "dev"}, Project{Workspace: "local"}, Show{Target: "web"}, Detail{Format: "json"}}, []string{"dispatch", "--profile", "dev", "project", "--workspace", "local", "show", "web", "detail", "--format", "json"}},
	}
	for _, test := range tests {
		got, err := driver.Render(test.frames...)
		if err != nil || !reflect.DeepEqual(got, test.want) {
			t.Errorf("Render(%v) = %v, %v; want %v", test.frames, got, err, test.want)
		}
	}
	for _, frames := range [][]any{nil, {Dispatch{}, Show{}}, {Dispatch{}, Project{}, OtherShow{}}, {Dispatch{}, Project{}, Show{}, OtherShow{}}} {
		if got, err := driver.Render(frames...); err == nil || got != nil {
			t.Errorf("Render(%v) = %v, %v; want an error and no result", frames, got, err)
		}
	}
}

func TestNewRejectsInvalidFrameRegistration(t *testing.T) {
	for name, value := range map[string]any{"": Dispatch{}, "nil": nil, "pointer": (*Dispatch)(nil)} {
		if _, err := New(Chart{}, map[string]any{name: value}); err == nil {
			t.Errorf("New registry %q = nil error; want invalid registration", name)
		}
	}
}

func TestNewUsesChartTypeNotValues(t *testing.T) {
	chart := Chart{Dispatch: Dispatch{Profile: "must not appear"}}
	driver, err := New(chart, map[string]any{
		"dispatch": Dispatch{}, "project": Project{}, "show": Show{},
		"detail": Detail{}, "other": Other{}, "other-show": OtherShow{},
	})
	if err != nil {
		t.Fatal(err)
	}
	args, err := driver.Render(Dispatch{}, Project{}, Show{Target: "api"})
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"dispatch", "--profile", "", "project", "--workspace", "", "show", "api"}
	if !reflect.DeepEqual(args, want) {
		t.Fatalf("Render = %v; want %v", args, want)
	}
}
