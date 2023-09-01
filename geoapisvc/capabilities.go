package main

import (
	"apisvc/lib/redisstore"
	"apisvc/lib/tile38store"
)

const (
	cell_site              = "%s_R%s_CELL_SITE"
	cell_data              = "%s_R%s_CELL_DATA"
	cell_data_fields       = "tech,band,height,azimuth,mt,et,model,hbw,vendor"
	site_cells             = "%s_R%s_SITE_CELLS"
	site_data              = "%s_R%s_SITE_DATA"
	site_data_fields       = "name,lat,lng,type"
	site_pos               = "%s_R%s_SITE_POS"
	events                 = "%s_R%s_EVENTS"
	event_cell_tile        = "%s_R%s_EVENT_CELL_TILE"
	event_cell_tile_fields = "cell,tile"
)

var (
	rs  = redisstore.New("localhost", 6379)
	t38 = tile38store.New("localhost", 9851)
)
