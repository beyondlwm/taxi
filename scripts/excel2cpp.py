# Copyright (C) 2018-present ichenq@outlook.com. All rights reserved.
# Distributed under the terms and conditions of the Apache License.
# See accompanying files LICENSE.

from __future__ import print_function
from __future__ import unicode_literals
from __future__ import division

import codecs
import json
import os
import sys
import csv
import time
import exportutil


# 类型映射
def map_cpp_type(typ):
    type_mapping = {
        'bool':     'bool',
        'int8':     'int8_t',
        'uint8':    'uint8_t',
        'int16':    'int16_t',
        'uint16':   'uint16_t',
        'int':      'int',
        'int32':    'int32_t',
        'uint32':   'uint32_t',
        'int64':    'int64_t',
        'uint64':   'uint64_t',
        'float':    'float',
        'float32':  'float',
        'float64':  'double',
        'enum':     'enum',
        'string':   'std::string',
    }
    if typ.startswith('array'):
        t = exportutil.array_element_type(typ)
        elem_type = type_mapping[t]
        return 'std::vector<%s>' % elem_type
    else:
        return type_mapping[typ]


def is_pod_type(typ):
    assert len(typ.strip()) > 0
    return typ.strip() != 'std::string' # only string is non-pod


def instance_data_name(name):
    return '_instance_%s_data' % name.lower()
 

# 为类型加上默认值
def name_with_default_value(field, typename):
    typename = typename.strip()
    line = ''
    if typename == 'bool':
        line = '%s = false;' % field['name']
    elif exportutil.is_integer_type(field['type_name']):
        line = '%s = 0;' % field['name']
    elif exportutil.is_floating_type(field['type_name']):
        line = '%s = 0.0;' % field['name']
    else:
        line = '%s;' % field['name']
    assert len(line) > 0
    return line

    
# 生成结构体定义
def gen_cpp_struct(struct):
    content = '// %s\n' % struct['comment']
    content += 'struct %s \n{\n' % struct['name']
    fields = struct['fields']
    if struct['options']['parse-kv-mode']:
        fields = exportutil.construct_struct_kv_fields(struct)
        
    max_name_len = exportutil.max_field_length(fields, 'name', None)
    max_type_len = exportutil.max_field_length(fields, 'original_type_name', map_cpp_type)
    for field in fields:
        typename = map_cpp_type(field['original_type_name'])
        assert typename != "", field['original_type_name']
        typename = exportutil.pad_spaces(typename, max_type_len+1)
        name = name_with_default_value(field, typename)
        name = exportutil.pad_spaces(name, max_name_len + 5)
        content += '    %s %s // %s\n' % (typename, name, field['comment'])
    
    return content

    
def get_struct_keys(struct, keyname):
    column_keys = struct['options'][keyname].split(',')
    assert len(column_keys) > 0, struct['name']
    key_tuples = []
    fields = struct['fields']
    for column in column_keys:
        idx, field = exportutil.get_field_by_column_index(struct, int(column))
        typename = map_cpp_type(field['original_type_name'])
        name = field['name']
        key_tuples.append((typename, name))
    return key_tuples
    
    
# 生成方法声明    
def gen_static_method_declare(struct):
    content = ''
    content += '    static int Load();\n'
    if struct['options']['parse-kv-mode']:
        content += '    static const %s* Instance();\n' % struct['name']
        return content
        
    get_keys = get_struct_keys(struct, 'get-keys')
    get_args = []
    for tpl in get_keys:
        typename = tpl[0]
        if not is_pod_type(typename):
            typename = 'const %s&' % typename
        get_args.append(typename + ' ' + tpl[1])      
    #content += '    static void UnMarshal(const std::vector<StringPiece>& row, std::vector<%s>* dataptr);\n' % struct['name']
    content += '    static const std::vector<%s>* GetData(); \n' % struct['name']
    content += '    static const %s* Get(%s);\n' % (struct['name'], ', '.join(get_args))
    
    if 'range-keys'  in struct['options']:
        range_keys = get_struct_keys(struct, 'range-keys')
        range_args = []
        for tpl in range_keys:
            typename = tpl[0]
            if not is_pod_type(typename):
                typename = 'const %s&' % typename
            range_args.append(typename + ' ' + tpl[1])  
        content += '    static std::vector<const %s*> GetRange(%s);\n' % (struct['name'], ', '.join(range_args))

    return content

    
