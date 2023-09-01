package tile38store

import (
	"fmt"
	"os"

	"github.com/tidwall/gjson"
	"github.com/xjem/t38c"
)

type GeoPointObject map[string]interface{}

// type GeoPointObject struct {
// 	Id     string                 `json:"id"`
// 	Lat    float64                `json:"lat"`
// 	Lng    float64                `json:"lng"`
// 	Fields map[string]interface{} `json:"fields"`
// }

// TILE38 STORAGE

type Tile38Store struct {
	host string
	port int
}

func New(host string, port int) *Tile38Store {
	return &Tile38Store{
		host: host,
		port: port,
	}
}

func (c *Tile38Store) newClient() (client *t38c.Client) {
	client, err := t38c.New(fmt.Sprintf("%s:%d", c.host, c.port))
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return
}

func (c *Tile38Store) releaseClient(client *t38c.Client) {
	client.Close()
}

func (c *Tile38Store) AllPoints(key string) (result []GeoPointObject, err error) {
	worstCaseResult := []GeoPointObject{}
	client := c.newClient()
	defer c.releaseClient(client)

	data, err := client.Execute("SCAN", key)
	if err != nil {
		return worstCaseResult, err
	}

	bytes := []byte(gjson.GetBytes(data, "objects").Raw)
	result = c.getGeoPointObjects(bytes, -1, "distance")
	return
}

func (c *Tile38Store) NearBy(key string, lat, lng, radius float64, limit int) (result []GeoPointObject, err error) {
	worstCaseResult := []GeoPointObject{}
	client := c.newClient()
	defer c.releaseClient(client)

	data, err := client.Execute("NEARBY", key, "DISTANCE", "POINT",
		fmt.Sprintf("%f", lat), fmt.Sprintf("%f", lng), fmt.Sprintf("%f", radius))

	if err != nil {
		return worstCaseResult, err
	}

	bytes := []byte(gjson.GetBytes(data, "objects").Raw)
	result = c.getGeoPointObjects(bytes, -1, "distance")
	return
}

func (c *Tile38Store) PointsIntersect(key, areaType string, data interface{}) (result []GeoPointObject, err error) {
	worstCaseResult := []GeoPointObject{}
	client := c.newClient()
	defer c.releaseClient(client)

	if areaType == "geojson" {
		data, err1 := client.Execute("INTERSECTS", key, "LIMIT", "1000000", "OBJECT", data.(string))
		if err1 != nil {
			return worstCaseResult, err1
		}
		result = c.getGeoPointObjects(data, -1)
	} else if areaType == "circle" {
		lat := gjson.Get(data.(string), "lat").String()
		lng := gjson.Get(data.(string), "lng").String()
		rad := gjson.Get(data.(string), "radius").String()
		data, err1 := client.Execute("INTERSECTS", key, "LIMIT", "1000000", "CIRCLE", lat, lng, rad)
		if err1 != nil {
			return worstCaseResult, err1
		}
		result = c.getGeoPointObjects(data, -1)
	} else {
		err = fmt.Errorf("unknown area type")
		result = worstCaseResult
	}
	return
}

func (c *Tile38Store) getGeoPointObjects(resultSet []byte, limit int, props ...string) (result []GeoPointObject) {
	cols := []string{}
	for _, c := range gjson.GetBytes(resultSet, "fields").Array() {
		cols = append(cols, c.String())
	}

	result = []GeoPointObject{}
	for j, o := range gjson.GetBytes(resultSet, "objects").Array() {
		if limit > -1 && j < limit {
			break
		}
		obj := o.Map()
		id := obj["id"].String()
		// item := GeoPointObject{
		// 	Id:     id,
		// 	Lat:    obj["object"].Map()["coordinates"].Array()[1].Float(),
		// 	Lng:    obj["object"].Map()["coordinates"].Array()[0].Float(),
		// 	Fields: map[string]interface{}{},
		// }

		item := GeoPointObject{
			"id":  id,
			"lat": obj["object"].Map()["coordinates"].Array()[1].Float(),
			"lng": obj["object"].Map()["coordinates"].Array()[0].Float(),
		}

		for i, c := range cols {
			item[c] = obj["fields"].Array()[i].Float()
		}

		for _, p := range props {
			item[p] = obj[p].Float()
		}
		result = append(result, item)
	}
	return
}
