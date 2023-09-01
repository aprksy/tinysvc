package main

import (
	"flag"
	"fmt"

	"github.com/gomodule/redigo/redis"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	tile38ip       = "127.0.0.1"
	tile38port     = 9851
	boundaryDir    = "boundaries"
	boundaryBucket = "boundaries"
	resultBucket   = ""
)

var (
	pool *redis.Pool
)

func main() {
	pAddress := flag.String("a", "0.0.0.0", "address to listen to")
	pPort := flag.Int("p", 9301, "port to listen to")

	// setup redis infrastructure
	pool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "127.0.0.1:6379")
		},
	}

	// load boundary files to Redis
	loadBoundaries(boundaryDir, boundaryBucket)

	// Echo instance
	e := echo.New()

	// Middleware
	// e.Use(middleware.Logger())
	// e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// cells
	e.POST("/cells/site", apiGetCellsSite)
	e.POST("/cells/details", apiGetCellsDetails)
	// sites
	e.POST("/sites/cells", apiGetSitesCells)
	e.POST("/sites/details", apiGetSitesDetails)
	e.POST("/sites/intersects", apiGetSitesIntersectBoundary)
	// events
	e.POST("/events/cell-tile", apiGetEventsCellTile)
	e.POST("/events/intersects", apiGetEventsIntersectBoundary)
	// boundaries
	e.PUT("/boundaries", apiAddBoundary)
	e.PUT("/boundaries/", apiAddBoundary)
	e.GET("/boundaries", apiGetBoundary)
	e.POST("/boundaries", apiUpdateBoundary)
	e.POST("/boundaries/", apiUpdateBoundary)
	e.DELETE("/boundaries", apiDeleteBoundary)
	e.DELETE("/boundaries/", apiDeleteBoundary)
	// boundary-facts
	e.POST("/bfact", apiGetBoundaryFacts)
	// simulation
	e.POST("/simulation/ui/rsrp", apiSimulateRsrp)
	e.POST("/simulation/ui/kpi", apiSimulateKpi)
	e.POST("/simulation/automated", apiSimulateRsrp)

	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%d", *pAddress, *pPort)))
}
