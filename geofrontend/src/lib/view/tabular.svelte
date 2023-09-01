<script>
    import { DataTable, Pagination } from "carbon-components-svelte";

    export let title='Table Title';
    export let description='Table description';
    export let header=[
        {key:"id", value:"ID"},
        {key:"name", value:"Name"},
        {key:"desc", value:"Description"},
    ];
    export let data=[
        {id:1, name:"Items 1", desc: "Item 1 description"},
        {id:2, name:"Items 2", desc: "Item 2 description"},
        {id:3, name:"Items 3", desc: "Item 3 description"},
    ];
    export let pagesize=10;
    export let visible=false
    let pageno=1;
</script>

<div style="width:100%; display:{visible?'block':'none'}">
    <DataTable
        title={title}
        description={description}
        zebra
        sortable        
        headers={header}
        pageSize={pagesize}
        page={pageno}
        rows={data}
    />
    <Pagination
        bind:pageSize={pagesize}
        bind:page={pageno}
        totalItems={data?data.length:0}
        pageSizeInputDisabled
    />
</div>

<style>
</style>