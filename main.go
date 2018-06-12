// Copyright (C) 2018-present ichenq@outlook.com. All rights reserved.
// Distributed under the terms and conditions of the Apache License.
// See accompanying files LICENSE.

package main

import (
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"time"

	_ "github.com/MakingGame/taxi/excel"
	"github.com/MakingGame/taxi/exporter"
	"github.com/MakingGame/taxi/importer"
	_ "github.com/MakingGame/taxi/mydb"
	"github.com/jessevdk/go-flags"
)

type Options struct {
	Version     bool   `short:"v" long:"version" description:"show version string"`
	Mode        string `short:"M" long:"mode" description:"mode of importer source"`
	ImportArgs  string `long:"import-args" description:"arguments of importer"`
	ExportArgs  string `long:"export-args" description:"arguments of exporter"`
	ExporterDir string `long:"export-dir" description:"exporter template directory"`
}

func NewOptions() *Options {
	return &Options{}
}

var Version = "0.1.1"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	var opts = NewOptions()
	if _, err := flags.Parse(opts); err != nil {
		if e, ok := err.(*flags.Error); ok {
			if e.Type == flags.ErrHelp {
				return
			}
		}
		log.Fatalf("flags.Parse: %v", err)
	}
	if opts.Version {
		fmt.Printf("ver %s, (built w/%s, OS/Arch: %s/%s)", Version, runtime.Version(),
			runtime.GOOS, runtime.GOARCH)
		return
	}

	worker := importer.ImporterByName(opts.Mode)
	if worker == nil {
		fmt.Printf("unrecognized imporeter mode [%v], run taxi --help to show how to use", opts.Mode)
		return
	}
	if err := worker.Init(opts.ImportArgs); err != nil {
		fmt.Printf("Init failed: %v\n", err)
		return
	}
	defer worker.Close()

	descriptors, err := worker.Import()
	if err != nil {
		fmt.Printf("Import failed: %v\n", err)
		return
	}

	if err := exporter.RunExport(opts.ExporterDir, opts.ExportArgs, descriptors); err != nil {
		fmt.Printf("Run export failed: %v\n", err)
		return
	}
	fmt.Println("program exit")
}
