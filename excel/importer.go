// Copyright (C) 2018-present ichenq@outlook.com. All rights reserved.
// Distributed under the terms and conditions of the Apache License.
// See accompanying files LICENSE.

package excel

import (
	"github.com/MakingGame/taxi/descriptor"
	"github.com/MakingGame/taxi/importer"
)

type ExcelImporter struct {
}

func (e *ExcelImporter) Name() string {
	return "excel"
}

func (e *ExcelImporter) Init(filepath string) error {
	return nil
}

func (e *ExcelImporter) Import() (*descriptor.ImportResult, error) {
	return nil, nil
}

func (e *ExcelImporter) Close() {

}

func init() {
	importer.Register(&ExcelImporter{})
}
