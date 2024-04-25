// Package stdlib provides a few macro types which are replaced or expanded to platform specific types
// by a specific renderer. The standard types end with a ! and each build-in renderer is aware of them
// (which may also mean, that it just ignored). You cannot extend the standard types, because the
// semantic may be very specific (e.g. marshalling/unmarshalling) and its usage may be scattered
// through an entire renderer.
package stdlib
