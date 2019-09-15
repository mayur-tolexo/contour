package favour

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsDefault(t *testing.T) {
	i := 0
	a := assert.New(t)
	a.True(IsDefault(i), "Int Default")
	a.True(IsDefault(0.0), "Float Default")
	a.True(IsDefault(""), "String Default")
	a.True(IsDefault(nil), "Interface Default")
}
