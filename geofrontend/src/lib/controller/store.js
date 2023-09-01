import { writable } from 'svelte/store';

export const storeRegionals = writable([
    { id: "01", text: "REGIONAL - 01" },
    { id: "02", text: "REGIONAL - 02" },
    { id: "03", text: "REGIONAL - 03" },
    { id: "04", text: "REGIONAL - 04" },
    { id: "05", text: "REGIONAL - 05" },
    { id: "06", text: "REGIONAL - 06" },
    { id: "07", text: "REGIONAL - 07" },
    { id: "08", text: "REGIONAL - 08" },
    { id: "09", text: "REGIONAL - 09" },
    { id: "10", text: "REGIONAL - 10" },
    { id: "11", text: "REGIONAL - 11" },
    { id: "12", text: "REGIONAL - 12" },
]);
export const storeRegionalsSelected = writable({});

export const storeSites = writable([]);
export const storeSitesSelected = writable({});

export const storeNearbyCenter = writable({});
export const storeNearbySites = writable([]);
export const storeNearbySitesLimit = writable(100);
export const storeNearbySitesRadius = writable(2000.0);

export const storeBoundaries = writable([]);
export const storeBoundaries_selected = writable([]);
export const storeBoundaryData = writable({});

export const storeSitesInBoundary = writable([]);
export const storeSiteIdsInBoundary = writable([]);
export const storeSiteNamesInBoundary = writable([]);
export const storeSiteDetails = writable([]);
export const storeSiteCells = writable([]);
export const storeSimulation = writable({});

export const storeServingCells = writable([]);
export const storeServingCells_selected = writable([]);