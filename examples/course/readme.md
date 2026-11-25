# Find a complete route through ambiguous choices

An immediate choice and a complete route answer different questions. This
program models a command tree with two `client` commands: the first has a
`status` child, while the second has a `version` child. Selecting `client`
immediately is ambiguous, but the full route to `version` is unique.

Run it from `examples/`:

```sh
go run ./course
```

Expected output:

```text
Next client: ambiguous (two immediate matches)
RequireOneRoute: tool -> client (Client B) -> version
AllRoutes: 2
  tool -> client (Client A)
  tool -> client (Client B)
FirstRoute: tool -> client (Client A)
```

## Walkthrough

In `main.go`, `Program` and `Command` define the grammar, and `App` supplies
the application shape to build. `named` selects a command by its `name`
Feature in the `commands` Lane. `run` compiles the grammar, builds the model,
and uses `Next` to show that the immediate `client` choice is ambiguous. It
then prepares a three-waypoint Course (`tool` → `client` → `version`) and calls
`RequireOneRoute`; a shorter Course demonstrates `AllRoutes` and `FirstRoute`.

Sail compiles the declared grammar, populates its Legend from the Chart's type
structure, and selects stored tokens for the requested route. The framework
chooses the selectors, decides which route result it needs, and handles the
selected application data; Course navigation does not inspect or mutate live
Chart values. See the [navigation manual](../../manual/navigation.md#prepare-and-search-a-course)
for the API contract and ordering details.
