// Copyright (C) 2018-present ichenq@outlook.com. All rights reserved.
// Distributed under the terms and conditions of the Apache License.
// See accompanying files LICENSE.

package mydb

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/MakingGame/taxi/descriptor"
)

// column of a mysql table
type MyColumn struct {
	Name    string
	Type    string
	Comment string
	IsNull  bool
	Key     string
}

func (c *MyColumn) ToEnumTypeName() string {
	switch {
	case strings.HasPrefix(c.Type, "tinyint"):
		if strings.HasPrefix(c.Type, "tinyint(1)") {
			return "bool"
		}
		if strings.Index(c.Type, "unsigned") > 0 {
			return "uint8"
		}
		return "int8"

	case strings.HasPrefix(c.Type, "smallint"):
		if strings.Index(c.Type, "unsigned") > 0 {
			return "uint16"
		}
		return "int16"

	case strings.HasPrefix(c.Type, "mediumint"),
		strings.HasPrefix(c.Type, "int"):
		if strings.Index(c.Type, "unsigned") > 0 {
			return "uint32"
		}
		return "int32"

	case strings.HasPrefix(c.Type, "bigint"):
		if strings.Index(c.Type, "unsigned") > 0 {
			return "uint64"
		}
		return "int64"

	case strings.HasPrefix(c.Type, "float"):
		return "float32"

	case strings.HasPrefix(c.Type, "double"),
		strings.HasPrefix(c.Type, "decimal"):
		return "float64"

	case strings.Index(c.Type, "blob") >= 0,
		strings.Index(c.Type, "binary") >= 0:
		return "bytes"

	case strings.HasPrefix(c.Type, "datetime"),
		strings.HasPrefix(c.Type, "timestamp"),
		strings.HasPrefix(c.Type, "date"),
		strings.HasPrefix(c.Type, "time"),
		strings.HasPrefix(c.Type, "year"):
		return "datetime"

	case strings.HasPrefix(c.Type, "char"),
		strings.HasPrefix(c.Type, "varchar"),
		strings.Index(c.Type, "text") >= 0,
		strings.Index(c.Type, "json") >= 0:
		return "string"

	default:
		log.Printf("unsupported mysql type %v\n", c.Type)
		return c.Type
	}
}

func (c *MyColumn) ToDescriptor() *descriptor.FieldDescriptor {
	var name = c.ToEnumTypeName()
	return &descriptor.FieldDescriptor{
		Name:             c.Name,
		CamelCaseName:    descriptor.CamelCase(c.Name),
		TypeName:         name,
		Type:             descriptor.NameToType(name),
		OriginalTypeName: c.Type,
		Comment:          c.Comment,
	}
}

// table of a mysql schema
type MyTable struct {
	Name        string
	Comment     string
	Columns     []*MyColumn
	PrimaryKeys []string
	UniqueKeys  []string
	IndexKeys   []string
}

func (m *MyTable) ToDescriptor() *descriptor.StructDescriptor {
	var desp = &descriptor.StructDescriptor{
		Name:          m.Name,
		CamelCaseName: descriptor.CamelCase(m.Name),
		Comment:       m.Comment,
		Options:       make(map[string]string),
	}
	if len(m.PrimaryKeys) > 0 {
		desp.Options["primary_keys"] = strings.Join(m.PrimaryKeys, ",")
	}
	if len(m.UniqueKeys) > 0 {
		desp.Options["unique_keys"] = strings.Join(m.UniqueKeys, ",")
	}
	if len(m.IndexKeys) > 0 {
		desp.Options["index_keys"] = strings.Join(m.IndexKeys, ",")
	}

	var prevField *descriptor.FieldDescriptor
	for _, col := range m.Columns {
		var field = col.ToDescriptor()
		if prevField != nil && descriptor.IsVectorFields(prevField, field) {
			prevField.IsVector = true
			field.IsVector = true
		}
		prevField = field
		desp.Fields = append(desp.Fields, field)
	}
	return desp
}

