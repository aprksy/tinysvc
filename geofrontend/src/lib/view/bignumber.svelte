<script>
    export let field='fieldName';
    export let value=0;
    export let color='#888';
    export let rightMost=false;
    export let width='200px';
    export let direction='none';
    export let arrowColors={
        'up': 'green',
        'down': 'red',
    };
    let rightMargin = rightMost ? 0: 10;
</script>

<div class="container" style="width:{width};margin-right:{rightMargin}px; overflow-y:auto;">
    <div class="head">
        <div class="field">{field}</div>
        <div 
            class="triangle-up" 
            style="border-color: 
                    transparent 
                    transparent 
                    {arrowColors.up} 
                    transparent; 
                display:{direction == 'up' ? 'block': 'none'}"
            />
        <div 
            class="triangle-down" 
            style="border-color: 
                    {arrowColors.down} 
                    transparent 
                    transparent 
                    transparent; 
                display:{direction == 'down' ? 'block': 'none'}"/>
    </div>
    <p class="value" style="color:{color};">{value}</p>
</div>

<style>
    .container {
        display: flex;
        flex-flow: column nowrap;
        padding: 5px;
        border: 1px solid #eee;
        background-color: #fff;
    }
    .head {
        display: flex;
        flex-flow: row nowrap;
        justify-content: space-between;
        width: 100%;
    }
    .field {
        color: #aaa;
        font-size: 14px;
        text-transform: uppercase;
        align-self: stretch;
    }
    .value {
        width: 100%;
        font-size: 40px;
        text-align: right;
        align-self: flex-end;
    }
    .triangle-up {
        width: 0;
        height: 0;
        border-style: solid;
        border-width: 0 12px 15px 12px;
    }
    .triangle-down {
        width: 0;
        height: 0;
        border-style: solid;
        border-width: 15px 12px 0 12px;
    }
</style>