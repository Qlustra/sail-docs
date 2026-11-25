package main

import (
	"fmt"
	"log"
	"strings"
)

func main() {
	driver, err := demoDriver()
	if err != nil {
		log.Fatal(err)
	}
	description, err := driver.Describe("dispatch", "project", "show")
	if err != nil {
		log.Fatal(err)
	}
	args, err := driver.Render(
		Dispatch{Profile: "prod"}, Project{Workspace: "demo"}, Show{Verbose: true, Target: "api"},
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", strings.Join(description, "\n"))
	fmt.Printf("[%s]\n", strings.Join(args, " "))
}