def gen_equal_stmt(prefix, struct, key):
    keys = get_struct_keys(struct, key)
    args = []
    for tpl in keys:
        args.append('%s%s == %s' % (prefix, tpl[1], tpl[1]))
    return ' && '.join(args)


# 定义静态变量
def gen_static_data_define(struct):
    content = ''
    varname = instance_data_name(struct['name'])
    if struct['options']['parse-kv-mode']:
        content += '    static %s* %s = nullptr;\n' % (struct['name'], varname)
    else:
        content += '    static std::vector<%s>* %s = nullptr;\n' % (struct['name'], varname)
    return content
    

#
def gen_static_data_method(struct):
    content = ''
    varname = instance_data_name(struct['name'])
    if struct['options']['parse-kv-mode']:
        content += 'const %s* %s::Instance()\n' % (struct['name'], struct['name'])
        content += '{\n'
        content += '    BEATS_ASSERT(%s != nullptr);\n' % varname
        content += '    return %s;\n' % varname
        content += '}\n\n'
    else:
        content += 'const std::vector<%s>* %s::GetData()\n' % (struct['name'], struct['name'])
        content += '{\n'
        content += '    BEATS_ASSERT(%s != nullptr);\n' % varname
        content += '    return %s;\n' % varname
        content += '}\n\n'
    return content
   
    
# 生成Get方法    
def gen_static_get_method(struct):
    content = ''
    if struct['options']['parse-kv-mode']:
        return content
        
    keys = get_struct_keys(struct, 'get-keys')
    args = []
    for tpl in keys:
        typename = tpl[0]
        if not is_pod_type(typename):
            typename = 'const %s&' % typename        
        args.append('%s %s' %  (typename, tpl[1]))
        
    content += 'const %s* %s::Get(%s)\n' % (struct['name'], struct['name'], ', '.join(args))
    content += '{\n'
    content += '    const vector<%s>* dataptr = GetData();\n' % struct['name']
    content += '    BEATS_ASSERT(dataptr != nullptr && dataptr->size() > 0);\n'
    content += '    for (size_t i = 0; i < dataptr->size(); i++)\n'
    content += '    {\n'
    content += '        if (%s)\n' % gen_equal_stmt('dataptr->at(i).', struct, 'get-keys')
    content += '        {\n'
    content += '            return &dataptr->at(i);\n'
    content += '        }\n'
    content += '    }\n'
    content += '    BEATS_ASSERT(false, "%s.Get: no item found");\n' % struct['name']
    content += '    return nullptr;\n'
    content += '}\n\n'
    return content


# 生成GetRange方法
def gen_static_range_method(struct):
    content = ''
    if struct['options']['parse-kv-mode']:
        return content    

    if 'range-keys' not in struct['options']: 
        return content 

    keys = get_struct_keys(struct, 'range-keys')
    args = []
    for tpl in keys:
        typename = tpl[0]
        if not is_pod_type(typename):
            typename = 'const %s&' % typename        
        args.append('%s %s' %  (typename, tpl[1]))
        
    content += 'std::vector<const %s*> %s::GetRange(%s)\n' % (struct['name'], struct['name'], ', '.join(args))
    content += '{\n'
    content += '    const vector<%s>* dataptr = GetData();\n' % struct['name']
    content += '    std::vector<const %s*> range;\n' % struct['name']
    content += '    BEATS_ASSERT(dataptr != nullptr && dataptr->size() > 0);\n'
    content += '    for (size_t i = 0; i < dataptr->size(); i++)\n'
    content += '    {\n'
    content += '        if (%s)\n' % gen_equal_stmt('dataptr->at(i).', struct, 'range-keys')
    content += '        {\n'
    content += '            range.push_back(&dataptr->at(i));\n'
    content += '        }\n'
    content += '        else \n'
    content += '        {\n'
    content += '            if (!range.empty()) \n'
    content += '                break;\n'
    content += '        }\n'
    content += '    }\n'
    content += '    BEATS_ASSERT(!range.empty(), "%s.GetRange: no item found");\n' % struct['name']
    content += '    return range;\n'
    content += '}\n\n'
    return content    
    
