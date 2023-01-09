package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContains(t *testing.T) {
	a := []string{"1", "2", "3"}
	b := "1"
	flag := Contains(a, b)
	assert.Equal(t, flag, true)
}

func TestIntersection(t *testing.T) {
	a := []string{"1", "2", "3"}
	b := []string{"2", "3", "4"}
	ret := Intersection(a, b)
	assert.Equal(t, ret, []string{"2", "3"})
}

func TestUnion(t *testing.T) {
	a := []string{"1", "2", "3"}
	b := []string{"2", "3", "4"}
	ret := Union(a, b)
	assert.Equal(t, ret, []string{"1", "2", "3", "4"})
}

func TestDifference(t *testing.T) {
	a := []string{"1", "2", "3"}
	b := []string{"2", "3", "4"}
	ret1 := Difference(a, b)
	assert.Equal(t, ret1, []string{"1"})
	ret2 := Difference(b, a)
	assert.Equal(t, ret2, []string{"4"})
}
