package geographic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Coordinates_Transform(t *testing.T) {
	assert := assert.New(t)

	lon, lat, err := Transform(2154, 1044593.0, 6298716.0)
	assert.NoError(err)
	assert.InDelta(7.2761920740520, lon, 0.00000001)
	assert.InDelta(43.703478617764, lat, 0.00000001)
}

func Test_Coordinates_Transform_EPSG_27572(t *testing.T) {
	assert := assert.New(t)

	lon, lat, err := Transform(27572, 617766.0, 2609021.0)
	assert.NoError(err)

	assert.InDelta(2.586260171787, lon, 0.00000001)
	assert.InDelta(50.475718687453, lat, 0.00000001)
}

func benchmark_EPSG27572(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Transform(27572, 617766.0, 2609021.0)
	}
}

func benchmark_EPSG2154(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Transform(2154, 1044593.0, 6298716.0)
	}
}

func BenchmarkEPSG27572(b *testing.B) { benchmark_EPSG27572(b) }
func BenchmarkEPSG2154(b *testing.B)  { benchmark_EPSG2154(b) }
