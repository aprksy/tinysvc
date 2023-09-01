package main

import (
	"apisvc/lib/tile38store"
	"fmt"
	"strings"
	"time"
)

func getEventsCellTile(date time.Time, region string, events ...string) ([]map[string]interface{}, error) {
	worstCaseResult := []map[string]interface{}{}
	retval, err := rs.Values(fmt.Sprintf(event_cell_tile, date.Format("20060102"), region), events)
	if err != nil {
		return worstCaseResult, err
	}
	result := worstCaseResult
	fields := strings.Split(event_cell_tile_fields, ",")
	for j, e := range retval {
		values := strings.Split(string(e), ",")
		data := map[string]interface{}{"id": j + 1, "cell": events[j]}
		for i := range fields {
			data[fields[i]] = values[i]
		}
		result = append(result, data)
	}
	return result, nil
}

func getEventsIntersectsBoundary(date time.Time, region string, areaType, boundary string) ([]tile38store.GeoPointObject, error) {
	worstCaseResult := []tile38store.GeoPointObject{}
	retval, err := t38.PointsIntersect(fmt.Sprintf(events, date.Format("20060102"), region), areaType, boundary)
	if err != nil {
		return worstCaseResult, err
	}
	result := retval
	return result, nil
}
