# gtl

Generic containers and algorithms for Go, in the spirit of C++'s STL and
Elixir's `Enum` and `Stream`.

The standard library's `slices`, `maps` and `iter` packages cover the basics.
gtl builds on them with a fuller set of algorithms, one consistent API across
containers, and containers the standard library doesn't have.

## Layout

- `container`: data structures, one package each (`slice`, `set`, `hashmap`,
  ...). Their methods are the eager algorithms: each call runs now and returns
  the same container type, or modifies it in place. Every container exposes
  `All()`, which starts a lazy pipeline, and `Collect()` to build one from a
  sequence.
- `algorithm`: lazy pipelines over sequences, written as method chains on a
  `Stream`. Nothing runs until the result is ranged over or collected, and a
  pipeline stops as soon as its consumer does.

Directories nest for organization only; call sites use the package name
alone, as in `set.New(...)` and `algorithm.From(...)`.

Anything with the shape of an `iter.Seq`, including standard slices and maps
and containers from other libraries, enters a pipeline through
`algorithm.From`.

## Example

```go
users := slice.New(ada, bob, cy)

// Eager: each step runs now and returns a Slice.
active := users.Keep(User.IsActive)

// Lazy: one pass, stops after three matches.
names := users.All().Keep(User.IsActive).Map(User.Name).Take(3).Collect()

// Containers interoperate through Stream.
ids := set.Collect(active.All().Map(User.ID))

// A plain []User gets the methods through Wrap, without copying.
raw := loadUsers()
slice.Wrap(&raw).Retain(User.IsActive)

// A standard sequence enters with From; Seq hands a Stream back out.
sorted := slices.Sorted(algorithm.From(maps.Keys(index)).Keep(isLive).Seq())
```

## Packages

<!-- gendoc:begin -->

### algorithm

Package algorithm provides pipelines over sequences that do no work until their result is consumed. ([README](algorithm/README.md))

#### Seq[T]

Seq is satisfied by every push-iterator function type over single values, including iter.Seq and Stream, so a function taking an S accepts either.

#### Seq2[K, V]

Seq2 is the Seq counterpart for pairs, satisfied by iter.Seq2 and Stream2.

#### Stream[T]

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

#### Stream2[K, V]

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

### container/hashmap

Package hashmap provides an unordered collection of key/value pairs. ([README](container/hashmap/README.md))

#### HashMap[K, V]

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

### container/set

Package set provides an unordered collection of unique values. ([README](container/set/README.md))

#### Set[T]

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

### container/slice

Package slice provides a slice with the library's algorithms as methods. ([README](container/slice/README.md))

#### Slice[T]

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

## Naming

Filtering is `Keep` and `Reject` rather than `filter`, which doesn't say which
side it keeps. Functions over key/value pairs (`iter.Seq2` in `lazy`, maps in
`eager`) carry a `2` suffix, as in the standard library.

## Requirements

Go 1.27 or later. The chainable methods rely on generic methods.

## Development

```
go test ./...
go vet ./...
go generate ./...   # refresh the function tables in the READMEs
```

The function tables in every README are generated from the sources. The test
suite fails while they are stale.
