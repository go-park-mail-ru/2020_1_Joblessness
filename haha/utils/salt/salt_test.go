package salt

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFlow(t *testing.T) {
	password := "password"
	wrongPas := "wrong"

	hash, err := HashAndSalt(password)
	assert.NoError(t, err)

	res := ComparePasswords(hash, password)
	assert.True(t, res)

	res = ComparePasswords(hash, wrongPas)
	assert.False(t, res)
}
