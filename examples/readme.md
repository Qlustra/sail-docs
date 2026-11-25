# Sail examples

These are independent programs. Start with the task you want to solve; each
walkthrough shows what to run, what to expect, and which decisions belong to
the framework. For a capability reference, see the [manual](../manual/readme.md).

## Run the examples

Use Go 1.27 or newer. This module uses the published Sail `v1.0.0`; no Go
workspace is needed. Run these commands from `examples/`:

```sh
go mod download github.com/qlustra/sail github.com/qlustra/intel
go run ./population
```

Run `go test ./...` from `examples/` to check the examples.

If you use [Task](https://taskfile.dev/), `task --list` shows one task per
example. For instance, `task population` runs the population example. The
tasks disable Go workspace overrides so they use the pinned release.

## Choose a task

| Task | Command |
| --- | --- |
| [Describe application settings without reading their live values](population/readme.md) | `go run ./population` |
| [Tell a bad grammar, a bad Chart, and an impossible route apart](validation/readme.md) | `go run ./validation` |
| [Split repeated members into ordered groups](lanes/readme.md) | `go run ./lanes` |
| [List nested declarations with inherited context](survey/readme.md) | `go run ./survey` |
| [Find the sole complete route through ambiguous siblings](course/readme.md) | `go run ./course` |
| [Update one field in a live Chart](live-field/readme.md) | `go run ./live-field` |
| [Preview and read an in-memory workspace catalog](catalog/readme.md) | `go run ./catalog` |
| [Describe a recursive command tree and render separate frames](command-frames/readme.md) | `go run ./command-frames` |
