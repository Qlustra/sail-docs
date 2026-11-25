# Navigate a populated Legend

Navigation lets a framework select and process occurrences in a populated
**Legend** without walking the application’s live Chart. A **Region** or
**Landmark** occurrence is a token: a Region opens a nested scope, while a
Landmark is a leaf. Lanes and Spots group tokens, but are not themselves
navigation steps. The Legend is the starting scope.

Use **Survey** to process the occurrences in a model, **Next** to choose one
immediate occurrence at a time, or a prepared **Course** to search for complete
routes through a chain of selections. The framework supplies the criteria and
decides what to do with the selected model data.

## Select tokens

A selector names the Go type of the token to select. It can also restrict the
incoming Lane and apply a predicate. For a job Landmark with a `Name` Feature:

```go
jobs := sail.Select[*Job](func(job *Job) bool {
    return job.Name.Get() == "compile"
}).InLane("jobs")
```

The Lane ID identifies a Lane use in the enclosing scope. Type and Lane
restrictions are checked before the predicate. Omit the predicate or `InLane`
when the framework does not need that restriction.

Selection tests the pointer to the stored token. `Select[*Job]()` can match a
Job occurrence; `Select[Job]()` does not match by making a value copy. An
interface selector can match when the stored pointer implements that interface,
and the interface holds that same pointer. These are pointers into the
populated Legend, not to Chart fields. They follow ordinary Go storage rules;
Sail does not keep a pointer attached if the framework moves or replaces the
model storage.

## Survey occurrences

`Survey` walks present tokens depth-first. Its order chooses the order of
siblings: `sail.ChartOrder` follows recognized Chart declaration order, while
`sail.LegendOrder` follows Lane and Spot declaration order. Repeated
occurrences within a Spot retain Chart order. The choice does not change which
occurrences population admitted.

Each rule pairs a selector with an encounter reducer. The reducer receives
position information, a prefix, and the selected token, and returns the prefix
for subsequent processing. For example, a rule can extend a path inherited by
descendants of a Region:

```go
rule := sail.Select[*Stage]().Encounter(
    func(at sail.Encounter, prefix string, stage *Stage) string {
        return prefix + "/" + stage.Name.Get()
    },
)
err := navigator.Survey(model, sail.ChartOrder, "", rule)
```

Here `Stage` is a framework-defined Region with a `Name` Feature. The encounter
also exposes `Depth()` and `LaneID()`; the Legend is depth zero and its
immediate tokens are depth one. Survey returns an error, not a collected result;
the framework can collect or render information in its reducers.

At a token, matching rules run in registration order: each sees the previous
rule’s returned prefix. The final prefix is passed to each child. A child’s
replacement does not rebind its siblings’ prefix, but this is ordinary Go
value passing, not a deep copy: mutations through shared pointers, maps, or
slice backing arrays can be visible elsewhere.

A selector filters encounters, not descent. A Region that does not match a
rule can still have descendants that do. Survey offers no pruning or early
halt control.

## Choose an immediate step

`Next` looks only at immediate tokens below the current position. It does not
search through a nonmatching Region. Starting from the Legend:

```go
selection, job, err := navigator.Next(model, jobs)
if err != nil {
    // Handle an invalid selection or navigation error.
} else if selection.Ambiguous() {
    // More than one immediate job matched.
} else if selection.Matched() {
    // Exactly one matched; job is its stored pointer.
} else {
    // No immediate job matched.
}
```

Only a unique selection can advance with `selection.Next(anotherSelector)`.
The returned token is zero when there is no match or the step is ambiguous;
check the selection rather than interpreting that zero as a result. A Landmark
has no children to select.

Immediate multiplicity is not the number of complete routes. For example, two
immediate Regions may match a first selector even though only one contains a
matching child. `Next` reports the first step as ambiguous; it does not look
ahead to resolve the ambiguity. Use a Course when the question is which whole
chains can complete.

## Prepare and search a Course

A **Plot** collects selectors as waypoints; `Prepare` checks the chain against
the compiled grammar and produces a reusable **Course**:

```go
plot := navigator.PlotCourse()
plot.Waypoint(stageSelector)
plot.Waypoint(jobSelector)
course, err := plot.Prepare()
if err != nil {
    return err
}
route, err := course.RequireOneRoute(model, sail.ChartOrder)
```

Each waypoint selects an immediate token below the preceding one, and can use
its own concrete or interface type, Lane restriction, and predicate. Preparation
checks selector configuration and whether the grammar admits a connected path.
It does not run predicates or inspect a populated Legend. In particular, a
possible path is not a promise that every candidate can continue, that a
predicate will match, or that any route will exist at runtime. Later edits to
the Plot do not change the prepared Course.

Course searches explore matching immediate steps depth-first, continuing with
siblings when a partial route cannot complete. `ChartOrder` and `LegendOrder`
choose sibling order as in Survey. Search outcomes count **complete routes**,
not intermediate matches:

- `FirstRoute` returns the first complete route in the chosen order, or
  `sail.ErrNoRoute`.
- `RequireOneRoute` requires exactly one, returning `sail.ErrNoRoute` or
  `sail.ErrAmbiguousRoute` otherwise.
- `AllRoutes` returns all complete routes in order; no matches is a successful
  empty result.

A route contains the stored token pointers in waypoint order. Its slice is
separate, but its elements still refer to the populated Legend. To reach the
corresponding live application field instead, use the token’s
[Chart location](chart-location.md).

[Previous: Compile and populate](compile-and-populate.md) · [Contents](readme.md) ·
[Next: Find the corresponding live Chart field](chart-location.md)
