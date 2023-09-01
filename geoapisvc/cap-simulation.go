package main

import (
	"apisvc/lib/tile38store"
	"math"
)

func computeRsrpAverageBySite(data map[string]interface{}, site string, tile string) map[string]float64 {
	tileEventSum := map[string]float64{}
	tileEventCount := map[string]int{}
	for siteName, siteStat := range data["service-site"].(map[string]map[string]interface{}) {
		if site != siteName {
			for _, index := range siteStat["indexes"].([]int) {
				events := data["events"].([]tile38store.GeoPointObject)[index]
				if value, ok := tileEventSum[events[tile].(string)]; !ok {
					tileEventSum[events[tile].(string)] = events["rsrp"].(float64)
					tileEventCount[events[tile].(string)] = 1
				} else {
					tileEventSum[events[tile].(string)] = value + events["rsrp"].(float64)
					tileEventCount[events[tile].(string)] = tileEventCount[events[tile].(string)] + 1
				}
			}
		}
	}
	result := map[string]float64{}
	for key := range tileEventCount {
		numerator := tileEventSum[key]
		denominator := 1
		if tileEventCount[key] == 0 {
			numerator = 0
		} else {
			denominator = tileEventCount[key]
		}
		value := numerator / float64(denominator)
		result[key] = value
	}
	return result
}

func computeKpiAverageBySite(data map[string]interface{}, site, tile, kpi string) map[string]float64 {
	tileEventSum := map[string]float64{}
	tileEventCount := map[string]int{}
	for siteName, siteStat := range data["service-site"].(map[string]map[string]interface{}) {
		if site != siteName {
			for _, index := range siteStat["indexes"].([]int) {
				events := data["events"].([]tile38store.GeoPointObject)[index]
				if events[kpi] != nil {
					tileId := events["tiles"].(map[string]interface{})[tile].(map[string]interface{})["id"].(string)
					// if value, ok := tileEventSum[events[tile].(string)]; !ok {
					if value, ok := tileEventSum[tileId]; !ok {
						// tileEventSum[events[tile].(string)] = events[kpi].(float64)
						// tileEventCount[events[tile].(string)] = 1
						tileEventSum[tileId] = events[kpi].(float64)
						tileEventCount[tileId] = 1
					} else {
						// tileEventSum[events[tile].(string)] = value + events[kpi].(float64)
						// tileEventCount[events[tile].(string)] = tileEventCount[events[tile].(string)] + 1
						tileEventSum[tileId] = value + events[kpi].(float64)
						tileEventCount[tileId] = tileEventCount[tileId] + 1
					}
				}
			}
		}
	}
	result := map[string]float64{}
	for key := range tileEventCount {
		numerator := tileEventSum[key]
		denominator := 1
		if tileEventCount[key] == 0 {
			numerator = 0
		} else {
			denominator = tileEventCount[key]
		}
		value := numerator / float64(denominator)
		result[key] = value
	}
	return result
}

func computeMinMax(data map[string]float64) (float64, float64) {
	min := 0.0
	max := -1000.0
	for _, value := range data {
		if value > max {
			max = value
		}
		if value < min {
			min = value
		}
	}
	return min, max
}

func groupByCategory(data map[string]float64, names []string, limits []float64) map[string]int {
	result := map[string]int{}
	intervals := [][2]float64{}

	for i := 0; i < len(limits)-1; i++ {
		intervals = append(intervals, [2]float64{limits[i], limits[i+1]})
	}

	for _, name := range names {
		result[name] = 0
	}

	for _, d := range data {
		for i, interval := range intervals {
			if interval[0] <= d && d < interval[1] {
				result[names[i]] = result[names[i]] + 1
				break
			}
		}
	}
	return result
}

func groupByRoundedValue(data map[string]float64) [][2]float64 {
	result := [][2]float64{}
	temp := map[float64]int{}
	for _, v := range data {
		if value, ok := temp[math.Round(v)]; !ok {
			temp[math.Round(v)] = 1
		} else {
			temp[math.Round(v)] = value + 1
		}
	}
	for k, v := range temp {
		result = append(result, [2]float64{k, float64(v)})
	}
	return result
}

func setTileCategory(data map[string]float64, names []string, limits []float64) map[string]map[string]interface{} {
	result := map[string]map[string]interface{}{}
	intervals := [][2]float64{}

	for i := 0; i < len(limits)-1; i++ {
		intervals = append(intervals, [2]float64{limits[i], limits[i+1]})
	}

	for k, v := range data {
		catIndex := 0
		for i, interval := range intervals {
			if interval[0] <= v && v < interval[1] {
				catIndex = i
				break
			}
		}
		result[k] = map[string]interface{}{
			"value":    v,
			"catIndex": catIndex,
			"catName":  names[catIndex],
		}
	}
	return result
}

