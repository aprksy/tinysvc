import { groupLimits } from "$lib/model/defaults";

export let c1 = "#3944c5"
export let c2 = "#2cb485"
export let c3 = "#ff7700"
export let c4 = "#c80064"

let colors = [c4, c3, c2, c1];

export function byCategory(cat) {
    switch (cat) {
        case "POOR": return c4;
        case "FAIR": return c3;
        case "GOOD": return c2;
        case "EXCELLENT": return c1;
        case "UNDEFINED": return 'black';
    }
}

export function byDelta(delta) {
    switch (delta) {
        case "UPGRADED": return c2;
        case "UNCHANGE": return c3;
        case "DEGRADED": return c4;
    }
}

export function byStatus(status) {
    switch (status) {
        case "SAFE": return c2;
        case "UNSAFE": return c4;
        case "FATAL": return 'black';
    }
}

export function fromValue(value) {
    let intervals = [];
    for (let i=0; i<groupLimits.length-1; i++) {
        intervals.push([groupLimits[i], groupLimits[i+1]]);
    }

    for (let i=0; i<intervals.length; i++) {
        if (intervals[i][0]<value && value<=intervals[i][1]) {
            return colors[i];
        }
    }
}

export function legends(kind) {
    let result = [];
    switch (kind) {
        case 'original':
            result.push({color:byCategory('POOR'), label: 'POOR'});
            result.push({color:byCategory('FAIR'), label: 'FAIR'});
            result.push({color:byCategory('GOOD'), label: 'GOOD'});
            result.push({color:byCategory('EXCELLENT'), label: 'EXCELLENT'});
            result.push({color:byCategory('UNDEFINED'), label: 'UNDEFINED'});
            break;
        case 'simulation':
            result.push({color:byCategory('POOR'), label: 'POOR'});
            result.push({color:byCategory('FAIR'), label: 'FAIR'});
            result.push({color:byCategory('GOOD'), label: 'GOOD'});
            result.push({color:byCategory('EXCELLENT'), label: 'EXCELLENT'});
            result.push({color:byCategory('UNDEFINED'), label: 'UNDEFINED'});
            break;
        case 'final-status':
            result.push({color:byStatus('SAFE'), label: 'SAFE'});
            result.push({color:byStatus('UNSAFE'), label: 'UNSAFE'});
            result.push({color:byStatus('FATAL'), label: 'FATAL'});
            break;
    
        default:
            break;
    }
    return result;
}