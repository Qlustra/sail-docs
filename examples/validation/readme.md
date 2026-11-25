# Tell validation boundaries apart

This example demonstrates three distinct checks: `Compile` rejects a Legend
whose two lanes reuse the same lane ID; `Build` accepts a Chart containing its
required option and rejects one without it; and `Course.Prepare` rejects a
second waypoint after a terminal Landmark. These failures happen at different
stages, so the example labels each result by stage rather than exposing
diagnostic wording.

Run it from `examples/`:

```sh
go run ./validation
```

Output:

```text
Compile: rejected duplicate lane IDs
Build: accepted Chart with required option
Build: rejected missing required option
Prepare: rejected waypoint after terminal Landmark
```

## Walkthrough

In [`main.go`](main.go), the `duplicateLanes` Legend declares two uses of lane ID `options`. Its
`Compile` call must return no Navigator and an error. The valid `language`
Legend compiles, and `validChart` supplies the required `option` Spot; the
program also checks that the resulting model records the Spot as present. An
empty Chart then fails `Build` because it has no occurrence of that required
Spot.

Finally, the program plots two waypoints selecting the `option` Landmark.
`Prepare` rejects that course because a Landmark is terminal: no next waypoint
can follow it. Each rejection is asserted by the presence of an error and the
absence of a successful result, never by parsing English diagnostic text.

## Policy

Diagnostics are intended for people and may change wording. The program prints
stable stage-and-condition labels and deliberately omits raw diagnostics.
The framework decides how to present failures and whether to impose stricter
policies; Sail performs the checks shown here. See the
[compile and populate manual](../../manual/compile-and-populate.md)
and [navigation manual](../../manual/navigation.md) for the corresponding
contracts.
