package main

import "fmt"

type text struct{ Value string }

func (t text) Read() string  { return t.Value }
func (text) Summary() string { return "A declaration-time summary." }

type workspace struct {
	Welcome text `catalog:"document,summary" name:"welcome"`
	Guides  struct {
		Install   text `catalog:"document" name:"install"`
		Reference struct {
			API text `catalog:"document" name:"api"`
		} `catalog:"folder" name:"reference"`
	} `catalog:"folder" name:"guides"`
	About text `catalog:"document" name:"about"`
}

// demonstrate owns this demo's display format; Catalog owns its preview and
// route policy.
func demonstrate() error {
	catalog, err := New(workspace{})
	if err != nil {
		return fmt.Errorf("build catalog: %w", err)
	}
	preview, err := catalog.Preview()
	if err != nil {
		return fmt.Errorf("preview catalog: %w", err)
	}
	fmt.Println("Catalog preview:")
	for _, entry := range preview {
		if entry.HasSummary {
			fmt.Printf("  %s — %s\n", entry.Path, entry.Summary)
		} else {
			fmt.Printf("  %s\n", entry.Path)
		}
	}

	second := workspace{}
	second.Guides.Install.Value = "selected workspace"
	value, err := catalog.Read(&second, "guides", "install")
	if err != nil {
		return fmt.Errorf("read selected document: %w", err)
	}
	fmt.Printf("Read guides/install: %q\n", value)
	return nil
}
