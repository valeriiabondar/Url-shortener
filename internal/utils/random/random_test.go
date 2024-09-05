package random

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewRandomAlias(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{
			name:   "test 1",
			length: 3,
		},
		{
			name:   "test 2",
			length: 5,
		},
		{
			name:   "test 3",
			length: 6,
		},
		{
			name:   "test 4",
			length: 10,
		}, {
			name:   "test 5",
			length: 20,
		}, {
			name:   "test 6",
			length: 30,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			str1 := NewRandomAlias(tt.length)
			time.Sleep(1 * time.Nanosecond)
			str2 := NewRandomAlias(tt.length)
			time.Sleep(1 * time.Nanosecond)

			t.Logf("Test: %s, Generated strings: str1 = %s, str2 = %s", tt.name, str1, str2)

			assert.Len(t, str1, tt.length)
			assert.Len(t, str2, tt.length)

			assert.NotEqual(t, str1, str2)
		})
	}
}
