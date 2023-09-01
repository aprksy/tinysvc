package main

import (
	"apisvc/lib/tile38store"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/tidwall/gjson"
)

func about(c echo.Context) error {
	c.JSON(http.StatusOK, "")
	return nil
}

func getValues(data []byte, fields ...string) []interface{} {
	worstCaseResult := []interface{}{}
	if len(fields) == 0 {
		return worstCaseResult
	}
	result := worstCaseResult
	for _, f := range fields {
		if f == "data" {
			result = append(result, gjson.GetBytes(data, f).Raw)
		} else {
			result = append(result, gjson.GetBytes(data, f).Value())
		}
	}
	return result
}

func getValuesAsMap(data []byte, fields ...string) map[string]interface{} {
	worstCaseResult := map[string]interface{}{}
	if len(fields) == 0 {
		return worstCaseResult
	}
	result := worstCaseResult
	for _, f := range fields {
		if f == "data" {
			result[f] = gjson.GetBytes(data, f).Raw
		} else {
			result[f] = gjson.GetBytes(data, f).Value()
		}
	}
	return result
}

func apiGetCellsSite(c echo.Context) error {
	// check request body
	data, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, "")
		return nil
	}

	// if not complete, throw error
	values := getValues(data, "date", "region", "cells")
	for _, v := range values {
		if v == nil {
			c.JSON(http.StatusBadRequest, "")
			return nil
		}
	}

	// if date not valid, throw error
	date, err := time.Parse("20060102", values[0].(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, "")
		return nil
	}

	region := values[1].(string)
	cells := strings.Split(values[2].(string), ",")

	// if error in reading redis, throw error
	sites, err := getCellsSite(date, region, cells...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "")
		return nil
	}

	// build response as map
	result := map[string]string{}
	for i := range cells {
		result[cells[i]] = sites[i]
	}

	// return response
	c.JSON(http.StatusOK, result)
	return nil
}

func apiGetCellsDetails(c echo.Context) error {
	// check request body
	data, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, "")
		return nil
	}

	// if not complete, throw error
	values := getValues(data, "date", "region", "cells")
	for _, v := range values {
		if v == nil {
			c.JSON(http.StatusBadRequest, "")
			return nil
		}
	}

	// if date not valid, throw error
	date, err := time.Parse("20060102", values[0].(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, "")
		return nil
	}

	region := values[1].(string)
	cells := strings.Split(values[2].(string), ",")

	// if error in reading redis, throw error
	details, err := getCellsData(date, region, cells...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "")
		return nil
	}

	// build response as map
	result := map[string]map[string]interface{}{}
	for i := range cells {
		result[cells[i]] = details[i]
	}

	// return response
	c.JSON(http.StatusOK, result)
	return nil
}

func apiGetSitesCells(c echo.Context) error {
	// check request body
	data, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, "")
		return nil
	}

	// if not complete, throw error
	values := getValues(data, "date", "region", "sites")
	for _, v := range values {
		if v == nil {
			c.JSON(http.StatusBadRequest, "")
			return nil
		}
	}

	// if date not valid, throw error
	date, err := time.Parse("20060102", values[0].(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, "")
		return nil
	}

	region := values[1].(string)
	sites := strings.Split(values[2].(string), ",")

	// if error in reading redis, throw error
	cells, err := getSiteCells(date, region, sites...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "")
		return nil
	}

	// build response as map
	result := map[string][]string{}
	for i := range sites {
		result[sites[i]] = cells[i]
	}

	// return response
	c.JSON(http.StatusOK, result)
	return nil
}

func apiGetSitesDetails(c echo.Context) error {
	// check request body
	data, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, "")
		return nil
	}

	// if not complete, throw error
	values := getValues(data, "date", "region", "sites")
	for _, v := range values {
		if v == nil {
			c.JSON(http.StatusBadRequest, "")
			return nil
		}
	}

	// if date not valid, throw error
	date, err := time.Parse("20060102", values[0].(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, "")
		return nil
	}

	region := values[1].(string)
	sites := strings.Split(values[2].(string), ",")

	// if error in reading redis, throw error
	details, err := getSitesData(date, region, sites...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "")
		return nil
	}

	// build response as map
	result := map[string]map[string]interface{}{}
	for i := range sites {
		result[sites[i]] = details[i]
	}

	// return response
	c.JSON(http.StatusOK, result)
	return nil
}

