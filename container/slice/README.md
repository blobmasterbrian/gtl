# slice

A `[]T` with the library's algorithms as methods, so a pipeline over a slice
can be written as a chain. It is the sequence container: ordered, indexable,
and growable, the counterpart of C++'s `vector`.

`queries.go` holds the methods that leave the slice untouched; `mutations.go`
holds the ones that modify it in place. `All` hands a chain over to package
`algorithm`. A slice that isn't a `Slice` gets the methods through
`Wrap(&s)`, which takes its address so that in-place methods update it.

<!-- gendoc:begin -->

## Slice[T]

Slice is a []T with the library's eager algorithms as methods, so a pipeline can be written as a chain.

| Method | Description | Complexity |
|---|---|---|
| `Collect` | Collect returns a slice containing the elements of src, which may be an iter.Seq or an algorithm.Stream. | Time O(n), space O(n). |
| `New` | New returns a slice containing items. | Time O(n), space O(n), sized up front. |
| `Wrap` | Wrap gives a slice that is not a Slice the container's methods without copying it. | Time O(1), no allocation. |
| `All` | All returns the elements as a lazy pipeline, in order. | Time O(1) to create and O(n) to walk, space O(1). |
| `Apply` | Apply replaces each element with f applied to it. | Time O(n), space O(1). |
| `Delete` | Delete removes the elements for which pred returns true and shortens the slice to fit. | Time O(n), space O(1). |
| `Each` | Each calls f for each element. | Time O(n), space O(1). |
| `Exists` | Exists reports whether pred returns true for any element. | Time O(n) in the worst case, stopping at the first match. Space O(1). |
| `Find` | Find returns the first element for which pred returns true. | Time O(n) in the worst case, stopping at the first match. Space O(1). |
| `Keep` | Keep returns the elements for which pred returns true. | Time O(n), space O(k). |
| `Map` | Map returns f applied to each element. | Time O(n), space O(n) in a single allocation. |
| `Reduce` | Reduce folds the slice into a single value, starting from init. | Time O(n), space O(1) beyond the accumulator. |
| `Reject` | Reject returns the elements for which pred returns false. | Time O(n), space O(k). |
| `Retain` | Retain keeps the elements for which pred returns true and shortens the slice to fit. | Time O(n), space O(1). |
| `Satisfies` | Satisfies reports whether pred returns true for every element. | Time O(n) in the worst case, stopping at the first failure. Space O(1). |

<!-- gendoc:end -->