# 生成Load方法    
def gen_static_load_method_kv_mode(struct):
    content = ''
    rows = struct['options']['datarows']
    keycol = struct['options']['key-column'] 
    valcol = struct['options']['value-column'] 
    typcol = int(struct['options']['value-type-column']) 
    assert keycol > 0 and valcol > 0 and typcol > 0
    
    keyidx, keyfield = exportutil.get_field_by_column_index(struct, keycol)
    validx, valfield = exportutil.get_field_by_column_index(struct, valcol)
    typeidx, typefield = exportutil.get_field_by_column_index(struct, typcol)

    content = 'int %s::Load()\n' % struct['name']
    content += '{\n'
    content += '    const char* csvpath = "/csv/%s.csv";\n' % struct['name'].lower()
    content += '    %s* dataptr = BEATS_NEW(%s, "autogen", csvpath);\n' % (struct['name'], struct['name'])
    content += '    const vector<vector<string>>& rows = CResourceManager::GetInstance()->ReadCsvToRows(csvpath);\n'
    content += '    BEATS_ASSERT(rows.size() >= %d && rows[0].size() >= %d);\n' % (len(rows), validx)
    idx = 0
    for row in rows:
        name = rows[idx][keyidx].strip()
        typename = rows[idx][typeidx].strip()
        typename = map_cpp_type(typename)
        content += '    if (!rows[%d][%d].empty())\n' % (idx, validx)
        content += '    {\n'
        if typename == 'std::string':
            content += '        dataptr->%s = rows[%d][%d];\n' % (name, idx, validx)
        else:
            content += '        dataptr->%s = to<%s>(rows[%d][%d]);\n' % (name, typename, idx, validx)
        content += '    }\n'
        idx += 1
    varname = instance_data_name(struct['name'])
    content += '    BEATS_SAFE_DELETE(%s);\n' % varname
    content += '    %s = dataptr;\n' % varname
    content += '    return 0;\n'
    content += '}\n\n'    
    return content
    

def gen_field_array_assign_stmt(field, idx, delim):
    content = ''
    elemt_type = map_cpp_type(exportutil.array_element_type(field['original_type_name']))
    content += '            const vector<string>& array = Split(row[%d], "%s");\n' % (idx, delim)
    content += '            for (size_t i = 0; i < array.size(); i++)\n'
    content += '            {\n'
    content += '                item.%s.push_back(to<%s>(array[i]));\n' % (field['name'], elemt_type)
    content += '            }\n'
    return content



def gen_all_field_assign_stmt(struct):
    content = ''
    idx = 0
    array_delim = '|'
    if 'array-delim' in struct['options']:
        array_delim = struct['options']['array-delim']


    for field in struct['fields']:
        origin_type = field['original_type_name']
        typename = map_cpp_type(origin_type)
        content += '        if (!row[%d].empty())\n' % idx
        content += '        {\n'
        if origin_type.startswith('array'):
            content += gen_field_array_assign_stmt(field, idx, array_delim)
        else:
            content += '            item.%s = to<%s>(row[%d]);\n' % (field['name'], typename, idx)
        content += '        }\n'
        idx +=1
    return content


# 生成UnMarshal方法
def gen_unmarshal_method(struct):
    content = ''
    name = struct['name']
    content += 'void %s::UnMarshal(const vector<StringPiece>& row, vector<%s>* dataptr)\n' % (name, name)
    content += '{\n'
    content += '    BEATS_ASSERT(row.size() >= %d);\n' % len(struct['fields'])
    content += '    %s item;\n' % name
    content += gen_all_field_assign_stmt(struct)
    content += '        dataptr->push_back(item);\n'
    content += '}\n\n'
    return content
    
    
# 生成Load方法
def gen_static_load_method(struct):
    content = ''
    if struct['options']['parse-kv-mode']:
        return gen_static_load_method_kv_mode(struct)
    
    varname = instance_data_name(struct['name'])
    content += 'int %s::Load()\n' % struct['name']
    content += '{\n'
    content += '    const char* csvpath = "/csv/%s.csv";\n' % struct['name'].lower()
    content += '    vector<%s>* dataptr = BEATS_NEW(vector<%s>, "autogen", csvpath);\n' %  (struct['name'], struct['name'])
    content += '    const vector<vector<string>>& rows = CResourceManager::GetInstance()->ReadCsvToRows(csvpath);\n'
    content += '    BEATS_ASSERT(rows.size() > 0);\n'
    content += '    for (size_t i = 0; i < rows.size(); i++)\n'
    content += '    {\n'
    content += '        const vector<string>& row = rows[i];\n'
    content += '        BEATS_ASSERT(row.size() >= %d);\n' % len(struct['fields'])
    content += '        %s item;\n' % struct['name']
    content += gen_all_field_assign_stmt(struct)
    content += '        dataptr->push_back(item);\n'
    content += '    }\n'
    content += '    BEATS_ASSERT(dataptr->size() > 0);\n'
    content += '    BEATS_SAFE_DELETE(%s);\n' % varname
    content += '    %s = dataptr;\n' % varname
    content += '    return 0;\n'
    content += '}\n\n'
    return content
    
    
