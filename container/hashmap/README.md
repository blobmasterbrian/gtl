# hashmap

An unordered collection of key/value pairs with unique keys, backed by a Go
map. `All` exposes the pairs as an `algorithm.Stream2`, `Keys` and `Values`
expose either half as an `algorithm.Stream`, and `Collect` builds a map from
any sequence of pairs. `Map` transforms each pair into a new pair; `MapKeys`
and `MapValues` transform one half and leave the other; `Collapse` turns each
pair into a single value and returns a `slice.Slice`.

`queries.go` holds the methods that leave the map untouched; `mutations.go`
holds the ones that modify it. A map that isn't a `HashMap` gets the methods
through `Wrap(m)`, which shares it rather than copying.

<!-- gendoc:begin -->

## HashMap[K, V]

HashMap is an unordered collection of key/value pairs with unique keys.

| Method | Description | Complexity |
|---|---|---|
| `Collect` | Collect returns a map containing the pairs of src, which may be an iter.Seq2 or an algorithm.Stream2. | Time O(n) expected, space O(n). |
| `New` | New returns a map containing the entries of every map in sources, a later map overriding an earlier one where keys repeat. | Time O(n), space O(n), sized up front. |
| `Wrap` | Wrap gives a map that is not a HashMap the container's methods without copying it. | Time O(1). |
| `All` | All returns the map's pairs as a lazy pipeline, in no particular order. | Time O(1) to create and O(n) to walk, space O(1). |
| `Collapse` | Collapse returns f applied to each entry, collapsing the entries to a slice in no particular order. | Time O(n), space O(n) in a single allocation. |
| `Contains` | Contains reports whether k is a key of the map. | Time O(1) expected. |
| `Delete` | Delete deletes the entries for which pred returns true. | Time O(n) expected, space O(1). |
| `Get` | Get returns the value stored for k and whether k is present. | Time O(1) expected. |
| `Insert` | Insert stores every pair of src, which may be an iter.Seq2 or an algorithm.Stream2, replacing any values already stored under their keys. | Time O(n) expected in the pairs of src, space O(n). |
| `Keep` | Keep returns a new map of the entries for which pred returns true. | Time O(n) expected, space O(k). |
| `Keys` | Keys returns the map's keys as a lazy pipeline, in no particular order. | Time O(1) to create and O(n) to walk, space O(1). |
| `Len` | Len returns the number of entries in the map. | Time O(1). |
| `Map` | Map returns a new map of f applied to each entry. | Time O(n) expected, space O(n), sized up front. |
| `MapKeys` | MapKeys returns a new map with f applied to each key and the values unchanged. | Time O(n) expected, space O(n), sized up front. |
| `MapValues` | MapValues returns a new map with the same keys and f applied to each value. | Time O(n) expected, space O(n), sized up front. |
| `Reject` | Reject returns a new map of the entries for which pred returns false. | Time O(n) expected, space O(k). |
| `Remove` | Remove removes the entry for k. | Time O(1) expected. |
| `Retain` | Retain keeps the entries for which pred returns true and deletes the rest. | Time O(n) expected, space O(1). |
| `Set` | Set stores v under k, replacing any value already stored there. | Time O(1) expected, amortized over growth. |
| `Values` | Values returns the map's values as a lazy pipeline, in no particular order. | Time O(1) to create and O(n) to walk, space O(1). |

<!-- gendoc:end -->
