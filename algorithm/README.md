# algorithm

Pipelines over sequences that do no work until their result is consumed,
written as method chains on a `Stream`. Reach for this package when the input
is large, comes from a container or a stream rather than a slice, or won't be
used in full.

`From` turns anything with the shape of an `iter.Seq` into a `Stream`, and
every container's `All()` returns one directly. `adapters.go` holds the
methods that return a new stream; `terminals.go` holds the ones that consume
it and return a value. `Stream2` is the same for pairs. `Seq()` hands a stream
back to code that takes the standard type.

The eager counterparts of these operations are methods on the containers.

<!-- gendoc:begin -->

## Seq[T]

Seq is satisfied by every push-iterator function type over single values, including iter.Seq and Stream, so a function taking an S accepts either.

## Seq2[K, V]

Seq2 is the Seq counterpart for pairs, satisfied by iter.Seq2 and Stream2.

## Stream[T]

Stream is a sequence with this package's pipeline operations as methods.

| Method | Description | Complexity |
|---|---|---|
| `From` | From wraps a sequence as a Stream. | Time O(1), no allocation. |
| `Collect` | Collect runs the pipeline and returns its elements as a slice. | Time O(n), space O(n). |
| `Count` | Count returns the number of elements. | Time O(n), space O(1). |
| `Drop` | Drop returns the stream without its first count elements. | Time O(n), space O(1). |
| `Each` | Each calls f for each element. | Time O(n), space O(1). |
| `Exists` | Exists reports whether pred returns true for any element. | Time O(n) in the worst case, stopping at the first match and stopping the source with it. Space O(1). |
| `Find` | Find returns the first element for which pred returns true. | Time O(n) in the worst case, stopping at the first match and stopping the source with it. Space O(1). |
| `Keep` | Keep returns the elements for which pred returns true. | Time O(n), space O(1). |
| `Map` | Map returns f applied to each element. | Time O(n), space O(1). |
| `Reduce` | Reduce folds the stream into a single value, starting from init. | Time O(n), space O(1) beyond the accumulator. |
| `Reject` | Reject returns the elements for which pred returns false. | Time O(n), space O(1). |
| `Satisfies` | Satisfies reports whether pred returns true for every element. | Time O(n) in the worst case, stopping at the first failure and stopping the source with it. Space O(1). |
| `Seq` | Seq returns the stream as a plain iter.Seq, for functions that take one. | Time O(1), no allocation. |
| `Take` | Take returns the first count elements, or all of them if there are fewer. | Time O(min(n, count)), space O(1). |

## Stream2[K, V]

Stream2 is the Stream counterpart for pairs, such as a map's entries or a slice's index/value pairs.

| Method | Description | Complexity |
|---|---|---|
| `From2` | From2 wraps a sequence of pairs as a Stream2. | Time O(1), no allocation. |
| `Collapse` | Collapse returns f applied to each pair, collapsing the pairs to single values. | Time O(n), space O(1). |
| `Count` | Count returns the number of pairs. | Time O(n), space O(1). |
| `Each` | Each calls f for each pair. | Time O(n), space O(1). |
| `Keep` | Keep returns the pairs for which pred returns true. | Time O(n), space O(1). |
| `Map` | Map returns f applied to each pair. | Time O(n), space O(1). |
| `Reject` | Reject returns the pairs for which pred returns false. | Time O(n), space O(1). |
| `Seq` | Seq returns the stream as a plain iter.Seq2, for functions that take one. | Time O(1), no allocation. |

<!-- gendoc:end -->
