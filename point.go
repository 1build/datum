package datum

import (
	"database/sql/driver"
	"errors"
	"fmt"

	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/ewkb"
	"github.com/volatiletech/null/convert"
)

// Point is a PostGIS point.
type Point struct {
	SRID int
	Lng  float64
	Lat  float64
}

// NewPoint creates a new point.
func NewPoint(srid int, lng float64, lat float64) Point {
	return Point{
		SRID: srid,
		Lng:  lng,
		Lat:  lat,
	}
}

// PointFromVal creates a new point from any value. If the value is nil, the point (0 0) is returned.
func PointFromVal(value interface{}, srid int) (*Point, error) {
	point := NewPoint(0, 0)
	var geometryBuffer []byte
	convert.ConvertAssign(&geometryBuffer, value)
	geometry, err := ewkb.Unmarshal(geometryBuffer)
	if err != nil {
		// no library seems to understand what PostGIS is giving us so I just return 0, 0
		return &point, nil
	}

	geomPoint, ok := geometry.(*geom.Point)
	if !ok {
		return nil, errors.New("geometry is not a point")
	}

	point = NewPoint(srid, geomPoint.Y(), geomPoint.X())

	return &point, nil
}

// Scan implements the Scanner interface.
func (p *Point) Scan(value interface{}, srid int) error {
	point, err := PointFromVal(value, srid)
	if err != nil {
		return err
	}
	*p = *point
	return nil
}

// Value implements the driver Value interface.
func (p Point) Value() (driver.Value, error) {
	return fmt.Sprintf("SRID=%d;POINT(%f %f)", p.SRID, p.Lng, p.Lat), nil
}