func apiGetSitesIntersectBoundary(c echo.Context) error {
	// check request body
	data, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, "")
		return nil
	}

	// if not complete, throw error
	fields := []string{"date", "region", "boundaryId", "areaType", "data"}
	values := getValuesAsMap(data, fields...)
	// values := getValues(data, "date", "region", "boundaryId", "areaType", "data")
	// for _, v := range values {
	// 	if v == nil {
	// 		c.JSON(http.StatusBadRequest, "")
	// 		return nil
	// 	}
	// }

	// if date not valid, throw error
	var date time.Time
	if v := values["date"]; v != nil {
		date, err = time.Parse("20060102", v.(string))
		if err != nil {
			c.JSON(http.StatusBadRequest, "")
			return nil
		}
	}

	region := ""
	if v := values["region"]; v != nil {
		region = v.(string)
	}

	boundaryId := ""
	if v := values["boundaryId"]; v != nil {
		boundaryId = v.(string)
	}

	areaType := ""
	if v := values["areaType"]; v != nil {
		areaType = v.(string)
	}

	boundaryData := ""
	if v := values["data"]; v != nil {
		boundaryData = v.(string)
	}

reevaluate:
	var result []tile38store.GeoPointObject
	if boundaryId != "" {
		_areaType, _data, err := getBoundary(boundaryBucket, boundaryId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "")
			return nil
		}
		geojson, _ := json.Marshal(_data)
		areaType = _areaType
		boundaryData = string(geojson)
		boundaryId = ""
		goto reevaluate
	} else if areaType != "" && boundaryData != "" {
		// TODO: check value validity
		result, err = getSitesIntersectsBoundary(date, region, areaType, boundaryData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "")
			return nil
		}
	} else {
		c.JSON(http.StatusBadRequest, "")
		return nil
	}
	c.JSON(http.StatusOK, result)
	return nil
}

func apiGetEventsCellTile(c echo.Context) error {
	// check request body
	data, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, "")
		return nil
	}

	// if not complete, throw error
	values := getValues(data, "date", "region", "events")
	for _, v := range values {
		if v == nil {
			c.JSON(http.StatusBadRequest, "")
			return nil
		}
	}

	// if date not valid, throw error
	date, err := time.Parse("20060102", values[0].(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, "")
		return nil
	}

	region := values[1].(string)
	events := strings.Split(values[2].(string), ",")

	// if error in reading redis, throw error
	details, err := getEventsCellTile(date, region, events...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "")
		return nil
	}

	// build response as map
	result := map[string]map[string]interface{}{}
	for i := range events {
		result[events[i]] = details[i]
	}

	// return response
	c.JSON(http.StatusOK, result)
	return nil
}

func apiGetEventsIntersectBoundary(c echo.Context) error {
	// check request body
	data, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, "")
		return nil
	}

	// if not complete, throw error
	values := getValues(data, "date", "region", "areaType", "data")
	for _, v := range values {
		if v == nil {
			c.JSON(http.StatusBadRequest, "")
			return nil
		}
	}

	// if date not valid, throw error
	date, err := time.Parse("20060102", values[0].(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, "")
		return nil
	}

	region := values[1].(string)
	areaType := values[2].(string)
	boundaryData := values[3].(string)

	// if error in reading redis, throw error
	result, err := getEventsIntersectsBoundary(date, region, areaType, boundaryData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "")
		return nil
	}

	// return response
	c.JSON(http.StatusOK, result)
	return nil
}

func apiAddBoundary(c echo.Context) error {
	// check request body
	data, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, "")
		return nil
	}

	// if not complete, throw error
	values := getValues(data, "id", "areaType", "data")
	for _, v := range values {
		if v == nil {
			c.JSON(http.StatusBadRequest, "")
			return nil
		}
	}

	id := values[0].(string)
	areaType := values[1].(string)
	boundaryData := values[2].(string)

	// if error in reading redis, throw error
	err = addBoundary(boundaryBucket, boundaryDir, id, areaType, boundaryData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "")
		return nil
	}

	// return response
	c.JSON(http.StatusOK, []byte{})
	return nil
}

func apiUpdateBoundary(c echo.Context) error {
	// check request body
	data, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, "")
		return nil
	}

	// if not complete, throw error
	values := getValues(data, "id", "areaType", "data")
	for _, v := range values {
		if v == nil {
			c.JSON(http.StatusBadRequest, "")
			return nil
		}
	}

	id := values[0].(string)
	areaType := values[1].(string)
	boundaryData := values[2].(string)

	// if error in reading redis, throw error
	err = updateBoundary(boundaryBucket, boundaryDir, id, areaType, boundaryData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "")
		return nil
	}

	// return response
	c.JSON(http.StatusOK, []byte{})
	return nil
}

func apiDeleteBoundary(c echo.Context) error {
	// check request body
	data, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, "")
		return nil
	}

	// if not complete, throw error
	values := getValues(data, "id")
	for _, v := range values {
		if v == nil {
			c.JSON(http.StatusBadRequest, "")
			return nil
		}
	}

	id := values[0].(string)

	// if error in reading redis, throw error
	err = removeBoundary(boundaryBucket, boundaryDir, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "")
		return nil
	}

	// return response
	c.JSON(http.StatusOK, []byte{})
	return nil
}