func setTileDelta(data0, data1 map[string]map[string]interface{}, minIdx int) map[string]map[string]interface{} {
	result := map[string]map[string]interface{}{}
	for k, d0 := range data0 {
		value := 1.0
		catName := "UNDEFINED"
		catIndex := -1
		delta := "UNDEFINED"
		status := "UNDEFINED"
		if d1, ok := data1[k]; ok {
			value = d1["value"].(float64)
			catName = d1["catName"].(string)
			catIndex = d1["catIndex"].(int)

			if d1["catIndex"].(int) < d0["catIndex"].(int) {
				delta = "DEGRADED"
			} else if d1["catIndex"].(int) > d0["catIndex"].(int) {
				delta = "UPGRADED"
			} else {
				delta = "UNCHANGE"
			}
		}
		if catIndex < 0 {
			status = "FATAL"
		} else if catIndex < minIdx {
			status = "UNSAFE"
		} else {
			status = "SAFE"
		}
		result[k] = map[string]interface{}{
			"value0":    d0["value"].(float64),
			"category0": d0["catName"].(string),
			"value1":    value,
			"category1": catName,
			"delta":     delta,
			"status":    status,
		}
	}
	return result
}

func setTilesStatusSummary(data map[string]map[string]interface{}) map[string]int {
	result := map[string]int{"SAFE": 0, "UNSAFE": 0, "FATAL": 0}
	for _, d := range data {
		switch d["status"].(string) {
		case "SAFE":
			result["SAFE"] = result["SAFE"] + 1
		case "UNSAFE":
			result["UNSAFE"] = result["UNSAFE"] + 1
		case "FATAL":
			result["FATAL"] = result["FATAL"] + 1
		}
	}
	return result
}

func setTilesDeltaSummary(data map[string]map[string]interface{}) map[string]int {
	result := map[string]int{"UPGRADED": 0, "DEGRADED": 0, "UNCHANGE": 0}
	for _, d := range data {
		switch d["delta"].(string) {
		case "UPGRADED":
			result["UPGRADED"] = result["UPGRADED"] + 1
		case "DEGRADED":
			result["DEGRADED"] = result["DEGRADED"] + 1
		case "UNCHANGE":
			result["UNCHANGE"] = result["UNCHANGE"] + 1
		}
	}
	return result
}

func simulateRemoveSite(data map[string]interface{}, sites []string, tile string,
	groupNames []string, groupLimits []float64, minCatIndex int) map[string]interface{} {
	result := map[string]interface{}{}
	tiles := map[string]map[string]float64{}
	for _, e := range data["events"].([]tile38store.GeoPointObject) {
		tiles[e[tile].(string)] = map[string]float64{}
		tiles[e[tile].(string)]["lat"] = e["clat"].(float64)
		tiles[e[tile].(string)]["lng"] = e["clng"].(float64)
	}
	inboundarySites := map[string]map[string]interface{}{}
	simulateAllSites := len(sites) == 0
	for _, e := range data["inboundary-sites"].([]tile38store.GeoPointObject) {
		inboundarySites[e["id"].(string)] = map[string]interface{}{}
		inboundarySites[e["id"].(string)]["lat"] = e["lat"].(float64)
		inboundarySites[e["id"].(string)]["lng"] = e["lng"].(float64)
		inboundarySites[e["id"].(string)]["id"] = e["id"].(string)
		inboundarySites[e["id"].(string)]["name"] = e["name"].(string)
		inboundarySites[e["id"].(string)]["type"] = e["type"].(string)
		if simulateAllSites {
			sites = append(sites, e["id"].(string))
		}
	}
	origTiles := computeRsrpAverageBySite(data, "original", tile)
	origMin, origMax := computeMinMax(origTiles)
	origCatGroups := groupByCategory(origTiles, groupNames, groupLimits)
	origIntGroups := groupByRoundedValue(origTiles)
	origTilesCat := setTileCategory(origTiles, groupNames, groupLimits)
	origTilesDelta := setTileDelta(origTilesCat, origTilesCat, minCatIndex)
	original := map[string]interface{}{
		"tiles":     origTilesDelta,
		"count":     len(origTiles),
		"min":       origMin,
		"max":       origMax,
		"category":  origCatGroups,
		"frequency": origIntGroups,
	}

	result["tiles"] = tiles
	result["sites"] = inboundarySites
	result["original"] = original

	simulations := map[string]map[string]interface{}{}
	for _, site := range sites {
		simTiles := computeRsrpAverageBySite(data, site, tile)
		simMin, simMax := computeMinMax(simTiles)
		simCatGroups := groupByCategory(simTiles, groupNames, groupLimits)
		simIntGroups := groupByRoundedValue(simTiles)
		simTilesCat := setTileCategory(simTiles, groupNames, groupLimits)
		simTilesDelta := setTileDelta(origTilesCat, simTilesCat, minCatIndex)
		simTilesDeltaSummary := setTilesDeltaSummary(simTilesDelta)
		simTilesStatusSummary := setTilesStatusSummary(simTilesDelta)
		simSite := map[string]interface{}{
			"tiles":         simTilesDelta,
			"count":         len(simTiles),
			"min":           simMin,
			"max":           simMax,
			"category":      simCatGroups,
			"frequency":     simIntGroups,
			"deltaSummary":  simTilesDeltaSummary,
			"statusSummary": simTilesStatusSummary,
		}
		simulations[site] = simSite
	}
	result["simulation"] = simulations

	return result
}

