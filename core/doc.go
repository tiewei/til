// Package core encapsulates the core logic for performing operations on parsed
// Bridge descriptions, such as resolving dependencies and translating to
// Kubernetes API objects.
package core

// NOTE(antoineco): a lot of the concepts found in this package were borrowed
// from a major refactoring of the HashiCorp Terraform core in 2015, which
// resulted in a much cleaner and testable software architecture that is still
// the foundation of the product today.
// Terraform relies heavily on graphs to apply its "plans", and back then
// developers had to rethink the design of this central piece in-depth and
// reduce its complexity in order to be able to keep maintaining it.
//
// For further context, please see
//  * pull-request that introduced the refactored code: https://github.com/hashicorp/terraform/pull/1010
//  * walkthough of those changes by Mitchell Hashimoto: https://www.youtube.com/watch?v=suiuaKgaQ74
