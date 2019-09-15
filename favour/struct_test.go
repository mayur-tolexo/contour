package favour

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStructTag(t *testing.T) {
	a := assert.New(t)

	type basic struct {
		Name   string  `json:"name"`
		Age    float64 `json:"age"`
		Number int     `json:"number"`
	}

	b := basic{Name: "Mayur", Age: 24, Number: 123123}

	v, err := StructTagValue(b, "")
	a.NoError(err)
	a.Equal("Mayur", v["name"])
	a.Equal(24.0, v["age"])
	a.Equal(int64(123123), v["number"])
}
