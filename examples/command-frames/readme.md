# Describe a command tree and render explicit frames

This example uses Sail to declare a reusable command grammar, then implements
the application-specific driver that describes a requested route and turns
caller-supplied values into argument elements. It does not parse input, quote a
shell command, or launch a process.

From `examples/`, run it with:

```sh
go run ./command-frames
```

Observed output:

```text
dispatch --profile <value>
project --workspace <value>
show --verbose <target>
[dispatch --profile prod project --workspace demo show --verbose api]
```

## Follow the source

1. [`grammar.go`](grammar.go) defines Sail's grammar shape: commands are
   repeated Spots at the root and under each command, while options are Spots
   in a command's body.
2. [`demo.go`](demo.go) supplies the tagged declaration: `dispatch` contains
   `project`, which contains `show`; `other` separately contains another
   command named `show`. `Chart` is the declaration value. The types registered
   with `demoDriver` also identify the types accepted as separate frames.
3. [`driver.go`](driver.go) has the reusable application-side logic. `New`
   compiles and populates the grammar and records registered frame types;
   `step` moves only to an immediate child. `Describe` follows the explicit
   path one level at a time, so `dispatch, project, show` selects the intended
   route, while `dispatch, show` fails rather than skipping `project`. It
   needs no configured frame values.
4. In the same file, `Render` walks that same recursive route, matches each
   command's declared frame type, and reads values from the corresponding
   caller-supplied frame. It appends options in declaration order. The true
   `Verbose` flag emits `--verbose`; false emits nothing. Tests also exercise
   different values, a deeper `detail` frame, and rejection of a frame from
   the unrelated branch.
5. [`main.go`](main.go) wires the demo together and prints the description and
   rendered argument elements shown above.

Sail compiles the tagged declaration and provides the populated model and
recursive navigation. A Sail `ChartLocation` is a root-relative witness to a
field in the exact `Chart` used for population; resolving one accesses that
field relative to that chart root. It does not identify or resolve the
separately supplied `Dispatch`, `Project`, or `Show` frame values in `Render`.
Those values are independent inputs, selected here by the driver's frame-type
registry and read by the driver's field mapping.

The driver owns the policies above the declaration: frame registration and
matching, route description, option-to-field mapping, usage notation, and
rendering values into argument elements. Sail does not parse arguments,
execute commands, or supply values for the caller.

Check the example and its regression tests from `examples/`:

```sh
go test ./command-frames
```
