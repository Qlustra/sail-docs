# Preview and read an in-memory workspace catalog

This runnable example interprets tagged struct declarations as named folders
and documents. It previews declarations and reads from a caller-selected live
workspace without filesystem access.

Run from `examples/`:

```sh
go run ./catalog
```

Observed output:

```text
Catalog preview:
  welcome — A declaration-time summary.
  guides/install
  guides/reference/api
  about
Read guides/install: "selected workspace"
```

## Follow the source

1. In [`catalog.go`](catalog.go), `model`, `entries`, `folder`, and `document`
   define the Sail grammar. `New` compiles and populates that grammar from a
   workspace type; `Catalog.Preview` surveys declarations without reading live
   document values.
2. In [`example.go`](example.go), `workspace` supplies the tagged declaration.
   A summary appears only when a document has the `summary` qualifier tag
   (`catalog:"document,summary"`) **and** its type implements `Summary()`.
   `summary.Text` maps that qualifier to the required `Descriptive` barrel.
3. `Catalog.Read` in `catalog.go` accepts separate exact path segments, selects
   one unique route, then resolves and reads the field from the workspace
   supplied by the caller.
4. `demonstrate` in `example.go` chooses how preview entries are displayed and
   which document to read; `main.go` handles command-line errors. The catalog
   framework owns its summary opt-in and unique-route policy, while this
   example owns its presentation. If membership were discovered at runtime,
   the framework would also need to own that discovery and its ordering. This
    example performs no filesystem I/O or runtime discovery.

Run the regression tests from `examples/` with `go test ./...`.
