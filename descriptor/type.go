// Copyright (C) 2018-present ichenq@outlook.com. All rights reserved.
// Distributed under the terms and conditions of the Apache License.
// See accompanying files LICENSE.

package descriptor

import (
	"strings"
)

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
	TypeEnum_Array:    "array",
	TypeEnum_Map:      "map",
	TypeEnum_Any:      "any",
}

var AbstractTypeNames = map[string]TypeEnum{
	"array": TypeEnum_Array,
	"map":   TypeEnum_Map,
	"any":   TypeEnum_Any,
}

// enum type to enum name
func TypeToName(typ TypeEnum) string {
	if name, found := enumNames[typ]; found {
		return name
	}
	return ""
}

func IsPrimitiveType(typename string) bool {
	for name, _ := range AbstractTypeNames {
		if name == typename {
			return false
		}
	}
	return true
}

// enum name to enum type
func NameToType(typename string) TypeEnum {
	for name, e := range AbstractTypeNames {
		if strings.Index(typename, name) >= 0 {
			return e
		}
	}
	for k, v := range enumNames {
		if v == typename {
			return k
		}
	}
	return TypeEnum_Unknown
}
