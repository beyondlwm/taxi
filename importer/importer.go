// Copyright (C) 2018-present ichenq@outlook.com. All rights reserved.
// Distributed under the terms and conditions of the Apache License.
// See accompanying files LICENSE.

package importer

import (
	"log"
	"sync"

	"github.com/MakingGame/taxi/descriptor"
)

type Importer interface {
	Init(string) error
	Name() string
	Import() (*descriptor.ImportResult, error)
	Close()
}

var (
	mut       sync.Mutex
	importers = make(map[string]Importer)
)

func Register(importer Importer) {
	mut.Lock()
	defer mut.Unlock()
	if importer == nil {
		panic("Register: nil importer registration")
	}
	var name = importer.Name()
	if _, dup := importers[name]; dup {
		log.Panicf("Register: duplicat importer[%v] registration", name)
	}
	importers[name] = importer
}

func ImporterByName(name string) Importer {
	mut.Lock()
	v := importers[name]
	mut.Unlock()
	return v
}
