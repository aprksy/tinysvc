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
    } from "carbon-components-svelte";
    import { Upload, View, TreeView, Area, AreaCustom, WatsonHealthCircleMeasurement, ChartNetwork, Reset, CheckboxChecked, OrderDetails } from "carbon-icons-svelte";
    import MapBoundary from "carbon-icons-svelte/lib/MapBoundary.svelte";
    import SettingsAdjust from "carbon-icons-svelte/lib/SettingsAdjust.svelte";
    import UserAvatarFilledAlt from "carbon-icons-svelte/lib/UserAvatarFilledAlt.svelte";
    import LeafletMap from "$lib/LeafletMap.svelte";
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

    function getBoundaries() {
        storeBoundaries.set([
            { id: "0", text: "R03-Rectangle-TebetBarat-01", tag: "4G", boundary:{"type":"FeatureCollection","features":[{"type":"Feature","properties":{},"geometry":{"type":"Polygon","coordinates":[[[106.83955192565918,-6.238983289054349],[106.86058044433592,-6.238983289054349],[106.86058044433592,-6.221449256187597],[106.83955192565918,-6.221449256187597],[106.83955192565918,-6.238983289054349]]]}}]}},
            { id: "1", text: "R03-Polygon-Halim-01", tag: "4G", boundary:{"type":"FeatureCollection","features":[{"type":"Feature","properties":{},"geometry":{"type":"Polygon","coordinates":[[[106.87491416931152,-6.259332270466314],[106.88023567199707,-6.26411011355174],[106.87946319580078,-6.269143865206112],[106.88401222229004,-6.270594259210957],[106.88796043395996,-6.264536704558167],[106.88718795776366,-6.2587350370072965],[106.88246726989746,-6.25566353986335],[106.87808990478516,-6.256175457307679],[106.87491416931152,-6.259332270466314]]]}}]}},
            { id: "2", text: "R03-Circle-Cibubur-01", tag: "4G", boundary: {} },
        ]);
    }

    let moduleName = "Usecase 07 - Geospatial based KPI Analysis"

    let layers = [
        // {"type":"FeatureCollection","features":[{"type":"Feature","properties":{},"geometry":{"type":"Polygon","coordinates":[[[106.87491416931152,-6.259332270466314],[106.88023567199707,-6.26411011355174],[106.87946319580078,-6.269143865206112],[106.88401222229004,-6.270594259210957],[106.88796043395996,-6.264536704558167],[106.88718795776366,-6.2587350370072965],[106.88246726989746,-6.25566353986335],[106.87808990478516,-6.256175457307679],[106.87491416931152,-6.259332270466314]]]}}]},
    ];

    let mainMapLayers;
    let editMapLayers;

</script>
  
<Header company="Telkomsel" platformName={moduleName} bind:isSideNavOpen>
    <svelte:fragment slot="skip-to-content">
        <SkipToContent />
    </svelte:fragment>
    <HeaderUtilities>
        <HeaderAction
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
        </HeaderAction>
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
                            mainMapLayers.redraw(layers);
                        }}
                        on:clear={e => mainMapLayers.clear()}
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
                <!-- <LeafletMap id="main-map" width={664} height={786} bind:layers={mainMapLayers} drawControls={true} /> -->
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
    on:close
    on:submit
>
    <div class="container row space-between" style="height:600px;">
        <div style="width:50%; height:100%">
            <div class="container col start">
                <FileUploaderDropContainer
                    labelText="Drag and drop files here or click to upload"
                    accept={[".json", ".zip", ".txt", "csv"]}
                    validateFiles={(files) => {
                        return files.filter((file) => true);
                    }}
                    on:change={(e) => {
                        uploadStatus.fileDate = e.detail[0].lastModified;
                        uploadStatus.fileSize = e.detail[0].size;
                        uploadStatus.fileName = e.detail[0].name;
                        uploadStatus.status = uploadStatus.fileSize > 1024 * 1024 ? "failed": "success";
                        console.log(uploadStatus);
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
        <div style="width:50%; height:100%; background-color:#ddd;"></div>
    </div>
</Modal>

<Modal
    size="lg"
    bind:open={MapBoundaryEditOpen}
    modalHeading="Edit boundary"
    primaryButtonText="Confirm Edit"
    secondaryButtonText="Cancel"
    on:click:button--secondary={() => (MapBoundaryEditOpen = false)}
    on:open
    on:close
    on:submit
>
    <div class="container row space-between" style="height:600px;">
        <div class="container col start" style="width:calc(50% - 20px); height:100%">
            <DataTable
                sortable
                title="Available boundaries"
                description="Boundaries available to use"
                headers={[
                    { key: "text", value: "Name" },
                    { key: "tag", value: "Tag" },
                ]}
                rows={$storeBoundaries}
                on:click:row={(e) => {
                    console.log(e.detail)
                    layers = [
                        e.detail.boundary,
                    ]
                    // editMapLayers.redraw(layers);
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
                <LeafletMap id="edit-map" width={670} height={550} bind:layers={editMapLayers} drawControls={true} />
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
  <p>This is a permanent action and cannot be undone.</p>
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