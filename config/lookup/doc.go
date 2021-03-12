// Package lookup contains interfaces for looking up data in a Bridge configuration.
//
// Those are kept separate from the "config" package so we can introduce
// opinions on top of the generic "domain" (the Bridge config), and to avoid
// circular dependencies.
package lookup
