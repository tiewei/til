package k8s

import "bridgedl/lang/k8s"

// This file contains aliases to variables and functions from the lang/k8s
// package that are commonly used in component implementations, so that those
// components only have to depend on a single import.

var NewDestination = k8s.NewDestination
var DestinationCty = k8s.DestinationCty

var ObjectReferenceCty = k8s.ObjectReferenceCty
var IsObjectReference = k8s.IsObjectReference
