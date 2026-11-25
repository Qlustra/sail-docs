# Compile and populate

Compilation asks whether the framework has authored a valid grammar.
Population asks whether a particular application declaration fits that grammar,
and produces a model if it does. A valid grammar need not accept every Chart;
a failed Chart does not invalidate its Navigator.

## Compile the Legend once

`sail.Compile[L]()` checks the authored Legend `L` and returns a reusable
Navigator. Compilation checks field roles, tag syntax, and relationships in
the grammar, such as overlapping Lane selectors. It does not need a Chart or
invoke its Barrel providers. A fatal error returns no Navigator.

`navigator.Warnings()` returns a snapshot of nonfatal advisories on a successful
Navigator. An eager Lane, for example, can leave a grammar viable but make a
later requirement depend on a field that ends the earlier Lane. Sail does not
print advisories or decide whether to treat them as fatal. See
[Lane policies](shape-a-legend.md#lane-policies) for what that arrangement means
for Charts.

## Supply Chart type evidence

| What you have | Entry point |
| --- | --- |
| The Chart type is known to the caller. | `navigator.Build[ReleaseChart]()` |
| The Chart type is determined by the caller’s value. | `navigator.Populate(chart)` |

Both interpret a struct type and return a fresh populated `*L`. `Build` does
not construct a Chart. `Populate` uses the argument’s dynamic type; it does
not read its field values or retain the instance. A nil interface supplies
no type, and a pointer to a Chart is not a struct root, even when non-nil.
Use a pointer later when [resolving a live field](chart-location.md), not to
populate the model.

## What Sail reads from a Chart

Sail considers nonblank, named properties in Go declaration order. Go’s
visibility rules still apply: promoted properties appear at the embedding’s
position, and shadowed or ambiguous promoted fields are not revived. An
embedding is not itself an occurrence even if it carries a grammar tag.
Selected properties must have value-only field paths and non-pointer carrier
types. A Region additionally needs a struct carrier so Sail can interpret
its nested declaration.

The Chart is **open**: fields without the grammar’s tag, unrecognized
families, and families outside the current scope do not move the Lane cursor.
This is what lets the application keep its own fields alongside declarations.
It does not excuse a missed required Spot, a recognized field out of Lane
order, or invalid modifiers on a participating family.

For example, this Chart property gives a pipeline grammar one candidate in
the `job` family:

```go
Compile struct{} `pipeline:"job" name:"compile"`
```

When a Lane admits it, Sail checks the selected Spot’s constraints and
populates one job, including its `name` Feature. An unrelated
field between two jobs does not end the Lane. But if a second job selects a
single Spot that is already full, Sail reports a cardinality violation rather
than retrying that field in a later Lane. A Landmark acceptance failure or a
failed eligible provider likewise fails population rather than changing the
match. See [Recognize and describe occurrences](recognize-and-describe.md)
for the selection and value rules.

## Read the outcome at the right phase

| Outcome | What to check |
| --- | --- |
| Legend compilation fails | Check the Legend’s marker and field declarations, selectors, and relationships between Lanes. There is no Navigator to reuse. |
| Compilation succeeds with warnings | Inspect the advisories and decide whether the framework accepts that grammar. A warning is not a result for any particular Chart. |
| Chart population fails | Check the field, its carrier type and modifiers, and the selected Spot or enclosing Lane. There is no partial model; the Navigator may be reused. |
| Chart population succeeds | Read ordinary Go model fields, including their presence and value accessors, or [navigate the populated Legend](navigation.md). |

Show diagnostic text to readers rather than parsing it as a data format.
The framework owns presentation and any additional domain validation, such
as rejecting duplicate job names.

[Previous: Recognize and describe occurrences](recognize-and-describe.md) ·
[Contents](readme.md) · [Next: Navigate a populated Legend](navigation.md)
