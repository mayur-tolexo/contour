package favour

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStructTag(t *testing.T) {
	a := assert.New(t)

	type basicJSON struct {
		Name   string  `json:"name"`
		Age    float64 `json:"age"`
		Number int     `json:"number"`
		Value  int     `sql:"value,omitempty"`
	}

	bJSON := basicJSON{Name: "Mayur", Age: 24, Number: 123123}

	v, err := StructTagValue(bJSON, JSONTag)
	a.NoError(err)
	a.Equal("Mayur", v["name"])
	a.Equal(24.0, v["age"])
	a.Equal(int64(123123), v["number"])
	a.Nil(v["value"])

	type basicSQL struct {
		Name   string  `sql:"name"`
		Age    float64 `sql:"age,omitempty"`
		Number int     `sql:"number"`
		Ignore string  `sql:"-"`
		Value  int     `sql:"value,omitempty"`
	}

	bSQL := basicSQL{Name: "Mayur", Age: 24, Number: 123123, Ignore: "Ignore"}

	v, err = StructTagValue(bSQL, SQLTag)
	a.NoError(err)
	a.Equal("Mayur", v["name"])
	a.Equal(24.0, v["age"])
	a.Equal(int64(123123), v["number"])
	a.Nil(v["ignore"])
	a.Nil(v["value"])
}

func TestEmbededStruct(t *testing.T) {
	a := assert.New(t)

	type Common struct {
		Number int `json:"number"`
		Value  int `sql:"value,omitempty"`
	}

	type basicJSON struct {
		Name string  `json:"name"`
		Age  float64 `json:"age"`
		Common
	}

	bJSON := basicJSON{Name: "Mayur", Age: 24}
	bJSON.Number = 123123

	v, err := StructTagValue(bJSON, JSONTag)
	a.NoError(err)
	a.Equal("Mayur", v["name"])
	a.Equal(24.0, v["age"])
	a.Equal(int64(123123), v["number"])
	a.Nil(v["value"])
}
