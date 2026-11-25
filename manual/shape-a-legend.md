# Shape a Legend

A Legend’s shape says which occurrences belong together, which can nest, and
which groups the Chart must supply in order. Define that structure on the same
Go fields the framework will read in the populated model.

## Scopes and occurrences

A **Region** is an occurrence that opens a nested scope. A **Landmark** is a
leaf occurrence. Either can carry descriptive values, but only a Region has
Lanes for further occurrences. The Legend itself is the root scope.

Each scope contains ordered **Lane uses**. A Lane groups permitted members;
each **Spot** in its body permits occurrences of a particular Region or
Landmark. For example, a pipeline can contain stages, each with jobs:

```go
type PipelineLegend struct {
    legend.Legend `sail:"legend" mark:"pipeline"`
    Stages        StageLane `sail:"lane"`
}

type StageLane struct {
    legend.Lane
    Items []Stage `sail:"spot" spot:"multi,required"`
}

type Stage struct {
    legend.Region `sail:"region" family:"stage"`
    Jobs          JobLane `sail:"lane"`
}

type JobLane struct {
    legend.Lane
    Items []Job `sail:"spot" spot:"multi,required"`
}

type Job struct {
    legend.Landmark `sail:"landmark" family:"job"`
    Name            legend.Feature `sail:"feature" symbol:"name"`
}
```

Here `pipeline` names the Chart tag Sail reads; `stage` and `job` are
families permitted in different scopes. A Chart can declare a stage as a
struct-valued property and jobs inside that struct:

```go
type ReleaseChart struct {
    Build struct {
        Compile struct{} `pipeline:"job" name:"compile"`
    } `pipeline:"stage"`
}
```

The stage and job are separate occurrences. `Name` describes the selected
job; it is not another Spot. The [next chapter](recognize-and-describe.md)
explains how tags select and describe occurrences.

## Reuse a Region at any depth

When each level has the same grammar, use one Region type rather than a
different model type per level. For example, a `Section` can contain more
sections by using a Lane with a repeated Spot for `Section`:

```go
type Section struct {
    legend.Region `sail:"region" family:"section"`
    Name          legend.Feature `sail:"feature" symbol:"name"`
    Children      SectionLane   `sail:"lane"`
}

type SectionLane struct {
    legend.Lane
    Sections []Section `sail:"spot" spot:"multi"`
}
```

A Legend can use `SectionLane` at its root as well. Sail compiles the
recursive grammar without unrolling it to a fixed depth. Each populated
section is still a separate occurrence, and the same `*Section` selector can
match at any level. The recursive Spot here is optional, so a section can
have no children. The Chart declares a finite tree of struct-valued fields;
the recursive Legend does not make the Chart a pointer graph.

## Spot cardinality

A `single` Spot is a value field and allows at most one occurrence. A `multi`
Spot is a slice and allows any number. Add `,required` to require at least one:

| Spot tag | Occurrences |
| --- | --- |
| `single` | zero or one |
| `single,required` | exactly one |
| `multi` | zero or more |
| `multi,required` | one or more |

Each matching Chart property contributes one occurrence. A `multi` Spot
collects occurrences; it does not expand a collection held by a Chart property.
Requiredness sets a minimum, not the point at which a Lane stops consuming.
When a Lane closes, Sail checks its required Spots. Another match to an already
full single is an error, not an instruction to advance to the next Lane.

## Lane uses and consumption

Members permitted by one Lane can appear in any order. The Lane uses in a
Legend or Region, however, are ordered. Sail consumes matching Chart
properties in Chart order in the current Lane. When a candidate does not
match a Spot there, Sail closes that Lane and tries the same candidate in the
next one, subject to the required minima. It also checks minima at the end of
the scope.

A Lane body can be reused in multiple fields. Each field is a separate use,
with its own position and counts. An optional `lane:"..."` ID and the
`exploration` policy belong to the field use, not the shared body. A Lane ID
must be unique among the Lane uses of its enclosing scope. The same ID can
be used in another scope and can be used to [filter encounters](navigation.md#select-tokens).

## Lane policies

`unambiguous` is the default. At compilation, Sail compares the Lane’s
selectors with following Lanes in the same scope, including the first later
Lane with a required Spot. It rejects overlapping selectors and does not
look beyond that required Lane. This is a grammar check, not lookahead during
population.

Use `exploration:"eager"` when an earlier Lane should consume its matches
even if a later Lane also admits them. For example, suppose an eager Lane
accepts jobs and blocks a later Lane that requires a job. If the only available
selector is `job`, the earlier Lane consumes every job and compilation rejects
the unreachable requirement. If the later Lane also accepts a `checkpoint`
that the eager Lane does not, a Chart can use a checkpoint to enter the later
Lane and then supply its job.

Sail can warn about that dependency without rejecting the grammar. Population
never backtracks to leave an earlier Lane or moves an already consumed job to
the later Lane. A field ignored by the grammar cannot act as the checkpoint.

These rules describe the grammar, not framework policy such as whether job
names are unique or whether a job should run. For exact marker and call
contracts, see the [Legend Go documentation](https://pkg.go.dev/github.com/qlustra/sail/legend)
and [Sail Go documentation](https://pkg.go.dev/github.com/qlustra/sail).

[Contents](readme.md) · [Next: Recognize and describe occurrences](recognize-and-describe.md)
