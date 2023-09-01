import * as color from "$lib/view/colors";

export const mainMap = {
    id: 'main-map',
    width: 664,
    height: 678,
    // center: [-6.175392, 106.827153], // Jakarta
    center: [3.597031, 98.678513], // Medan
    zoom: 13,
    drawControls: false,
    drawnItems: null,
    tileLayerGroupActual: null,
    tileLayerGroupSimulation: null,
    tileLayerGroupDelta: null,
    map: null,
}

export const addMap = {
    id: 'add-map',
    width: 670,
    height: 600,
    // center: [-6.175392, 106.827153], // Jakarta
    center: [3.597031, 98.678513], // Medan
    zoom: 13,
    drawControls: false,
    drawnItems: null,
    map: null,
}

export const editMap = {
    id: 'edit-map',
    width: 670,
    height: 550,
    // center: [-6.175392, 106.827153], // Jakarta
    center: [3.597031, 98.678513], // Medan
    zoom: 13,
    drawControls: true,
    drawnItems: null,
    map: null,
}

export const deleteMap = {
    id: 'delete-map',
    width: 670,
    height: 600,
    // center: [-6.175392, 106.827153], // Jakarta
    center: [3.597031, 98.678513], // Medan
    zoom: 13,
    drawControls: false,
    drawnItems: null,
    map: null,
}

export const beforeMap = {
    id: 'before-map',
    width: 390,
    height: 400,
    // center: [-6.175392, 106.827153], // Jakarta
    center: [3.597031, 98.678513], // Medan
    zoom: 13,
    drawControls: false,
    drawnItems: null,
    tileLayerGroupActual: null,
    tileLayerGroupSimulation: null,
    tileLayerGroupDelta: null,
    map: null,
}

export const afterMap = {
    id: 'after-map',
    width: 390,
    height: 400,
    // center: [-6.175392, 106.827153], // Jakarta
    center: [3.597031, 98.678513], // Medan
    zoom: 13,
    drawControls: false,
    drawnItems: null,
    tileLayerGroupActual: null,
    tileLayerGroupSimulation: null,
    tileLayerGroupDelta: null,
    map: null,
}

export const deltaMap = {
    id: 'delta-map',
    width: 520,
    height: 400,
    // center: [-6.175392, 106.827153], // Jakarta
    center: [3.597031, 98.678513], // Medan
    zoom: 13,
    drawControls: false,
    drawnItems: null,
    tileLayerGroupActual: null,
    tileLayerGroupSimulation: null,
    tileLayerGroupDelta: null,
    map: null,
}

let myIcon;

export function clearMap(obj) {
    obj.drawnItems.clearLayers();
}

export function clearAllMaps() {
    clearMap(mainMap);
    clearMap(beforeMap);
    clearMap(afterMap);
}

