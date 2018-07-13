// Copyright (C) 2018-present ichenq@outlook.com. All rights reserved.
// Distributed under the terms and conditions of the Apache License.
// See accompanying files LICENSE.

package mydb

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/MakingGame/taxi/descriptor"
	"github.com/MakingGame/taxi/importer"
	"github.com/MakingGame/taxi/version"
	_ "github.com/go-sql-driver/mysql"
)

type MySQLImporter struct {
	db   MyDB
	conn *sql.DB
}

func (m *MySQLImporter) Name() string {
	return "mysql"
}

func (m *MySQLImporter) Init(args string) error {
	opts, err := importer.ParseArgs(args)
	if err != nil {
		return err
	}

	var dsn = m.makeDSN(opts)
	log.Printf("DSN: %s\n", dsn)
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	if err := conn.Ping(); err != nil {
		return err
	}
	m.conn = conn
	m.db.Init(conn, opts["db"], opts["table"])
	return nil
}

func (m *MySQLImporter) Close() {
	m.conn.Close()
}

// DSN format: user:pwd@tcp(host:port)/db?charset=utf8&interpolateParams=true&parseTime=true&loc=Local
func (m *MySQLImporter) makeDSN(kv map[string]string) string {
	// make a copy
	var opts = make(map[string]string, len(kv))
	for k, v := range kv {
		opts[k] = v
	}
	// data source name
	var user = opts["user"]
	var passwd = opts["passwd"]
	var host = opts["host"]
	if host == "" {
		host = "localhost"
	}
	var protocol = opts["protocol"]
	if protocol == "" {
		protocol = "tcp"
	}
	var port = opts["port"]
	if port == "" {
		port = "3306"
	}
	var db = opts["db"]

	// delete no-parameters
	delete(opts, "user")
	delete(opts, "passwd")
	delete(opts, "protocol")
	delete(opts, "host")
	delete(opts, "port")
	delete(opts, "db")

	// set default values
	var params = map[string]string{}
	var unsetDefaultValues = map[string]string{
		"charset":           "utf8",
		"interpolateParams": "true",
		"parseTime":         "true",
		"loc":               "Local",
	}
	for k, v := range unsetDefaultValues {
		if _, found := opts[k]; !found {
			params[k] = v
		}
	}

	// accepted DSN param
	var acceptedParams = map[string]string{
		"parameters":       "Parameters",
		"address":          "Address",
		"maxAllowedPacket": "maxAllowedPacket",
		"multiStatements":  "multiStatements",
		"readTimeout":      "readTimeout",
		"writeTimeout":     "writeTimeout",
		"timeout":          "timeout",
	}
	for k, v := range acceptedParams {
		if s, found := opts[k]; found {
			params[v] = s
		}
	}

	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%s:%s@%s(%s:%s)/%s?", user, passwd, protocol, host, port, db)
	var atEnd = 0
	for k, v := range params {
		atEnd++
		if atEnd < len(params) {
			fmt.Fprintf(&buf, "%s=%s&", k, v)
		} else {
			fmt.Fprintf(&buf, "%s=%s", k, v)
		}
	}
	return buf.String()
}

func (m *MySQLImporter) Import() (*descriptor.ImportResult, error) {
	if err := m.db.Load(); err != nil {
		return nil, err
	}
	var descriptors []*descriptor.StructDescriptor
	for _, tbl := range m.db.Tables {
		descriptors = append(descriptors, tbl.ToDescriptor())
	}
	var result = &descriptor.ImportResult{
		Version:     version.Version,
		Comment:     "mysql",
		Timestamp:   descriptor.FormatTime(time.Now()),
		Descriptors: descriptors,
	}
	return result, nil
}

func init() {
	importer.Register(&MySQLImporter{})
}
