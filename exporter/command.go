// Copyright (C) 2018-present ichenq@outlook.com. All rights reserved.
// Distributed under the terms and conditions of the Apache License.
// See accompanying files LICENSE.

package exporter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/MakingGame/taxi/descriptor"
)

const (
	MainScriptSeperator = "exporter"
	asciiAlphabet       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

var RegisteredInterpreters = map[string]string{
	".py": "python",
	".js": "node",
}

func RandAlphbetString(length int) string {
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		idx := rand.Int() % len(asciiAlphabet)
		result[i] = asciiAlphabet[idx]
	}
	return string(result)
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

func MakeOneTempFile() string {
	for maxtry := 20; maxtry > 0; maxtry-- {
		var filename = fmt.Sprintf("%s/taxi_%s.json", os.TempDir(), RandAlphbetString(8))
		if f, err := os.Open(filename); err != nil {
			return filename
		} else {
			f.Close()
		}
	}
	return ""
}

func StoreResultToTempFile(result *descriptor.ImportResult) (string, error) {
	var filepath = MakeOneTempFile()
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

func RunExport(dir, params string, result *descriptor.ImportResult) error {
	if dir == "" {
		return fmt.Errorf("export dir is empty, no exported executed")
	}

	filepath, err := StoreResultToTempFile(result)
	if err != nil {
		return err
	}
	defer os.Remove(filepath)

	log.Printf("write descriptor to file %s\n", filepath)

	// marshal request to command line argument string
	var request = &descriptor.ExportRequest{
		Version:  "1.0.1",
		Format:   "json",
		Filepath: filepath,
	}
	data, err := json.Marshal(request)
	if err != nil {
		return err
	}
	var argument = string(data)
	var scripts = EnumerateExporterScripts(dir)
	for _, script := range scripts {
		if err := RunScriptCommand(script, argument, params); err != nil {
			log.Printf("%s %s: %v\n", script, argument, err)
		}
	}

	return nil
}
