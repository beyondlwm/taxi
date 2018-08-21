// Copyright (C) 2018-present ichenq@outlook.com. All rights reserved.
// Distributed under the terms and conditions of the Apache License.
// See accompanying files LICENSE.

package descriptor

import (
	"strconv"
	"time"
	"unicode"
	"unicode/utf8"
)

const timestampLayout = "2006-01-02 15:04:05.999"

func FormatTime(t time.Time) string {
	return t.Format(timestampLayout)
}

// 首字符大写
func Capitalize(name string) string {
	r, i := utf8.DecodeRuneInString(name)
	return string(unicode.ToUpper(r)) + name[i:]
}

//
func FindLastIndexNotNumber(s string) int {
	for i := len(s); i > 1; i-- {
		if s[i-1] >= '0' && s[i-1] <= '9' {
			continue
		}
		return i - 1
	}
	return -1
}

//
func StringVectorIndex(name string) int {
	var i = len(name) - 1
	for i >= 0 && name[i] >= '0' && name[i] <= '9' {
		i -= 1
	}
	i += 1
	if v, e := strconv.Atoi(name[i:]); e == nil {
		return v
	}
	return 0
}

// 找出一群字符串里的最长公共子串
func CommonPrefix(ss ...string) string {
	if len(ss) == 0 {
		return ""
	}
	var prefix = ss[0]
	for _, s := range ss {
		if len(s) < len(prefix) {
			prefix = prefix[:len(s)]
		}
		if prefix == "" {
			return ""
		}
		for i := 0; i < len(prefix); i++ {
			if prefix[i] != s[i] {
				prefix = prefix[:i]
				break
			}
		}
	}
	return prefix
}

//是否是相似的列（归为数组）
func IsVectorFields(prev, cur *FieldDescriptor) bool {
	if prev.OriginalTypeName == cur.OriginalTypeName {
		var prefix = CommonPrefix(prev.Name, cur.Name)
		var err error
		if _, err = strconv.Atoi(prev.Name[len(prefix):]); err != nil {
			return false
		}
		if _, err = strconv.Atoi(cur.Name[len(prefix):]); err != nil {
			return false
		}
		return StringVectorIndex(prev.Name)+1 == StringVectorIndex(cur.Name)
	}
	return false
}
