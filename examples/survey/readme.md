# Survey nested occurrences

This example walks a workflow's nested steps and sibling notes. Run from
`examples/` with `go run ./survey`.

## Walkthrough

In [`main.go`](main.go), `Workflow` declares an `Entries` Lane with `Branch` and
`Note` Spots. Each `Branch` has a `Steps` Lane of `Step` occurrences. `Chart`
supplies two `notify` branches, with `draft` and `deliver` steps, plus a root
`memo` note. `main` compiles `Workflow` and builds its populated Legend from
`Chart`.

The `branchPath` encounter rule appends each branch name to the path; the `step`
and `note` rules print occurrences and pass the path onward unchanged. `main`
runs these rules twice, selecting `ChartOrder` and `LegendOrder` respectively:

```text
Chart order
step /notify/draft
note /memo
step /notify/deliver
Legend order
step /notify/draft
step /notify/deliver
note /memo
```

Chart order follows recognized declarations in `Chart`; Legend order follows
the Lane and Spot order in `Workflow`. The example's framework author defines
the Legend grammar and decides what these workflow entries mean or what action
to take. Sail interprets the tags to populate that Legend and provides Survey's
traversal, ordering, and encounter callbacks; it does not interpret workflow
policy or execute steps. See the [manual's navigation chapter](../../manual/navigation.md#survey-occurrences)
for the Survey contract.
