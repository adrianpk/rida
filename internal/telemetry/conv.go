package telemetry

import (
	"net/http"
	"strconv"
)

func NewQuery(r *http.Request) (Query, error) {
	q := r.URL.Query()
	var qry Query
	var err error

	qry.Area.MinLat, err = parseFloat(q.Get("minLat"))
	if err != nil {
		return qry, err
	}

	qry.Area.MinLng, err = parseFloat(q.Get("minLng"))
	if err != nil {
		return qry, err
	}

	qry.Area.MaxLat, err = parseFloat(q.Get("maxLat"))
	if err != nil {
		return qry, err
	}

	qry.Area.MaxLng, err = parseFloat(q.Get("maxLng"))
	if err != nil {
		return qry, err
	}

	qry.Status = Status(q.Get("status"))

	return qry, nil
}

func parseFloat(val string) (float64, error) {
	return strconv.ParseFloat(val, 64)
}
