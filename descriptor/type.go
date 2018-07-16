// Copyright (C) 2018-present ichenq@outlook.com. All rights reserved.
// Distributed under the terms and conditions of the Apache License.
// See accompanying files LICENSE.

package descriptor

import (
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

// 首字符大写
func Capitalize(name string) string {
	r, i := utf8.DecodeRuneInString(name)
	return string(unicode.ToUpper(r)) + name[i:]
}

const timestampLayout = "2006-01-02 15:04:05.999"

func FormatTime(t time.Time) string {
	return t.Format(timestampLayout)
}

var enumNames = map[TypeEnum]string{
	TypeEnum_Nil:      "nil",
	TypeEnum_Bool:     "bool",
	TypeEnum_Int8:     "int8",
	TypeEnum_Uint8:    "uint8",
	TypeEnum_Int16:    "int16",
	TypeEnum_Uint16:   "uint16",
	TypeEnum_Int:      "int",
	TypeEnum_Int32:    "int32",
	TypeEnum_Uint32:   "uint32",
	TypeEnum_Int64:    "int64",
	TypeEnum_Uint64:   "uint64",
	TypeEnum_Float:    "float",
	TypeEnum_Float32:  "float32",
	TypeEnum_Float64:  "float64",
	TypeEnum_String:   "string",
	TypeEnum_Enum:     "enum",
	TypeEnum_Bytes:    "bytes",
	TypeEnum_DateTime: "datetime",
	TypeEnum_Json:     "json",
	TypeEnum_Array:    "array",
	TypeEnum_Any:      "any",
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
	if strings.Index(typename, "array") >= 0 {
		return TypeEnum_Array
	}
	for k, v := range enumNames {
		if v == typename {
			return k
		}
	}
	return TypeEnum_Unknown
}
