package geographic

import (
	"github.com/everystreet/go-proj/v6/proj"
	"github.com/wroge/wgs84"
)

func Transform(srsName int, x, y float64) (lat, lon float64, e error) {
	if srsName == 27572 {
		var xy proj.XY
		xy.X = x
		xy.Y = y
		e = proj.CRSToCRS(
			"EPSG:27572",
			"+proj=latlong",
			func(pj proj.Projection) {
				proj.TransformForward(pj, &xy)
			})

		return xy.X, xy.Y, e
	}

	epsg := wgs84.EPSG()
	lon, lat, _, e = epsg.SafeTransform(srsName, 4326)(x, y, 0)

	return lon, lat, e
}