// mysql database
type MyDB struct {
	DBName    string
	TableName string
	Tables    []*MyTable
	db        *sql.DB
}

func (m *MyDB) Init(db *sql.DB, dbname, tablename string) {
	m.db = db
	m.DBName = dbname
	m.TableName = tablename
}

func (m *MyDB) Load() error {
	var dbnames = []string{}
	if m.DBName != "" {
		dbnames = append(dbnames, m.DBName)
	} else {
		names, err := m.ShowSchemas()
		if err != nil {
			return err
		}
		dbnames = names
	}
	if len(dbnames) == 0 {
		log.Printf("empty database\n")
		return nil
	}
	for _, dbname := range dbnames {
		log.Printf("load database %s\n", dbname)
		if _, err := m.db.Exec("USE " + dbname); err != nil {
			return err
		}
		var tablenames = []string{}
		if dbname == m.DBName && m.TableName != "" {
			tablenames = append(tablenames, m.TableName)
		} else {
			names, err := m.ShowTables()
			if err != nil {
				return err
			}
			tablenames = names
		}
		for _, name := range tablenames {
			log.Printf("load table %s.%s", dbname, name)
			tbl, err := m.LoadTable(name)
			if err != nil {
				return err
			}
			m.Tables = append(m.Tables, tbl)
		}
	}

	return nil
}

func (m *MyDB) ShowSchemas() ([]string, error) {
	rows, err := m.db.Query("SHOW DATABASES")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var names = []string{}
	for rows.Next() {
		var name string
		if err = rows.Scan(&name); err != nil {
			return nil, err
		}
		names = append(names, name)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return names, nil
}

func (m *MyDB) ShowTables() ([]string, error) {
	rows, err := m.db.Query("SHOW TABLES")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var names = []string{}
	for rows.Next() {
		var name string
		if err = rows.Scan(&name); err != nil {
			return nil, err
		}
		names = append(names, name)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return names, nil
}

func (m *MyDB) LoadTable(tablename string) (*MyTable, error) {
	var name, engine, rowFormat, collation, createOptions, comment string
	var version, rows, avgRowLen, dataLen, maxDataLen, indexLen, dataFree interface{}
	var autoIncr, createTime, updateTime, checkTime, checksum interface{}

	var row = m.db.QueryRow(fmt.Sprintf("SHOW TABLE STATUS LIKE '%s'", tablename))
	if err := row.Scan(&name, &engine, &version, &rowFormat, &rows, &avgRowLen, &dataLen, &maxDataLen, &indexLen, &dataFree,
		&autoIncr, &createTime, &updateTime, &checkTime, &collation, &checksum, &createOptions, &comment); err != nil {
		return nil, err
	}

	var table = &MyTable{
		Name:    name,
		Comment: comment,
	}
	if table.Comment == "" {
		table.Comment = "  " // fix protobuf empty string
	}
	if err := m.LoadColumns(table); err != nil {
		return nil, err
	}
	return table, nil
}

func (m *MyDB) LoadColumns(table *MyTable) error {
	rows, err := m.db.Query("SHOW FULL COLUMNS FROM " + table.Name)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var fieldname, fieldtype, nullable, key, extra, priv, comment string
		var collation, defval interface{}
		if err = rows.Scan(&fieldname, &fieldtype, &collation, &nullable, &key, &defval, &extra, &priv, &comment); err != nil {
			return err
		}
		var column = &MyColumn{
			Name:    fieldname,
			Type:    fieldtype,
			Comment: comment,
			Key:     key,
			IsNull:  nullable == "YES",
		}
		if column.Comment == "" {
			column.Comment = " " // fix protobuf empty string
		}
		table.Columns = append(table.Columns, column)
		if key == "PRI" {
			table.PrimaryKeys = append(table.PrimaryKeys, fieldname)
		} else if key == "UNI" {
			table.UniqueKeys = append(table.UniqueKeys, fieldname)
		} else if key == "MUL" {
			table.IndexKeys = append(table.IndexKeys, fieldname)
		}
	}
	return nil
}
