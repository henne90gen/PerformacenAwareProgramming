package main

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func withinTolerance(t *testing.T, a, b float64) {
	e := 1e-8
	if a == b {
		return
	}

	d := math.Abs(a - b)

	if b == 0 {
		require.Less(t, d, e)
	}

	require.Less(t, (d / math.Abs(b)), e)
}

func TestHaversine(t *testing.T) {
	result := Haversine(-62.896773782897057004, -40.186985933141897931, -60.232742837180452966, -24.993796147592455270, 6372.8)
	require.Equal(t, 1707.9848326312076, result)

	result = Haversine(21.025119, -53.677203, 62.415201, -72.597892, 6372.8)
	withinTolerance(t, 2849.0521357751214, result)

	result = Haversine(7.046473, -53.998459, -44.023025, -61.530947, 6372.8)
	withinTolerance(t, 3055.523222, result)
}
