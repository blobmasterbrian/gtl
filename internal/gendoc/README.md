# gendoc

Generates the function tables in every `README.md`, including the top-level
one, from the package sources, so the tables cannot drift from the code or
from each other. `go generate ./...` runs it; its test fails while any table
is stale, any exported function or method lacks a complexity note, or a
wrapper method's note disagrees with the function it wraps.
