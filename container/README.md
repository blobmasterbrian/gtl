# container

Data structures, one package each. Every container exposes `All()`, which
returns an `algorithm.Stream` (a `Stream2` for pairs), and `Collect()` to
build one from any sequence, and offers the library's algorithms as methods
that return the same container type. Within each package, `queries.go` holds the methods that leave the
container untouched and `mutations.go` the ones that modify it.

The directory is organization only: import paths nest under `container/`,
but call sites use the package name alone, as in `set.New(...)`.
