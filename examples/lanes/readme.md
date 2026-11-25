# Split repeated members into lanes

Reuse the same job lane before and after a required checkpoint. A present
required boundary separates otherwise overlapping selectors; an optional
boundary does not separate them under the default `unambiguous`
policy. Setting the earlier lane to `exploration:"eager"` changes how it
consumes matching jobs, but does not make an optional boundary required.

Run from `examples/`:

```sh
go run ./lanes
```

Output:

```text
required boundary: Before "A", boundary present true, After "B"
missing boundary: population failed
optional boundary + default reuse: compilation rejected
optional boundary + eager Before: compiles; unseparated population fails
checkpoint dependency advisories: 1
eager checkpoint: Before 2 jobs, After "C", checkpoint present true
missing checkpoint boundary: required later job population failed
```

In [`main.go`](main.go), compare the Legend with a required separating Lane
to the one with an optional separator. A Chart with the required separator
populates both uses of the shared job Lane; omitting it fails population.
The checkpoint scenario shows eager exploration assigning
the first two jobs to the multi Spot, then the checkpoint and subsequent job
to the later lane. The required later job makes that route structurally
reachable only with a checkpoint; `Warnings()` reports the dependency as a
nonfatal advisory, printed here by length rather than by assuming a concrete
warning representation. The checkpoint Spot is optional (`single`), while the
later jobs Spot is required (`multi,required`): the latter requires at least
one job, not a checkpoint. The checkpoint is nevertheless a useful entry boundary
for a valid Chart because the eager lane consumes matching jobs before it.

The framework author chooses lane order, exploration and occurrence requirements;
Sail enforces those choices. The framework owns what a checkpoint means and how
to present advisories. The example only records its presence.
