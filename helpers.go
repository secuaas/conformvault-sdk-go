package conformvault

import "strconv"

// String returns a pointer to the given string value.
func String(v string) *string { return &v }

// Bool returns a pointer to the given bool value.
func Bool(v bool) *bool { return &v }

// Int returns a pointer to the given int value.
func Int(v int) *int { return &v }

// itoa converts an int to a string (shorthand for strconv.Itoa).
func itoa(v int) string { return strconv.Itoa(v) }
