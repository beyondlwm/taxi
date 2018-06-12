// Copyright (C) 2018-present ichenq@outlook.com. All rights reserved.
// Distributed under the terms and conditions of the Apache License.
// See accompanying files LICENSE.

package mydb

import (
	"fmt"
	"testing"
	"database/sql"


)

var (
	host   = "localhost"
	port   = 3306
	user   = "root"
	passwd = "holyshit"
)

func TestMySQLDB_Load(t *testing.T) {
	var format = "%s:%s@tcp(%s:%d)/%s?charset=utf8&interpolateParams=true&parseTime=true&loc=Local"
	var dsn string = fmt.Sprintf(format, user, passwd, host, port, "")
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("open mysql: %v", err)
	}
	defer conn.Close()
	var db MyDB
	db.Init(conn, "test", "")
	if err := db.Load(); err != nil {
		t.Fatalf("%v", err)
	}
	fmt.Printf("%v\n", db.Tables)
}
