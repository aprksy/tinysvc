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
    import * as defaults from "../lib/model/defaults";
    import * as mainUi from "../lib/controller/mainUi";
    import * as maps from "../../src/lib/view/maps";
    import * as colors from "../../src/lib/view/colors";
    import * as chart from "../../src/lib/view/chartData";
    import * as client from "../../src/lib/controller/geoapiClient";
    import Bignumber from "$lib/view/bignumber.svelte";

    let isSideNavOpen = false;
    let moduleName = "Usecase 10 - Blacksite v1.1";
    let shouldFilterItem = mainUi.shouldFilterItem;
    let boundaries;
    let siteNames;
    let drawOptions;
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

    function onSimulationCompleted(data) {
        siteNames = data.siteNames;
        drawOptions = {
            dismantledSite: '',
            allSites: data.simulationResult['sites'],
            allTiles: data.simulationResult['tiles'],
            boundaryData: data.boundaryData,
            simData: data.simulationResult['original'],
        }
        simulationResult = data.simulationResult;
        // map.simOptions.dismantledSite = "";
        // map.simOptions.allSites = result["sites"];
        // map.simOptions.allTiles = result["tiles"];
        // map.simOptions.boundaryData = boundaryData;
        // map.simOptions.simData = result["original"];
    }

    function computeStatistics(data0, data1) {
        console.log(data1)
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
                'value': data1.category.Excellent,
                'change': data1.category.Excellent - data0.category.Excellent > 0 ? 
                    'up': data1.category.Excellent - data0.category.Excellent == 0 ? 'none': 'down',
            },
            'catGood': {
                'value': data1.category.Good,
                'change': data1.category.Good - data0.category.Good > 0 ? 
                    'up': data1.category.Good - data0.category.Good == 0 ? 'none': 'down',
            },
            'catFair': {
                'value': data1.category.Fair,
                'change': data1.category.Fair - data0.category.Fair > 0 ? 
                    'up': data1.category.Fair - data0.category.Fair == 0 ? 'none': 'down',
            },
            'catPoor': {
                'value': data1.category.Poor,
                'change': data1.category.Poor - data0.category.Poor > 0 ? 
                    'up': data1.category.Poor - data0.category.Poor == 0 ? 'none': 'down',
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

    onMount(async () => {
        // get boundaries
        boundaries = await client.fetchBoundaries();

        let mainMapTilesActual = L.layerGroup();
        let mainMapTilesSimulation = L.layerGroup();
        // getBoundaries()
        if(browser) {
            maps.setupMap(maps.mainMap);
            // setupMap(addMap);
            // setupMap(editMap);
            // setupMap(deleteMap);
            // maps.setupMap(maps.beforeMap);
            // maps.setupMap(maps.afterMap);
            // maps.setupMap(maps.deltaMap);
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
                </div>
                <!-- external boundary source panel -->
                <div class="container row space-between" style="width:calc(100%); align-items: center; height:40px;">
                    <div style="width:calc(50% - 5px)">
                        <ComboBox
                            size="sm"
                            placeholder="Select from registered boundary"
                            items={boundaries}
                            on:select={async (e) => {
                                let boundaryId = e.detail.selectedItem.text;
                                let params = {
                                    date: '20220520',
                                    region: '01',
                                    boundaryId: boundaryId,
                                }
                                await client.performSimulation(params, onSimulationCompleted);
                                drawOptions.mapObj = maps.mainMap;
                                maps.drawSimulationCategory(drawOptions);
                                drawOptions.simData = simulationResult["original"];
                                computeStatistics(simulationResult["original"], drawOptions.simData);
                            }}
                            on:clear={(e) => {
                                maps.clearMap(maps.mainMap);
                            }}
                        />
                    </div>
                    <div> - or - </div>
                    <Button size="sm" kind="tertiary">Load from File</Button>
                </div>
            </div>
            <!-- map panel -->
            <div class="container row start" style="width:100%; height:680px;">
                <div id="main-map" style="width:{maps.mainMap.width}px; height:{maps.mainMap.height}px;"></div>
            </div>
            <!-- legend panel -->
            <div class="container row space-between" style="width:calc(100% - 28px); height:60px; align-items: center; margin: 0 14px;">
                <div>Legend</div>
            </div>
        </div>
    
        <!-- Fine tune panel -->
        <div class="container col start border-right" style="width:calc(33% - 1px); height:100%;">
            <div class="container col" style="width:calc(100%); height:110px; padding: 14px;">
                <div style="width:100%; padding-bottom:14px;">
                    <ComboBox
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
                            maps.drawSimulationCategory(drawOptions);
                            computeStatistics(simulationResult["original"], drawOptions.simData);
                            console.log(simulationStat)
                        }}
                        on:clear={(e) => {
                            drawOptions.dismantledSite = '';
                            drawOptions.simData = simulationResult["original"];
                            maps.drawSimulationCategory(drawOptions);
                            computeStatistics(simulationResult["original"], drawOptions.simData);
                        }}
                    />
                </div>
                <div style="display:flex; flex-flow:row nowrap; width:100%; padding:5px 0;">
                    <Bignumber 
                        field="tiles" 
                        value={simulationStat.count.value} 
                        color="blue" 
                        direction={simulationStat.count.change}
                        width="calc((100% - 20px)/3)"/>
                    <Bignumber 
                        field="rsrp max" 
                        value={simulationStat.max.value} 
                        color="green" 
                        direction={simulationStat.max.change}
                        width="calc((100% - 20px)/3)"/>
                    <Bignumber 
                        field="rsrp min" 
                        value={simulationStat.min.value} 
                        color="red" 
                        direction={simulationStat.min.change}
                        rightMost 
                        width="calc((100% - 20px)/3)"/>
                </div>
                <div style="display:flex; flex-flow:row nowrap; width:100%; padding:5px 0;">
                    <Bignumber 
                        field="Excellent" 
                        value={simulationStat.catExcellent.value} 
                        color="blue" 
                        direction={simulationStat.catExcellent.change}
                        width="calc((100% - 30px)/4)"/>
                    <Bignumber 
                        field="Good" 
                        value={simulationStat.catGood.value} 
                        color="green" 
                        direction={simulationStat.catGood.change}
                        width="calc((100% - 30px)/4)"/>
                    <Bignumber 
                        field="Fair" 
                        value={simulationStat.catFair.value} 
                        color="red" 
                        direction={simulationStat.catFair.change}
                        width="calc((100% - 30px)/4)"/>
                    <Bignumber 
                        field="Poor" 
                        value={simulationStat.catPoor.value} 
                        color="red" 
                        direction={simulationStat.catPoor.change}
                        rightMost 
                        width="calc((100% - 30px)/4)"/>
                </div>
                <div style="display:flex; flex-flow:row nowrap; width:100%; padding:5px 0;">
                    <Bignumber 
                        field="upgraded" 
                        value={simulationStat.deltaUpgraded.value} 
                        color="blue" 
                        direction={simulationStat.deltaUpgraded.change}
                        width="calc((100% - 20px)/3)"/>
                    <Bignumber 
                        field="unchanged" 
                        value={simulationStat.deltaUnchanged.value} 
                        color="green" 
                        direction={simulationStat.deltaUnchanged.change}
                        width="calc((100% - 20px)/3)"/>
                    <Bignumber 
                        field="degraded" 
                        value={simulationStat.deltaDegraded.value} 
                        color="red" 
                        direction={simulationStat.deltaDegraded.change}
                        rightMost
                        width="calc((100% - 20px)/3)"/>
                </div>
                <div style="display:flex; flex-flow:row nowrap; width:100%; padding:5px 0;">
                    <Bignumber 
                        field="safe" 
                        value={simulationStat.statusSafe.value} 
                        color="blue" 
                        direction={simulationStat.statusSafe.change}
                        width="calc((100% - 20px)/3)"/>
                    <Bignumber 
                        field="unsafe" 
                        value={simulationStat.statusUnsafe.value} 
                        color="green" 
                        direction={simulationStat.statusUnsafe.change}
                        width="calc((100% - 20px)/3)"/>
                    <Bignumber 
                        field="fatal" 
                        value={simulationStat.statusFatal.value} 
                        color="red" 
                        direction={simulationStat.statusFatal.change}
                        rightMost
                        width="calc((100% - 20px)/3)"/>
                </div>
            </div>
            
            <!-- cell list panel -->
            <!--div class="container row start border-bottom" style="calc(100% - 60px); height:380px; padding: 0 14px;">
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
                        bind:pageSize={mainUi.table2PageSize}
                        bind:page={mainUi.table2CurrentPage}
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
                        bind:pageSize={mainUi.table2PageSize}
                        bind:page={mainUi.table2CurrentPage}
                        totalItems={$storeServingCells.length}
                        pageSizes={[10,15,20,50]}
                        pageInputDisabled
                    />
                </div>
            </div-->
            <!-- command panel -->
            <!--div class="container row space-between" style="width:calc(100%); height:60px; align-items: center; padding: 0 14px;">
                <Button size="sm" kind="tertiary">Save scenario</Button>
                <div>
                    <Button size="sm" kind="tertiary" on:click={() => {
                            // storeServingCells_selected.set([]);
                            // simulationTiles = drawSimulatedTiles([mainMap, afterMap], resultSet, $storeServingCells_selected)
                            // drawDeltaTiles([mainMap, deltaMap], actualTiles, simulationTiles)
                            // donutSimulation = [];
                            // donutSimDelta = [];
                            // barGroupActualVsSimulation = [];
                        }
                    }>Reset</Button>
                    <Button size="sm" kind="primary" on:click={() => {
                            // simulationTiles = drawSimulatedTiles([mainMap, afterMap], resultSet, $storeServingCells_selected)
                            // donutSimulation = getDonutChartData(simulationTiles, "category", categoryNames)
                            // let deltaTiles = drawDeltaTiles([mainMap, deltaMap], actualTiles, simulationTiles)
                            // donutSimDelta = getDonutChartData(deltaTiles, "status", statusNames)
                            // barGroupActualVsSimulation = getBarGroupChartData([actualTiles, simulationTiles], "category", 
                            //     categoryNames, ["Actual", "Simulation"])
                            // tableTilesDelta = getTableData(deltaTiles, ["tileId", "avgRSRP0", "avgRSRP1", "statusText"],
                            //     (row) => {return row["status"] < 1})
                            // console.log(tableTilesDelta)
                        }
                    }>Redraw</Button>
                </div>
            </div-->
        </div>
    
        <!-- Result panel -->
        <div class="container col start" style="width:calc(34%); height:100%;">
            <div class="container row" style="width:calc(100%); height:282px; align-items: center; padding: 14px;">
                <div style="margin-top:0px; width:220px;">
                    <DonutChart 
                        data={chart.donutSimDelta}
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
                                    "Unchanged": colors.c3,
                                    "Safe": colors.c2,
                                    "Unsafe": colors.c4
                                }
                            },
                            "toolbar":{
                                "enabled": false,
                            },
                            "height": "200px",
                            "data": {
                                "loading": chart.donutSimDelta.length == 0,
                            },
                        }}
                    />
                </div>
                <div style="margin-top:0px; width:400px;">
                    <BarChartGrouped
                    data={chart.barGroupActualVsSimulation}
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
                                    "Actual": colors.c1,
                                    "Simulation": colors.c2,
                                }
                            },
                            "toolbar":{
                                "enabled": false,
                            },
                            "data": {
                                "loading": chart.barGroupActualVsSimulation.length == 0,
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
                        rows={chart.tableTilesDelta}
                        bind:pageSize={mainUi.table3PageSize}
                        bind:page={mainUi.table3CurrentPage}
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
                                    <div style="width:100%; padding:5px; text-align:center; text-transform:uppercase; background-color:{colors.c4}">
                                        {cell.value}
                                    </div>
                                {:else}
                                    <div style="width:100%; padding:5px; text-align:center; text-transform:uppercase; background-color:{colors.c2}">
                                        {cell.value}
                                    </div>
                                {/if}
                            {:else}
                                <div style="width:300px; text-overflow:ellipsis;">{cell.value}</div>
                            {/if}
                        </svelte:fragment>
                    </DataTable>
                    {#if chart.tableTilesDelta.length == 0}
                        <div style="width:100%; text-align:center;">
                            <div style="font-weight:500;">No data</div>
                            <div>check one or more cell on the middle panel,<br>and click redraw</div>
                        </div>
                    {/if}
                    <Pagination
                        bind:pageSize={mainUi.table3PageSize}
                        bind:page={mainUi.table3CurrentPage}
                        totalItems={chart.tableTilesDelta.length}
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