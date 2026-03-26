package program

import (
	"cmp"
	"fmt"
	"os"

	"go.n16f.net/uuid"
)

func OptionalEnvironmentVariable(name, defaultValue string) string {
	return cmp.Or(os.Getenv(name), defaultValue)
}

func EnvironmentVariable(name string) string {
	value := os.Getenv(name)
	if value == "" {
		Abort("missing or empty environment variable %q", name)
	}

	return value
}

func OptionalBooleanEnvironmentVariable(name string, defaultValue bool) bool {
	s := OptionalEnvironmentVariable(name, fmt.Sprintf("%v", defaultValue))

	var value bool

	switch s {
	case "true":
		value = true
	case "false":
		value = false
	default:
		Abort("invalid environment variable %q: invalid boolean %q", name, s)
	}

	return value
}

func UUIDEnvironmentVariable(name string) uuid.UUID {
	s := EnvironmentVariable(name)

	var id uuid.UUID
	if err := id.Parse(s); err != nil {
		Abort("invalid environment variable %q: invalid UUID: %v", name, err)
	}

	return id
}
