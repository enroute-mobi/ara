package geographic

import (
	"fmt"
	"math"

	"github.com/wroge/wgs84/v2"
)

func Transform(srsName int, x, y float64) (lon, lat float64, e error) {
	lon, lat, _ = wgs84.Transform(wgs84.EPSG(srsName), wgs84.EPSG(4326))(x, y, 0)

	if math.IsNaN(lon) || math.IsNaN(lat) {
		return 0, 0, fmt.Errorf("unsupported coordinate reference system EPSG:%d", srsName)
	}

	return lon, lat, nil
}
