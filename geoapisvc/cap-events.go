package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/tidwall/gjson"
)

/*
boundary file content shall be as follows:
{
	"id": "boundary_id",
	"areaType": "circle/geojson",
	"data": **data-definition**
}
this also for the key's value in Redis
*/

func loadBoundaries(dir, bucket string) {
	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Printf("unable to load boundary files: %s\n", err.Error())
		return
	}
	count := 0
	for _, f := range files {
		if !f.IsDir() {
			if path.Ext(f.Name()) == ".json" {
				fname := path.Join(dir, f.Name())
				content, err := os.ReadFile(fname)
				if err == nil {
					err1 := loadBoundary(bucket, content)
					if err1 == nil {
						count++
					}
				}
			}
		}
	}
}

func loadBoundary(bucket string, content []byte) error {
	id := gjson.GetBytes(content, "id").String()
	return rs.SetField(bucket, id, content)
}

func addBoundary(bucket, dir, id, areaType string, boundary interface{}) error {
	if !validStrId(id) {
		return fmt.Errorf("invalid id")
	}
	// content := map[string]interface{}{
	// 	"id": id,
	// 	"data": map[string]interface{}{
	// 		"areaType": areaType,
	// 		"data":     boundary,
	// 	},
	// }
	content := map[string]interface{}{
		"id":       id,
		"areaType": areaType,
		"data":     boundary,
	}

	data, err := json.Marshal(content)
	if err != nil {
		return err
	}
	os.WriteFile(path.Join(dir, id+".json"), data, 0644)
	loadBoundary(bucket, data)
	return nil
}

func getBoundaries(bucket string) (result []string, err error) {
	result, err = rs.AllFields(bucket)
	if err != nil {
		return []string{}, err
	}
	return
}

func getBoundary(bucket, id string) (areaType string, data interface{}, err error) {
	result, err := rs.Value(bucket, id)
	if err != nil {
		return "", "", err
	}
	var jsonObj map[string]interface{}
	rawStr, _ := strconv.Unquote(gjson.GetBytes(result, "data").Raw)
	json.Unmarshal([]byte(rawStr), &jsonObj)
	return gjson.GetBytes(result, "areaType").String(), jsonObj, nil
}

func updateBoundary(bucket, dir, id, areaType, boundary string) error {
	return addBoundary(bucket, dir, id, areaType, boundary)
}

func removeBoundary(bucket, dir, id string) error {
	rs.DelField(bucket, id)
	os.Remove(path.Join(dir, id+".json"))
	return nil
}

func getBoundaryFactsFromId(date time.Time, region string, boundaryBucket, boundaryId, resultBucket, resultId string,
	storeToCache bool) (map[string]interface{}, error) {
	result := map[string]interface{}{}
	areaType, data, err := getBoundary(boundaryBucket, boundaryId)
	if err != nil {
		return result, err
	}
	geojson, _ := json.Marshal(data)
	result, err = getBoundaryFactsFromObject(date, region, areaType, string(geojson))
	if err != nil {
		return result, err
	}
	if storeToCache {
		result["resultId"] = resultId
		storeBoundaryFactsToCache(resultBucket, resultId, result)
	}
	return result, nil
}
