<script>
    import {
        Button,
        Header,
        HeaderUtilities,
        HeaderAction,
        HeaderGlobalAction,
        HeaderPanelLinks,
        HeaderPanelDivider,
        HeaderPanelLink,
        SideNav,
        SideNavItems,
        SideNavMenu,
        SideNavMenuItem,
        SideNavLink,
        SkipToContent,
        Content,
        Grid,
        Row,
        Column,
        Tile,
    } from "carbon-components-svelte";
    import { ArrowRight } from "carbon-icons-svelte";
    import SettingsAdjust from "carbon-icons-svelte/lib/SettingsAdjust.svelte";
    import UserAvatarFilledAlt from "carbon-icons-svelte/lib/UserAvatarFilledAlt.svelte";
  
    let isSideNavOpen = false;
    let isOpen1 = false;
    let isOpen2 = false;
</script>
  
<Header company="Telkomsel" platformName="Mediation Server - Geospatial Usecases" bind:isSideNavOpen>
    <svelte:fragment slot="skip-to-content">
        <SkipToContent />
    </svelte:fragment>
    <HeaderUtilities>
        <!-- <HeaderGlobalAction aria-label="Settings" icon={SettingsAdjust} /> -->
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
        <!-- <HeaderAction bind:isOpen={isOpen2}>
        <HeaderPanelLinks>
            <HeaderPanelDivider>Switcher subject 1</HeaderPanelDivider>
            <HeaderPanelLink>Switcher item 1</HeaderPanelLink>
            <HeaderPanelDivider>Switcher subject 2</HeaderPanelDivider>
            <HeaderPanelLink>Switcher item 1</HeaderPanelLink>
            <HeaderPanelLink>Switcher item 2</HeaderPanelLink>
            <HeaderPanelLink>Switcher item 3</HeaderPanelLink>
            <HeaderPanelLink>Switcher item 4</HeaderPanelLink>
            <HeaderPanelLink>Switcher item 5</HeaderPanelLink>
        </HeaderPanelLinks>
        </HeaderAction> -->
    </HeaderUtilities>
</Header>
  
<!-- <SideNav bind:isOpen={isSideNavOpen}>
    <SideNavItems>
        <SideNavLink text="Link 1" />
        <SideNavLink text="Link 2" />
        <SideNavLink text="Link 3" />
        <SideNavMenu text="Menu">
        <SideNavMenuItem href="/" text="Link 1" />
        <SideNavMenuItem href="/" text="Link 2" />
        <SideNavMenuItem href="/" text="Link 3" />
        </SideNavMenu>
    </SideNavItems>
</SideNav> -->
  
<Content>
    <div class="container col self-center">
        <div style="padding:30px; text-align:center; font-size:26px; 
        padding-bottom:30px; border-bottom:1px solid #ddd; margin-bottom:30px;">Select module to get started</div>
        <div class="container row wrap start modules-container self-center">
            <div class="module margin30">
                <div class="box-title">Usecase 07 - Geospatial based KPI Analysis</div>
                <div class="box-subtitle">Description</div>
                <div class="box-content">Lorem ipsum dolor sit amet, consectetur adipiscing elit, 
                    sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, 
                    quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. </div>
                <div class="box-subtitle">Analysis tools</div>
                <div class="box-content">Maps, Table, Charts</div>
                <div class="box-subtitle">Data sources</div>
                <div class="box-content">Cell events, sitedb, all-regions</div>
                <div class="box-control container row end">
                    <Button size="field" icon={ArrowRight} on:click={()=>{window.open("/uc07-1");}}>OPEN</Button>
                </div>
            </div>
            <div class="module margin30">
                <div class="box-title">Usecase 09 - Custom tile size KPI Aggregation</div>
                <div class="box-subtitle">Description</div>
                <div class="box-content">Lorem ipsum dolor sit amet, consectetur adipiscing elit, 
                    sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, 
                    quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. </div>
                <div class="box-subtitle">Analysis tools</div>
                <div class="box-content">Maps, Table, Charts</div>
                <div class="box-subtitle">Data sources</div>
                <div class="box-content">Cell events, sitedb, all-regions</div>
                <div class="box-control container row end">
                    <Button size="field" icon={ArrowRight} >OPEN</Button>
                </div>
            </div>
            <div class="module">
                <div class="box-title">Usecase 10 - Blacksite Simulation</div>
                <div class="box-subtitle">Description</div>
                <div class="box-content">Lorem ipsum dolor sit amet, consectetur adipiscing elit, 
                    sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, 
                    quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. </div>
                <div class="box-subtitle">Analysis tools</div>
                <div class="box-content">Maps, Table, Charts</div>
                <div class="box-subtitle">Data sources</div>
                <div class="box-content">Cell events, sitedb, all-regions</div>
                <div class="box-control container row end">
                    <Button size="field" icon={ArrowRight} on:click={()=>{window.open("/uc10_3");}}>OPEN</Button>
                </div>
            </div>
        </div>
    </div>
</Content>

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

.self-center {
    align-self: center;
}

.modules-container {
    width: 80%;
}

.module {
    width: calc(33% - 20px);
    height: 400px;
    padding: 20px;
    background-color: #f8f8f8;
}

.box-title {
    padding: 20px 0;
    font-size: 18px;
    font-weight: 500;
}

.box-subtitle {
    padding-top: 10px;
    margin-bottom: 6px;
    margin-top: 10px;
    font-weight: 500;
    border-top: 1px solid #ddd;
}

.box-content {
    line-height: 20px;
}

.box-control {
    padding: 20px 0;
}

.margin20 {
    margin-right: 20px;
}

.margin30 {
    margin-right: 30px;
}
</style>