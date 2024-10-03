package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	usr := User{}
	resUser := Create(int64(33), float64(55), float64(77))
	assert.IsType(t, &usr, resUser, "TYPES NOT EQUEL")
}

func TestUpdate(t *testing.T) {
	usr := User{}
	res := usr.Update()
	if !res {
		t.Error("Expected true, but got false")
	}
}
