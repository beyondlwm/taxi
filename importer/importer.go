// Copyright (C) 2018-present ichenq@outlook.com. All rights reserved.
// Distributed under the terms and conditions of the Apache License.
// See accompanying files LICENSE.

package importer

import (
	"fmt"
	"log"
	"strings"
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

func ParseArgs(args string) (map[string]string, error) {
	if args == "" {
		return nil, fmt.Errorf("ParseArgs: empty arguments")
	}
	var kvlist = strings.Split(args, ",")
	var opts = make(map[string]string)
	for _, item := range kvlist {
		kv := strings.Split(item, "=")
		if len(kv) != 2 {
			return nil, fmt.Errorf("ParseArgs: invalid key-value pair arguments: %s", item)
		}
		opts[kv[0]] = kv[1]
	}
	return opts, nil
}
