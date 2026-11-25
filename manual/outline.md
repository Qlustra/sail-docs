# Sail manual — candidate outline

## 1. Sail’s world

Start with the two people writing Go code. The **framework author** defines
what an application can declare. The **application author** makes those
declarations in a **Chart**. The framework’s **Legend** is both the model it
wants to work with and the place it declares the grammar. Sail compiles that
grammar into a reusable **Navigator**, then uses it to populate a fresh
instance of the same Legend type from a Chart type.

```text
Legend type ── compile ──► Navigator
Chart type + Navigator ── populate ──► populated Legend (same Legend type)
```

Use one declaration from the pilot to make the vocabulary concrete: the
Legend declares `pipeline` as the tag key and `job` as a permitted family;
the Chart says `pipeline:"job" name:"compile"`; the framework reads a job
named `compile` from the populated Legend. Introduce **occurrence** as that
particular job in the model, as distinct from the rule that permits jobs.

Close with the boundary the rest of the manual relies on: the Chart type is
not the live Chart instance, and neither is runtime input such as CLI
arguments. Sail interprets the declaration; the framework decides what to
do with the result. Introduce the remaining terms where they first matter,
rather than putting a full glossary up front.

## 2. Shape a Legend

- **Regions and Landmarks:** a nested scope versus a leaf. Both can carry
  descriptive values; only Regions contain further occurrences.
- **Lanes and Spots:** ordered groups versus permitted members. Cover single,
  repeated, and required Spots, and the distinction between a reusable Lane
  body and each field that uses it.
- **Lane policies:** what `unambiguous` checks and what `eager` consumption
  commits to. Distinguish order between Lanes from order within a Lane.

## 3. Recognize and describe occurrences

- **Mark, family, and discriminators:** which Chart tag Sail reads, which
  role a field declares, and how an exact discriminator set selects a variant.
- **Features, Qualifiers, and Barrels:** a tag value, a modifier’s presence,
  and a typed value from a declaration provider. Explain how a Qualifier can
  gate Barrel capture without selecting a different variant.
- **Landmark acceptance and presence:** a type rule checks a selected carrier;
  it does not choose the Spot. Distinguish absence from a present empty, zero,
  or nil value when reading the model.

## 4. Compile and populate

- `Compile` validates the Legend; warnings are advisories on the Navigator,
  not evidence that a particular Chart will pass.
- `Build` and `Populate` supply Chart type evidence in different ways and
  both produce fresh Legends. Cover Go-visible properties and promotion,
  value-only paths, and the open-Chart rule: out-of-scope fields are ignored,
  but a selected field’s processing failure is not.
- Follow one declaration through selection and into the model. Show where
  failures arise and what the framework can do with them, without treating
  internal diagnostic strings as a compatibility contract.

## 5. Navigate a populated Legend

- **Shared vocabulary:** Region and Landmark occurrences are tokens. Lanes
  and Spots group them but are not steps. Selectors match stored pointers
  (including via interfaces), optionally by Lane and predicate.
- **Survey:** visit tokens in Chart or Legend sibling order; compose rules
  and pass reducer prefixes to descendants. Filtering an encounter does not
  prune the walk.
- **Next:** test one immediate step and tell no match, one match, and
  ambiguity apart. Only a unique selection can advance.
- **Prepared Courses:** a Plot describes a chain of selectors; preparation
  checks the grammar without running predicates. First, unique, and all
  route searches count *complete routes*, not matches at intermediate steps.

## 6. Find the corresponding live Chart field

- A `ChartLocation` records the exact root type and a field path, not a live
  reference. Its path is in the Chart, not the populated Legend.
- `Resolve` checks a non-nil pointer to the exact Chart root type and gives
  checked access to the corresponding field. Contrast this with a Feature
  read or a Barrel provider call: it is access to live storage, not a
  declaration read, argument parser, or dispatcher.

The existing [compile-and-populate pilot](compile-and-populate.md) is a
test of the writing style. It samples material from the opening chapters;
it need not become a chapter verbatim. The manual should explain how the
capabilities fit together without requiring a runnable example. Keep exact
call contracts in Go docs and task walkthroughs in the examples.