# 导出为C++头文件
def export_header_content(struct, params):
    content = ''
    content += gen_cpp_struct(struct)
    content += '\n'
    content += gen_static_method_declare(struct)
    content += '};\n\n'
    return content
    
    
# 导出为C++源文件    
def export_cpp_content(struct, params):
    content = ''
    content += gen_static_data_method(struct)
    content += gen_static_get_method(struct)
    content += gen_static_range_method(struct)
    content += gen_static_load_method(struct)
    content += '\n'
    return content
    

# 生成全局Load/Clear函数
def gen_global_method(descriptors):
    content = ''
    clear_method_content = '// load all configurations\nvoid ClearAllAutogenConfig()\n{\n'
    load_method_content = '// clear all configuration\nvoid LoadAllAutogenConfig()\n{\n'
    for struct in descriptors:
        load_method_content += '    %s::Load();\n' % struct['name']
        clear_method_content += '    BEATS_SAFE_DELETE(%s);\n' % instance_data_name(struct['name'])
    load_method_content += '}\n\n'   
    clear_method_content += '}\n\n' 
    content += load_method_content
    content += clear_method_content
    return content


# 执行导出
def run_export(info, params):
    h_include_headers = [
        '#include <stdint.h>',
        '#include <string>',
        '#include <vector>',
    ]    
    curtime = exportutil.current_time_str()
    header_content = '// This file is auto-generated by taxi at %s, DO NOT EDIT!\n\n#pragma once\n\n' % curtime
    header_content += '\n'.join(h_include_headers) + '\n\n'
    header_content += '// load all configurations\nvoid LoadAllAutogenConfig();\n\n'
    header_content += '// clear all configurations\nvoid ClearAllAutogenConfig();\n\n'

    cpp_include_headers = [
        '#include "stdafx.h"',
        '#include <stddef.h>',
        '#include <memory>',
        '#include "AutogenConfig.h"',
        '#include "Utility/Conv.h"',
        '#include "Resource/ResourceManager.h"',
    ]
    cpp_content = '// This file is auto-generated by taxi at %s, DO NOT EDIT!\n\n' % curtime
    cpp_content += '\n'.join(cpp_include_headers) + '\n\n'
    cpp_content += 'using namespace std;\n\n'

    data_only = False
    if 'data-only' in params:
        data_only = True

    descriptors = info['descriptors']
    class_content = ''
    for struct in descriptors:
        exportutil.setup_comment(struct)
        exportutil.setup_key_value_mode(struct)
        datarows = exportutil.validate_data_file(struct, params)
        struct['options']['datarows'] = datarows
        if not data_only:
            header_content += export_header_content(struct, params)
            class_content += export_cpp_content(struct, params)

    if data_only:
        return
        
    static_define_content = 'namespace {\n'
    for struct in descriptors:
        static_define_content += gen_static_data_define(struct)
    static_define_content += '}\n\n'


    outdir = '.'
    if 'outsrc-dir' in params:
        outdir = params['outsrc-dir']
    filename = outdir + '/AutogenConfig.h'
    exportutil.write_content_to_file(filename, header_content, 'gbk')
    
    cpp_content += static_define_content 
    cpp_content += gen_global_method(descriptors)
    cpp_content += class_content
    filename = outdir + '/AutogenConfig.cpp'
    exportutil.write_content_to_file(filename, cpp_content, 'gbk') 
    

def main():
    print(sys.argv)
    if len(sys.argv) < 2:
        print('no arguments specified')
        sys.exit(1)
    
    params = {}
    if len(sys.argv) >= 3:
        params = exportutil.parse_argv_kv()

    info = json.loads(sys.argv[1])
    f = codecs.open(info['filepath'], 'r', 'utf8')
    obj = json.loads(f.read())
    f.close()

    run_export(obj, params)


if __name__ == '__main__':
    main()
    
