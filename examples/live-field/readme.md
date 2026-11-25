# Update one live Chart field

Use a populated Landmark's `ChartLocation` to resolve its corresponding field
in a live Chart instance. The location records the exact Chart type and field
path; resolving it returns access to live storage, not a copy of the populated
Legend or its captured tag value.

## Walkthrough

In [`main.go`](main.go), `model` defines the Legend and `field` its Landmark;
`chart.First` and `chart.Second` are the tagged Chart properties. The path
`[1 0]` enters the embedding at index `1`, then its `First` field at index `0`.
`demonstrate`
compiles the model, uses `Navigator.Build[chart]()` as type evidence, and gets
the first Landmark's `ChartLocation`. It resolves that location against a live
`chart`, obtains the field address, and writes through it. It also tries a
distinct `otherChart` with the same layout to show exact root-type checking.
The function returns observed values; [`main_test.go`](main_test.go) asserts
the path, selected-field mutation, unchanged sibling, and rejected type.

The boundary is deliberate: Sail resolves checked access to live storage; it
does not decide what the field means or dispatch framework behavior. See
[Find the corresponding live Chart field](../../manual/chart-location.md) for
the API contract.

Run this example from `examples/`:

```sh
go run ./live-field
```

Expected output:

```text
type evidence: main.chart, field path [1 0] (name "first")
live Resolve: First "before" -> "after"; Second remains "sibling"
different Chart type rejected
```

The printed values show the operation; the test verifies its assertions.
