// Copyright (C) 2018-present ichenq@outlook.com. All rights reserved.
// Distributed under the terms and conditions of the Apache License.
// See accompanying files LICENSE.

package exporter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/MakingGame/taxi/descriptor"
	"github.com/MakingGame/taxi/version"
)

const (
	MainScriptSeperator = "exporter"
)

var RegisteredInterpreters = map[string]string{
	".py": "python",
	".js": "node",
}

func EnumerateExporterScripts(dir string) []string {
	var scripts []string
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("walk dir[%s]: %v\n", dir, err)
			return err
		}
		if info.IsDir() {
			return nil
		}
		var ext = filepath.Ext(path)

		interpreter := RegisteredInterpreters[ext]
		if interpreter == "" {
			return nil
		}
		var base = filepath.Base(path)
		if strings.Index(base, MainScriptSeperator) >= 0 {
			scripts = append(scripts, path)
		}
		return nil
	})
	return scripts
}

func StoreResultToTempFile(result *descriptor.ImportResult) (string, error) {
	var filepath = descriptor.MakeOneTempFile("taxi", ".json")
	if filepath == "" {
		return "", fmt.Errorf("cannot create temporary file")
	}

	// write json text to file
	f, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return "", err
	}
	defer f.Close()
	data, err := json.Marshal(result)
	if err != nil {
		return "", err
	}
	f.Write(data)
	f.Sync()
	return filepath, nil
}

func RunScriptCommand(script, argument, params string) error {
	var interpreter = RegisteredInterpreters[filepath.Ext(script)]
	if interpreter == "" {
		log.Fatalf("invalid interpreter of %s", script)
	}
	var output bytes.Buffer
	var cmd = exec.Command(interpreter, script, argument, params)
	cmd.Stdout = &output
	cmd.Stderr = &output
	if err := cmd.Run(); err != nil {
		log.Printf("run %s %s: %v\n%s", interpreter, script, err, output.String())
	} else {
		log.Printf("run %s %s:\n%s\n", interpreter, script, output.String())
	}

	return nil
}

func RunExport(path, dir, params string, result *descriptor.ImportResult) error {
	if dir == "" && path == "" {
		return fmt.Errorf("exporter path is empty, no exporter executed")
	}

	filepath, err := StoreResultToTempFile(result)
	if err != nil {
		return err
	}
	defer os.Remove(filepath)

	log.Printf("write descriptor to file %s\n", filepath)

	// marshal request to command line argument string
	var request = &descriptor.ExportRequest{
		Version:  version.Version,
		Format:   "json",
		Filepath: filepath,
	}
	data, err := json.Marshal(request)
	if err != nil {
		return err
	}
	var argument = string(data)
	var scripts []string
	if dir != "" {
		scripts = EnumerateExporterScripts(dir)
	}
	if path != "" {
		scripts = append(scripts, path)
	}
	for _, script := range scripts {
		if err := RunScriptCommand(script, argument, params); err != nil {
			log.Printf("%s %s: %v\n", script, argument, err)
		}
	}

	return nil
}
