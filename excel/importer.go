// Copyright (C) 2018-present ichenq@outlook.com. All rights reserved.
// Distributed under the terms and conditions of the Apache License.
// See accompanying files LICENSE.

package excel

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/MakingGame/taxi/descriptor"
	"github.com/MakingGame/taxi/importer"
	"github.com/MakingGame/taxi/version"
)

type ExcelImporter struct {
	filelist         []string
	doc              *excelize.File
	meta             map[string]string
	currentSheetName string
	dataRows         [][]string
}

func enumerateExcelFiles(dir string) []string {
	var files []string
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("enumerateExcelFiles: %s,  %v\n", dir, err)
			return err
		}
		if info.IsDir() {
			return nil
		}
		var ext = filepath.Ext(path)
		if ext == ".xlsx" {
			files = append(files, path)
		}
		return nil
	})
	return files
}

func (e *ExcelImporter) Name() string {
	return "excel"
}

func (e *ExcelImporter) Init(args string) error {
	e.meta = map[string]string{}
	opts, err := importer.ParseArgs(args)
	if err != nil {
		return err
	}

	var files []string
	if dir := opts["filedir"]; dir != "" {
		files = enumerateExcelFiles(dir)
	}
	if filepath := opts["filename"]; filepath != "" {
		files = append(files, filepath)
	}
	e.filelist = files
	return nil
}

func (e *ExcelImporter) parseMeta() error {
	var rows = e.doc.GetRows(PredefMetaSheet)
	if len(rows) > 0 {
		for _, row := range rows {
			if len(row) >= 2 {
				var key = strings.TrimSpace(row[0])
				var value = strings.TrimSpace(row[1])
				if key != "" && value != "" {
					e.meta[key] = value // key-value pair
				}
			}
		}
	}

	if e.meta[PredefStructTypeColumn] == "" {
		e.meta[PredefStructTypeColumn] = "1"
	}
	if e.meta[PredefStructNameColumn] == "" {
		e.meta[PredefStructNameColumn] = "2"
	}
	if e.meta[PredefCommentColumn] == "" {
		e.meta[PredefCommentColumn] = "3"
	}
	if e.meta[PredefDataStartColumn] == "" {
		e.meta[PredefDataStartColumn] = "4"
	}
	if len(e.meta["keys"]) == 0 {
		e.meta["keys"] = "1" // default first column is key
	}
	fmt.Printf("parsed sheet meta %v\n", e.meta)
	return nil
}

func (e *ExcelImporter) parseSheet(sheetName string) (*descriptor.StructDescriptor, error) {
	e.dataRows = nil
	e.currentSheetName = sheetName

	var rows = e.doc.GetRows(sheetName)
	if len(rows) == 0 {
		return nil, fmt.Errorf("sheet is empty")
	}

	// validate meta index
	typeColumnIndex, err := strconv.Atoi(e.meta[PredefStructTypeColumn])
	if err != nil {
		fmt.Printf("parseSheet: parse %s failed\n", PredefStructTypeColumn)
		return nil, err
	}
	if typeColumnIndex >= len(rows) {
		return nil, fmt.Errorf("type column index overflow, %d/%d", typeColumnIndex, len(rows))
	}
	nameColumnIndex, err := strconv.Atoi(e.meta[PredefStructNameColumn])
	if err != nil {
		fmt.Printf("parseSheet: parse %s failed\n", PredefStructNameColumn)
		return nil, err
	}
	if nameColumnIndex >= len(rows) {
		return nil, fmt.Errorf("name column index overflow, %d/%d", nameColumnIndex, len(rows))
	}
	dataStartColumnIndex, err := strconv.Atoi(e.meta[PredefDataStartColumn])
	if err != nil {
		fmt.Printf("parseSheet: parse %s failed\n", PredefDataStartColumn)
		return nil, err
	}
	if dataStartColumnIndex >= len(rows) || dataStartColumnIndex <= typeColumnIndex || dataStartColumnIndex <= nameColumnIndex {
		return nil, fmt.Errorf("data start column index overflow, %d/%d", dataStartColumnIndex, len(rows))
	}
	var dataEndColumnIndex = len(rows)
	if e.meta[PredefDataStartColumn] != "" {
		index, err := strconv.Atoi(e.meta[PredefDataStartColumn])
		if err != nil {
			return nil, err
		}
		if index >= len(rows) || index < dataStartColumnIndex {
			return nil, fmt.Errorf("data end column index overflow, %d", dataEndColumnIndex)
		}
		if index > dataStartColumnIndex {
			dataEndColumnIndex = index
		}
	}
	var des = e.parseSheetData(rows, typeColumnIndex, nameColumnIndex, dataStartColumnIndex, dataEndColumnIndex)
	return des, nil
}

