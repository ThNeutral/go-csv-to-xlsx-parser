package main

import (
	"encoding/csv"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"

	"github.com/tealeg/xlsx"
)

const (
	inputFolder   = "./input"
	outputFolder  = "./output"
	csvExtension  = ".csv"
	xslxExtension = ".xlsx"
)

func main() {
	var paths []string
	filepath.WalkDir(inputFolder, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(d.Name()) == csvExtension {
			paths = append(paths, path)
		}
		return nil
	})
	var wg sync.WaitGroup
	for pathIndex, path := range paths {
		wg.Add(1)
		go func(pathIndex int, path string) {
			defer wg.Done()

			inputFile, err := os.Open(path)
			if err != nil {
				fmt.Printf("Error: %v", err)
				return
			}
			defer inputFile.Close()

			csvReader := csv.NewReader(inputFile)
			records, err := csvReader.ReadAll()
			if err != nil {
				fmt.Printf("Error: %v", err)
				return
			}

			outputFile := xlsx.NewFile()
			sheet, err := outputFile.AddSheet("Sheet1")
			if err != nil {
				fmt.Printf("Error: %v", err)
				return
			}

			for _, rowIterable := range records {
				row := sheet.AddRow()
				for _, cellValue := range rowIterable {
					cell := row.AddCell()
					cell.Value = cellValue
				}
			}

			err = outputFile.Save(outputFolder + "/" + fmt.Sprint(pathIndex) + xslxExtension)
			if err != nil {
				fmt.Printf("Error: %v", err)
				return
			}
		}(pathIndex, path)
	}
	wg.Wait()
}
