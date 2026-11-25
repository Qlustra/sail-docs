# Describe application settings without reading their live values

This example uses a framework-owned `Legend` as both grammar and result model,
then applies that grammar to an application's ordinary Go `Chart`. `Build`
uses only the Chart's type; `Populate` accepts a value to supply its dynamic
type, but does not inspect its configured fields. The two calls return separate
models, and the provider runs on a zero-valued declaration receiver.

Run from `examples/`:

```sh
go run ./population
```

Output:

```text
Service: reporter
  scalar: port
  enum: format
    text (declared default: true)
    json (declared default: false)
Build:    Render reports (port on declaration receiver: 0)
Populate: Render reports (port on declaration receiver: 0)
```

In [`main.go`](main.go), start with the app’s `server`: its unmodified
`service:"setting"` populates the scalar Spot, while `service:"setting,enum"`
selects the enum Region. Inside that Region, `optionLandmark.Accept` checks for
`bool` fields, and the `default` Qualifier marks the declared default. Neither
the option name nor that default depends on the application’s live boolean
value.

Follow `showPopulation` to see `Build` and `Populate` read the same tags into
separate models. Both descriptions show port zero: Sail invokes `About` on
a zero-initialized `server`, not the configured `Port: 8080`. See the manual’s
[compile and populate](../../manual/compile-and-populate.md) chapter for the
two entry points.

The framework owns the grammar, the meaning of a default and how to use the
result; Sail interprets the declarations and enforces the authored rules.
