# set

An unordered collection of unique values, backed by a Go map. It is the
reference for how containers in this library are shaped: `All` exposes the
values as an `algorithm.Stream`, `Collect` builds a set from any sequence,
and the algorithms are available as methods that return another set.

`queries.go` holds the methods that leave the set untouched; `mutations.go`
holds the ones that modify it.

<!-- gendoc:begin -->

## Set[T]

Set is an unordered collection of unique values.

| Method | Description | Complexity |
|---|---|---|
| `Collect` | Collect returns a set containing the elements of src, which may be an iter.Seq or an algorithm.Stream. | Time O(n) expected, space O(n). |
| `New` | New returns a set containing items. | Time O(n), space O(n), sized up front. |
| `Add` | Add adds v to the set. | Time O(1) expected, amortized over growth. |
| `All` | All returns the set's values as a lazy pipeline, in no particular order. | Time O(1) to create and O(n) to walk, space O(1). |
| `Contains` | Contains reports whether v is in the set. | Time O(1) expected. |
| `Delete` | Delete removes the values for which pred returns true. | Time O(n) expected, space O(1). |
| `Insert` | Insert adds every element of src, which may be an iter.Seq or an algorithm.Stream, to the set. | Time O(n) expected in the elements of src, space O(n). |
| `Keep` | Keep returns a new set of the values for which pred returns true. | Time O(n) expected, space O(k). |
| `Len` | Len returns the number of values in the set. | Time O(1). |
| `Map` | Map returns a new set of f applied to each value. | Time O(n) expected, space O(k), with k at most n after deduplication. |
| `Reject` | Reject returns a new set of the values for which pred returns false. | Time O(n) expected, space O(k). |
| `Remove` | Remove removes v from the set. | Time O(1) expected. |
| `Retain` | Retain keeps the values for which pred returns true and removes the rest. | Time O(n) expected, space O(1). |

<!-- gendoc:end -->
