# gtl

GTL is a STL (Standard Template Library) implementation for Golang.  
  
The Go Template Library is a set of Golang template structs to provide common programming data structures and functions.
It is a library of containers, algorithms, and functions. It provides data structures such as vectors, sets, etc and
provides enumeration operations such as map, reduce, filter for them.  
  
This library seeks to conform with the clear and conciseness of go convention. As such, certain naming conventions were
abandoned for clarity. Ex: there are two filter functions, `keep` and `reject` which are more explicit than `filter`.  
  
A standard template library should seek to achieve high performance on the general-purpose structures and functions
it provides. As such, this library will seek to achieve the best performance while still being generic, foregoing things
like using an element in a loop in favor of accessing elements via index.  
