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

	"github.com/MakingGame/taxi/descriptor"
	"github.com/MakingGame/taxi/version"
)

var RegisteredInterpreters = map[string]string{
	".py":  "python",
	".js":  "node",
	".lua": "lua",
	".cs":  "donet run",
	".go":  "go run",
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
		scripts = append(scripts, path)
		return nil
	})
	return scripts
}

func StoreResultToTempFile(result *descriptor.ImportResult) (string, error) {
	var filepath = descriptor.MakeOneTempFile("taxi_meta", ".json")
	if filepath == "" {
		return "", fmt.Errorf("StoreResultToTempFile: cannot create temporary file")
	}

	// write json text to file
	f, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return "", err
	}
	defer f.Close()
	data, err := json.MarshalIndent(result, "", "  ")
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
		log.Fatalf("RunScriptCommand: invalid interpreter of %s", script)
	}
	var output bytes.Buffer
	fmt.Printf("run %s %s", interpreter, script)
	var cmd = exec.Command(interpreter, script, argument, params)
	cmd.Stdout = &output
	cmd.Stderr = &output
	if err := cmd.Run(); err != nil {
		fmt.Printf("run %s %s: %v\n%s", interpreter, script, err, output.String())
	} else {
		fmt.Printf("run %s %s:\n%s\n", interpreter, script, output.String())
	}

	return nil
}

func RunExport(filepath, dir, params string, result *descriptor.ImportResult) error {
	if dir == "" && filepath == "" {
		return fmt.Errorf("RunExport: exporter path is empty, no exporter executed")
	}

	filename, err := StoreResultToTempFile(result)
	if err != nil {
		return err
	}

	fmt.Printf("write descriptor to file %s\n", filename)

	// marshal request to command line argument string
	var request = &descriptor.ExportRequest{
		Version:  version.Version,
		Format:   "json",
		Filepath: filename,
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
	if filepath != "" {
		scripts = append(scripts, filepath)
	}
	if len(scripts) == 0 {
		return fmt.Errorf("RunExport: no export script found")
	}
	for _, script := range scripts {
		if err := RunScriptCommand(script, argument, params); err != nil {
			fmt.Printf("%s %s: %v\n", script, argument, err)
		}
	}

	return nil
}
