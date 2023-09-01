<script>
    export let layout='';
    export let items = [
        {color:'red', label:'Item 1'},
        {color:'green', label:'Item 2'},
        {color:'blue', label:'Item 3'},
    ];
</script>

<div class="container" style="flex-flow:{layout=='vertical'?'column':'row'} nowrap;">
    {#each items as item}
        <div class="item">
            <div class="color-box" style="background-color:{item.color}"></div>
            <div class="label">{item.label}</div>
        </div>
    {/each}
</div>

<style>
    .container {
        display: flex;
        flex-flow: row nowrap;
        padding: 5px;
        background-color: #fff;
    }
    .item {
        display: flex;
        flex-flow: row nowrap;
        justify-content: flex-start;
        margin-right: 20px;
    }
    .color-box {
        color: black;
        width: 15px;
        height: 15px;
        margin-right: 12px;
        border-radius: 50%;
    }
    .label {
        height: 20px;
        font-size: 14px;
        text-transform: uppercase;
    }
</style>