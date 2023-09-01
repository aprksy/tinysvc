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
    } from "carbon-components-svelte";
    import { Upload, View, TreeView, Area, AreaCustom, WatsonHealthCircleMeasurement, ChartNetwork, Reset, CheckboxChecked, OrderDetails } from "carbon-icons-svelte";
    import MapBoundary from "carbon-icons-svelte/lib/MapBoundary.svelte";
    import SettingsAdjust from "carbon-icons-svelte/lib/SettingsAdjust.svelte";
    import UserAvatarFilledAlt from "carbon-icons-svelte/lib/UserAvatarFilledAlt.svelte";
    import { onMount } from 'svelte';
    import { browser } from '$app/env';
    import { 
        storeRegionals, 
        storeBoundaries,
    } from "../lib/controller/store.js";
  
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

    function shouldFilterItem(item, value) {
        if (!value) return true;
        return item.text.toLowerCase().includes(value.toLowerCase());
    }

    // async function getBoundaries() {
	// 	const res = await fetch('http://127.0.0.1:9301/boundaries', {
	// 		method: 'GET',
    //         headers: {
    //             "Content-Type": "application/json",
    //         },
	// 		body: JSON.stringify(boundary)
	// 	})
	// 	result = await res.json()
	// }

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



    let moduleName = "Usecase 07 - Geospatial based KPI Analysis"

    let layers = [
        // {"type":"FeatureCollection","features":[{"type":"Feature","properties":{},"geometry":{"type":"Polygon","coordinates":[[[106.87491416931152,-6.259332270466314],[106.88023567199707,-6.26411011355174],[106.87946319580078,-6.269143865206112],[106.88401222229004,-6.270594259210957],[106.88796043395996,-6.264536704558167],[106.88718795776366,-6.2587350370072965],[106.88246726989746,-6.25566353986335],[106.87808990478516,-6.256175457307679],[106.87491416931152,-6.259332270466314]]]}}]},
    ];

    let mainMapLayers;
    let editMapLayers;

    const mainMap = {
        id: 'main-map',
        width: 664,
        height: 786,
        center: [-6.175392, 106.827153],
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
             // sel.type = type;
                switch (type) {
                    case 'polyline':
                        layer.editing.latlngs.forEach(element => {
                            // sel.coords.push(element);
                        });
                        break;
                    case 'polygon':
                        layer.editing.latlngs[0][0].forEach(element => {
                            // sel.coords.push(element);
                        });
                        break;
                    case 'rectangle':
                        layer.editing._shape._latlngs[0].forEach(element => {
                            // sel.coords.push(element);
                        });
                        break;
                    case 'circle':
                        // sel.coords.push(layer.editing._shape._latlng);
                        // sel.radius = layer.editing._shape._mRadius;
                        break;
                    default:
                }
                let gjson = layer.toGeoJSON()
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
        <div class="container col space-between">
            <div class="container col start">
                <div class="container row start">
                    
                </div>
                <div class="container row start"></div>
            </div>
            <div class="container row start">

            </div>
        </div>
    </div>
</div>

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