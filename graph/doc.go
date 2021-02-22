// Package graph provides primitives for building and traversing directed
// graphs.
//
// The implementation is purposely generic and could be used for representing
// any kind of data. In practice, the main application of those graphs is via
// the child "build" package for representing the different messaging
// components of a Bridge, as well as the data flow between those components.
package graph
