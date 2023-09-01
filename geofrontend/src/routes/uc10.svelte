<svelte:head>
    <link rel="stylesheet" href="https://unpkg.com/leaflet@1.7.1/dist/leaflet.css"
    integrity="sha512-xodZBNTC5n17Xt2atTPuE1HxjVMSvLVW9ocqUKLsCC5CXdbqCmblAshOMAS6/keqq/sMZMZ19scR4PsZChSR7A=="
    crossorigin=""/>
    <script src="https://unpkg.com/leaflet@1.7.1/dist/leaflet.js"
    integrity="sha512-XQoYMqMTK8LvdxXYG3nZ448hOEQiglfqkJs1NOQV44cWnUrBc8PkAOcXy20w0vlaXaVUearIOBhiXZ5V3ynxwA=="
    crossorigin=""></script>
    <link rel="stylesheet" href="https://unpkg.com/leaflet-draw@0.4.1/dist/leaflet.draw.css" />
    <script src="https://unpkg.com/leaflet-draw@0.4.1/dist/leaflet.draw.js"></script>
</svelte:head>
<script>
// @ts-nocheck

    import {
        Button,
        Header,
        HeaderUtilities,
        HeaderAction,
        HeaderPanelLinks,
        HeaderPanelDivider,
        HeaderPanelLink,
        SkipToContent,
        ComboBox,
        Modal,
        OverflowMenu, 
        OverflowMenuItem,
        Popover,
        DatePicker,
        DatePickerInput,
        TileGroup,
        RadioTile,
        TextInput,
        Accordion,
        AccordionItem,
        FileUploader,
        FileUploaderDropContainer,
        ButtonSet,
        FileUploaderItem,
        DataTable,
        Toolbar,
        ToolbarContent,
        ToolbarSearch,
        ToolbarMenu,
        ToolbarMenuItem,
        SelectItem,
        FluidForm,
        Form,
        Tab,
        Tabs,
        TabContent,
        Search,
        Slider,
        MultiSelect,
        Link,
        Pagination,
    } from "carbon-components-svelte";
    import { Upload, View, TreeView, Area, AreaCustom, WatsonHealthCircleMeasurement, ChartNetwork, Reset, CheckboxChecked, OrderDetails, Launch, Settings, Number_0 } from "carbon-icons-svelte";
    import Drawer from 'svelte-drawer-component';
    import MapBoundary from "carbon-icons-svelte/lib/MapBoundary.svelte";
    import SettingsAdjust from "carbon-icons-svelte/lib/SettingsAdjust.svelte";
    import UserAvatarFilledAlt from "carbon-icons-svelte/lib/UserAvatarFilledAlt.svelte";
    import { onMount } from 'svelte';
    import { browser } from '$app/env';
    import { 
        storeRegionals, 
        storeBoundaries,
        storeSites,
        storeRegionalsSelected,
        storeNearbyCenter,
        storeNearbySitesLimit,
        storeNearbySitesRadius,
        storeNearbySites,
        storeServingCells,
        storeServingCells_selected,
    } from "../lib/controller/store.js";

    let table1PageSize = 10;
    let table1CurrentPage = 1;

    let table2PageSize = 15;
    let table2CurrentPage = 1;
  
    let drawerOpen = false;
    let isSideNavOpen = false;
    let isOpen1 = false;
    let legendOpen = false;
    let MapBoundaryAddOpen = false;
    let MapBoundaryEditOpen = false;
    let MapBoundaryDeleteOpen = false;
    let AutomationOpen = false;

    let uploadStatus = {
        status: "",
        fileName: "",
        fileDate: 0,
        fileSize: 0,
    }

    let resultSet;
    let tilesActual = {};
    let tilesSimulation = {};
    let mainMapTilesActual;
    let mainMapTilesSimulation;

    function shouldFilterItem(item, value) {
        if (!value) return true;
        return item.text.toLowerCase().includes(value.toLowerCase());
    }

    async function getSites() {
        const res = await fetch('http://localhost:9301/sites?region=' + $storeRegionalsSelected.id, {
			method: 'GET',
		})
		let result = await res.json()
        storeSites.set(result)
    }

    async function getNearbySites(region, lat, lng, limit, radius) {
        const res = await fetch('http://localhost:9301/sites/nearby?' + 
            'region=' + region + '&' +
            'lat=' + lng + '&' +
            'lng=' + lat + '&' +
            'limit=' + limit + '&' +
            'radius=' + radius, {
			method: 'GET',
		})
		let result = await res.json()
        result.forEach((e) => {
            let parts = e.text.split('_')
            e['shortText'] = parts[0].length > 1 ? parts[0] : parts[0] + "_" + parts[1];
            e.distance = e.distance.toFixed(1)
        })
        storeNearbySites.set(result)
        addNearbyMarkers(markerGroup, $storeNearbySites, $storeNearbySitesLimit)
    }

    async function getBoundaries() {
		const res = await fetch('http://localhost:9301/boundaries', {
			method: 'GET',
		})
		let resp = await res.json()
        let result = []
        resp.forEach((e, i) => {
            let src = JSON.parse(atob(e))
            let bound = {
                id: i,
                text: src["name"],
                type: src["type"],
                tag: src["tags"],
                boundary: JSON.parse(src["data"]),
            }
            result.push(bound)
        })
        storeBoundaries.set(result)
	}

    async function addBoundary(data) {
		const res = await fetch('http://localhost:9301/boundaries', {
			method: 'PUT',
            headers: {
               "Content-Type": "application/json",
            },
		    body: JSON.stringify(data)
		})
        if (res.ok) {
            getBoundaries()
        } else {

        }
	}

    let moduleName = "Usecase 10 - Blacksite"

    let layers = [
        // {"type":"FeatureCollection","features":[{"type":"Feature","properties":{},"geometry":{"type":"Polygon","coordinates":[[[106.87491416931152,-6.259332270466314],[106.88023567199707,-6.26411011355174],[106.87946319580078,-6.269143865206112],[106.88401222229004,-6.270594259210957],[106.88796043395996,-6.264536704558167],[106.88718795776366,-6.2587350370072965],[106.88246726989746,-6.25566353986335],[106.87808990478516,-6.256175457307679],[106.87491416931152,-6.259332270466314]]]}}]},
    ];

    let mainMapLayers;
    let editMapLayers;

    const mainMap = {
        id: 'main-map',
        width: 664,
        height: 786,
        // center: [-6.175392, 106.827153], // Jakarta
        center: [3.597031, 98.678513], // Medan
        zoom: 13,
        drawControls: true,
        drawnItems: null,
        map: null,
    }

    const addMap = {
        id: 'add-map',
        width: 670,
        height: 600,
        center: [-6.175392, 106.827153],
        zoom: 13,
        drawControls: false,
        drawnItems: null,
        map: null,
    }

    const editMap = {
        id: 'edit-map',
        width: 670,
        height: 550,
        center: [-6.175392, 106.827153],
        zoom: 13,
        drawControls: true,
        drawnItems: null,
        map: null,
    }

    const deleteMap = {
        id: 'delete-map',
        width: 670,
        height: 600,
        center: [-6.175392, 106.827153],
        zoom: 13,
        drawControls: false,
        drawnItems: null,
        map: null,
    }

    function clearMap(obj) {
        obj.drawnItems.clearLayers();
    }

    function redrawMap(obj, data) {
        if (data) {
            clearMap(obj);
            data.forEach(el => {
                L.geoJson(el, {
                    onEachFeature: function (feature, layer) {
                        obj.drawnItems.addLayer(layer);
                        var bounds = layer.getBounds();
                        obj.map.fitBounds(bounds);
                    }
                });
            });
        }
    }

    let myIcon;

    let markerGroup;
    function clearMarkers(obj) {
        obj.drawnItems.removeLayer(markerGroup);
    }

    function addMainMarker(obj, lat, lng, radius, otherMarkers) {
        clearMarkers(obj);
        let marker = L.marker([lat, lng])
        let circle = L.circle([lat, lng], {radius: radius}).setStyle({
            weight: 1, fillOpacity: 0.1
        })
        let lg = L.layerGroup([circle, marker])
        addNearbyMarkers(lg, otherMarkers, 10);
        obj.drawnItems.addLayer(lg)

        let bound = circle.getBounds();
        obj.map.fitBounds(bound);

        markerGroup = lg
    }

    function addNearbyMarkers(lgroup, data, max) {
        myIcon = L.icon({
            iconUrl: 'map-marker-icon-gray.png',
            iconSize: [18, 30], // size of the icon
        });
        data.forEach((e, i) => {
            if (i <= max) {
                let marker = L.marker([e.lng, e.lat], {icon: myIcon});
                marker.bindPopup(e.text);
                lgroup.addLayer(marker);
            }
        })
    }

    function getCategoryByValue(value) {
        let l1 = -97
        let l2 = -102
        let l3 = -110

        let c1 = 0
        let c2 = 1
        let c3 = 2
        let c4 = 3

        if (value >= l1) {
            return c1
        } else if (value >= l2) {
            return c2
        } else if (value >= l3) {
            return c3
        } else {
            return c4
        }
    }

    function getColorByCategory(value) {
        let c1 = "#0066ff"
        let c2 = "#44aa00"
        let c3 = "#ff9900"
        let c4 = "#ff0000"
        let colors = [c1, c2, c3, c4];

        return colors[value]
    }

    function mapLength(data) {
        let count = 0
        for (const [key, value] of Object.entries(data)) {
            count++;
        }
        return count;
    }

    function sumValue(arrayOfObject, field) {
        let result = 0
        for (const [key, value] of Object.entries(arrayOfObject)) {
            result = result + value[field]
        }
        return result;
    }

    function proportionValue(arrayOfObject, field1, field2) {
        let result1 = 0
        for (const [key, value] of Object.entries(arrayOfObject)) {
            result1 = result1 + value[field1]
        }
        let result2 = 0
        for (const [key, value] of Object.entries(arrayOfObject)) {
            result2 = result2 + value[field2]
        }
        return result1/result2;
    }

    function getServingCells(data) {
        let cells = [];
        let id=0;
        for (const [cn, cd] of Object.entries(data)) {
            cells.push({
                id: cn,
                name: cn,
                tileCount: mapLength(cd),
                eventCount: sumValue(cd, "count"),
                avgRsrp: proportionValue(cd, "sum", "count").toFixed(4),
            })
            id++;
        }
        storeServingCells.set(cells)
    }

    function drawTiles(mapObj, data) {
        let lg;
        let tileData;
        let stroke = false;
        let color = "#eee";
        let radius = 18.5;
        mainMapTilesSimulation.clearLayers();
        tilesSimulation = {};
        if ($storeServingCells_selected.length == 0) {
            mainMapTilesActual.clearLayers();
            lg = mainMapTilesActual;
            tileData = data.tiles
            for (const [tileId, tileObj] of Object.entries(tileData)) {
                let tile = tileObj;
                tile["average"] = tileObj.sum/tileObj.count;
                tile["category"] = getCategoryByValue(tile["average"]);
                tilesActual[tileId] = tile;
            }
            tileData = tilesActual;
        } else {
            tileData = {};
            stroke = true;
            radius = 10;
            lg = mainMapTilesSimulation;
            for (const [cellName, cellObj] of Object.entries(data.cells)) {
                let skip = false;
                for (let i = 0; i < $storeServingCells_selected.length; i++) {
                    if (cellName == $storeServingCells_selected[i]) {
                        skip = true;
                        break;
                    }
                }
                if (skip) {
                    continue
                }

                for (const [tileId, tileObj] of Object.entries(cellObj)) {
                    if (!(tileId in tileData)) {
                        tileData[tileId] = {
                            count: 0,
                            lat: 0,
                            lng: 0,
                            servingCellCount: 0,
                            sum: 0
                        }
                    }
                    let tile = tileData[tileId];
                    tile["count"] = tile["count"] + tileObj["count"];
                    tile["sum"] = tile["sum"] + tileObj["sum"];
                    tile["lat"] = tileObj["lat"];
                    tile["lng"] = tileObj["lng"];
                    tileData[tileId] = tile;
                }
            }
            for (const [tileId, tileObj] of Object.entries(tileData)) {
                let tile = tileObj;
                tile["average"] = tileObj.sum/tileObj.count;
                tile["category"] = getCategoryByValue(tile["average"]);
                tilesSimulation[tileId] = tile;
            }
            tileData = tilesSimulation;
        }

        let tilesDelta = [];
        for (const [tileId, tileObj] of Object.entries(tilesSimulation)) {
            if (tilesActual[tileId].category != tilesSimulation[tileId].category) {
                tilesDelta.push(tileObj);
            }
        }

        for (const [tileId, tileObj] of Object.entries(tileData)) {
            let circle = L.circle([tileObj.lat, tileObj.lng], {radius: radius}).setStyle({
                stroke: stroke,
                weight: 1, 
                color: color,
                fillOpacity: 0.6,
                fillColor: getColorByCategory(tileObj.category),
            })
            lg.addLayer(circle)
        }

        tilesDelta.forEach((e) => {
            let circle = L.circle([e.lat, e.lng], {radius: 18.5}).setStyle({
                stroke: true,
                weight: 2, 
                color: "black",
                fill: false,
            });
            lg.addLayer(circle);
        });
        mapObj.drawnItems.addLayer(lg)
    }

    async function sendReqsToGeosvr(regional, type, boundary) {
		const res = await fetch('http://localhost:9301/cell-tiles/raw', {
			method: 'POST',
            headers: {
                "Content-Type": "application/json",
            },
			body: JSON.stringify({
                "region": regional,
                "type": type,
                "data": boundary,
            })
		})
		let result = await res.json()
        resultSet = result;

        drawTiles(mainMap, result)
        getServingCells(result.cells)
	}

    function setupMap(obj) {
        obj.map = L.map(obj.id).setView(obj.center, obj.zoom);
        L.tileLayer('http://192.168.200.25/hot/{z}/{x}/{y}.png', {
            attribution: '© <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
        }).addTo(obj.map);
        obj.drawnItems = new L.FeatureGroup();
        obj.map.addLayer(obj.drawnItems);

        if (obj.drawControls) {
            var selectionType = '';
            var drawControl = new L.Control.Draw({
                position: 'topright',
                draw: {
                    polyline: {
                        metric: true
                    },
                    polygon: {
                        allowIntersection: true,
                        showArea: false,
                        drawError: {
                            //color: '#b00b00',
                            timeout: 1000
                        },
                        shapeOptions: {
                            color: '#0000ff',
                            fill: false,
                            fillOpacity: 0.1,
                            weight: 2,
                            
                        }
                    },
                    rectangle: {
                        shapeOptions: {
                            showArea: false,
                            fill: false,
                            fillOpacity: 0.1,
                            weight: 2,
                            color: '#0000ff',
                        }
                    },
                    circle: {
                        shapeOptions: {
                            showArea: false,
                            fill: false,
                            fillOpacity: 0.1,
                            weight: 2,
                            color: '#0000ff',
                        }
                    },
                    marker: false
                },
                edit: {
                    featureGroup: obj.drawnItems,
                    remove: true
                }
            });
        obj.map.addControl(drawControl);
        obj.map.on('draw:created', function (e) {
                var type = e.layerType,
                    layer = e.layer;
             if (type === 'marker') {
                    layer.bindPopup('A popup!');
                }
                obj.drawnItems.addLayer(layer);
                obj.map.fitBounds(layer.getBounds());

                let result;
             // sel.type = type;
                switch (type) {
                    case 'polyline':
                        // layer.editing.latlngs.forEach(element => {
                        //     sel.coords.push(element);
                        // });
                        break;
                    case 'polygon':
                        sendReqsToGeosvr("01", "geojson", layer.toGeoJSON())
                        // layer.editing.latlngs[0][0].forEach(element => {
                        //    sel.coords.push(element);
                        // });
                        break;
                    case 'rectangle':
                        sendReqsToGeosvr("01", "geojson", layer.toGeoJSON())
                        // layer.editing._shape._latlngs[0].forEach(element => {
                        //     sel.coords.push(element);
                        // });
                        break;
                    case 'circle':
                        // sel.coords.push(layer.editing._shape._latlng);
                        // sel.radius = layer.editing._shape._mRadius;
                        break;
                    default:
                }
                // let gjson = layer.toGeoJSON()
                // doPost(gjson)
            });
            obj.map.on('draw:edited', function (e) {
                var layers = e.layers;
                var countOfEditedLayers = 0;
                layers.eachLayer(function(layer) {
                    countOfEditedLayers++;
                    let gjson = layer.toGeoJSON()
                    // doPost(gjson)
                });
            });
        }
    }      

    onMount(async () => {
        mainMapTilesActual = L.layerGroup();
        mainMapTilesSimulation = L.layerGroup();
        getBoundaries()
        if(browser) {
            setupMap(mainMap);
            setupMap(addMap);
            setupMap(editMap);
            setupMap(deleteMap);
        }
    });
