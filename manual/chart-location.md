# Find the corresponding live Chart field

A populated Legend describes occurrences from a Chart **type**; it does not
hold the application’s live Chart instance. When a framework needs the actual
field corresponding to an occurrence, the token’s **ChartLocation** bridges
that gap. It records where the field is, and the framework can resolve that
location against a suitable instance later.

This is distinct from reading a Feature’s captured tag value or a Barrel’s
captured provider result. Those describe declarations in the populated model;
resolving a location gives checked access to live storage. Sail does not parse
runtime arguments, choose an operation, or dispatch it merely by resolving a
field.

## A location is a snapshot of correspondence

Present Region and Landmark occurrences expose `ChartLocation()`. Check its
result: an absent occurrence has no location. A successful location records
the exact root Chart Go type and a root-relative Go field-index path.
`ChartType()` gives an independent type interrogator, and `Index()` returns a
copy of the path. The path counts ignored Chart fields and includes embedding
steps where applicable. It is not an index into a Spot’s model slice.

The same location can be resolved against different instances of its recorded
root type. Model occurrence presence, a Feature’s presence, and the contents of
a live field are separate facts.

## Resolve against an instance

`Resolve` requires a non-nil pointer directly to the exact recorded Chart root
type. A struct value, a nil pointer, a pointer to a pointer, or a distinct
defined type with the same fields is not that root. A Go type alias denotes
the same type. Resolution checks access to the field and returns a handle to
the actual field in the supplied instance, not a copy of the model token.

For example, suppose `job` is a Landmark selected during [navigation](navigation.md)
and its Chart property is declared as `Compile string`. A framework can use
the location to change the live string independently of any captured tag:

```go
location, check := job.ChartLocation()
if !check.OK() {
    return fmt.Errorf("locate job: %s", check.Why())
}
field, check := location.Resolve(&chart)
if !check.OK() {
    return fmt.Errorf("resolve job: %s", check.Why())
}
address, check := field.Address()
if !check.OK() {
    return fmt.Errorf("address job: %s", check.Why())
}
destination, ok := address.(*string)
if !ok {
    return fmt.Errorf("job field: expected *string, got %T", address)
}
*destination = "release"
```

`Address` exposes the address as `any`, so the framework checks the expected
type before writing. A resolved struct field can instead be examined through
the handle’s `AsStruct` operation. Any subsequent method invocation and its
policy belong to the framework.

Resolution does not allocate missing intermediate pointers and can fail when a
path is inaccessible. The location and model do not retain the instance, while
the resolved handle refers to its live storage. The framework owns the usual
Go lifetime and synchronization considerations. Location resolution is not a
live traversal of the Chart; navigation continues to operate on populated Legend
tokens.

[Previous: Navigate a populated Legend](navigation.md) · [Contents](readme.md)
