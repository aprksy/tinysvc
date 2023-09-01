package main

import (
	"apisvc/lib/tile38store"
	"encoding/json"
	"sync"
	"time"
)

func getCellStat(events []tile38store.GeoPointObject, cellTile []map[string]interface{}, tileKind string) map[string]map[string]interface{} {
	tiles := map[string]map[string]bool{}
	cellIndex := map[string]map[string]interface{}{}
	var cellIdx []int
	for i, e := range events {
		cellName := cellTile[i]["cell"].(string)

		var indexes []int
		if cellIndex[cellName] != nil {
			cellIdx = cellIndex[cellName]["indexes"].([]int)
			cellIdx = append(cellIdx, i)
			indexes = append(indexes, cellIdx...)
			cellIndex[cellName]["eventCount"] = cellIndex[cellName]["eventCount"].(int) + 1
			cellIndex[cellName]["indexes"] = indexes
		} else {
			tiles[cellName] = map[string]bool{}
			cellIndex[cellName] = map[string]interface{}{}
			cellIndex[cellName]["indexes"] = []int{i}
			cellIndex[cellName]["eventCount"] = 1
		}
		tiles[cellName][e[tileKind].(string)] = true
		cellIndex[cellName]["tileCount"] = len(tiles[cellName])
		cellIndex[cellName]["site"] = e["site"].(string)
	}

	return cellIndex
}

func getSiteStat(cellStats map[string]map[string]interface{}) map[string]map[string]interface{} {
	siteIndex := map[string]map[string]interface{}{}
	var siteIdx []int
	for _, v := range cellStats {
		site := v["site"].(string)
		var indexes []int
		if siteIndex[site] != nil {
			siteIdx = siteIndex[site]["indexes"].([]int)
			siteIdx = append(siteIdx, v["indexes"].([]int)...)
			indexes = append(indexes, siteIdx...)
			siteIndex[site]["eventCount"] = siteIndex[site]["eventCount"].(int) + v["eventCount"].(int)
			siteIndex[site]["tileCount"] = siteIndex[site]["tileCount"].(int) + v["tileCount"].(int)
			siteIndex[site]["indexes"] = indexes
		} else {
			siteIndex[site] = map[string]interface{}{}
			siteIndex[site]["indexes"] = v["indexes"].([]int)
			siteIndex[site]["eventCount"] = v["eventCount"].(int)
			siteIndex[site]["tileCount"] = v["tileCount"].(int)
		}
	}
	return siteIndex
}

func getBoundaryFactsFromObject(date time.Time, region string, areaType, boundary string) (map[string]interface{}, error) {
	result := map[string]interface{}{}

	var (
		wg            sync.WaitGroup
		sites, events []tile38store.GeoPointObject
		err1, err2    error
	)
	wg.Add(2)
	go func() {
		sites, err1 = getSitesIntersectsBoundary(date, region, areaType, boundary)
		wg.Done()
	}()
	go func() {
		events, err2 = getEventsIntersectsBoundary(date, region, areaType, boundary)
		wg.Done()
	}()
	wg.Wait()

	// prepare siteIds for next queries
	siteIds := []string{}
	for _, e := range sites {
		siteIds = append(siteIds, e["id"].(string))
	}

	var (
		siteCells [][]string
		siteData  []map[string]interface{}
	)

	wg.Add(2)
	go func() {
		siteCells, _ = getSiteCells(date, region, siteIds...)
		wg.Done()
	}()
	go func() {
		siteData, _ = getSitesData(date, region, siteIds...)
		wg.Done()
	}()
	wg.Wait()

	for i, s := range sites {
		s["cells"] = siteCells[i]
		s["name"] = siteData[i]["name"].(string)
		s["type"] = siteData[i]["type"].(string)
		sites[i] = s
	}

	eventIds := []string{}
	for _, e := range events {
		eventIds = append(eventIds, e["id"].(string))
	}

	// FIXME: fix error handling
	cellTile, _ := getEventsCellTile(date, region, eventIds...)
	cellIds := []string{}
	for _, ct := range cellTile {
		cellIds = append(cellIds, ct["cell"].(string))
	}

	siteFromCell, _ := getCellsSite(date, region, cellIds...)
	for i, e := range events {
		cellName := cellTile[i]["cell"].(string)
		siteId := siteFromCell[i]
		e["cell"] = cellName
		e["tileCovmo"] = cellTile[i]["tile"].(string)
		e["site"] = siteId

		// add normal geohash9
		gh7, _ := GeohashEncode(e["lat"].(float64), e["lng"].(float64), 7)
		gh7center, _ := GeohashDecode(gh7)
		gh9, _ := GeohashEncode(e["lat"].(float64), e["lng"].(float64), 9)
		gh9center, _ := GeohashDecode(gh9)
		gh9_2x2, gh9_4x4, gh9_2x2c, gh9_4x4c, _ := Geohash9_ext(gh9)
		e["tiles"] = map[string]interface{}{
			"tileCovmo": map[string]interface{}{
				"id":     e["tileCovmo"].(string),
				"lat":    e["clat"].(float64),
				"lng":    e["clng"].(float64),
				"radius": 37.0 / 2,
			},
			"tileGeohash7": map[string]interface{}{
				"id":     gh7,
				"lat":    gh7center.Lat,
				"lng":    gh7center.Lng,
				"radius": 152.4 / 2,
			},
			"tileGeohash9": map[string]interface{}{
				"id":     gh9,
				"lat":    gh9center.Lat,
				"lng":    gh9center.Lng,
				"radius": 4.77 / 2,
			},
			"tileGeohash9_2x2": map[string]interface{}{
				"id":     gh9_2x2,
				"lat":    gh9_2x2c.Lat,
				"lng":    gh9_2x2c.Lng,
				"radius": 4.77,
			},
			"tileGeohash9_4x4": map[string]interface{}{
				"id":     gh9_4x4,
				"lat":    gh9_4x4c.Lat,
				"lng":    gh9_4x4c.Lng,
				"radius": 4.77 * 2,
			},
		}
		events[i] = e
	}

	cellIndex := getCellStat(events, cellTile, "tileCovmo")

	siteIndex := getSiteStat(cellIndex)

	result["inboundary-sites"] = sites
	result["events"] = events
	result["service-site"] = siteIndex
	result["service-cell"] = cellIndex
	// result["cellTile"] = cellTile

	var err error
	if err != nil {
		err = err1
	} else if err2 != nil {
		err = err2
	}
	return result, err
}

func storeBoundaryFactsToCache(bucket, resultId string, data map[string]interface{}) (string, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	err = rs.SetField(bucket, "RESULT."+resultId, bytes)
	// TODO: set expiration
	if err != nil {
		return "", err
	}
	return resultId, nil
}
