# Sail manual

Sail requires Go 1.27 or newer. It uses generic methods, so a call such as
`navigator.Build[ReleaseChart]()` supplies a type argument to a method.

## Sail’s world

A framework needs a model it can work with. It may also need to say what an
application is allowed to declare. Sail lets the framework author express both
in one Go type: the **Legend**. Its tagged fields describe the grammar; after
population, those same fields hold the model the framework reads.

The application author writes a **Chart**: a Go struct that uses the framework’s
tags. The framework compiles its Legend once into a **Navigator**. The Navigator
holds the grammar, not a set of application values. Given a Chart type, it
applies that grammar and fills a fresh instance of the Legend type.

```text
Legend type ── compile ──► Navigator
Chart type + Navigator ── populate ──► populated Legend (same Legend type)
```

Consider a pipeline framework. Its Legend says to look under the tag key
`pipeline`, permits the family `job`, and says a job can carry a name. The
application declares:

```go
type ReleaseChart struct {
    Compile struct{} `pipeline:"job" name:"compile"`
}
```

The resulting Legend contains a job named `compile`. The rule that permits
jobs is part of the authored grammar; this particular populated job is an
**occurrence**. If the Chart declared a second job in a repeated position,
the framework would read two occurrences of the same kind, not two rules.

The Chart declaration is also distinct from a **live Chart instance**. A live
instance can hold configured receivers or mutable state; Sail can populate
the Legend without one. Runtime inputs, such as CLI arguments or a request to
run the `compile` job, are something else again. The framework decides how
to act on them and whether it needs the corresponding live field. Sail does
not run jobs, parse invocations, or decide what these tags mean to a user.

## Read in capability order

| Chapter | What it explains |
| --- | --- |
| [Shape a Legend](shape-a-legend.md) | Nested scopes and leaves, Lane uses, Spot cardinality, and how Lanes consume occurrences. |
| [Recognize and describe occurrences](recognize-and-describe.md) | Selecting by family and variant; reading tags, modifier presence, and typed declared values; checking a Landmark’s carrier. |
| [Compile and populate](compile-and-populate.md) | Validate the grammar, interpret Chart declarations, and tell which phase a failure belongs to. |
| [Navigate a populated Legend](navigation.md) | Survey encounters, take a unique immediate step, or search for complete routes. |
| [Find the corresponding live Chart field](chart-location.md) | Keep a field’s location separate from its storage and resolve it against a live instance. |

Use the manual to choose and combine capabilities. Use `go doc` in your own
module for exact signatures and call contracts in the Sail version you use;
the [Sail API](https://pkg.go.dev/github.com/qlustra/sail) and
[Legend marker API](https://pkg.go.dev/github.com/qlustra/sail/legend) are
reference entry points. [Runnable examples](../examples/readme.md) are optional:
use them to try a task, not to fill in a missing precondition in the manual.
