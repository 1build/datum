package datum

import (
	"database/sql/driver"
	"encoding/hex"

	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/ewkb"
	"github.com/twpayne/go-geom/encoding/wkbcommon"
)

// NullPoint --
type NullPoint struct {
	Point *geom.Point
}

// Scan scans from a []byte.
func (gc *NullPoint) Scan(src interface{}) error {
	if src == nil {
		return nil
	}

	s, ok := src.(string)
	if !ok {
		return ewkb.ErrExpectedByteSlice{Value: src}
	}

	b, err := hex.DecodeString(s)
	if err != nil {
		return ewkb.ErrExpectedByteSlice{Value: src}
	}

	got, err := ewkb.Unmarshal(b)
	if err != nil {
		return err
	}

	gc1, ok := got.(*geom.Point)
	if !ok {
		return wkbcommon.ErrUnexpectedType{Got: gc1, Want: gc}
	}

	gc.Point = gc1
	return nil
}

// Value returns the WKB encoding of gc.
func (gc NullPoint) Value() (driver.Value, error) {
	return value(gc.Point)
}
