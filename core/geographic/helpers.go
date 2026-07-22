package geographic

import (
	"fmt"
	"math"

	"github.com/wroge/wgs84/v2"
)

func Transform(srsName int, x, y float64) (lon, lat float64, e error) {
	transform, err := wgs84.Transform(srsName, 4326)
	if err != nil {
		return 0, 0, fmt.Errorf("unsupported coordinate reference system EPSG:%d", srsName)
	}

	lon, lat, _, _ = transform(x, y, 0)

	if math.IsNaN(lon) || math.IsNaN(lat) {
		return 0, 0, fmt.Errorf("unsupported coordinate reference system EPSG:%d", srsName)
	}

	return lon, lat, nil
}
