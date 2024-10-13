package program

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitCommandName(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		s     string
		names []string
	}{
		{"",
			[]string{}},
		{"  	",
			[]string{}},
		{"foo",
			[]string{"foo"}},
		{"	 	foo ",
			[]string{"foo"}},
		{"foo bar",
			[]string{"foo", "bar"}},
		{" foo	bar		baz	",
			[]string{"foo", "bar", "baz"}},
	}

	for _, test := range tests {
		label := fmt.Sprintf("%q", test.s)
		names := splitCommandName(test.s)
		assert.Equal(test.names, names, label)
	}
}
