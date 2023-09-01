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
        Dropdown,
    } from "carbon-components-svelte";
    import { Upload, View, Download, Area, AreaCustom, WatsonHealthCircleMeasurement, ChartNetwork, Reset, CheckboxChecked, OrderDetails, Launch, Settings, Number_0, JoinRight, Label } from "carbon-icons-svelte";
    import Drawer from 'svelte-drawer-component';
    import MapBoundary from "carbon-icons-svelte/lib/MapBoundary.svelte";
    import SettingsAdjust from "carbon-icons-svelte/lib/SettingsAdjust.svelte";
    import UserAvatarFilledAlt from "carbon-icons-svelte/lib/UserAvatarFilledAlt.svelte";
    import "carbon-components/css/carbon-components.min.css";
    import { BarChartGrouped, BarChartSimple, DonutChart } from "@carbon/charts-svelte";
    import "@carbon/charts/styles.min.css";
    import { onMount } from 'svelte';
    import { browser } from '$app/env';
    import * as defaults from "../lib/model/defaults";
    import * as mainUi from "../lib/controller/mainUi";
    import * as maps from "../../src/lib/view/maps";
    import * as colors from "../../src/lib/view/colors";
    import * as chart from "../../src/lib/view/chartData";
    import * as client from "../../src/lib/controller/geoapiClient";
    import * as tabular from "../../src/lib/controller/tabular";
    import Bignumber from "$lib/view/bignumber.svelte";
    import Colorlegends from "$lib/view/colorlegends.svelte";
    import Tabular from "$lib/view/tabular.svelte";
    import Multitabular from "$lib/view/multitabular.svelte";

    let isSideNavOpen = false;
    let moduleName = "Usecase 10 - Blacksite v1.1";
    let shouldFilterItem = mainUi.shouldFilterItem;
    let boundaries;
    let boundaryId = '';
    let mapPage=1;
    let selectedDate="2022-05-20";
    let selectedRegion="01";
    let selectedGridSize="tileCovmo";
    let siteNames;
    let drawOptions;
    let exportPrefix = '20220520-boundaryId';
    let simulationResult;
    let simulationStat = {
        count: {value: 'N/A', change: 'N/A'},
        min: {value: 'N/A', change: 'N/A'},
        max: {value: 'N/A', change: 'N/A'},
        catExcellent: {value: 'N/A', change: 'N/A'},
        catGood: {value: 'N/A', change: 'N/A'},
        catFair: {value: 'N/A', change: 'N/A'},
        catPoor: {value: 'N/A', change: 'N/A'},
        deltaUpgraded: {value: 'N/A', change: 'N/A'},
        deltaDegraded: {value: 'N/A', change: 'N/A'},
        deltaUnchanged: {value: 'N/A', change: 'N/A'},
        statusSafe: {value: 'N/A', change: 'N/A'},
        statusUnsafe: {value: 'N/A', change: 'N/A'},
        statusFatal: {value: 'N/A', change: 'N/A'},
    };
    let chartsData = {
        donut: {
            tilesKpi: [],
            delta: [],
            status: [],
        },
        barGroup: {
            compareKpi: [],
            distribution: [],
        },
        bar: {
            distributionOriginal: [],
            distributionSimulation: [],
            colors: {},
        },
    }
    let tabularIndex = -1;
    let tableSites = {
        title: 'Sites in boundary',
        description: 'Sites which is located in the boundary',
        header: [
            {key: 'no', value: 'No'},
            {key: 'id', value: 'ID'},
            {key: 'name', value: 'Name'},
            {key: 'lat', value: 'Latitude'},
            {key: 'lng', value: 'Longitude'},
            {key: 'type', value: 'Type'},
        ],
    }
    let tableTilesOri = {
        title: 'Tiles before dismantle',
        description: 'Tiles value and category before dismantle simulation',
        header: [
            {key: 'no', value: 'No'},
            {key: 'id', value: 'ID'},
            {key: 'lat', value: 'Latitude'},
            {key: 'lng', value: 'Longitude'},
            {key: 'value', value: 'RSRP'},
            {key: 'category', value: 'Category'},
        ],
    }
    let tableSimulation = [];

    let tables = [];

    function onSimulationCompleted(data) {
        siteNames = data.siteNames;
        drawOptions = {
            dismantledSite: '',
            tileRadius: data.simulationResult['tileRadius'],
            allSites: data.simulationResult['sites'],
            allTiles: data.simulationResult['tiles'],
            boundaryData: data.boundaryData,
            simData: data.simulationResult['original'],
            oriData: data.simulationResult['original'],
        }
        simulationResult = data.simulationResult;

        // generate tables data
        tableSites['data'] = tabular.generateSitesTable(simulationResult)
        tableTilesOri['data'] = tabular.generateTilesOriTable(simulationResult)
        
        tables.push(
            {id:0, text: 'Sites in boundary', tabularData:tableSites},
            {id:1, text: 'Tiles before simulation', tabularData:tableTilesOri},
        )

        let i = 2;
        for (const [siteid, value] of Object.entries(simulationResult.simulation)) {
            let tableSim = {
                title: 'Site ' + siteid +' dismantling',
                description: 'Result from site ' + siteid + ' dismantling simulation',
                header: [
                    {key: 'no', value: 'No'},
                    {key: 'id', value: 'ID'},
                    {key: 'value0', value: 'RSRP 0'},
                    {key: 'category0', value: 'Category 0'},
                    {key: 'value1', value: 'RSRP 1'},
                    {key: 'category1', value: 'Category 1'},
                    {key: 'delta', value: 'Delta'},
                    {key: 'status', value: 'Status'},
                ],
                data: tabular.generateSimulationTables(simulationResult, siteid),
            }
            tables.push({id:i, text: tableSim['title'], tabularData:tableSim});
            i++;
        }
    }

    function exportToFile(fn, data, format) {
        let dataFile;
        let mime = 'text/plain';
        if (format == 'csv') {
            let header = [];
            let content = [];
            data.tabularData.header.forEach((e) => {
                header.push(e.value);
            });
            data.tabularData.data.forEach((e) => {
                let row = [];
                for (const [key, value] of Object.entries(e)) {
                    row.push(value);
                }
                content.push(row.join(','));
            });
            dataFile = header.join(',') + '\n' + content.join('\n')
            mime = 'text/csv';
        } else if (format == 'json') {
            mime = 'application/json';
            dataFile = JSON.stringify(data);
        }

        var blob = new Blob([dataFile], { type: mime + ';charset=utf-8;' });
        if (navigator.msSaveBlob) { // IE 10+
            navigator.msSaveBlob(blob, fn + '.' + format);
        } else {
            var link = document.createElement("a");
            if (link.download !== undefined) { // feature detection
                // Browsers that support HTML5 download attribute
                var url = URL.createObjectURL(blob);
                link.setAttribute("href", url);
                link.setAttribute("download", fn + '.' + format);
                link.style.visibility = 'hidden';
                document.body.appendChild(link);
                link.click();
                document.body.removeChild(link);
            }
        }
    }

    function computeStatistics(data0, data1) {
        let result = {
            'count': {
                'value': data1.count,
                'change': data1.count - data0.count > 0 ? 'up': data1.count - data0.count == 0 ? 'none': 'down',
            },
            'min': {
                'value': data1.min,
                'change': data1.min - data0.min > 0 ? 'up': data1.min - data0.min == 0 ? 'none': 'down',
            },
            'max': {
                'value': data1.max,
                'change': data1.max - data0.max > 0 ? 'up': data1.max - data0.max == 0 ? 'none': 'down',
            },
            'catExcellent': {
                'value': data1.category.EXCELLENT,
                'change': data1.category.EXCELLENT - data0.category.EXCELLENT > 0 ? 
                    'up': data1.category.EXCELLENT - data0.category.EXCELLENT == 0 ? 'none': 'down',
            },
            'catGood': {
                'value': data1.category.GOOD,
                'change': data1.category.GOOD - data0.category.GOOD > 0 ? 
                    'up': data1.category.GOOD - data0.category.GOOD == 0 ? 'none': 'down',
            },
            'catFair': {
                'value': data1.category.FAIR,
                'change': data1.category.FAIR - data0.category.FAIR > 0 ? 
                    'up': data1.category.FAIR - data0.category.FAIR == 0 ? 'none': 'down',
            },
            'catPoor': {
                'value': data1.category.POOR,
                'change': data1.category.POOR - data0.category.POOR > 0 ? 
                    'up': data1.category.POOR - data0.category.POOR == 0 ? 'none': 'down',
            },
            'deltaUpgraded': {
                'value': data1.deltaSummary ? data1.deltaSummary.UPGRADED: 'N/A',
                'change': 'none',
            },
            'deltaUnchanged': {
                'value': data1.deltaSummary ? data1.deltaSummary.UNCHANGE: 'N/A',
                'change': 'none',
            },
            'deltaDegraded': {
                'value': data1.deltaSummary ? data1.deltaSummary.DEGRADED: 'N/A',
                'change': 'none',
            },
            'statusSafe': {
                'value': data1.deltaSummary ? data1.statusSummary.SAFE: 'N/A',
                'change': 'none',
            },
            'statusUnsafe': {
                'value': data1.deltaSummary ? data1.statusSummary.UNSAFE: 'N/A',
                'change': 'none',
            },
            'statusFatal': {
                'value': data1.deltaSummary ? data1.statusSummary.FATAL: 'N/A',
                'change': 'none',
            },
        }
        simulationStat = result;
        return result;
    }

    function setupDistributionData(data) {
        let result = [];
        data.frequency.forEach(e => {
            result.push({"group": e[0], "value": e[1]});
        });
        return result;
    }

    function setupDistributionColor(data) {
        let result = {};
        data.frequency.forEach(e => {
            let key = e[0];
            let value = colors.fromValue(key);
            result[key.toString()] = value;
        });
        console.log(result);
        return result;
    }

    function setupChartsData(data0, data1) {
        chartsData = {
            donut: {
                tilesKpi: [
                    {"group": "EXCELLENT", "value": data1.category.EXCELLENT},
                    {"group": "GOOD", "value": data1.category.GOOD},
                    {"group": "FAIR", "value": data1.category.FAIR},
                    {"group": "POOR", "value": data1.category.POOR},
                ],
                delta: [
                    {"group": "UPGRADE", "value": data1.deltaSummary ? data1.deltaSummary.UPGRADED: 0},
                    {"group": "UNCHANGE", "value": data1.deltaSummary ? data1.deltaSummary.UNCHANGE: 0},
                    {"group": "DEGRADE", "value": data1.deltaSummary ? data1.deltaSummary.DEGRADED: 0},
                ],
                status: [
                    {"group": "SAFE", "value": data1.deltaSummary ? data1.statusSummary.SAFE: 0},
                    {"group": "UNSAFE", "value": data1.deltaSummary ? data1.statusSummary.UNSAFE: 0},
                    {"group": "FATAL", "value": data1.deltaSummary ? data1.statusSummary.FATAL: 0},
                ],
            },
            barGroup: {
                compareKpi: [
                    {"group":"AFTER", "key": "EXCELLENT", "value": data1.category.EXCELLENT},
                    {"group":"AFTER", "key": "GOOD", "value": data1.category.GOOD},
                    {"group":"AFTER", "key": "FAIR", "value": data1.category.FAIR},
                    {"group":"AFTER", "key": "POOR", "value": data1.category.POOR},
                    {"group":"BEFORE", "key": "EXCELLENT", "value": data0.category.EXCELLENT},
                    {"group":"BEFORE", "key": "GOOD", "value": data0.category.GOOD},
                    {"group":"BEFORE", "key": "FAIR", "value": data0.category.FAIR},
                    {"group":"BEFORE", "key": "POOR", "value": data0.category.POOR},
                ],
            },
            bar: {
                distributionOriginal: setupDistributionData(data0),
                distributionSimulation: setupDistributionData(data1),
                colors: setupDistributionColor(data1),
            },
        }
        //console.log(chartsData.bar.colors)
    }

    onMount(async () => {
        // get boundaries
        boundaries = await client.fetchBoundaries();

        let mainMapTilesActual = L.layerGroup();
        let mainMapTilesSimulation = L.layerGroup();
        // getBoundaries()
        if(browser) {
            maps.setupMap(maps.mainMap);
            maps.setupMap(maps.addMap);
            // setupMap(editMap);
            // setupMap(deleteMap);
            maps.setupMap(maps.beforeMap);
            maps.setupMap(maps.afterMap);
            // maps.setupMap(maps.deltaMap);
        }
        mapPage=0;
    });

    let mapBoundaryAddOpen = false;
    let uploadStatus = {
        fileDate: 0,
        fileSize: 0,
        fileName: '',
        status: '',
    }
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
<div class="container col start" style="width:100vw; height:calc(100vh - 50px); margin-top:50px; background-color:#fafafa;">
    <!-- Toolbar and other non workflow functions -->
    <!-- <div class="control row border-bottom" style="width:100%; height:60px"></div> -->
    <div class="container row start" style="width:100vw; height:100%; background-color:transparent;">
        <!-- Requirement panel -->
        <div class="container col start border-right" style="width:calc(33% - 1px); height:100%; background-color:#fff;">
            <div class="container col" style="width:calc(100%); height:110px; padding: 14px;">
                <!-- date & regional panel -->
                <div class="container row space-between" style="width:calc(100%); height:40px;">
                    <div style="width:50%">
                        <DatePicker 
                            light
                            datePickerType="single" 
                            bind:value={selectedDate} 
                            dateFormat="Y-m-d">
                            <DatePickerInput size="sm" placeholder="mm/dd/yyyy" />
                        </DatePicker>
                    </div>
                    <div style="width:25%">
                        <Dropdown
                            type="inline"
                            size="sm"
                            placeholder="Select regional"
                            bind:selectedId={selectedRegion}
                            items={defaults.regionals}
                            on:select={(e) => {
                                //storeRegionalsSelected.set(e.detail.selectedItem)
                                // getSites()
                            }}
                            on:clear={(e) => {
                                //storeRegionalsSelected.set({})
                                //storeSites.set([]);
                            }}
                        />
                    </div>
                    <div style="width:25%">
                        <Dropdown
                            type="inline"
                            size="sm"
                            placeholder="Select grid size"
                            bind:selectedId={selectedGridSize}
                            items={defaults.gridSizes}
                            on:select={(e) => {
                                //storeRegionalsSelected.set(e.detail.selectedItem)
                                // getSites()
                            }}
                            on:clear={(e) => {
                                //storeRegionalsSelected.set({})
                                //storeSites.set([]);
                            }}
                        />
                    </div>
                </div>
                <!-- external boundary source panel -->
                <div class="container row space-between" style="width:calc(100%); align-items: center; height:40px;">
                    <div class="container row space-between" style="width:calc(50% - 0); align-items:center;">
                        <ComboBox
                            light
                            size="sm"
                            placeholder="Select from registered boundary"
                            items={boundaries}
                            on:select={async (e) => {
                                boundaryId = e.detail.selectedItem.text;
                                let params = {
                                    date: selectedDate.replace(/-/g, ''),
                                    region: selectedRegion,
                                    boundaryId: boundaryId,
                                    tileField: selectedGridSize,
                                    kpi: 'rsrp',
                                }
                                await client.performSimulation(params, onSimulationCompleted);
                                drawOptions.mapObj = maps.mainMap;
                                if (boundaryId != '') {
                                    if (mapPage == 0 || mapPage == 2) {
                                        maps.drawSimulationCategory(drawOptions);
                                    } else if (mapPage == 1) {
                                        maps.drawOnMultimap(drawOptions);
                                    }
                                }
                                drawOptions.simData = simulationResult["original"];
                                computeStatistics(simulationResult["original"], drawOptions.simData);
                                setupChartsData(simulationResult["original"], drawOptions.simData);
                                tabularIndex = tabularIndex != -1 ? tabularIndex: 0;
                            }}
                            on:clear={(e) => {
                                maps.clearAllMaps();
                                computeStatistics({}, {});
                                setupChartsData({}, {});
                                e.detail.selectedItem = null;
                            }}
                        />
                        <Button 
                            size="small" 
                            iconDescription="Custom Boundary" 
                            icon={AreaCustom} 
                            on:click={e => {
                                mapBoundaryAddOpen=true;
                            }}
                        />
                    </div>
                    <div style="width:12px"></div>
                    <div style="width:calc(50% + 0px)">
                        <ComboBox
                            light
                            size="sm"
                            placeholder="Select site to simulate dismantle"
                            items={siteNames}
                            on:select={(e) => {
                                drawOptions.dismantledSite = e.detail.selectedId;
                                if (e.detail.selectedId != '') {
                                    drawOptions.simData = simulationResult["simulation"][e.detail.selectedId];
                                } else {
                                    drawOptions.simData = simulationResult["original"];
                                }
                                if (boundaryId != '') {
                                    if (mapPage == 0 || mapPage == 2) {
                                        maps.drawSimulationCategory(drawOptions);
                                    } else if (mapPage == 1) {
                                        maps.drawOnMultimap(drawOptions);
                                    }
                                }
                                computeStatistics(simulationResult["original"], drawOptions.simData);
                                setupChartsData(simulationResult["original"], drawOptions.simData);
                            }}
                            on:clear={(e) => {
                                drawOptions.dismantledSite = '';
                                drawOptions.simData = simulationResult["original"];
                                if (boundaryId != '') {
                                    if (mapPage == 0 || mapPage == 2) {
                                        maps.drawSimulationCategory(drawOptions);
                                    } else if (mapPage == 1) {
                                        maps.drawOnMultimap(drawOptions);
                                    }
                                }
                                computeStatistics(simulationResult["original"], drawOptions.simData);
                                setupChartsData(simulationResult["original"], drawOptions.simData);
                            }}
                        />
                    </div>
                </div>
            </div>
            <!-- map panel -->
            <div class="container row start" style="width:100%; height:680px;">
                <div id="main-map" style="width:{maps.mainMap.width}px; height:{maps.mainMap.height}px;"></div>
            </div>
            <!-- legend panel -->
            <div class="container row space-between" style="width:calc(100% - 28px); height:60px; align-items: center; margin: 0 14px;">
                <Colorlegends items={colors.legends('original')}/>
            </div>
        </div>
    
        <!-- Fine tune panel -->
        <div class="container col start" style="width:calc(67% - 12px); height:100%; background-color:#fafafa; overflow-y:clip;">
            <Tabs 
                bind:selected={mapPage}
                on:change={(e) => {
                    if (boundaryId != '') {
                        if (mapPage == 0 || mapPage == 2) {
                            maps.drawSimulationCategory(drawOptions);
                        } else if (mapPage == 1) {
                            maps.drawOnMultimap(drawOptions);
                        }
                    }
                }}
            >
                <Tab label="Dashboard" />
                <Tab label="Multimap" />
                <Tab label="Tables" />
                <svelte:fragment slot="content">
                    <TabContent>
                        <div class="container col" style="width:calc(100%); height:100%;">
                            <div class="container col start" style="width:100%; height:calc(100% - 60px);">
                                <div style="display:flex; flex-flow:row nowrap; width:100%;">
                                    <Bignumber 
                                        field="tiles" 
                                        value={simulationStat.count.value} 
                                        color="gray" 
                                        direction={simulationStat.count.change}
                                        width="calc((100% - 40px)/5)"/>
                                    <Bignumber 
                                        field="EXCELLENT" 
                                        value={simulationStat.catExcellent.value} 
                                        color={colors.byCategory("EXCELLENT")} 
                                        direction={simulationStat.catExcellent.change}
                                        width="calc((100% - 40px)/5)"/>
                                    <Bignumber 
                                        field="GOOD" 
                                        value={simulationStat.catGood.value} 
                                        color={colors.byCategory("GOOD")} 
                                        direction={simulationStat.catGood.change}
                                        width="calc((100% - 40px)/5)"/>
                                    <Bignumber 
                                        field="FAIR" 
                                        value={simulationStat.catFair.value} 
                                        color={colors.byCategory("FAIR")} 
                                        direction={simulationStat.catFair.change}
                                        width="calc((100% - 40px)/5)"/>
                                    <Bignumber 
                                        field="POOR" 
                                        value={simulationStat.catPoor.value} 
                                        color={colors.byCategory("POOR")} 
                                        direction={simulationStat.catPoor.change}
                                        rightMost 
                                        width="calc((100% - 40px)/5)"/>
                                </div>
            
                                <div class="container row start" style="width:100%; height:430px; background-color:#fff; margin:12px 0; border:1px solid #eee;">
                                    <div class="container col start" style="width:400px; padding:12px;">
                                        <DonutChart 
                                            data={chartsData.donut.tilesKpi}
                                            options={{
                                                "data": {
                                                },
                                                "resizable": true,
                                                "donut": {
                                                    "center": {
                                                        "label": "Tiles"
                                                    },
                                                    "alignment": "center",
                                                },
                                                "legend": {
                                                    "alignment": "center",
                                                },
                                                "color": {
                                                    "scale": {
                                                        "EXCELLENT": colors.c1,
                                                        "GOOD": colors.c2,
                                                        "FAIR": colors.c3,
                                                        "POOR": colors.c4
                                                    }
                                                },
                                                "toolbar":{
                                                    "enabled": false,
                                                },
                                                "height": "385px",
                                                "data": {
                                                    "loading": chartsData.donut.tilesKpi.length == 0,
                                                },
                                            }}
                                        />
                                    </div>
                                    <div class="container col start" style="width:300px; padding:12px;">
                                        <BarChartGrouped
                                            data={chartsData.barGroup.compareKpi}
                                                options={{
                                                    "axes": {
                                                        "bottom": {
                                                            "mapsTo": "value"
                                                        },
                                                        "left": {
                                                            "scaleType": "labels",
                                                            "mapsTo": "key"
                                                        }
                                                    },
                                                    "bars": {
                                                        "width": 40,
                                                    },
                                                    "height": "385px",
                                                    "color": {
                                                        "scale": {
                                                            "BEFORE": colors.c1,
                                                            "AFTER": colors.c2,
                                                        }
                                                    },
                                                    "toolbar":{
                                                        "enabled": false,
                                                    },
                                                    "data": {
                                                        "loading": chartsData.barGroup.compareKpi == 0,
                                                    },
                                                }}
                                            />
                                    </div>
                                    <div class="container col start" style="width:calc(100% - 700px); border-right:1px solid #eee; padding:12px;">
                                        <BarChartSimple
                                            data={chartsData.bar.distributionSimulation}
                                                options={{
                                                    "axes": {
                                                        "left": {
                                                            "mapsTo": "value"
                                                        },
                                                        "bottom": {
                                                            "scaleType": "linear",
                                                            "mapsTo": "group",
                                                            "includeZero": false,
                                                        }
                                                    },
                                                    "legend": {
                                                        "enabled": false,
                                                    },
                                                    "height": "385px",
                                                    "getFillColor": function(group) {
                                                        return colors.fromValue(group);
                                                    },
                                                    "toolbar":{
                                                        "enabled": false,
                                                    },
                                                    "data": {
                                                        "loading": chartsData.bar.distributionSimulation == 0,
                                                    },
                                                }}
                                            />
                                    </div>
                                </div>
            
                                <div class="container row start" style="width:100%; height:300px;">
                                    <div class="container row start" style="width:calc(50% - 6px); height:282px; margin-right:12px">
                                        <div style="display:flex; flex-flow:column nowrap; width:200px; margin-right:12px;">
                                            <Bignumber 
                                                field="UPGRADE" 
                                                value={simulationStat.deltaUpgraded.value} 
                                                color={colors.byDelta('UPGRADED')}
                                                direction={simulationStat.deltaUpgraded.change}
                                                width="200px"/>
                                            <div style="height:12px"/>
                                            <Bignumber 
                                                field="UNCHANGE" 
                                                value={simulationStat.deltaUnchanged.value} 
                                                color={colors.byDelta('UNCHANGE')}
                                                direction={simulationStat.deltaUnchanged.change}
                                                width="200px"/>
                                            <div style="height:12px"/>
                                            <Bignumber 
                                                field="DEGRADE" 
                                                value={simulationStat.deltaDegraded.value} 
                                                color={colors.byDelta('DEGRADED')}
                                                direction={simulationStat.deltaDegraded.change}
                                                width="200px"/>
                                        </div>
                                        <div style="width:calc(100% - 212px); padding:5px 0; background-color:#fff; border:1px solid #eee;">
                                            <DonutChart 
                                                data={chartsData.donut.delta}
                                                options={{
                                                    "data": {
                                                    },
                                                    "resizable": true,
                                                    "donut": {
                                                        "center": {
                                                            "label": "Tiles"
                                                        },
                                                        "alignment": "center",
                                                    },
                                                    "legend": {
                                                        "alignment": "center",
                                                    },
                                                    "color": {
                                                        "scale": {
                                                            "UPGRADE": colors.c2,
                                                            "UNCHANGE": colors.c3,
                                                            "DEGRADE": colors.c4,
                                                        }
                                                    },
                                                    "toolbar":{
                                                        "enabled": false,
                                                    },
                                                    "height": "250px",
                                                    "data": {
                                                        "loading": chartsData.donut.delta.length == 0,
                                                    },
                                                }}
                                            />
                                        </div>
                                    </div>
                                    <div class="container row start" style="width:calc(50% - 6px); height:282px;">
                                        <div style="display:flex; flex-flow:column nowrap; width:200px; margin-right:12px;">
                                            <Bignumber 
                                                field="SAFE" 
                                                value={simulationStat.statusSafe.value} 
                                                color={colors.byStatus('SAFE')}
                                                direction={simulationStat.statusSafe.change}
                                                width="200px"/>
                                            <div style="height:12px"/>
                                            <Bignumber 
                                                field="UNSAFE" 
                                                value={simulationStat.statusUnsafe.value} 
                                                color={colors.byStatus('UNSAFE')}
                                                direction={simulationStat.statusUnsafe.change}
                                                width="200px"/>
                                            <div style="height:12px"/>
                                            <Bignumber 
                                                field="FATAL" 
                                                value={simulationStat.statusFatal.value} 
                                                color={colors.byStatus('FATAL')}
                                                direction={simulationStat.statusFatal.change}
                                                width="200px"/>
                                        </div>
                                        <div style="width:calc(100% - 212px); padding:5px 0; background-color:#fff; border:1px solid #eee;">
                                            <DonutChart 
                                                data={chartsData.donut.status}
                                                options={{
                                                    "data": {
                                                    },
                                                    "resizable": true,
                                                    "donut": {
                                                        "center": {
                                                            "label": "Tiles"
                                                        },
                                                        "alignment": "center",
                                                    },
                                                    "legend": {
                                                        "alignment": "center",
                                                    },
                                                    "color": {
                                                        "scale": {
                                                            "SAFE": colors.c2,
                                                            "UNSAFE": colors.c4,
                                                            "FATAL": "black",
                                                        }
                                                    },
                                                    "toolbar":{
                                                        "enabled": false,
                                                    },
                                                    "height": "250px",
                                                    "data": {
                                                        "loading": chartsData.donut.status.length == 0,
                                                    },
                                                }}
                                            />
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </TabContent>
                    <TabContent>
                        <div class="container col start" style="border:1px solid #eee;">
                            <div class="container row start" style="width:100%; height:54px; background-color:white;">
                                <div style="width:calc(50% - 6px); margin-right:12px; padding:15px;">
                                    <div>AFTER 
                                        <span style="font-weight:600; color:#08a;">{drawOptions ? drawOptions.dismantledSite:'(no site selected)'}</span> 
                                            DISMANTLING
                                    </div>
                                </div>
                                <div style="width:calc(50% - 6px); padding:15px;">
                                    <div>
                                        <span style="font-weight:600; color:#08a;">{drawOptions ? drawOptions.dismantledSite:'(no site selected)'}</span> 
                                        FINAL STATUS
                                    </div>
                                </div>
                            </div>
                            <div class="container row start" style="width:100%; height:calc(100% - 75px); background-color:white;">
                                <div id="before-map" style="width:calc(50% - 6px); height:676px; margin-right:12px;"></div>
                                <div id="after-map" style="width:calc(50% - 6px); height:676px;"></div>
                            </div>
                            <div class="container row start" style="width:100%; height:71px; background-color:white;">
                                <div style="width:calc(50% - 6px); margin-right:12px; padding:15px;">
                                    <Colorlegends items={colors.legends('simulation')}/>
                                </div>
                                <div style="width:calc(50% - 6px); padding:15px;">
                                    <Colorlegends items={colors.legends('final-status')}/>
                                </div>
                            </div>
                        </div>
                    </TabContent>
                    <TabContent>
                        <div class="container col start" style="background-color:white; height:803px; border:1px solid #eee; padding:12px;">
                            <div class="container row space-between" style="width:100%; align-items:center; border-bottom:1px solid #eee; padding-bottom:12px;">
                                <div class="container col stretch" style="width:400px;">
                                    <ComboBox
                                        size="sm"
                                        placeholder="Select table"
                                        items={tables}
                                        selectedId={tabularIndex}
                                        on:select={async (e) => {
                                            tabularIndex = e.detail.selectedId;
                                        }}
                                        on:clear={(e) => {
                                            tabularIndex = -1
                                        }}
                                    />
                                </div>
                                <div class="container row end" style="width:750px;">
                                    <OverflowMenu open flipped  style="width: auto;">
                                        <div slot="menu" class="container row" style="padding: 1rem;">
                                            <div style="margin-right:12px;">Download</div>
                                            <Download />
                                        </div>
                                        <OverflowMenuItem 
                                            text="This table (CSV)" 
                                            on:click={(e) => {
                                                let fn = exportPrefix + '-' + 
                                                    tables[tabularIndex].tabularData.title.split(' ').join('_').toLowerCase();
                                                exportToFile(fn, tables[tabularIndex], 'csv');
                                            }}
                                        />
                                        <OverflowMenuItem text="All CSV (zipped)" hasDivider/>
                                        <OverflowMenuItem text="All (zipped)"/>
                                        <OverflowMenuItem 
                                            text="Raw (JSON)" 
                                            hasDivider
                                            on:click={(e) => {
                                                console.log(simulationResult);
                                                let fn = exportPrefix + '-simulation_result';
                                                exportToFile(fn, simulationResult, 'json');
                                            }}
                                        />
                                      </OverflowMenu>
                                </div>
                            </div>
                            <Multitabular>
                                {#each tables as table}
                                    <Tabular 
                                        title={table.tabularData.title} 
                                        description={table.tabularData.description} 
                                        header={table.tabularData.header}
                                        data={table.tabularData.data}
                                        visible={tabularIndex==table.id}
                                    />
                                {/each}
                            </Multitabular>
                        </div>
                    </TabContent>
                </svelte:fragment>
            </Tabs>
        </div>
    </div>
</div>

<Modal
    passiveModal
    size="lg"
    bind:open={mapBoundaryAddOpen}
    modalHeading="Manage Boundaries"
    on:click:button--secondary={() => (mapBoundaryAddOpen = false)}
    on:open
    on:close={(e) => {
        maps.clearMap(maps.addMap)
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
    <Tabs 
    >
        <Tab label="Draw Map" />
        <Tab label="Upload Map" />
        <Tab label="Edit Map" />
        <Tab label="Delete Map" />
        <svelte:fragment slot="content">
            <TabContent></TabContent>
            <TabContent></TabContent>
            <TabContent></TabContent>
            <TabContent></TabContent>
        </svelte:fragment>
    </Tabs>
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
                                maps.redrawMap(maps.addMap, [JSON.parse(result)]);
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
            <div id="add-map" style="width:{maps.addMap.width}px; height:{maps.addMap.height}px;"></div>
        </div>
    </div>
</Modal>

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