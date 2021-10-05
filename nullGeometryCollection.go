package postgis

import (
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"strings"

	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/ewkb"
	"github.com/twpayne/go-geom/encoding/wkbcommon"
)

// NullGeometryCollection --
type NullGeometryCollection struct {
	GeometryCollection *geom.GeometryCollection
}

// Scan scans from a []byte.
func (gc *NullGeometryCollection) Scan(src interface{}) error {
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

	gc1, ok := got.(*geom.GeometryCollection)
	if !ok {
		return wkbcommon.ErrUnexpectedType{Got: gc1, Want: gc}
	}

	gc.GeometryCollection = gc1
	return nil
}

// Value returns the WKB encoding of gc.
func (gc NullGeometryCollection) Value() (driver.Value, error) {
	return value(gc.GeometryCollection)
}

func value(g geom.T) (driver.Value, error) {
	// The go-geom types are slightly weird. Geometry objects we call this function with will actually be pointers -
	// some of them can be nil pointers. Because the type here assumes it's not a pointer, the only way to establish
	// whether the object actually exists is by casting it to string and seeing what's inside.
	b, _ := json.Marshal(g)
	sb := &strings.Builder{}
	if string(b) != "null" {
		if err := ewkb.Write(sb, ewkb.NDR, g); err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}

	return hex.EncodeToString([]byte(sb.String())), nil
}