func apiGetBoundary(c echo.Context) error {
	// check query params
	id := c.Request().URL.Query().Get("id")

	if id == "" {
		result, err := getBoundaries(boundaryBucket)
		if err != nil {
			if err.Error() == "redigo: nil returned" {
				c.JSON(http.StatusNoContent, "")
			} else {
				c.JSON(http.StatusInternalServerError, "")
			}
			return nil
		}
		c.JSON(http.StatusOK, result)
		return nil
	}

	// if error in reading redis, throw error
	areaType, content, err := getBoundary(boundaryBucket, id)
	if err != nil {
		if err.Error() == "redigo: nil returned" {
			c.JSON(http.StatusNoContent, "")
		} else {
			c.JSON(http.StatusInternalServerError, "")
		}
		return nil
	}

	// var data map[string]interface{}
	// err = json.Unmarshal([]byte(content), &data)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	response := map[string]interface{}{
		"id":       id,
		"areaType": areaType,
		"data":     content,
	}

	// return response
	c.JSON(http.StatusOK, response)
	return nil
}

func apiGetBoundaryFacts(c echo.Context) error {
	// check request body
	data, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, "")
		return nil
	}

	// if not complete, throw error
	fields := []string{"date", "region", "boundaryId", "cacheId"}
	values := getValuesAsMap(data, fields...)
	for _, field := range fields {
		if _, ok := values[field]; field != "cacheId" && !ok {
			c.JSON(http.StatusBadRequest, "")
			return nil
		}
	}

	// if date not valid, throw error
	var date time.Time
	if v := values["date"]; v != nil {
		date, err = time.Parse("20060102", v.(string))
		if err != nil {
			c.JSON(http.StatusBadRequest, "")
			return nil
		}
	}

	region := ""
	if v := values["region"]; v != nil {
		region = v.(string)
	}

	boundaryId := ""
	if v := values["boundaryId"]; v != nil {
		boundaryId = v.(string)
	}

	cacheId := ""
	if v := values["cacheId"]; v != nil {
		cacheId = v.(string)
	}

	// if cacheId exist, get data from cache

	var result map[string]interface{}
	resultId := cacheId
	if !date.IsZero() && boundaryId != "" && region != "" {
		resultId = randStr(16)
		result, err = getBoundaryFactsFromId(date, region, boundaryBucket, boundaryId, resultBucket, resultId, true)
	} else if cacheId != "" {
		_temp, err := rs.Value(resultBucket, "RESULT."+cacheId)
		if err != nil {
			fmt.Println(err.Error())
			c.JSON(http.StatusInternalServerError, "")
			return nil
		}
		err = json.Unmarshal([]byte(_temp), &result)
		if err != nil {
			fmt.Println(err.Error())
			c.JSON(http.StatusInternalServerError, "")
			return nil
		}
	} else {
		c.JSON(http.StatusBadRequest, "")
		return nil
	}
	// if error in reading redis, throw error
	if err != nil {
		c.JSON(http.StatusInternalServerError, "")
		return nil
	}

	// return response
	c.JSON(http.StatusOK, result)
	return nil
}

func apiSimulateRsrp(c echo.Context) error {
	// check request body
	data, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, "")
		return nil
	}

	// if not complete, throw error
	fields := []string{"date", "region", "boundaryId", "cacheId", "sites", "tile", "groupNames",
		"groupLimits", "minCatIndex"}
	values := getValuesAsMap(data, fields...)
	for _, field := range fields {
		if _, ok := values[field]; field != "cacheId" && !ok {
			c.JSON(http.StatusBadRequest, "")
			return nil
		}
	}

	// if date not valid, throw error
	var date time.Time
	if v := values["date"]; v != nil {
		date, err = time.Parse("20060102", v.(string))
		if err != nil {
			c.JSON(http.StatusBadRequest, "")
			return nil
		}
	}

	region := ""
	if v := values["region"]; v != nil {
		region = v.(string)
	}

	boundaryId := ""
	if v := values["boundaryId"]; v != nil {
		boundaryId = v.(string)
	}

	cacheId := ""
	if v := values["cacheId"]; v != nil {
		cacheId = v.(string)
	}

	sites := []string{}
	if v := values["sites"]; v != nil {
		for _, s := range v.([]interface{}) {
			sites = append(sites, s.(string))
		}
	}

	tile := "tileCovmo"
	if v := values["tile"]; v != nil {
		tile = v.(string)
	}

	groupNames := []string{"POOR", "FAIR", "GOOD", "EXCELLENT"}
	if v := values["groupNames"]; v != nil {
		groupNames = []string{}
		for _, e := range v.([]interface{}) {
			groupNames = append(groupNames, e.(string))
		}
	}

	groupLimits := []float64{-200, -110, -102, -97, 0}
	if v := values["groupLimits"]; v != nil {
		groupLimits = []float64{}
		for _, e := range v.([]interface{}) {
			groupLimits = append(groupLimits, e.(float64))
		}
	}

	minCatIndex := 1
	if v := values["minCatIndex"]; v != nil {
		minCatIndex = int(v.(float64))
	}

	// if cacheId exist, get data from cache
	var facts map[string]interface{}
	resultId := cacheId
	if !date.IsZero() && boundaryId != "" && region != "" {
		resultId = randStr(16)
		facts, err = getBoundaryFactsFromId(date, region, boundaryBucket, boundaryId, resultBucket, resultId, true)
	} else if cacheId != "" {
		_temp, err := rs.Value(resultBucket, "RESULT."+cacheId)
		if err != nil {
			fmt.Println(err.Error())
			c.JSON(http.StatusInternalServerError, "")
			return nil
		}
		err = json.Unmarshal([]byte(_temp), &facts)
		if err != nil {
			fmt.Println(err.Error())
			c.JSON(http.StatusInternalServerError, "")
			return nil
		}
	} else {
		c.JSON(http.StatusBadRequest, "")
		return nil
	}
	// if error in reading redis, throw error
	if err != nil {
		c.JSON(http.StatusInternalServerError, "")
		return nil
	}

	result := simulateRemoveSite(facts, sites, tile, groupNames, groupLimits, minCatIndex)

	// return response
	c.JSON(http.StatusOK, result)
	return nil
}

