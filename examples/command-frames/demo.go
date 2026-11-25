package main

// These structs are the example's declared command and frame types. Their
// nested fields form the Chart; the caller later supplies distinct frame values.
type Dispatch struct {
	Profile string  `cli:"option" name:"--profile" field:"Profile" kind:"value"`
	Project Project `cli:"command" name:"project" frame:"project"`
}

type Project struct {
	Workspace string `cli:"option" name:"--workspace" field:"Workspace" kind:"value"`
	Show      Show   `cli:"command" name:"show" frame:"show"`
}

type Show struct {
	Verbose bool   `cli:"option" name:"--verbose" field:"Verbose" kind:"flag"`
	Target  string `cli:"option" name:"target" field:"Target" kind:"arg"`
	Detail  Detail `cli:"command" name:"detail" frame:"detail"`
}

type Detail struct {
	Format string `cli:"option" name:"--format" field:"Format" kind:"value"`
}

type Other struct {
	Show OtherShow `cli:"command" name:"show" frame:"other-show"`
}

type OtherShow struct {
	Detail bool `cli:"option" name:"--detail" field:"Detail" kind:"flag"`
}

type Chart struct {
	Dispatch Dispatch `cli:"command" name:"dispatch" frame:"dispatch"`
	Other    Other    `cli:"command" name:"other" frame:"other"`
}

func demoDriver() (*Driver, error) {
	return New(Chart{}, map[string]any{
		"dispatch": Dispatch{}, "project": Project{}, "show": Show{},
		"detail": Detail{}, "other": Other{}, "other-show": OtherShow{},
	})
}
