# Datum [![GoDoc](https://godoc.org/github.com/1build/datum?status.svg)](https://godoc.org/github.com/1build/datum)

Datum is a set of tools for serializing geospatial primitives with database/sql & sqlboiler. It's currently limited to `point`, `nullPoint`, and `nullGeometryCollection`.


## Example

Below is a minimal example of using datum to serialize a geospatial point to a record with sqlboiler generated structs.


```go
    import (
        "gitub.com/myorg/path/to/sqlboiler/models"
        "github.com/1build/datum"
    )

    const srid := 4326

    record := &models.Supplier{
        Name: input.Name,
        Location: datum.Point{
            SRID: srid,
            Lat:  input.Lat,
            Lng:  input.Lng,
        },
    }

    if err := record.Insert(ctx, tx, boil.Infer()); err != nil {
        warning := fmt.Sprintf("[Supplier.Repository.CreateSupplier]: Couldn't create new supplier: %s", input.Name)
        logger.Errorf(ctx, warning)

        return nil, errors.New(warning)
    }
```

