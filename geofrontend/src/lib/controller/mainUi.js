let table1PageSize = 10;
let table1CurrentPage = 1;
export let table2PageSize = 10;
export let table2CurrentPage = 1;
export let table3PageSize = 10;
export let table3CurrentPage = 1;

export function shouldFilterItem(item, value) {
    if (!value) return true;
    return item.text.toLowerCase().includes(value.toLowerCase());
}
