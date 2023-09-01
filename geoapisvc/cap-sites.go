package main

import (
	"apisvc/lib/tile38store"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func getSiteCells(date time.Time, region string, sites ...string) ([][]string, error) {
	worstCaseResult := [][]string{}
	retval, err := rs.Values(fmt.Sprintf(site_cells, date.Format("20060102"), region), sites)
	if err != nil {
		return worstCaseResult, err
	}
	result := worstCaseResult
	for _, e := range retval {
		cells := strings.Split(string(e), ",")
		result = append(result, cells)
	}
	return result, nil
}

func getSitesData(date time.Time, region string, sites ...string) ([]map[string]interface{}, error) {
	worstCaseResult := []map[string]interface{}{}
	retval, err := rs.Values(fmt.Sprintf(site_data, date.Format("20060102"), region), sites)
	if err != nil {
		return worstCaseResult, err
	}
	result := worstCaseResult
	fields := strings.Split(site_data_fields, ",")
	for j, e := range retval {
		values := strings.Split(string(e), ",")
		data := map[string]interface{}{"id": j + 1, "site": sites[j]}
		for i := range fields {
			if i == 1 || i == 2 {
				data[fields[i]], _ = strconv.ParseFloat(values[i], 64)
			} else {
				data[fields[i]] = values[i]
			}
		}
		result = append(result, data)
	}
	return result, nil
}

func getSitesIntersectsBoundary(date time.Time, region string, areaType, boundary string) ([]tile38store.GeoPointObject, error) {
	worstCaseResult := []tile38store.GeoPointObject{}
	retval, err := t38.PointsIntersect(fmt.Sprintf(site_pos, date.Format("20060102"), region), areaType, boundary)
	if err != nil {
		return worstCaseResult, err
	}
	result := retval
	return result, nil
}