export function redrawMap(obj, data) {
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

export function drawBoundary(mapObj, boundaryData) {
    if (boundaryData) {
        let style = {
            stroke: true,
            // color: "#000",
            opacity: 1.0,
            weight: 3.0,
            fillOpacity: 0.1,
            dashArray: "10, 8",
        }
        if (boundaryData.areaType == "geojson") {
            let boundary = L.geoJson(boundaryData.data, {style: style}).addTo(mapObj.map);
            mapObj.drawnItems.addLayer(boundary);
            // mapData.data.features.forEach(el => {
            //     L.geoJson(el, {
            //         onEachFeature: function (feature, layer) {
            //             mapObj.drawnItems.addLayer(layer);
            //         }
            //     });
            // });
            var bounds = boundary.getBounds();
            mapObj.map.fitBounds(bounds);
        } else if (boundaryData.areaType == "circle") {
            let opt = style;
            opt["radius"] = boundaryData.data.radius
            let boundary = L.circle([boundaryData.data.lat, boundaryData.data.lng], opt).addTo(mapObj.map);
            mapObj.drawnItems.addLayer(boundary);
            var bounds = boundary.getBounds();
            mapObj.map.fitBounds(bounds);
        }
    }
}

export function drawTileKpi(mapObj, tileValues, tileLoc, tileRadius) {
    if (tileValues) {
        for (const [key, value] of Object.entries(tileValues)) {
            let lat = tileLoc[key].lat
            let lng = tileLoc[key].lng

            let tile = L.circle([lat, lng], {radius: tileRadius}).setStyle({
                stroke: false,
                weight: 1, 
                color: "#000",
                fillOpacity: 0.8,
                fillColor: color.byCategory(value.category1),
            });
            mapObj.drawnItems.addLayer(tile);
        }
    }
}

export function drawTileChanges(mapObj, tileValues, tileLoc, tileRadius, changeType) {
    if (tileValues) {
        for (const [key, value] of Object.entries(tileValues)) {
            let lat = tileLoc[key].lat
            let lng = tileLoc[key].lng

            let tileColor = color.byStatus(value.status);
            if (changeType == 'delta') {
                tileColor = color.byDelta(value.delta);
            }

            let tile = L.circle([lat, lng], {radius: tileRadius}).setStyle({
                stroke: false,
                weight: 1, 
                color: "#000",
                fillOpacity: 0.8,
                fillColor: tileColor,
            });
            mapObj.drawnItems.addLayer(tile);
        }
    }
}

export function drawSites(mapObj, sites, dismantledSites) {
    myIcon = L.icon({
        iconUrl: 'map-marker-icon-gray.png',
        iconSize: [18, 30], // size of the icon
    });
    // data.forEach((e, i) => {
    //     let marker = L.marker([e.lng, e.lat], {icon: myIcon});
    //     marker.bindPopup(e.text);
    //     mapObj.drawnItems.addLayer(marker);
    // })
    for (const [key, value] of Object.entries(sites)) {
        if (!dismantledSites.includes(key)) {
            let marker = L.marker([value.lat, value.lng]);
            marker.bindPopup(value.name);
            mapObj.drawnItems.addLayer(marker);
        }
    }
}

export function drawSimulationCategory(opt) {
    clearMap(opt.mapObj);
    drawBoundary(opt.mapObj, opt.boundaryData);
    drawTileKpi(opt.mapObj, opt.simData.tiles, opt.allTiles, opt.tileRadius);
    drawSites(opt.mapObj, opt.allSites, [opt.dismantledSite]);
}

export function drawOnMultimap(opt) {
    // clear all maps
    clearMap(mainMap);
    clearMap(beforeMap);
    clearMap(afterMap);
    // draw boundary on each map
    drawBoundary(mainMap, opt.boundaryData)
    drawBoundary(beforeMap, opt.boundaryData)
    drawBoundary(afterMap, opt.boundaryData)

    // draw original tiles
    drawTileKpi(mainMap, opt.oriData.tiles, opt.allTiles, opt.tileRadius);
    // draw simulated tiles
    drawTileKpi(beforeMap, opt.simData.tiles, opt.allTiles, opt.tileRadius);
    // draw changes
    // TODO: draw upgrade/unchange/degrade
    drawTileChanges(afterMap, opt.simData.tiles, opt.allTiles, opt.tileRadius, 'status');
    // draw sites
    drawSites(mainMap, opt.allSites, []);
    drawSites(beforeMap, opt.allSites, [opt.dismantledSite]);
    drawSites(afterMap, opt.allSites, [opt.dismantledSite]);
}

export let markerGroup;
function clearMarkers(obj) {
    obj.drawnItems.removeLayer(markerGroup);
}

export function setupMap(obj) {
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

            // beforeMap.map.fitBounds(layer.getBounds());
            // afterMap.map.fitBounds(layer.getBounds());
            // deltaMap.map.fitBounds(layer.getBounds());

            let result;
         // sel.type = type;
            switch (type) {
                case 'polyline':
                    // layer.editing.latlngs.forEach(element => {
                    //     sel.coords.push(element);
                    // });
                    break;
                case 'polygon':
                    // sendReqsToGeosvr("01", "geojson", layer.toGeoJSON())
                    // layer.editing.latlngs[0][0].forEach(element => {
                    //    sel.coords.push(element);
                    // });
                    break;
                case 'rectangle':
                    // sendReqsToGeosvr("01", "geojson", layer.toGeoJSON())
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