func simulateRemoveSiteKpi(data map[string]interface{}, sites []string, tile, kpi string,
	groupNames []string, groupLimits []float64, minCatIndex int) map[string]interface{} {
	result := map[string]interface{}{}
	tiles := map[string]map[string]float64{}
	for _, e := range data["events"].([]tile38store.GeoPointObject) {
		etile := e["tiles"].(map[string]interface{})[tile].(map[string]interface{})
		tiles[etile["id"].(string)] = map[string]float64{}
		tiles[etile["id"].(string)]["lat"] = etile["lat"].(float64)
		tiles[etile["id"].(string)]["lng"] = etile["lng"].(float64)
	}
	inboundarySites := map[string]map[string]interface{}{}
	simulateAllSites := len(sites) == 0
	for _, e := range data["inboundary-sites"].([]tile38store.GeoPointObject) {
		inboundarySites[e["id"].(string)] = map[string]interface{}{}
		inboundarySites[e["id"].(string)]["lat"] = e["lat"].(float64)
		inboundarySites[e["id"].(string)]["lng"] = e["lng"].(float64)
		inboundarySites[e["id"].(string)]["id"] = e["id"].(string)
		inboundarySites[e["id"].(string)]["name"] = e["name"].(string)
		inboundarySites[e["id"].(string)]["type"] = e["type"].(string)
		if simulateAllSites {
			sites = append(sites, e["id"].(string))
		}
	}
	origTiles := computeKpiAverageBySite(data, "original", tile, kpi)
	origMin, origMax := computeMinMax(origTiles)
	origCatGroups := groupByCategory(origTiles, groupNames, groupLimits)
	origIntGroups := groupByRoundedValue(origTiles)
	origTilesCat := setTileCategory(origTiles, groupNames, groupLimits)
	origTilesDelta := setTileDelta(origTilesCat, origTilesCat, minCatIndex)
	original := map[string]interface{}{
		"tiles":     origTilesDelta,
		"count":     len(origTiles),
		"min":       origMin,
		"max":       origMax,
		"category":  origCatGroups,
		"frequency": origIntGroups,
	}

	result["tiles"] = tiles
	result["sites"] = inboundarySites
	if len(data["events"].([]tile38store.GeoPointObject)) > 0 {
		event0 := data["events"].([]tile38store.GeoPointObject)[0]
		result["tileRadius"] = event0["tiles"].(map[string]interface{})[tile].(map[string]interface{})["radius"]
	}
	result["original"] = original

	simulations := map[string]map[string]interface{}{}
	for _, site := range sites {
		simTiles := computeKpiAverageBySite(data, site, tile, kpi)
		simMin, simMax := computeMinMax(simTiles)
		simCatGroups := groupByCategory(simTiles, groupNames, groupLimits)
		simIntGroups := groupByRoundedValue(simTiles)
		simTilesCat := setTileCategory(simTiles, groupNames, groupLimits)
		simTilesDelta := setTileDelta(origTilesCat, simTilesCat, minCatIndex)
		simTilesDeltaSummary := setTilesDeltaSummary(simTilesDelta)
		simTilesStatusSummary := setTilesStatusSummary(simTilesDelta)
		simSite := map[string]interface{}{
			"tiles":         simTilesDelta,
			"count":         len(simTiles),
			"min":           simMin,
			"max":           simMax,
			"category":      simCatGroups,
			"frequency":     simIntGroups,
			"deltaSummary":  simTilesDeltaSummary,
			"statusSummary": simTilesStatusSummary,
		}
		simulations[site] = simSite
	}
	result["simulation"] = simulations

	return result
}
