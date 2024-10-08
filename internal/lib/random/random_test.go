package random

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRandomString(t *testing.T) {
	testSizes := []int{1, 5, 10, 20, 30, 50}
	for index, size := range testSizes {
		t.Run(
			fmt.Sprintf("Test №%d: string size = %d", index+1, size),
			func(t *testing.T) {
				str1 := NewRandomString(size)
				str2 := NewRandomString(size)

				assert.Len(t, str1, size)
				assert.Len(t, str2, size)

				assert.NotEqual(t, str1, str2)
			},
		)
	}
}
