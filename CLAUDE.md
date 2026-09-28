# gtl

Conventions for working in this repository.

## Documentation

- Every directory gets a `README.md` describing what it holds and why it is a
  separate unit. The package doc comment, not the README, documents the API.
- The function tables in every `README.md`, including the top-level one, are
  generated from the package sources by `go generate ./...`. Regenerate after
  adding, removing or renaming an exported function, method or type, or after
  changing its doc comment's first sentence or complexity note. `go test ./...`
  fails while a table is stale. Never edit the generated regions by hand.
- Every exported function and method documents its complexity in a paragraph
  that begins with `Time`. A method that wraps the package-level function of
  its own name states the same complexity as that function. `gendoc` fails on
  a missing or mismatched note.
