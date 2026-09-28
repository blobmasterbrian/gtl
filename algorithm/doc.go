// Package algorithm provides pipelines over sequences that do no work until
// their result is consumed.
//
// From wraps any push-iterator function, including an iter.Seq and every
// container's All, as a Stream; From2 does the same for pairs. A Stream's
// methods are either adapters, which return a new Stream without doing any
// work, or terminals, which consume the stream and return a value. Adapters
// compose: nothing runs until a terminal runs or the stream is ranged over,
// and ranging stops the whole pipeline as soon as the consumer breaks out.
// No method modifies its input; in-place operations belong to the containers.
//
// A source is ranged over at most once by any method, so single-use sources
// are safe everywhere. A returned stream is reusable exactly when its source
// is.
//
// Seq and Seq2 are constraints satisfied by iter.Seq, Stream and any other
// type of the same shape, for functions that accept either without a
// conversion. Seq on a Stream converts it back for functions that take only
// the standard type.
//
// Complexity notes give the cost of consuming a result in full, with n the
// number of elements in the source and k the number in the result. Building
// an adapter is O(1) in time and space, one closure, whatever the source
// size. Compared with a container's eager methods, a pipeline of d stages
// makes one pass in O(d) extra space rather than d passes with d-1
// intermediate results.
package algorithm
