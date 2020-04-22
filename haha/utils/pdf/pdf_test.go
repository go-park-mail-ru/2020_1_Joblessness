package pdf

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenderParser(t *testing.T) {
	for i := 1; i < 11; i++ {
		res := genderToStr(string(i))
		assert.NotEqual(t, "", res)
	}
}
