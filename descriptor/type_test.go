// Copyright (C) 2018-present ichenq@outlook.com. All rights reserved.
// Distributed under the terms and conditions of the Apache License.
// See accompanying files LICENSE.

package descriptor

import "testing"

func TestCapitalize(t *testing.T) {
	var s = "abcdefg"
	var b = Capitalize(s)
	if b != "Abcdefg" {
		t.Fatalf("invalid capitalize %v", b)
	}
}

func TestNameToType(t *testing.T) {
	var typenames = []string{
		"nil",
		"bool",
		"int8",
		"uint8",
		"int16",
		"uint16",
		"int32",
		"uint32",
		"int64",
		"uint64",
		"float32",
		"float64",
		"string",
		"bytes",
		"datetime",
	}
	for _, name := range typenames {
		NameToType(name)
	}
}

func TestTypeToName(t *testing.T) {
	var types = []TypeEnum{
		TypeEnum_Nil,
		TypeEnum_Bool,
		TypeEnum_Int8,
		TypeEnum_Uint8,
		TypeEnum_Int16,
		TypeEnum_Uint16,
		TypeEnum_Int32,
		TypeEnum_Uint32,
		TypeEnum_Int64,
		TypeEnum_Uint64,
		TypeEnum_Float32,
		TypeEnum_Float64,
		TypeEnum_String,
		TypeEnum_Bytes,
		TypeEnum_DateTime,
	}
	for _, typ := range types {
		TypeToName(typ)
	}
}
