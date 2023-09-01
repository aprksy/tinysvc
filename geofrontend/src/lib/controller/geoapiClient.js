// import {
//     storeBoundaries,
//     storeBoundaries_selected,
//     storeBoundaryData,
//     storeSitesInBoundary,
//     storeSiteIdsInBoundary,
//     storeSiteDetails,
//     storeSiteCells,
//     storeSimulation,
//     storeSiteNamesInBoundary,
// } from "$lib/store";
// import * as map from "$lib/maps";

import * as urls from '$lib/controller/urls';

// let boundaryData;
// storeBoundaryData.subscribe(value => {
// 	boundaryData = value;
// });

// let simulationData;
// storeSimulation.subscribe(value => {
// 	simulationData = value;
// });

export async function fetchBoundaries() {
    const res = await fetch(urls.allBoundaries, {
        method: 'GET',
        headers: {
            "Content-Type": "application/json",
        }
    })
    let result = await res.json();
    // transform to {"id":XX, "text": "YYY"}
    let fresult = [];
    result.forEach((e, i) => {
        let e1 = {
            "id": i,
            "text": e,
        }
        fresult.push(e1);
    });
    // storeBoundaries.set(fresult);
    return fresult;
}

export async function fetchBoundary(id) {
    const res = await fetch(urls.oneBoundary + id, {
        method: 'GET',
        headers: {
            "Content-Type": "application/json",
        }
    })
    let result = await res.json();
    // storeBoundaryData.set(result);
    return result;
}


export async function fetchSiteIntersects(date, region, boundaryId) {
    let siteIds = [];
    const res = await fetch(urls.siteIntersects, {
        method: 'POST',
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({
            "date": date, 
            "region": region, 
            "boundaryId": boundaryId
        })
    })
    let result = await res.json();
    result.forEach(e => siteIds.push(e.id));
    // storeSitesInBoundary.set(result);
    // storeSiteIdsInBoundary.set(ids);
    // fetchSiteDetails(date, region, ids);
    // fetchSiteCells(date, region, ids);
    // doSimulation(date, region, boundaryId, tileField, groupNames, groupLimits,
    //     minCatIndex, ids);
    return siteIds;
}

export async function fetchSiteDetails(date, region, ids) {
    let sites = ids.join(',');
    const res = await fetch(urls.siteDetails, {
        method: 'POST',
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({
            "date": date, 
            "region": region, 
            "sites": sites,
        })
    })
    let result = await res.json();
    
    // storeSiteNamesInBoundary.set(siteNames);
    // storeSiteDetails.set(result);
    return result;
}

export function getSiteNamesInBoundary(sites) {
    let result = [];
    for (const [key, value] of Object.entries(sites)) {
        result.push({id: value['site'], text: value['name']});
    }
    return result;
}

export async function fetchSiteCells(date, region, ids) {
    let sites = ids.join(',');
    const res = await fetch(urls.siteCells, {
        method: 'POST',
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({
            "date": date, 
            "region": region, 
            "sites": sites,
        })
    })
    let result = await res.json();
    // storeSiteCells.set(result);
    return result
}

export async function doSimulation(date, region, boundaryId, tile, kpi, gNames, gLimits, minCatIdx, sites) {
    const res = await fetch(urls.simulation, {
        method: 'POST',
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({
            "date": date, 
            "region": region, 
            "boundaryId": boundaryId,
            "tile": tile,
            "kpi": kpi,
            "groupNames": gNames,
            "groupLimits": gLimits,
            "sites": sites,
            "minCatIndex": minCatIdx,
        })
    })
    let result = await res.json();
    // storeSimulation.set(result);
    // site 'none' to be dismantled
    // map.simOptions.dismantledSite = "";
    // map.simOptions.allSites = result["sites"];
    // map.simOptions.allTiles = result["tiles"];
    // map.simOptions.boundaryData = boundaryData;
    // map.simOptions.simData = result["original"];

    // map.drawSimulationCategory(map.simOptions);
    return result;
}

export async function performSimulation(params, callback) {
    let boundaryData = await fetchBoundary(params['boundaryId']);
    let siteIds = await fetchSiteIntersects(params['date'], params['region'], params['boundaryId']);

    let p = params
    let siteDetails = await fetchSiteDetails(p['date'], p['region'], siteIds);
    let siteNames = getSiteNamesInBoundary(siteDetails);
    let siteCells = await fetchSiteCells(p['date'], p['region'], siteIds);
    let simulationResult = await doSimulation(p['date'], p['region'], p['boundaryId'], p['tileField'], 
        p['kpi'], p['groupNames'], p['groupLimits'], p['minCatIndex'], siteIds);

    if (callback) {
        callback({
            boundaryData: boundaryData,
            siteIds: siteIds,
            siteDetails: siteDetails,
            siteNames: siteNames,
            simulationResult: simulationResult,
        });
    }
}