</script>
  
<Header company="Telkomsel" platformName={moduleName} bind:isSideNavOpen>
    <svelte:fragment slot="skip-to-content">
        <SkipToContent />
    </svelte:fragment>
    <HeaderUtilities>
        <!-- <HeaderAction
            text="Administrator"
            bind:isOpen={isOpen1}
            icon={UserAvatarFilledAlt}
            closeIcon={UserAvatarFilledAlt}
        >
            <HeaderPanelLinks>
                <HeaderPanelDivider>Switcher subject 1</HeaderPanelDivider>
                <HeaderPanelLink>Switcher item 1</HeaderPanelLink>
                <HeaderPanelLink>Switcher item 2</HeaderPanelLink>
                <HeaderPanelLink>Switcher item 3</HeaderPanelLink>
                <HeaderPanelLink>Switcher item 4</HeaderPanelLink>
                <HeaderPanelDivider>Switcher subject 2</HeaderPanelDivider>
                <HeaderPanelLink>Switcher item 1</HeaderPanelLink>
                <HeaderPanelLink>Switcher item 2</HeaderPanelLink>
                <HeaderPanelDivider>Switcher subject 3</HeaderPanelDivider>
                <HeaderPanelLink>Switcher item 1</HeaderPanelLink>
            </HeaderPanelLinks>
        </HeaderAction> -->
    </HeaderUtilities>