func (e *ExcelImporter) parseSheetData(rows [][]string, typeColumnIndex, nameColumnIndex, dataStartColumnIndex, dataEndColumnIndex int) *descriptor.StructDescriptor {
	var class descriptor.StructDescriptor
	class.Comment = e.meta["comment"]
	if class.Comment == "" {
		class.Comment = " "
	}

	// class name
	var className = e.currentSheetName
	if e.meta[PredefClassName] != "" {
		className = e.meta[PredefClassName]
	}
	class.Name = className
	class.CamelCaseName = descriptor.CamelCase(className)

	var commentIndex = -1
	if e.meta[PredefCommentColumn] != "" {
		index, _ := strconv.Atoi(e.meta[PredefCommentColumn])
		if index > 0 {
			commentIndex = index - 1
		}
	}

	var typeRow = rows[typeColumnIndex-1]
	var namesRow = rows[nameColumnIndex-1]
	for i := 0; i < len(typeRow); i++ {
		if typeRow[i] == "" || namesRow[i] == "" { // skip empty
			continue
		}
		var field descriptor.FieldDescriptor
		field.Name = strings.TrimSpace(namesRow[i])
		field.CamelCaseName = descriptor.CamelCase(field.Name)
		field.TypeName = strings.TrimSpace(typeRow[i])
		field.OriginalTypeName = field.TypeName
		field.Type = descriptor.NameToType(typeRow[i])
		if field.Type == descriptor.TypeEnum_Unknown {
			log.Panicf("parseSheetData:detected unkown type: %s, %v\n", typeRow[i], field)
		}
		if commentIndex > 0 {
			field.Comment = rows[commentIndex][i]
		}
		if field.Comment == "" {
			field.Comment = " "
		}
		class.Fields = append(class.Fields, &field)
	}
	e.dataRows = rows[dataStartColumnIndex-1 : dataEndColumnIndex-1]
	var filename = descriptor.MakeOneTempFile("taxi", ".csv")
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		log.Panicf("parseSheetData: open file %s failed: %v\n", filename, err)
	}
	defer f.Close()
	var writer = csv.NewWriter(f)
	for _, row := range e.dataRows {
		if err := writer.Write(row); err != nil {
			log.Panicf("parseSheetData: write data to %s failed: %v", filename, err)
		}
	}
	class.Options = e.meta
	class.Options["datafile"] = filename
	fmt.Printf("write csv data file to %s\n", filename)
	return &class
}

func (e *ExcelImporter) imporeOneFile(result *descriptor.ImportResult) error {
	if err := e.parseMeta(); err != nil {
		return err
	}
	if e.meta["sheet"] != "" {
		var sheetName = e.meta["sheet"]
		des, err := e.parseSheet(sheetName)
		if err != nil {
			fmt.Printf("imporeOneFile: parse sheet %s failed\n", sheetName)
			return err
		}
		result.Descriptors = append(result.Descriptors, des)
		return nil
	}
	sheetMap := e.doc.GetSheetMap()
	if e.meta["parse-mode"] != "" {
		var mode = e.meta["parse-mode"]
		if mode == "active-only" {
			var sheetIndex = e.doc.GetActiveSheetIndex()
			if sheetIndex > 0 {
				var sheetName = sheetMap[sheetIndex]
				des, err := e.parseSheet(sheetName)
				if err != nil {
					fmt.Printf("imporeOneFile: parse sheet %s failed\n", sheetName)
					return err
				}
				result.Descriptors = append(result.Descriptors, des)
			} else {
				return fmt.Errorf("ExcelImporter: no active sheet")
			}
		} else if mode == "all" {
			for _, sheetName := range sheetMap {
				if sheetName != PredefMetaSheet {
					des, err := e.parseSheet(sheetName)
					if err != nil {
						fmt.Printf("imporeOneFile: parse sheet %s failed\n", sheetName)
						return err
					}
					result.Descriptors = append(result.Descriptors, des)
				}
			}
		} else {
			return fmt.Errorf("unsupported parse mode %s", mode)
		}
		return nil
	}
	// parse first sheet
	var sheetName = sheetMap[1]
	des, err := e.parseSheet(sheetName)
	if err != nil {
		fmt.Printf("imporeOneFile: parse sheet %s failed\n", sheetName)
		return err
	}
	result.Descriptors = append(result.Descriptors, des)
	return nil
}

func (e *ExcelImporter) Import() (*descriptor.ImportResult, error) {
	var result = &descriptor.ImportResult{
		Version:   version.Version,
		Comment:   "excel",
		Timestamp: descriptor.FormatTime(time.Now()),
	}
	if len(e.filelist) == 0 {
		return nil, fmt.Errorf("no excel file specified")
	}
	for _, filename := range e.filelist {
		fmt.Printf("start parse file %s\n", filename)
		doc, err := excelize.OpenFile(filename)
		if err != nil {
			fmt.Printf("ExcelImporter: OpenFile, %s\n", filename)
			return nil, err
		}
		e.doc = doc
		if err := e.imporeOneFile(result); err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (e *ExcelImporter) Close() {
	e.doc = nil
}

func init() {
	importer.Register(&ExcelImporter{})
}
