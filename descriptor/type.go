// Copyright (C) 2018-present ichenq@outlook.com. All rights reserved.
// Distributed under the terms and conditions of the Apache License.
// See accompanying files LICENSE.

package descriptor

import (
	"unicode"
	"unicode/utf8"
)

// 首字符大写
func Capitalize(name string) string {
	r, i := utf8.DecodeRuneInString(name)
	return string(unicode.ToUpper(r)) + name[i:]
}

var enumNames = map[TypeEnum]string{
	TypeEnum_Nil:      "nil",
	TypeEnum_Bool:     "bool",
	TypeEnum_Int8:     "int8",
	TypeEnum_Uint8:    "uint8",
	TypeEnum_Int16:    "int16",
	TypeEnum_Uint16:   "uint16",
	TypeEnum_Int32:    "int32",
	TypeEnum_Uint32:   "uint32",
	TypeEnum_Int64:    "int64",
	TypeEnum_Uint64:   "uint64",
	TypeEnum_Float32:  "float32",
	TypeEnum_Float64:  "float64",
	TypeEnum_String:   "string",
	TypeEnum_Bytes:    "bytes",
	TypeEnum_DateTime: "datetime",
}

var nameEnums = map[string]TypeEnum{
	"nil":      TypeEnum_Nil,
	"bool":     TypeEnum_Bool,
	"int8":     TypeEnum_Int8,
	"uint8":    TypeEnum_Uint8,
	"int16":    TypeEnum_Int16,
	"uint16":   TypeEnum_Uint16,
	"int32":    TypeEnum_Int32,
	"uint32":   TypeEnum_Uint32,
	"int64":    TypeEnum_Int64,
	"uint64":   TypeEnum_Uint64,
	"float32":  TypeEnum_Float32,
	"float64":  TypeEnum_Float64,
	"string":   TypeEnum_String,
	"bytes":    TypeEnum_Bytes,
	"datetime": TypeEnum_DateTime,
}

// enum type to enum name
func TypeToName(typ TypeEnum) string {
	if name, found := enumNames[typ]; found {
		return name
	}
	return ""
}

// enum name to enum type
func NameToType(typename string) TypeEnum {
	if e, found := nameEnums[typename]; found {
		return e
	}
	return TypeEnum_Unknown
}
