package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
	"os"
	"strings"
	"time"
)

var fileName = "./testdata/Gas Wisdom_Noon Report (V33).xlsx"

func main() {
	body, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	start := time.Now()
	sheetName, rows, err := readXLSX(body, 0)
	if err != nil {
		panic(err)
	}
	fmt.Printf("with excelize library: %s\n", time.Since(start))
	fmt.Printf("sheet name: %s", sheetName)

	_ = rows
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
		maxBlankColumns     = 10
		currentBlankColumns = 0
	)

	for cols.Next() && currentBlankColumns < maxBlankColumns { // stop iterating if max blank column limit eccided

		values, err := cols.Rows()
		if err != nil {
			return "", nil, fmt.Errorf("failed to get values from column: %w", err)
		}

		if isColumnBlank(values) {
			currentBlankColumns++
			continue
		} else { // if column isn't blank - break the sequence
			currentBlankColumns = 0
		}
	}

	return sheetName, excludeEmptyRows(rows), nil
}

func isColumnBlank(column []string) bool {
	for _, cell := range column {
		if strings.TrimSpace(cell) != "" {
			return false
		}
	}

	return true
}

func excludeEmptyRows(sourceRows [][]string) (rows [][]string) {
	rows = make([][]string, 0, len(sourceRows))
	for _, sRow := range sourceRows {
		if len(sRow) != 0 {
			rows = append(rows, sRow)
		}
	}

	return
}

func asRecords(rows [][]string) ([]Record, error) {
	switch {
	case len(rows) < 2:
		return nil, errors.New("")
	}
}
