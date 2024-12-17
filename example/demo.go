// Copyright 2024 Alex Athanasopoulos.
// SPDX-License-Identifier: Apache-2.0

package example

// Demo - the interface for remote calls
type Demo interface {
	// Repeat returns a []string consisting repeated elements
	// It demonstrates  inputs and outputs of various types.
	// It is also an easy way to generate large responses.
	Repeat(s string, count int) ([]string, error)

	// Time takes no arguments, returns multiple values, and does not return an error.
	Time() (hours, minutes, seconds int)

	// TimeStruct returns a struct
	TimeStruct() (Time, error)

	// TimePointer returns a struct pointer
	TimePointer() *Time

	// Wait waits the specified number of seconds
	// It can be used to test a long-running response.
	Wait(seconds int) error

	// Ping should return no errors
	// can be used to verify that communication is good
	Ping() error

	// Error returns a string and an error.
	// Tests error handling.
	Error() (string, error)
}

type Time struct {
	Hour, Minute, Second int
}
