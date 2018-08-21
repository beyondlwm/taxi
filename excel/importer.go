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

	"github.com/MakingGame/taxi/descriptor"
	"github.com/MakingGame/taxi/importer"
	"github.com/MakingGame/taxi/version"
	"github.com/tealeg/xlsx"
)

type ExcelImporter struct {
	filelist     []string
	doc          *xlsx.File
	meta         map[string]string
	currentSheet *xlsx.Sheet
	dataRows     [][]string
}

//文件名中有以下字符将被忽略
var ignoreExcelPattern = []string{
	"~$",
	"-TNP-",
	" - 副本",
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

		for _, name := range ignoreExcelPattern {
			if strings.Index(info.Name(), name) >= 0 {
				return nil
			}
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
	var skipNames []string
	if len(opts["skip"]) > 0 {
		skipNames = strings.Split(opts["skip"], ",")
	}
	for _, filename := range files {
		var found = false
		for _, v := range skipNames {
			if strings.Index(filename, v) >= 0 {
				found = true
				break
			}
		}
		if !found {
			e.filelist = append(e.filelist, filename)
		}
	}
	return nil
}

func (e *ExcelImporter) getSheetRows(sheet *xlsx.Sheet) [][]string {
	var textRows [][]string
	for _, row := range sheet.Rows {
		var textRow = make([]string, 0, len(row.Cells))
		for _, cell := range row.Cells {
			var text = strings.TrimSpace(cell.String())
			textRow = append(textRow, text)
		}
		textRows = append(textRows, textRow)
	}
	return textRows
}

func (e *ExcelImporter) parseMeta() error {
	e.meta = map[string]string{}
	var sheet = e.doc.Sheet[PredefMetaSheet]
	if sheet != nil {
		var rows = e.getSheetRows(sheet)
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

func (e *ExcelImporter) parseSheet(sheet *xlsx.Sheet) (*descriptor.StructDescriptor, error) {
	e.dataRows = nil
	e.currentSheet = sheet

	var rows = e.getSheetRows(sheet)
	if len(rows) == 0 {
		return nil, fmt.Errorf("sheet %v is empty", sheet.Name)
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
	var className = e.currentSheet.Name
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
	var nameDict = map[string]bool{}

	var prevField *descriptor.FieldDescriptor
	for i := 0; i < len(typeRow); i++ {
		if typeRow[i] == "" || namesRow[i] == "" { // skip empty
			continue
		}

		var field descriptor.FieldDescriptor
		field.Name = strings.TrimSpace(namesRow[i])
		field.CamelCaseName = descriptor.CamelCase(field.Name)
		field.OriginalTypeName = strings.TrimSpace(typeRow[i])
		field.Type = descriptor.NameToType(field.OriginalTypeName)
		field.TypeName = descriptor.TypeToName(field.Type)
		field.ColumnIndex = uint32(i + 1)

		if _, found := nameDict[field.Name]; found {
			log.Panicf("duplicate name defined, %s", field.Name)
		}

		if prevField != nil && descriptor.IsVectorFields(prevField, &field) {
			prevField.IsVector = true
			field.IsVector = true
		}
		prevField = &field

		if field.Type == descriptor.TypeEnum_Unknown || field.TypeName == "" {
			log.Panicf("parseSheetData: detected unkown type: %s, %v\n", field.OriginalTypeName, field)
		}
		if commentIndex > 0 {
			field.Comment = rows[commentIndex][i]
		}
		if field.Comment == "" {
			field.Comment = " "
		}
		nameDict[field.Name] = true
		class.Fields = append(class.Fields, &field)
	}
	var datarows = rows[dataStartColumnIndex-1 : dataEndColumnIndex]
	e.validateSheetRows(&class, datarows)

	fmt.Printf("total %d rows\n", len(e.dataRows))
	class.Options = e.meta
	e.writeCsvData(&class)
	return &class
}

// skip empty type/name column
func (e *ExcelImporter) validateSheetRows(class *descriptor.StructDescriptor, rows [][]string) {
	// pad empty row
	var maxRowLen = len(class.Fields)
	for _, row := range rows {
		if len(row) > maxRowLen {
			maxRowLen = len(row)
		}
	}
	for i, row := range rows {
		for j := len(row); j < maxRowLen; j++ {
			row = append(row, "")
		}
		rows[i] = row
	}

	var datarows = [][]string{}
	for _, row := range rows {
		var datarow = make([]string, 0, len(class.Fields))
		for _, field := range class.Fields {
			var text = row[field.ColumnIndex-1]
			datarow = append(datarow, text)
		}
		datarows = append(datarows, datarow)
	}

	e.dataRows = datarows
}

//写入数据到csv文件
func (e *ExcelImporter) writeCsvData(class *descriptor.StructDescriptor) {
	var filename = descriptor.MakeOneTempFile("taxi_"+class.Name, ".csv")
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		log.Panicf("parseSheetData: open file %s failed: %v\n", filename, err)
	}
	defer f.Close()
	var writer = csv.NewWriter(f)
	if err := writer.WriteAll(e.dataRows); err != nil {
		log.Panicf("parseSheetData: WriteAll, write csv data to %s failed: %v", filename, err)
	}
	if err := writer.Error(); err != nil {
		log.Panicf("parseSheetData: write csv data to %s failed: %v", filename, err)
	}
	class.Options["datafile"] = filename
	fmt.Printf("write csv data file to %s\n", filename)
}

func (e *ExcelImporter) imporeOneFile(result *descriptor.ImportResult) error {
	if err := e.parseMeta(); err != nil {
		return err
	}
	// parse first sheet
	for _, sheet := range e.doc.Sheets {
		des, err := e.parseSheet(sheet)
		if err != nil {
			fmt.Printf("imporeOneFile: parse sheet %s failed\n", sheet.Name)
			return err
		}
		fmt.Printf("parsed %s options: %v", sheet.Name, des.Options)
		result.Descriptors = append(result.Descriptors, des)
		break
	}
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
		doc, err := xlsx.OpenFile(filename)
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
	e.meta = nil
	e.currentSheet = nil
	e.filelist = nil
	e.dataRows = nil
}

func init() {
	importer.Register(&ExcelImporter{})
}
