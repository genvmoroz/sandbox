package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	excelize "github.com/xuri/excelize/v2"
)

var fileName = "./testdata/Gas Wisdom_Noon Report (V33).xlsx"

func main() {
	body, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	start := time.Now()
	sheetName, columns, err := readXLSX(body, 0)
	if err != nil {
		panic(err)
	}
	fmt.Printf("with excelize library: %s\n", time.Since(start))
	fmt.Printf("sheet name: %s", sheetName)

	reports, err := asReports(columns)
	if err != nil {
		panic(err)
	}

	content, err := json.MarshalIndent(reports, "", "\t")
	if err != nil {
		panic(err)
	}

	if err = os.WriteFile("out.json", content, os.ModePerm); err != nil {
		panic(err)
	}
}

// readXLSX reads specified sheet of provided xlsx file,
// returns:
// - name of sheet
// - rows which contains any data
// - error if any occurred
func readXLSX(body []byte, sheetIndex uint) (string, [][]string, error) {
	file, err := excelize.OpenReader(bytes.NewReader(body))
	if err != nil {
		return "", nil, fmt.Errorf("failed to read xlsx file body: %w", err)
	}

	switch {
	case file.SheetCount == 0:
		return "", nil, errors.New("xlsx file contains no sheets")
	case sheetIndex >= uint(file.SheetCount):
		return "", nil, fmt.Errorf("no sheet %d available, please select a sheet between 0 and %d", sheetIndex, file.SheetCount-1)
	}

	sheetName := file.GetSheetName(int(sheetIndex))

	cols, err := file.Cols(sheetName)
	if err != nil {
		return "", nil, fmt.Errorf("failed to read rows from sheet %s: %w", sheetName, err)
	}

	var (
		// once iterating through the file will meet 10 continuously blank
		// columns it will stop iterating since it seems to be the end of the file
		maxBlankColumns     = 3
		currentBlankColumns = 0
	)

	columns := make([][]string, 0)

	for {
		values, err := cols.Rows()
		if err != nil {
			return "", nil, fmt.Errorf("failed to get values from column: %w", err)
		}

		if isColumnBlank(values) {
			currentBlankColumns++
		} else { // if column isn't blank - break the sequence and append the column
			currentBlankColumns = 0
			columns = append(columns, values)
		}

		// stop iterating if max blank column limit exceeded
		if !cols.Next() || currentBlankColumns >= maxBlankColumns {
			break
		}
	}

	return sheetName, columns, nil
}

func isColumnBlank(column []string) bool {
	for _, cell := range column {
		if strings.TrimSpace(cell) != "" {
			return false
		}
	}

	return true
}

func asReports(columns [][]string) ([]Report, error) {
	switch {
	case len(columns) < 3:
		return nil, errors.New("should contain more that two columns, first is `Keys`, " +
			"second is `Units`, the rest for the values")
	}

	keys := columns[0]
	units := columns[1]

	values := columns[2:]

	titleIndex := findTitleIndex(keys)
	switch {
	case titleIndex == -1:
		return nil, errors.New("unexpected structure of the file, unable to find title index")
	case titleIndex >= len(keys)-1:
		return nil, errors.New("no keys under the title")
	}

	reports := make([]Report, 0, len(values))

	for _, cursor := range values {
		report, err := asReport(keys, units, cursor, titleIndex)
		if err != nil {
			return nil, err
		}

		reports = append(reports, report)
	}

	return reports, nil
}

func asReport(keys, units, values []string, titleIndex int) (Report, error) {
	report := Report{}

	date := strings.TrimSpace(valueFromByIndex(values, titleIndex))
	if len(date) != 0 {
		report.Date = date
	} else {
		report.Error += fmt.Sprintf("[ReportDate is blank]") // sample of error message
	}

	for i := titleIndex + 1; i < len(keys); i++ {
		var (
			key   = keys[i]
			unit  = valueFromByIndex(units, i)
			value = valueFromByIndex(values, i)
		)

		field := Field{
			Raw:   fmt.Sprintf("%s %s %s", key, unit, value),
			Key:   strings.TrimSpace(strings.TrimRight(strings.TrimLeft(key, "1234567890)"), ":")),
			Value: strings.TrimSpace(value),
			Units: strings.TrimSpace(unit),
		}
		report.Fields = append(report.Fields, field)
	}

	return report, nil
}

func valueFromByIndex(array []string, index int) string {
	if index >= len(array) {
		return "" // no value for the index
	}

	return array[index]
}

func findTitleIndex(keys []string) int {
	for index, key := range keys {
		if strings.EqualFold(strings.TrimSpace(key), "NOON REPORT") {
			return index
		}
	}
	return -1
}
