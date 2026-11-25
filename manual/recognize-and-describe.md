# Recognize and describe occurrences

Sail selects a Spot using a Chart property’s family and exact variant, then
describes the occurrence in the populated Legend. A Feature, Qualifier, or
Barrel adds information to that occurrence; it does not help choose the Spot.

## Mark, family, and variant

The Legend's `mark` names the Chart tag key. The first token in that tag is
the family; subsequent comma-separated tokens are modifiers. A Region or
Landmark marker declares its family and, optionally, a discriminator set:

```go
type PlainJob struct {
    legend.Landmark `sail:"landmark" family:"job"`
    Cached          legend.Qualifier `sail:"qualifier" symbol:"cached"`
}

type ScheduledJob struct {
    legend.Landmark `sail:"landmark" family:"job" discriminate:"scheduled"`
}

type JobLane struct {
    legend.Lane
    Plain     []PlainJob     `sail:"spot" spot:"multi"`
    Scheduled []ScheduledJob `sail:"spot" spot:"multi"`
}
```

In a scope using this Lane and a Legend marked `pipeline`, these properties
select the indicated variants:

```go
Compile struct{} `pipeline:"job,cached"`    // PlainJob, Cached present
Nightly struct{} `pipeline:"job,scheduled"` // ScheduledJob
Test    struct{} `pipeline:"job"`           // PlainJob, Cached absent
```

The family identifies the role. The discriminator set identifies a variant
of that family: selection requires the exact set, not just one matching
discriminator. The plain variant has the empty set. Modifier order does not
change the set. Modifiers are case-sensitive, are not trimmed, and must be
nonempty and unique. A Qualifier such as `cached` records presence on the
selected variant; it is not a discriminator and does not select another Spot.
`pipeline:"job,scheduled,cached"` selects the scheduled variant and then fails:
that variant does not declare `cached`. It does not fall back to the plain
variant. A modifier cannot be both a discriminator and a Qualifier anywhere
in the same family.

## Features record tag values

A Feature names a separate Chart tag key. It records the string under that
key on the selected property, rather than selecting the occurrence:

```go
// In a job Landmark:
Name legend.Feature `sail:"feature" symbol:"name"`

// On a Chart:
Compile struct{} `pipeline:"job" name:"compile"`
```

The populated job's `Name.Get()` is `"compile"`. `Name.Present()` distinguishes
an absent `name` tag from `name:""`. The occurrence's own `Present()` answers
a different question: whether that Region or Landmark occurrence was recorded.
Neither presence check by itself imposes a nonempty-name policy.

## Barrels capture typed declaration values

A **Barrel** captures a provider method’s result in the model. Unlike a
Feature’s tag string, the result can be any type assignable to the Barrel’s
stored value type:

```go
type Descriptive interface {
    About() string
}

// In a Region or Landmark:
Description legend.Barrel[string, Descriptive] `sail:"barrel" populate:"About"`

type CompileJob struct{}

func (CompileJob) About() string { return "Compile the project" }

// On a Chart:
Compile CompileJob `pipeline:"job"`
```

The Barrel declares the stored value type, the provider interface, and the
method to call. Use an exported method with no arguments and one assignable
result. Sail calls the provider on a fresh zero-valued surrogate of the
carrier, not on the live Chart property, so the provider must not rely on
configured instance state. `Get()` reads the value; `Present()` distinguishes
an unavailable capture from a captured zero or nil. An optional Barrel may be
absent when its provider is unavailable. Adding `barrel:"required"` requires
an available provider when capture is active, not a nonzero result. Once a
provider is eligible, a failed call is a population error even for an optional
Barrel.

A Qualifier can gate capture. A bare `legend.Qualifier` records presence;
a Qualifier body can instead contain Barrels. For example:

```go
type DefaultQualifier struct {
    legend.Qualifier `sail:"qualifier"`
    Value legend.Barrel[any, any] `sail:"barrel" populate:"Default" barrel:"required"`
}

// In a Landmark:
Default DefaultQualifier `sail:"qualifier" symbol:"default"`

type Remote string
func (Remote) Default() Remote { return "origin" }

// On a Chart:
Origin Remote `pipeline:"job,default"`
```

When `default` is present, Sail captures the provider's `Remote` value in the
Barrel. When it is absent, capture is inactive, even though the Barrel is
required. Here `any` permits lookup of the named method and storage of its
result; a nonempty provider interface, by contrast, uses ordinary Go method
sets and exact signatures.

### Inspect a captured value

Use `Get()` when the framework only needs the stored value. When it matters
how that value was declared, a populated Barrel also offers three checked reads:

| If you need to know… | Read |
| --- | --- |
| Which method supplied the value | `Method()` gives the provider method without a bound receiver. |
| What the provider declared and returned | `Result()` preserves the declared result type and the captured value. In the example above, the declared result type is `Remote`, while `Get()` returns the value in an `any` view. |
| What type the Chart property declared | `CarrierType()` gives the original property type, not the surrogate receiver type or the provider’s return type. |

Check `check.OK()` before using any of these results. If no capture occurred,
for example because a Qualifier was inactive, the inspection is unavailable.
A captured zero or nil still has inspection metadata. None of these reads
invokes the provider again or addresses a live Chart field.

## Acceptance follows selection

A Landmark can implement `Accept() typerule.RuleAPI` to constrain its Chart
carrier. Sail uses Intel’s [`github.com/qlustra/intel/rules/typerule`](https://pkg.go.dev/github.com/qlustra/intel/rules/typerule)
package directly for this hook. For example, a `Toggle` Landmark can require
the exact Go type `bool`:

```go
func (Toggle) Accept() typerule.RuleAPI {
    return typerule.Type().Is[bool]()
}
```

This accepts `bool` and its aliases, not a distinct defined type with
underlying type `bool`. Sail checks the original declared Chart property type
**after** the family and discriminator set select the Landmark. It does not
test a Feature value, a Barrel result, or the live field’s contents. A rejection
fails population; it does not make Sail try another Spot. Regions instead open
struct-valued nested scopes.

Selection, admission, and description are distinct from access to live Chart
storage. See the [Legend Go documentation](https://pkg.go.dev/github.com/qlustra/sail/legend)
and [Sail Go documentation](https://pkg.go.dev/github.com/qlustra/sail) for
the exact APIs and contracts.

[Previous: Shape a Legend](shape-a-legend.md) · [Contents](readme.md) ·
[Next: Compile and populate](compile-and-populate.md)