func apiSimulateKpi(c echo.Context) error {
	// check request body
	data, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, "")
		return nil
	}

	// if not complete, throw error
	fields := []string{"date", "region", "boundaryId", "cacheId", "sites", "tile", "kpi", "groupNames",
		"groupLimits", "minCatIndex"}
	values := getValuesAsMap(data, fields...)
	for _, field := range fields {
		if _, ok := values[field]; field != "cacheId" && !ok {
			c.JSON(http.StatusBadRequest, "")
			return nil
		}
	}

	// if date not valid, throw error
	var date time.Time
	if v := values["date"]; v != nil {
		date, err = time.Parse("20060102", v.(string))
		if err != nil {
			c.JSON(http.StatusBadRequest, "")
			return nil
		}
	}

	region := ""
	if v := values["region"]; v != nil {
		region = v.(string)
	}

	boundaryId := ""
	if v := values["boundaryId"]; v != nil {
		boundaryId = v.(string)
	}

	cacheId := ""
	if v := values["cacheId"]; v != nil {
		cacheId = v.(string)
	}

	sites := []string{}
	if v := values["sites"]; v != nil {
		for _, s := range v.([]interface{}) {
			sites = append(sites, s.(string))
		}
	}

	tile := "tileCovmo"
	if v := values["tile"]; v != nil {
		tile = v.(string)
	}

	kpi := "rsrp"
	if v := values["kpi"]; v != nil {
		kpi = v.(string)
	}

	groupNames := []string{"POOR", "FAIR", "GOOD", "EXCELLENT"}
	if v := values["groupNames"]; v != nil {
		groupNames = []string{}
		for _, e := range v.([]interface{}) {
			groupNames = append(groupNames, e.(string))
		}
	}

	groupLimits := []float64{-200, -110, -102, -97, 0}
	if v := values["groupLimits"]; v != nil {
		groupLimits = []float64{}
		for _, e := range v.([]interface{}) {
			groupLimits = append(groupLimits, e.(float64))
		}
	}

	minCatIndex := 1
	if v := values["minCatIndex"]; v != nil {
		minCatIndex = int(v.(float64))
	}

	// if cacheId exist, get data from cache
	var facts map[string]interface{}
	resultId := cacheId
	if !date.IsZero() && boundaryId != "" && region != "" {
		resultId = randStr(16)
		facts, err = getBoundaryFactsFromId(date, region, boundaryBucket, boundaryId, resultBucket, resultId, true)
	} else if cacheId != "" {
		_temp, err := rs.Value(resultBucket, "RESULT."+cacheId)
		if err != nil {
			fmt.Println(err.Error())
			c.JSON(http.StatusInternalServerError, "")
			return nil
		}
		err = json.Unmarshal([]byte(_temp), &facts)
		if err != nil {
			fmt.Println(err.Error())
			c.JSON(http.StatusInternalServerError, "")
			return nil
		}
	} else {
		c.JSON(http.StatusBadRequest, "")
		return nil
	}
	// if error in reading redis, throw error
	if err != nil {
		c.JSON(http.StatusInternalServerError, "")
		return nil
	}

	result := simulateRemoveSiteKpi(facts, sites, tile, kpi, groupNames, groupLimits, minCatIndex)

	// return response
	c.JSON(http.StatusOK, result)
	return nil
}