</Header>
  
<div class="container col start" style="width:100%; height:calc(100vh - 50px); margin-top:50px;">
    <div class="container row space-between border-bottom" id="global-toolbar" style="padding-top:2px;">
        <div class="container row space-between border-right" style="width:calc(35% - 4px); height:calc(100% - 20px); margin:10px;">
            <div class="container row end" style="width:190px; height:calc(100% - 20px);">
                <ComboBox
                    size="sm"
                    placeholder="Select regional"
                    items={$storeRegionals}
                    {shouldFilterItem}
                    on:select={(e) => {
                        storeRegionalsSelected.set(e.detail.selectedItem)
                        getSites()
                    }}
                    on:clear={(e) => {
                        storeRegionalsSelected.set({})
                        storeSites.set([]);
                    }}
                />
            </div>
            <div class="container row end" style="width:450px; height:calc(100% - 20px); margin-right:10px;">
                <div class="container col stretch" style="width:430px;">
                    <ComboBox
                        size="sm"
                        placeholder="Select boundary"
                        items={$storeBoundaries}
                        {shouldFilterItem}
                        on:select={(e) => {
                            layers = [
                                $storeBoundaries[e.detail.selectedId].boundary,
                            ]
                            redrawMap(mainMap, layers);
                        }}
                        on:clear={e => clearMap(mainMap)}
                    />
                </div>
                <OverflowMenu icon={MapBoundary} size="sm">
                    <OverflowMenuItem text="Refresh boundaries" on:click={() => {getBoundaries()}} />
                    <OverflowMenuItem text="Add boundary" on:click={() => MapBoundaryAddOpen = true} />
                    <OverflowMenuItem text="Edit boundary" on:click={() => MapBoundaryEditOpen = true} />
                    <OverflowMenuItem danger text="Delete boundary" on:click={() => MapBoundaryDeleteOpen = true} />
                </OverflowMenu>
            </div>
        </div>
        <div class="container row space-between border-right" style="width:calc(64% - 0px); height:calc(100% - 20px); margin:10px;">
            <div class="container row start">
                <p style="font-size:12px; padding-top:6px; color: #888; margin-right: 20px;">Data range: </p>
                <DatePicker
                    datePickerType="range"
                    on:change
                    valueFrom={"03/09/2022"}
                    valueTo={"03/09/2022"}
                    >
                    <DatePickerInput size="sm" placeholder="mm/dd/yyyy" />
                    <DatePickerInput size="sm" placeholder="mm/dd/yyyy" />
                </DatePicker>
                <!-- <Button on:click={() => drawerOpen = true}>Drawer</Button> -->
                <Button on:click={() => {drawTiles(mainMap, resultSet)}}>Drawer</Button>
            </div>
            <div class="container row end border-left" style="padding-left=20px;">
                <Button
                    kind="primary"
                    size="small"
                    icon={SettingsAdjust}
                    on:click={() => (AutomationOpen = true)}
                >Automation Settings</Button>
            </div>
        </div>
    </div>
    <div class="container row start" id="workspace">
        <div class="container col start border-right" style="width:35%; height:100%;">
            <div class="border-bottom" style="height:calc(100% - 50px);">
                <div id="main-map" style="width:{mainMap.width}px; height:{mainMap.height}px;"></div>
            </div>
            <div class="container row space-between" style="height:50px;">
                <div style="height:calc(100% - 20px); margin:10px; padding:10px;">Indicator Value:</div>
                <div style="height:calc(100% - 20px);" data-outline>
                    <Button kind="ghost" size="small" on:click={()=>{legendOpen = true}}>
                        Show legends
                    </Button>
                </div>
            </div>
            <Popover
                relative={false}
                caret
                align="top-right"
                open={legendOpen}
                closeOnOutsideClick
                on:click:outside={() => {
                    legendOpen = false;
                    console.log("on:click:outside");
                }}
            >
                <div style="left:300px;bottom:50px; padding:20px; width:250px; height:200px;">Content</div>
            </Popover>
        </div>
        <div class="container col start" style="width:100%; height:100%;">
            <Tabs type="container">
                <Tab label="Workspace" />
                <Tab label="Analytics" />
                <svelte:fragment slot="content">
                    <TabContent>
                        <div class="container row start" style="width:100%;">
                            <div class="container col start margin20 stretch" style="width:347px; overflow:hidden;">
                                <ComboBox
                                    size="sm"
                                    placeholder="Search site..."
                                    items={$storeSites}
                                    {shouldFilterItem}
                                    on:select={(e) => {
                                        storeNearbyCenter.set(e.detail.selectedItem);
                                        getNearbySites(
                                            $storeRegionalsSelected.id, 
                                            $storeNearbyCenter.lat,
                                            $storeNearbyCenter.lng,
                                            $storeNearbySitesLimit,
                                            $storeNearbySitesRadius
                                        );
                                        addMainMarker(
                                            mainMap, 
                                            $storeNearbyCenter.lng, 
                                            $storeNearbyCenter.lat, 
                                            $storeNearbySitesRadius,
                                            $storeNearbySites
                                        );
                                    }}
                                    on:clear={(e) => {
                                        storeNearbyCenter.set({});
                                        storeNearbySites.set([]);
                                        clearMarkers(mainMap);
                                    }}
                                />
                                <div style="min-height:20px;" />
                                <DataTable
                                    size="medium"
                                    sortable
                                    headers={[
                                        { key: "shortText", value: "Site Id" },
                                        { key: "distance", value: "Distance (m)" },
                                    ]}
                                    rows={$storeNearbySites}
                                    pageSize={table1PageSize}
                                    page={table1CurrentPage}
                                >
                                    <Toolbar>
                                        <ToolbarContent>
                                            <ToolbarSearch persistent value="" shouldFilterRows />
                                            <Button kind="tertiary" iconDescription="Tooltip text" icon={Settings} />
                                        </ToolbarContent>
                                    </Toolbar>
                                    <svelte:fragment slot="header" let:header>
                                        {#if header.key === "distance"}
                                            <div style="width:100%; text-align:right;">{header.value}</div>
                                        {:else}
                                            {header.value}
                                        {/if}
                                    </svelte:fragment>
                                    <svelte:fragment slot="cell" let:row let:cell>
                                        {#if cell.key === "distance"}
                                            <div style="width:100%; text-align:right;">{cell.value}</div>
                                        {:else}
                                            {cell.value}
                                        {/if}
                                    </svelte:fragment>
                                </DataTable>
                                <Pagination
                                    bind:pageSize={table1CurrentPage}
                                    bind:page={table1CurrentPage}
                                    totalItems={$storeNearbySites.length}
                                    pageSizeInputDisabled
                                    pageInputDisabled
                                />
                            </div>
                            <div class="container col start" style="width:820px;">
                                <ComboBox
                                    size="sm"
                                    placeholder="Minimum unit"
                                    items={[{id:1, text:"Tiles"}, {id:2, text: "Percents"}]}
                                    selectedId={2}
                                />
                                <div style="min-height:20px;" />
                                <DataTable
                                    selectable
                                    sortable
                                    size="medium"
                                    headers={[
                                        { key: "name", value: "Cell name", minwidth: "300px" },
                                        { key: "tileCount", value: "Tiles", width: "60px" },
                                        { key: "eventCount", value: "Events", width: "60px" },
                                        { key: "avgRsrp", value: "Avg. RSRP", width: "100px" },
                                    ]}
                                    rows={$storeServingCells}
                                    bind:selectedRowIds={$storeServingCells_selected}
                                    bind:pageSize={table2PageSize}
                                    bind:page={table2CurrentPage}
                                >
                                    <svelte:fragment slot="header" let:header>
                                        {#if header.key === "avgRsrp" || header.key === "eventCount" || header.key === "tileCount" }
                                            <div style="width:100%; text-align:right;">{header.value}</div>
                                        {:else}
                                            {header.value}
                                        {/if}
                                    </svelte:fragment>
                                    <svelte:fragment slot="cell" let:row let:cell>
                                        {#if cell.key === "avgRsrp" || cell.key === "eventCount" || cell.key === "tileCount" }
                                            <div style="width:100%; text-align:right;">{cell.value}</div>
                                        {:else}
                                            {cell.value}
                                        {/if}
                                    </svelte:fragment>
                                </DataTable>
                                <Pagination
                                    bind:pageSize={table2PageSize}
                                    bind:page={table2CurrentPage}
                                    totalItems={$storeServingCells.length}
                                    pageSizeInputDisabled
                                    pageInputDisabled
                                />
                            </div>
                            <!-- <div class="container col start" style="width:470px;">
                                <div class="container row start" style="width:100%">
                                    <div style="width:49%">
                                        <MultiSelect
                                            size="sm"
                                            label="Select technologies"
                                            items={[
                                                { id: "0", text: "Slack" },
                                                { id: "1", text: "Email" },
                                                { id: "2", text: "Fax" },
                                            ]}
                                        />
                                    </div>
                                    <div style="width:2%">
                                        
                                    </div>
                                    <div style="width:49%">
                                        <MultiSelect
                                            size="sm"
                                            label="Select bands"
                                            items={[
                                                { id: "0", text: "Slack" },
                                                { id: "1", text: "Email" },
                                                { id: "2", text: "Fax" },
                                            ]}
                                        />
                                    </div>
                                </div>
                                <div style="min-height:20px;" />
                                <DataTable
                                    selectable
                                    size="short"
                                    headers={[
                                        { key: "name", value: "Site name" },
                                        { key: "protocol", value: "Distance" },
                                    ]}
                                    rows={[
                                        {
                                            id: "a",
                                            name: "Load Balancer 3",
                                            protocol: "HTTP",
                                        },
                                        {
                                            id: "b",
                                            name: "Load Balancer 1",
                                            protocol: "HTTP",
                                        },
                                    ]}
                                />
                            </div> -->
                        </div>
                    </TabContent>
                    <TabContent>
                        <div class="container col start">
                            <div class="container row start">
                            
                            </div>
                            <div class="container row start"></div>
                        </div>
                        <div class="container row start">
                        </div>
                    </TabContent>
                </svelte:fragment>
            </Tabs>
        </div>
    </div>
</div>

<Drawer bind:open={drawerOpen} size='{1895 - mainMap.width}px' placement='right' on:clickAway={() => drawerOpen = false}>
    <Button on:click={() => drawerOpen = false}>Close</Button>
</Drawer>

<Modal
    size="lg"
    bind:open={MapBoundaryAddOpen}
    modalHeading="Create New Boundary"
    primaryButtonText="Confirm Creation"
    secondaryButtonText="Cancel"
    on:click:button--secondary={() => (MapBoundaryAddOpen = false)}
    on:open
    on:close={(e) => {
        clearMap(addMap)
    }}
    on:submit={(e) => {
        e.preventDefault();
        let name = document.getElementById("addBound-name").value
        let type = document.getElementById("addBound-type").value
        let tags = document.getElementById("addBound-tags").value
        let files = document.getElementById("addBound-file").files

        let fr = new FileReader();
        let newJson = {
            name: name,
            tags: tags,
            type: type,
        }
        fr.onload = function(e1) { 
            var result = e1.target.result;
            // var file = JSON.stringify(result, null, 2);
            newJson["file"] = btoa(unescape(encodeURIComponent(result)))
            addBoundary(newJson)
        }
        fr.readAsText(files.item(0));
    }}
>
    <div class="container row space-between" style="height:600px;">
        <div style="width:50%; height:100%">
            <div class="container col start stretch" style="padding-right:20px;">
                <TextInput id="addBound-name" labelText="Boundary name" placeholder="Enter boundary name..." required />
                <div style="min-height:20px;"/>
                <ComboBox
                    id="addBound-type"
                    titleText="Format type"
                    placeholder="Select format type"
                    items={[
                        {id: 0, text: "geojson", desc: "standard GeoJSON format"},
                        {id: 1, text: "mapinfo", desc: "standard MapInfo format" },
                        {id: 2, text: "simple-shape", desc: "support for Circle and Bounding box"},
                        {id: 3, text: "pg-coordlist", desc: "Polygon using coordinate list"},
                        {id: 4, text: "pl-coordlist", desc: "Polyline using coordinate list"},
                    ]}
                    on:select={(e) => {
                        
                    }}
                    on:clear={e => clearMap(mainMap)}
                    let:item
                >
                    <div>
                        <strong>{item.text}</strong>
                    </div>
                    <div>
                        {item.desc}
                    </div>
                </ComboBox>
                <div style="min-height:20px;"/>
                <TextInput id="addBound-tags" labelText="Tags" placeholder="Multiple tags separate by space" />
                <div style="min-height:20px;"/>
                <FileUploaderDropContainer
                    id="addBound-file"
                    labelText="Drag and drop files here or click to upload"
                    accept={[".json", ".zip", ".txt", "csv"]}
                    validateFiles={(files) => {
                        return files.filter((file) => true);
                    }}
                    on:change={(e) => {
                        let type = document.getElementById("addBound-type").value
                        let files = document.getElementById("addBound-file").files

                        let fr = new FileReader();
                        fr.onload = function(e1) { 
                            var result = e1.target.result;
                            console.log(result)
                            if (type == "geojson") {
                                redrawMap(addMap, [JSON.parse(result)]);
                            }
                        }
                        fr.readAsText(files.item(0));
                        
                        uploadStatus.fileDate = e.detail[0].lastModified;
                        uploadStatus.fileSize = e.detail[0].size;
                        uploadStatus.fileName = e.detail[0].name;
                        uploadStatus.status = uploadStatus.fileSize > 1024 * 1024 ? "failed": "success";

                    }}
                />
                {#if uploadStatus.status=="failed"}
                    <FileUploaderItem
                        invalid
                        id="readme"
                        name={uploadStatus.fileName}
                        errorSubject="File size exceeds 1.0MB limit"
                        errorBody="Please select a new file."
                        status="edit"
                        on:delete
                    />
                {/if}
                {#if uploadStatus.status=="success"}
                    <FileUploaderItem
                        id="readme"
                        name={uploadStatus.fileName}
                        status="complete"

                        on:delete
                    />
                {/if}
            </div>
        </div>
        <div style="width:50%; height:100%; background-color:#ddd;">
            <div id="add-map" style="width:{addMap.width}px; height:{addMap.height}px;"></div>
        </div>
    </div>
</Modal>

<Modal
    size="lg"
    bind:open={MapBoundaryEditOpen}
    modalHeading="Edit boundary"
    primaryButtonText="Close"
    
    on:click:button--secondary={() => (MapBoundaryEditOpen = false)}
    on:open
    on:close={(e) => {
        clearMap(editMap)
    }}
>
    <div class="container row space-between" style="height:600px;">
        <div class="container col start" style="width:calc(50% - 20px); height:100%">
            <DataTable
                sortable
                title="Available boundaries"
                description="Select boundary for editing, click 'save' once completed."
                headers={[
                    { key: "text", value: "Name" },
                    { key: "type", value: "Type" },
                    { key: "tag", value: "Tag" },
                ]}
                rows={$storeBoundaries}
                on:click:row={(e) => {
                    layers = [
                        e.detail.boundary,
                    ]
                    redrawMap(editMap, layers);
                }}
                >
                <Toolbar>
                    <ToolbarContent>
                    <ToolbarSearch persistent value="" shouldFilterRows />
                    </ToolbarContent>
                </Toolbar>
            </DataTable>
        </div>
        <div class="container col start" style="width:50%; height:100%;">
            <div class="map" style="width:100%; height:calc(100% - 51px); background-color:#ddd;">
                <!-- <LeafletMap id="edit-map" width={670} height={550} bind:layers={editMapLayers} drawControls={true} /> -->
                <div id="edit-map" style="width:{editMap.width}px; height:{editMap.height}px;"></div>
            </div>
            <div class="container row end" style="width:100%; height:51px; padding:10px 0;">
                <Button
                    kind="tertiary"
                    size="small"
                    icon={Reset}
                    on:click={() => (AutomationOpen = true)}
                >Reset</Button>
                <div style="min-width:10px"></div>
                <Button
                    kind="primary"
                    size="small"
                    icon={CheckboxChecked}
                    on:click={() => (AutomationOpen = true)}
                >Accept</Button>
            </div>
        </div>
    </div>
</Modal>

<Modal
    size="lg"
    danger
    bind:open={MapBoundaryDeleteOpen}
    modalHeading="Delete boundary"
    primaryButtonText="Confirm Delete"
    secondaryButtonText="Cancel"
    on:click:button--secondary={() => (MapBoundaryDeleteOpen = false)}
    on:open
    on:close
    on:submit
>
<div class="container row space-between" style="height:600px;">
    <div class="container col start" style="width:calc(50% - 20px); height:100%">
        <DataTable
            sortable
            title="Available boundaries"
            description="Select boundary for editing, click 'save' once completed."
            headers={[
                { key: "text", value: "Name" },
                { key: "type", value: "Type" },
                { key: "tag", value: "Tag" },
            ]}
            rows={$storeBoundaries}
            on:click:row={(e) => {
                layers = [
                    e.detail.boundary,
                ]
                redrawMap(deleteMap, layers);
            }}
            >
            <Toolbar>
                <ToolbarContent>
                <ToolbarSearch persistent value="" shouldFilterRows />
                </ToolbarContent>
            </Toolbar>
        </DataTable>
    </div>
    <div class="container col start" style="width:50%; height:100%;">
        <div class="map" style="width:100%; height:calc(100% - 51px); background-color:#ddd;">
            <div id="delete-map" style="width:{deleteMap.width}px; height:{deleteMap.height}px;"></div>
        </div>
    </div>
</div>
</Modal>

<Modal
  bind:open={AutomationOpen}
  modalHeading="Automation Settings"
  primaryButtonText="Confirm Save"
  secondaryButtonText="Cancel"
  on:click:button--secondary={() => (AutomationOpen = false)}
  on:open
  on:close
  on:submit
>
  <p>This is a permanent action and cannot be undone.</p>
</Modal>

<style>
    @import 'https://unpkg.com/leaflet@1.7.1/dist/leaflet.css';

    :global(.bx--list-box__menu-item, .bx--list-box__menu-item__option) {
        height: auto;
    }

    .container {
        display: flex;
        flex-wrap: nowrap;
    }

    .row {
        flex-direction: row;
    }

    .col {
        flex-direction: column;
    }

    .wrap {
        flex-wrap: wrap;
    }

    .start {
        justify-content: flex-start;
    }

    .end {
        justify-content: flex-end;
    }

    .center {
        justify-content: center;
    }

    .space-between {
        justify-content: space-between;
    }

    .space-around {
        justify-content: space-around;
    }

    .stretch {
        align-items: stretch;
    }

    .self-center {
        align-self: center;
    }

    .margin20 {
        margin-right: 20px;
    }

    .margin30 {
        margin-right: 30px;
    }

    .border-left {
        border-left: 1px solid #ddd;
    }

    .border-right {
        border-right: 1px solid #ddd;
    }

    .border-bottom {
        border-bottom: 1px solid #ddd;
    }

    #global-toolbar {
        height:60px;
    }

    #workspace {
        height:calc(100% - 61px);
    }

</style>