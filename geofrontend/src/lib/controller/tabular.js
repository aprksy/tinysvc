
export function generateSitesTable(data) {
    let result = [];
    let i = 1;
    for (const [siteid, value] of Object.entries(data.sites)) {
        result.push({
            no: i,
            id: value.id,
            name: value.name,
            lat: value.lat,
            lng: value.lng,
            type: value.type,
        });
        i++;
    }
    return result;
}

export function generateTilesOriTable(data) {
    let result = [];
    let i = 1;
    for (const [tileid, value] of Object.entries(data.original.tiles)) {
        result.push({
            no: i,
            id: tileid,
            lat: data.tiles[tileid].lat,
            lng: data.tiles[tileid].lng,
            value: value.value1,
            category: value.category1,
        });
        i++;
    }
    return result;
}

export function generateSimulationTables(data, siteid) {
    let result = [];
    let i = 1;
    for (const [tileid, value] of Object.entries(data.original.tiles)) {
        let tile = data.simulation[siteid].tiles[tileid];
        result.push({
            no: i,
            id: tileid,
            // lat: data.tiles[tileid].lat,
            // lng: data.tiles[tileid].lng,
            value0: value.value1,
            category0: value.category0,
            value1: tile?tile.value1:'n/a',
            category1: tile?tile.category1:'UNDEFINED',
            delta: tile?tile.delta:'DEGRADED',
            status: tile?tile.status:'FATAL',
        });
        i++;
    }
    return result;
}