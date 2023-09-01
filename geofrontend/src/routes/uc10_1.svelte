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
    import { Upload, View, TreeView, Area, AreaCustom, WatsonHealthCircleMeasurement, ChartNetwork, Reset, CheckboxChecked, OrderDetails, Launch, Settings, Number_0, JoinRight } from "carbon-icons-svelte";
    import Drawer from 'svelte-drawer-component';
    import MapBoundary from "carbon-icons-svelte/lib/MapBoundary.svelte";
    import SettingsAdjust from "carbon-icons-svelte/lib/SettingsAdjust.svelte";
    import UserAvatarFilledAlt from "carbon-icons-svelte/lib/UserAvatarFilledAlt.svelte";
    import "carbon-components/css/carbon-components.min.css";
    import { BarChartGrouped, DonutChart } from "@carbon/charts-svelte";
    import "@carbon/charts/styles.min.css";
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
    
    import Index from "./index.svelte";

    let table1PageSize = 10;
    let table1CurrentPage = 1;

    let table2PageSize = 10;
    let table2CurrentPage = 1;

    let table3PageSize = 10;
    let table3CurrentPage = 1;

    let donutTiles = [];
    let donutLoading = donutTiles.length;

    let donutActual = [];
    let donutSimulation = [];
    let donutDelta = [];
    let barGroupActualVsSimulation = [];
  
    let drawerOpen = false;
    let isSideNavOpen = false;
    let isOpen1 = false;
    let legendOpen = false;
    let MapBoundaryAddOpen = false;
    let MapBoundaryEditOpen = false;
    let MapBoundaryDeleteOpen = false;
    let SubpageDetailAnalysis = false;
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
    let tilesDelta = [];
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

    let moduleName = "Usecase 10 - Blacksite v1.1"

    let layers = [
        // {"type":"FeatureCollection","features":[{"type":"Feature","properties":{},"geometry":{"type":"Polygon","coordinates":[[[106.87491416931152,-6.259332270466314],[106.88023567199707,-6.26411011355174],[106.87946319580078,-6.269143865206112],[106.88401222229004,-6.270594259210957],[106.88796043395996,-6.264536704558167],[106.88718795776366,-6.2587350370072965],[106.88246726989746,-6.25566353986335],[106.87808990478516,-6.256175457307679],[106.87491416931152,-6.259332270466314]]]}}]},
    ];

    let mainMapLayers;
    let editMapLayers;

    const mainMap = {
        id: 'main-map',
        width: 664,
        height: 678,
        // center: [-6.175392, 106.827153], // Jakarta
        center: [3.597031, 98.678513], // Medan
        zoom: 13,
        drawControls: true,
        drawnItems: null,
        tileLayerGroupActual: null,
        tileLayerGroupSimulation: null,
        tileLayerGroupDelta: null,
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

    const beforeMap = {
        id: 'before-map',
        width: 390,
        height: 400,
        center: [-6.175392, 106.827153],
        zoom: 13,
        drawControls: false,
        drawnItems: null,
        tileLayerGroupActual: null,
        tileLayerGroupSimulation: null,
        tileLayerGroupDelta: null,
        map: null,
    }

    const afterMap = {
        id: 'after-map',
        width: 390,
        height: 400,
        center: [-6.175392, 106.827153],
        zoom: 13,
        drawControls: false,
        drawnItems: null,
        tileLayerGroupActual: null,
        tileLayerGroupSimulation: null,
        tileLayerGroupDelta: null,
        map: null,
    }

    const deltaMap = {
        id: 'delta-map',
        width: 520,
        height: 400,
        center: [-6.175392, 106.827153],
        zoom: 13,
        drawControls: false,
        drawnItems: null,
        tileLayerGroupActual: null,
        tileLayerGroupSimulation: null,
        tileLayerGroupDelta: null,
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

    let c1 = "#1984c5"
    let c2 = "#6cd4c5"
    let c3 = "orange"
    let c4 = "#c80064"
    let categoryNames = ["Excellent", "Good", "Fair", "Poor"];
    let statusNames = ["Unchanged", "Safe", "Unsafe"];

    function getColorByCategory(value) {
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

    let barChartActual = [];
    let barChartSimulation = [];
    let barChartCompare = [];
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

        let unchangedTiles = {
            group: "Unchanged",
            value: 0
        }
        let safeTiles = {
            group: "Safe",
            value: 0
        }
        let unsafeTiles = {
            group: "Unsafe",
            value: 0
        }
        tilesDelta = [];
        for (const [tileId, tileObj] of Object.entries(tilesSimulation)) {
            if (tilesActual[tileId].category != tilesSimulation[tileId].category) {
                let status = (tilesSimulation[tileId].category - tilesActual[tileId].category < 0) ||
                        (tilesSimulation[tileId].category - tilesActual[tileId].category > 0 && tilesSimulation[tileId].category < 3) ?
                        "SAFE": "UNSAFE";
                tilesDelta.push({
                    "id": tileId,
                    "tileId": tileId,
                    "lat": tileObj.lat,
                    "lng": tileObj.lng,
                    "cat0": tilesActual[tileId].category,
                    "cat1": tilesSimulation[tileId].category,
                    "deltaCat": tilesSimulation[tileId].category - tilesActual[tileId].category,
                    "avgRSRP0": (tilesActual[tileId].average).toFixed(2),
                    "avgRSRP1": (tilesSimulation[tileId].average).toFixed(2),
                    "deltaAvgRSRP": (tilesSimulation[tileId].average - tilesActual[tileId].average).toFixed(2),
                    "status": status,
                });
                if (status == "SAFE") {
                    safeTiles.value = safeTiles.value + 1;
                } else {
                    unsafeTiles.value = unsafeTiles.value + 1;
                }
            } else {
                unchangedTiles.value = unchangedTiles.value + 1;
            }
        }
        donutTiles = [unchangedTiles, safeTiles, unsafeTiles];

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
        // beforeMap.drawnItems.addLayer(mainMapTilesActual)
        // afterMap.drawnItems.addLayer(mainMapTilesSimulation)
    }

    function drawActualTiles(mapObjs, data) {
        let stroke = false;
        let color = "#eee";
        let radius = 18.5;
        mainMapTilesActual.clearLayers();
        let tileData = data.tiles
        for (const [tileId, tileObj] of Object.entries(tileData)) {
            let tile = tileObj;
            tile["average"] = tileObj.sum/tileObj.count;
            tile["category"] = getCategoryByValue(tile["average"]);
            tilesActual[tileId] = tile;
        }
        tileData = tilesActual;
        mapObjs.forEach((e) => {
            let lg = L.layerGroup();
            for (const [tileId, tileObj] of Object.entries(tileData)) {
                let circle = L.circle([tileObj.lat, tileObj.lng], {radius: radius}).setStyle({
                    stroke: stroke,
                    weight: 1, 
                    color: color,
                    fillOpacity: 0.8,
                    fillColor: getColorByCategory(tileObj.category),
                })
                lg.addLayer(circle)
            }
            e.drawnItems.addLayer(lg);
        })
        return tileData
    }

    function drawSimulatedTiles(mapObjs, data, removedCells) {
        let tileData = {};
        if (removedCells.length > 0) {
            for (const [cellName, cellObj] of Object.entries(data.cells)) {
                let skip = false;
                for (let i = 0; i < removedCells.length; i++) {
                    if (cellName == removedCells[i]) {
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

        mapObjs.forEach((e) => {
            let stroke = false;
            let color = "#eee";
            let radius = 18.5;
            let opacity = 0.8;
            if (e == mainMap) {
                stroke = true;
                radius = 10;
                opacity = 1.0;
            }
            e.drawnItems.removeLayer(e.tileLayerGroupSimulation);
            let lg = L.layerGroup();
            if (removedCells.length > 0) {
                for (const [tileId, tileObj] of Object.entries(tileData)) {
                    let circle = L.circle([tileObj.lat, tileObj.lng], {radius: radius}).setStyle({
                        stroke: stroke,
                        weight: 1, 
                        color: color,
                        fillOpacity: opacity,
                        fillColor: getColorByCategory(tileObj.category),
                    })
                    lg.addLayer(circle)
                }
            }
            e.drawnItems.addLayer(lg);
            e.tileLayerGroupSimulation = lg;
        })
        return tileData
    }

    function drawDeltaTiles(mapObjs, actualData, simulationData) {
        let tileData = {};
        if (simulationData != {}) {
            for (const [tileId, tileObj] of Object.entries(actualData)) {
                let tile = {
                        "id": tileId,
                        "tileId": tileId,
                        "lat": tileObj.lat,
                        "lng": tileObj.lng,
                        "cat0": actualData[tileId].category,
                        // "cat1": simulationData[tileId].category,
                        // "deltaCat": simulationData[tileId].category - actualData[tileId].category,
                        "avgRSRP0": (actualData[tileId].average).toFixed(2),
                        // "avgRSRP1": (simulationData[tileId].average).toFixed(2),
                        // "deltaAvgRSRP": (simulationData[tileId].average - actualData[tileId].average).toFixed(2),
                        "status": 0,
                    }
                if (tileId in simulationData) {
                    if (actualData[tileId].category != simulationData[tileId].category) {
                        let status = (simulationData[tileId].category - actualData[tileId].category < 0) ||
                            (simulationData[tileId].category - actualData[tileId].category > 0 && simulationData[tileId].category < 3) ?
                            1: 2;
                        tile.cat0 = actualData[tileId].category;
                        tile.cat1 = simulationData[tileId].category;
                        tile.deltaCat = simulationData[tileId].category - actualData[tileId].category;
                        tile.avgRSRP0 = (actualData[tileId].average).toFixed(2);
                        tile.avgRSRP1 = (simulationData[tileId].average).toFixed(2);
                        tile.deltaAvgRSRP = (simulationData[tileId].average - actualData[tileId].average).toFixed(2);
                        tile.status = status;
                    }
                } else {
                    tile.status = 2;
                }
                tile.statusText = statusNames[tile.status];
                tileData[tileId] = tile;
            }
        }
        
        mapObjs.forEach((e) => {
            e.drawnItems.removeLayer(e.tileLayerGroupDelta);
        });
        if ($storeServingCells_selected.length > 0) {
            mapObjs.forEach((e) => {
                let stroke = false;
                let radius = 18.5;
                let opacity = 0.8;
                let fill = true;
                let color = "black";
                let fillColor = "blue";
                if (e == mainMap) {
                    radius = 18.5;
                    fill = false;
                }
                let lg = L.layerGroup();
                for (const [tileId, tileObj] of Object.entries(tileData)) {
                    if (tileObj.status == 0) {
                        fillColor = c3;
                        stroke = false;
                    } else if (tileObj.status == 1) {
                        fillColor = c1;
                        stroke = true && e == mainMap;
                    } else if (tileObj.status == 2) {
                        fillColor = c4;
                        stroke = true && e == mainMap;
                    } 
                    let circle = L.circle([tileObj.lat, tileObj.lng], {radius: radius}).setStyle({
                        stroke: stroke,
                        weight: 1, 
                        color: color,
                        fill: fill,
                        fillOpacity: opacity,
                        fillColor: fillColor,
                    })
                    lg.addLayer(circle)
                }
                e.drawnItems.addLayer(lg);
                e.tileLayerGroupDelta = lg;
            })
        }
        return tileData
    }

    function getTableData(data, fields, filterFunc) {
        let id = 0;
        let result = [];
        for (const [key, val] of Object.entries(data)) {
            if (filterFunc(val)) {
                continue
            }
            let item = {
                id: id,
            }
            for (let i= 0; i<fields.length; i++) {
                item[fields[i]] = val[fields[i]]
            }
            id++;
            result.push(item);
        }
        return result;
    }

    function getDonutChartData(data, field, groupNames) {
        let items = [];
        for (let i = 0; i < groupNames.length; i++) {
            items.push({
                group: groupNames[i],
                value: 0,
            })
        }
        for (const [tileId, tileObj] of Object.entries(data)) {
            let tile = items[tileObj[field]]
            tile.value = tile.value + 1;
            items[tileObj[field]] = tile;
        }
        return items;
    }

    function getBarGroupChartData(data, field, keyNames, groups) {
        function processData(data, f, kn, g) {
            let items = [];
            for (let i = 0; i < kn.length; i++) {
                items.push({
                    group: g,
                    key: kn[i],
                    value: 0,
                })
            }
            for (const [tileId, tileObj] of Object.entries(data)) {
                let tile = items[tileObj[f]]
                tile.value = tile.value + 1;
                items[tileObj[f]] = tile;
            }
            return items;
        }

        let result = [];
        for (let i=0; i < groups.length; i++) {
            let temp = processData(data[i], field, keyNames, groups[i]);
            result.push(...temp)
        }
        return result
    }

    function getIntegralDistribution(data, field) {
        let result = {};
        for (const [key, val] of Object.entries(data)) {
            // result[val[field]]
        }
    }

    let actualTiles;
    let simulationTiles;
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

        // drawTiles(mainMap, result)
        actualTiles = drawActualTiles([mainMap, beforeMap], result)
        donutActual = getDonutChartData(actualTiles, "category", categoryNames);
        getServingCells(result.cells)
	}

    function setupMap(obj) {
        obj.map = L.map(obj.id).setView(obj.center, obj.zoom);
        L.tileLayer('http://192.168.200.25/hot/{z}/{x}/{y}.png', {
            attribution: '© <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
        }).addTo(obj.map);
        obj.drawnItems = new L.FeatureGroup();
        // beforeMap.drawnItems = new L.FeatureGroup();
        // afterMap.drawnItems = new L.FeatureGroup();
        // deleteMap.drawnItems = new L.FeatureGroup();

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

                beforeMap.map.fitBounds(layer.getBounds());
                afterMap.map.fitBounds(layer.getBounds());
                deltaMap.map.fitBounds(layer.getBounds());

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
            // setupMap(addMap);
            // setupMap(editMap);
            // setupMap(deleteMap);
            setupMap(beforeMap);
            setupMap(afterMap);
            setupMap(deltaMap);
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
<div class="container col start" style="width:100vw; height:calc(100vh - 50px); margin-top:50px;">
    <!-- Toolbar and other non workflow functions -->
    <div class="control row border-bottom" style="width:100%; height:60px"></div>
    <div class="container row start" style="width:100vw; height:100%;">
        <!-- Requirement panel -->
        <div class="container col start border-right" style="width:calc(33% - 1px); height:100%;">
            <div class="container col" style="width:calc(100%); height:110px; padding: 14px;">
                <!-- date & regional panel -->
                <div class="container row space-between" style="width:calc(100%); height:40px;">
                    <div style="width:50%">
                        <DatePicker datePickerType="single" on:change>
                            <DatePickerInput size="sm" placeholder="mm/dd/yyyy" />
                        </DatePicker>
                    </div>
                    <div style="width:50%">
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
                </div>
                <!-- external boundary source panel -->
                <div class="container row space-between" style="width:calc(100%); align-items: center; height:40px;">
                    <div style="width:calc(50% - 5px)">
                        <ComboBox
                            size="sm"
                            placeholder="Select from registered boundary"
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
                    <div> - or - </div>
                    <Button size="sm" kind="tertiary">Load from File</Button>
                </div>
            </div>
            <!-- map panel -->
            <div class="container row start" style="width:100%; height:680px;">
                <div id="main-map" style="width:{mainMap.width}px; height:{mainMap.height}px;"></div>
            </div>
            <!-- legend panel -->
            <div class="container row space-between" style="width:calc(100% - 28px); height:60px; align-items: center; margin: 0 14px;">
                <div>Legend</div>
            </div>
        </div>
    
        <!-- Fine tune panel -->
        <div class="container col start border-right" style="width:calc(33% - 1px); height:100%;">
            <!-- select site -->
            <div class="container col" style="width:calc(100%); height:108px; align-items: center; padding: 14px;">
                <div style="width:100%; height:40px;">
                    <MultiSelect
                        size="sm"
                        label="Filter by site(s)"
                        items={[
                            { id: "0", text: "Slack" },
                            { id: "1", text: "Email" },
                            { id: "2", text: "Fax" },
                        ]}
                        />
                </div>
                <!-- select technology & band -->
                <div class="container row space-between" style="width:calc(100%); height:40px; align-items: center;">
                    <div style="width:50%">
                        <MultiSelect
                            size="sm"
                            label="Filter by technology(s)"
                            items={[
                                { id: "0", text: "Slack" },
                                { id: "1", text: "Email" },
                                { id: "2", text: "Fax" },
                            ]}
                            />
                    </div>
                    <div style="width:10px"></div>
                    <div style="width:50%">
                        <MultiSelect
                            size="sm"
                            label="Filter by band(s)"
                            items={[
                                { id: "0", text: "Slack" },
                                { id: "1", text: "Email" },
                                { id: "2", text: "Fax" },
                            ]}
                            />
                    </div>
                </div>
            </div>
            <!-- cell list panel -->
            <div class="container row start border-bottom" style="calc(100% - 60px); height:680px; padding: 0 14px;">
                <div class="container col space-between" style="width:100%; overflow-y:scroll; overflow-x: hidden; background-color:#eee;">
                    <DataTable
                        selectable
                        sortable
                        size="medium"
                        headers={[
                            { key: "name", value: "Cell name", minWidth: "120px" },
                            { key: "tileCount", value: "Tiles", width: "60px" },
                            { key: "avgRsrp", value: "Avg. RSRP", width: "80px" },
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
                                <div style="width:100%; padding:5px 0; text-align:right;">{cell.value}</div>
                            {:else}
                                <div style="width:200px; padding:5px 0; white-space: nowrap; text-overflow:clip;">{cell.value}</div>
                            {/if}
                        </svelte:fragment>
                    </DataTable>
                    {#if $storeServingCells.length == 0}
                        <div style="width:100%; text-align:center;">
                            <div style="font-weight:500;">No data</div>
                            <div>complete fields on the left panel,<br>and select an area on the map</div>
                        </div>
                    {/if}
                    <Pagination
                        bind:pageSize={table2PageSize}
                        bind:page={table2CurrentPage}
                        totalItems={$storeServingCells.length}
                        pageSizes={[10,15,20,50]}
                        pageInputDisabled
                    />
                </div>
            </div>
            <!-- command panel -->
            <div class="container row space-between" style="width:calc(100%); height:60px; align-items: center; padding: 0 14px;">
                <Button size="sm" kind="tertiary">Save scenario</Button>
                <div>
                    <Button size="sm" kind="tertiary" on:click={() => {
                            storeServingCells_selected.set([]);
                            simulationTiles = drawSimulatedTiles([mainMap, afterMap], resultSet, $storeServingCells_selected)
                            drawDeltaTiles([mainMap, deltaMap], actualTiles, simulationTiles)
                            donutSimulation = [];
                            donutDelta = [];
                            barGroupActualVsSimulation = [];
                        }
                    }>Reset</Button>
                    <Button size="sm" kind="primary" on:click={() => {
                            simulationTiles = drawSimulatedTiles([mainMap, afterMap], resultSet, $storeServingCells_selected)
                            donutSimulation = getDonutChartData(simulationTiles, "category", categoryNames)
                            let deltaTiles = drawDeltaTiles([mainMap, deltaMap], actualTiles, simulationTiles)
                            donutDelta = getDonutChartData(deltaTiles, "status", statusNames)
                            barGroupActualVsSimulation = getBarGroupChartData([actualTiles, simulationTiles], "category", 
                                categoryNames, ["Actual", "Simulation"])
                            tilesDelta = getTableData(deltaTiles, ["tileId", "avgRSRP0", "avgRSRP1", "statusText"],
                                (row) => {return row["status"] < 1})
                            console.log(tilesDelta)
                        }
                    }>Redraw</Button>
                </div>
            </div>
        </div>
    
        <!-- Result panel -->
        <div class="container col start" style="width:calc(34%); height:100%;">
            <div class="container row" style="width:calc(100%); height:282px; align-items: center; padding: 14px;">
                <div style="margin-top:0px; width:220px;">
                    <DonutChart 
                        data={donutDelta}
                        options={{
                            "data": {
                            },
                            "resizable": true,
                            "donut": {
                                "center": {
                                    "label": "Tiles"
                                },
                            },
                            "color": {
                                "scale": {
                                    "Unchanged": c3,
                                    "Safe": c2,
                                    "Unsafe": c4
                                }
                            },
                            "toolbar":{
                                "enabled": false,
                            },
                            "height": "200px",
                            "data": {
                                "loading": donutDelta.length == 0,
                            },
                        }}
                    />
                </div>
                <div style="margin-top:0px; width:400px;">
                    <BarChartGrouped
                    data={barGroupActualVsSimulation}
                        options={{
                            "axes": {
                                "left": {
                                    "mapsTo": "value"
                                },
                                "bottom": {
                                    "scaleType": "labels",
                                    "mapsTo": "key"
                                }
                            },
                            "height": "200px",
                            "color": {
                                "scale": {
                                    "Actual": c1,
                                    "Simulation": c2,
                                }
                            },
                            "toolbar":{
                                "enabled": false,
                            },
                            "data": {
                                "loading": barGroupActualVsSimulation.length == 0,
                            },
                        }}
                    />
                </div>
            </div>
            <!-- tile list panel -->
            <div class="container row start border-bottom" style="calc(100%); height:506px; padding: 0 14px;">
                <div class="container col space-between" style="width:100%; background-color:#eee; overflow-y:scroll; overflow-x:hidden;">
                    <DataTable
                        sortable
                        filterable
                        size="medium"
                        headers={[
                            { key: "tileId", value: "Tile Id", width: "80px"},
                            { key: "avgRSRP0", value: "RSRP Before", width: "80px"},
                            { key: "avgRSRP1", value: "RSRP After", width: "80px"},
                            { key: "statusText", value: "Status", width: "80px"},
                        ]}
                        rows={tilesDelta}
                        bind:pageSize={table3PageSize}
                        bind:page={table3CurrentPage}
                    >
                        <svelte:fragment slot="header" let:header>
                            {#if header.key === "before" || header.key === "after" || header.key === "delta" }
                                <div style="width:100%; text-align:right;">{header.value}</div>
                            {:else}
                                {header.value}
                            {/if}
                        </svelte:fragment>
                        <svelte:fragment slot="cell" let:row let:cell>
                            {#if cell.key === "avgRSRP0" || cell.key === "avgRSRP1" }
                                <div style="width:100%; text-align:right; padding:5px; background-color:{getColorByCategory(getCategoryByValue(cell.value))}">
                                    {cell.value}
                                </div>
                            {:else if cell.key === "statusText"}
                                {#if cell.value === "Unsafe"}
                                    <div style="width:100%; padding:5px; text-align:center; text-transform:uppercase; background-color:{c4}">
                                        {cell.value}
                                    </div>
                                {:else}
                                    <div style="width:100%; padding:5px; text-align:center; text-transform:uppercase; background-color:{c2}">
                                        {cell.value}
                                    </div>
                                {/if}
                            {:else}
                                <div style="width:300px; text-overflow:ellipsis;">{cell.value}</div>
                            {/if}
                        </svelte:fragment>
                    </DataTable>
                    {#if tilesDelta.length == 0}
                        <div style="width:100%; text-align:center;">
                            <div style="font-weight:500;">No data</div>
                            <div>check one or more cell on the middle panel,<br>and click redraw</div>
                        </div>
                    {/if}
                    <Pagination
                        bind:pageSize={table3PageSize}
                        bind:page={table3CurrentPage}
                        totalItems={tilesDelta.length}
                        pageSizes={[10,15,20,50]}
                        pageInputDisabled
                    />
                </div>
            </div>
            <!-- command panel -->
            <div class="container row space-between" style="width:calc(100%); height:60px; align-items: center; padding: 0 14px;">
                <div></div>
                <div>
                    <Button size="sm" kind="primary" on:click={() => SubpageDetailAnalysis = true}>Show details</Button>
                </div>
            </div>
        </div>
    </div>
</div>

<Modal
    size="lg"
    bind:open={SubpageDetailAnalysis}
    modalHeading="Detailed Analysis"
    primaryButtonText="Confirm Save Scenario"
    secondaryButtonText="Cancel"
    on:click:button--secondary={() => (SubpageDetailAnalysis = false)}
    on:open
    on:close
    on:submit
>
    <div class="container row space-between" style="height:720px; overflow-y:auto;">
        <!-- left panel -->
        <div class="container col start" style="width:calc(60% - 0px); height:100%;">
            <!-- maps panel -->
            <div class="container row start" style="width:calc(100% - 0px); height:500px;">
                <!-- before map -->
                <div class="container col start" style="width:calc(50%);">
                    <div style="height:35px; padding:10px;">BEFORE</div>
                    <div id="before-map" style="width:calc(100% - 0px); height:{beforeMap.height}px;"></div>
                </div>
                <!-- map separator -->
                <div style="width:14px;"></div>
                <!-- after map -->
                <div class="container col start" style="width:calc(50%);">
                    <div style="height:35px; padding:10px;">AFTER</div>
                    <div id="after-map" style="width:calc(100% - 0px); height:{afterMap.height}px;"></div>
                </div>
            </div>
            <!-- donut charts -->
            <div style="width:100%; padding:10px 0; text-align:center; text-transform:uppercase;">Proportion: RSRP Category</div>
            <div class="container row start" style="width:calc(100% - 0px); height:300px;">
                <div class="container col start" style="width:calc(50% - 7px);">
                    <div style="margin-top:0px; width:100%; padding-top:20px; height:300px;">
                        <DonutChart 
                            data={donutActual}
                            options={{
                                "resizable": true,
                                "donut": {
                                    "center": {
                                        "label": "Tiles"
                                    },
                                },
                                "color": {
                                    "scale": {
                                        "Excellent": c1,
                                        "Good": c2,
                                        "Fair": c3,
                                        "Poor": c4,
                                    }
                                },
                                "toolbar":{
                                    "enabled": false,
                                },
                                "legend": {
                                    "position": "right",
                                },
                                "height": "230px",
                                "data": {
                                    "loading": donutActual.length == 0,
                                },
                            }}
                        />
                    </div>
                </div>
                <div style="width:14px;"></div>
                <div class="container col start" style="width:calc(50% - 7px);">
                    <div style="margin-top:0px; width:100%; padding-top:20px; height:300px;">
                        <DonutChart 
                            data={donutSimulation}
                            options={{
                                "resizable": true,
                                "donut": {
                                    "center": {
                                        "label": "Tiles"
                                    },
                                },
                                "color": {
                                    "scale": {
                                        "Excellent": c1,
                                        "Good": c2,
                                        "Fair": c3,
                                        "Poor": c4,
                                    }
                                },
                                "toolbar":{
                                    "enabled": false,
                                },
                                "legend": {
                                    "position": "right",
                                },
                                "height": "230px",
                                "data": {
                                    "loading": donutSimulation.length == 0,
                                },
                            }}
                        />
                    </div>
                </div>
            </div>
            <!-- bar group chart -->
            <div class="border-top" style="width:100%; padding:10px 0; text-align:center; text-transform:uppercase;">
                Comparison: Number of Tiles in Category
            </div>
            <div style="margin-top:0px; width:calc(100% - 14px);">
                <BarChartGrouped
                data={barGroupActualVsSimulation}
                    options={{
                        "axes": {
                            "left": {
                                "mapsTo": "value"
                            },
                            "bottom": {
                                "scaleType": "labels",
                                "mapsTo": "key"
                            }
                        },
                        "height": "200px",
                        "color": {
                            "scale": {
                                "Actual": c1,
                                "Simulation": c2,
                            }
                        },
                        "toolbar":{
                            "enabled": false,
                        },
                        "data": {
                            "loading": barGroupActualVsSimulation.length == 0,
                        },
                    }}
                />
            </div>
            <div class="container col start" style="width:calc(60% - 14px); height:100%; overflow-y: auto;">
                
            </div>
        </div>
        <!-- column separator -->
        <div style="width:14px; height:100%"></div>
        <!-- right oanel -->
        <div class="container col start" style="width:40%; height:100%;">
            <div class="container col start" style="width:100%; height:1040px; background-color:#eee;">
                <div style="height:35px; padding:10px;">DELTA PROFILE</div>
                <div id="delta-map" style="width:{deltaMap.width}px; height:{deltaMap.height}px;"></div>
                <div style="width:100%; padding:10px 0; text-align:center; text-transform:uppercase;">Proportion: Category Delta</div>
                <div class="container row center" style="margin-top:0px; width:100%; padding-top:20px; height:300px;">
                    <DonutChart 
                        data={donutDelta}
                        options={{
                            "resizable": true,
                            "donut": {
                                "center": {
                                    "label": "Tiles"
                                },
                            },
                            "color": {
                                "scale": {
                                    "Unchanged": c3,
                                    "Safe": c1,
                                    "Unsafe": c4,
                                }
                            },
                            "toolbar":{
                                "enabled": false,
                            },
                            "legend": {
                                "position": "right",
                            },
                            "height": "320px",
                            "width": "320px",
                            "data": {
                                "loading": donutDelta.length == 0,
                            },
                        }}
                    />
                </div>
                <div style="margin-top:0px; width:calc(100% - 0px);">
                    <BarChartGrouped
                    data={[
                            {
                                "group": "Before",
                                "key": "Excellent",
                                "value": 65000
                            },
                            {
                                "group": "Before",
                                "key": "Good",
                                "value": -29123
                            },
                            {
                                "group": "Before",
                                "key": "Fair",
                                "value": -35213
                            },
                            {
                                "group": "Before",
                                "key": "Poor",
                                "value": 51213
                            },
                            {
                                "group": "After",
                                "key": "Excellent",
                                "value": 32432
                            },
                            {
                                "group": "After",
                                "key": "Good",
                                "value": -21312
                            },
                            {
                                "group": "After",
                                "key": "Fair",
                                "value": -56456
                            },
                            {
                                "group": "After",
                                "key": "Poor",
                                "value": -21312
                            }
                        ]}
                        options={{
                            "axes": {
                                "left": {
                                    "mapsTo": "value"
                                },
                                "bottom": {
                                    "scaleType": "labels",
                                    "mapsTo": "key"
                                }
                            },
                            "height": "200px",
                            "color": {
                                "scale": {
                                    "Before": c1,
                                    "After": c2,
                                }
                            },
                            "toolbar":{
                                "enabled": false,
                            },
                            "data": {
                                "loading": true,
                            },
                        }}
                    />
                </div>
            </div>
        </div>
    </div>
</Modal>

<div class="app-drawer">
    <Drawer bind:open={drawerOpen} size='950px' placement='bottom' on:clickAway={() => drawerOpen = false}>
        <Button on:click={() => drawerOpen = false}>Close</Button>
    </Drawer>
</div>

<style>
    @import 'https://unpkg.com/leaflet@1.7.1/dist/leaflet.css';

    :global(.bx--list-box__menu-item, .bx--list-box__menu-item__option) {
        height: auto;
    }

    .app-drawer :global(.drawer .overlay) {
        background: rgba(100, 100, 100, 0.5);
    }

    .app-drawer :global(.drawer .panel) {
        color: white;
        z-index: 99999;
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

    .border-top {
        border-top: 1px solid #ddd;
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