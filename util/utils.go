package util

// AddrOf takes a value and returns its address as a pointer.
func AddrOf[T any](literal T) *T { return &literal }
