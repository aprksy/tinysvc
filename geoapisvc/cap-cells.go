package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func getCellsSite(date time.Time, region string, cells ...string) ([]string, error) {
	worstCaseResult := []string{}
	retval, err := rs.Values(fmt.Sprintf(cell_site, date.Format("20060102"), region), cells)
	if err != nil {
		return worstCaseResult, err
	}
	result := worstCaseResult
	for _, e := range retval {
		result = append(result, string(e))
	}
	return result, nil
}

func getCellsData(date time.Time, region string, cells ...string) ([]map[string]interface{}, error) {
	worstCaseResult := []map[string]interface{}{}
	retval, err := rs.Values(fmt.Sprintf(cell_data, date.Format("20060102"), region), cells)
	if err != nil {
		return worstCaseResult, err
	}
	intValues := []int{2, 3, 4, 5, 7}
	result := worstCaseResult
	fields := strings.Split(cell_data_fields, ",")
	for j, e := range retval {
		values := strings.Split(string(e), ",")
		data := map[string]interface{}{"id": j + 1, "cell": cells[j]}
		for i := range fields {
			if inIntSlice(i, intValues) {
				data[fields[i]], _ = strconv.ParseFloat(values[i], 64)
			} else {
				data[fields[i]] = values[i]
			}
		}
		result = append(result, data)
	}
	return result, nil
}
