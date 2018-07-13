// Copyright (C) 2018-present ichenq@outlook.com. All rights reserved.
// Distributed under the terms and conditions of the Apache License.
// See accompanying files LICENSE.

package excel

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/MakingGame/taxi/descriptor"
	"github.com/MakingGame/taxi/importer"
	"github.com/MakingGame/taxi/version"
)

type ExcelImporter struct {
	filename         string
	doc              *excelize.File
	meta             map[string]string
	currentSheetName string
	dataRows         [][]string
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
	var filename = opts["filename"]
	if filename == "" {
		return fmt.Errorf("empty excel filename")
	}
	doc, err := excelize.OpenFile(filename)
	if err != nil {
		return err
	}
	e.filename = filename
	e.doc = doc
	return nil
}

func (e *ExcelImporter) parseMeta() error {
	var rows = e.doc.GetRows(PredefMetaSheet)
	if len(rows) == 0 {
		return fmt.Errorf("no meta sheet found")
	}
	for _, row := range rows {
		if len(row) >= 2 {
			var key = strings.TrimSpace(row[0])
			var value = strings.TrimSpace(row[1])
			e.meta[key] = value // key-value pair
		}
	}
	if e.meta[PredefStructTypeColumn] == "" {
		return fmt.Errorf("struct type column not defined")
	}
	if e.meta[PredefStructNameColumn] == "" {
		return fmt.Errorf("struct name column not defined")
	}
	if e.meta[PredefDataStartColumn] == "" {
		return fmt.Errorf("struct data column not defined")
	}
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
		return nil, err
	}
	if typeColumnIndex >= len(rows) {
		return nil, fmt.Errorf("type column index overflow, %d/%d", typeColumnIndex, len(rows))
	}
	nameColumnIndex, err := strconv.Atoi(e.meta[PredefStructNameColumn])
	if err != nil {
		return nil, err
	}
	if nameColumnIndex >= len(rows) {
		return nil, fmt.Errorf("name column index overflow, %d/%d", nameColumnIndex, len(rows))
	}
	dataStartColumnIndex, err := strconv.Atoi(e.meta[PredefDataStartColumn])
	if err != nil {
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
	println(typeColumnIndex, nameColumnIndex, dataStartColumnIndex, dataEndColumnIndex)
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
		field.Name = namesRow[i]
		field.CamelCaseName = descriptor.CamelCase(field.Name)
		field.TypeName = typeRow[i]
		field.OriginalTypeName = field.TypeName
		field.Type = descriptor.NameToType(typeRow[i])
		if field.Type == descriptor.TypeEnum_Unknown {
			log.Panicf("detected unkown type: %s, %v", typeRow[i], field)
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
		log.Panicf("open file %s failed: %v", filename, err)
	}
	defer f.Close()
	var writer = csv.NewWriter(f)
	for _, row := range e.dataRows {
		if err := writer.Write(row); err != nil {
			log.Panicf("write data to %s failed: %v", filename, err)
		}
	}
	class.Options["datafile"] = filename
	return &class
}

func (e *ExcelImporter) Import() (*descriptor.ImportResult, error) {
	if err := e.parseMeta(); err != nil {
		return nil, err
	}
	var result = &descriptor.ImportResult{
		Version:   version.Version,
		Comment:   "excel",
		Timestamp: descriptor.FormatTime(time.Now()),
		Options:   e.meta,
	}
	if e.meta["sheet"] != "" {
		var sheetName = e.meta["sheet"]
		des, err := e.parseSheet(sheetName)
		if err != nil {
			fmt.Printf("parse sheet %s failed", sheetName)
			return nil, err
		}
		result.Descriptors = append(result.Descriptors, des)
		return result, nil
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
					fmt.Printf("parse sheet %s failed", sheetName)
					return nil, err
				}
				result.Descriptors = append(result.Descriptors, des)
			} else {
				return nil, fmt.Errorf("ExcelImporter: no active sheet")
			}
		} else if mode == "all" {
			for _, sheetName := range sheetMap {
				if sheetName != PredefMetaSheet {
					des, err := e.parseSheet(sheetName)
					if err != nil {
						fmt.Printf("parse sheet %s failed", sheetName)
						return nil, err
					}
					result.Descriptors = append(result.Descriptors, des)
				}
			}
		} else {
			return nil, fmt.Errorf("unsupported parse mode %s", mode)
		}
		return result, nil
	}
	// parse first sheet
	var sheetName = sheetMap[1]
	des, err := e.parseSheet(sheetName)
	if err != nil {
		fmt.Printf("parse sheet %s failed", sheetName)
		return nil, err
	}
	result.Descriptors = append(result.Descriptors, des)
	return result, nil
}

func (e *ExcelImporter) Close() {
	e.doc = nil
}

func init() {
	importer.Register(&ExcelImporter{})
}
