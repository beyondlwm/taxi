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
import exportutil


def map_go_type(typ):
    type_mapping = {
        'bool':     'bool',
        'int8':     'int8',
        'uint8':    'uint8',
        'int16':    'int16',
        'uint16':   'uint16',
        'int':      'int',
        'int32':    'int32',
        'uint32':   'uint32',
        'int64':    'int64',
        'uint64':   'uint64',
        'float':    'float32',
        'float32':  'float32',
        'float64':  'float64',
        'enum':     'int',
        'string':   'string',
    }
    if typ.startswith('array'):
        t = exportutil.array_element_type(typ)
        elem_type = type_mapping[t]
        return '[]' + elem_type
    return type_mapping[typ]


def const_key_name(name):
    return 'Key%sName' % name

def gen_go_struct(struct):
    fields = struct['fields']
    if struct['options']['parse-kv-mode']:
        fields = exportutil.construct_struct_kv_fields(struct)

    content = '// %s\n' % struct['comment']
    content += 'type %s struct\n{\n' % struct['camel_case_name']
    for field in fields:
        typename = map_go_type(field['original_type_name'])
        assert typename != "", field['original_type_name']
        content += '    %s %s // %s\n' % (field['camel_case_name'], typename, field['comment'])
    return content


#解析布尔    
def gen_parse_bool(text):
    content = '\t\tval, err := strconv.ParseBool(%s)\n' % text
    content += '\t\tif err != nil {\n'
    content += '\t\t\treturn fmt.Errorf("parse %s [%%v], %%v", %s, err)\n' % (text, text)
    content += '\t\t}\n'
    return content
    
    
# 解析整数    
def gen_parse_integer(text):
    content = '\t\tval, err := strconv.Atoi(%s)\n' % text
    content += '\t\tif err != nil {\n'
    content += '\t\t\treturn fmt.Errorf("parse %s [%%v], %%v", %s, err)\n' % (text, text)
    content += '\t\t}\n'
    return content

    
# 解析浮点数
def gen_parse_float(text):
    content = '\t\tval, err := strconv.ParseFloat(%s, 64)\n' % text
    content += '\t\tif err != nil {\n'
    content += '\t\t\treturn fmt.Errorf("parse %s [%%v], %%v", %s, err)\n' % (text, text)
    content += '\t\t}\n'
    return content
    

def gen_field_assgin_stmt(prefix, name, typename):
    content = ''
    if typename == 'bool':
        content += gen_parse_bool(prefix)
        content += '\t\tp.%s = %s(val)\n}\n' % (name, typename)
    elif exportutil.is_integer_type(typename):
        content += gen_parse_integer(prefix)
        content += '\t\tp.%s = %s(val)\n}\n' % (name, typename)
    elif exportutil.is_floating_type(typename):
        content += gen_parse_float(prefix)
        content += '\t\tp.%s = %s(val)\n}\n' % (name, typename)
    else:
        content += '\t\tp.%s = %s\n}\n' % (name, prefix)
    return content
    

# 数组格式赋值
def gen_field_array_assign_stmt(field, idx, array_delim):
    content = ''    
    name = field['name']
    elem_type = exportutil.array_element_type(field['original_type_name'])
    elem_type = map_go_type(elem_type)
    content += '\tfor _, item := range strings.Split(row[%d], "%s") {\n' % (idx, array_delim)
    prefix = 'item'
    
    if elem_type == 'bool':
        content += gen_parse_bool(prefix)
        content += '\tp.%s = append(p.%s, %s(%s))\n' % (name, name, elem_type, 'val')
    elif exportutil.is_integer_type(elem_type):
        content += gen_parse_integer(prefix)
        content += '\tp.%s = append(p.%s, %s(%s))\n' % (name, name, elem_type, 'val')
    elif exportutil.is_floating_type(elem_type):
        content += gen_parse_float(prefix)
        content += '\tp.%s = append(p.%s, %s(%s))\n' % (name, name, elem_type, 'val')
    else:
        content += '\tp.%s = append(p.%s, %s)\n' % (name, name, prefix)
    content += '\t}\n'
    return content
    
    
# 生成kv模式的parse方法    
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

    content += 'func (p *%s) ParseFromRows(rows [][]string) error {\n' % struct['camel_case_name']
    content += '\tif len(rows) < %d {\n' % len(rows)
    content += '\t\tlog.Panicf("%s:row length out of index, %%d < %d", len(rows))\n' % (struct['name'], len(rows))
    content += '\t}\n'
    
    idx = 0
    for row in rows:
        content += '\tif rows[%d][%d] != "" {\n' % (idx, validx)
        name = rows[idx][keyidx].strip()
        typename = rows[idx][typeidx].strip()
        typename = map_go_type(typename)     
        name = exportutil.camel_case(name)
        prefix = 'rows[%d][%d]' % (idx, validx)
        content += gen_field_assgin_stmt(prefix, name, typename)
        idx += 1
    content += '    return nil\n'
    content += '}\n\n'    
    return content

#     
def gen_parse_method(struct):
    if struct['options']['parse-kv-mode']:
        return gen_static_load_method_kv_mode(struct)

    
    array_delim = '|'
    if 'array-delim' in struct['options']:
        array_delim = struct['options']['array-delim']
        
    content = ''
    content += 'func (p *%s) ParseFromRow(row []string) error {\n' % struct['camel_case_name']
    content += '\tif len(row) < %d {\n' % len(struct['fields'])
    content += '\t\tlog.Panicf("%s: row length out of index %%d", len(row))\n' % struct['name']
    content += '\t}\n'
    idx = 0
    for field in struct['fields']:
        content += '\tif row[%d] != "" {\n' % idx
        origin_type_name = field['original_type_name']
        typename = map_go_type(origin_type_name)
        name = field['camel_case_name']
        prefix = 'row[%d]' % idx
        if origin_type_name.startswith('array'):
            content += gen_field_array_assign_stmt(field, idx, array_delim)
            content += '\t}\n'
        else:
            content += gen_field_assgin_stmt(prefix, name, typename)
        idx += 1
    content += 'return nil\n'
    content += '}\n\n'
    return content


def gen_load_method_kv(struct):
    content = ''
    content += 'func Load%s(loader DataSourceLoader) (*%s, error) {\n' % (struct['name'], struct['name'])
    content += '\tbuf, err := loader.LoadDataByKey(%s)\n' % const_key_name(struct['name'])
    content += '\tif err != nil {\n'
    content += '\treturn nil, err\n'
    content += '\t}\n'
    content += '\tr := csv.NewReader(buf)\n'
    content += '\trows, err := r.ReadAll()\n'
    content += '\tif err != nil {\n'
    content += '\t    log.Errorf("%s: csv read all, %%v", err)\n' % struct['name']
    content += '\t    return nil, err\n'
    content += '\t}\n'
    content += '\tvar item %s\n' % struct['name']
    content += '\tif err := item.ParseFromRows(rows); err != nil {\n'
    content += '\t    log.Errorf("%s: parse row %%d, %%v", len(rows), err)\n' % struct['name']
    content += '\t    return nil, err\n'
    content += '\t}\n'
    content += 'return &item, nil\n'
    content += '}\n\n'
    return content



def gen_load_method(struct):
    content = ''
    if struct['options']['parse-kv-mode']:
        return gen_load_method_kv(struct)
            
    content += 'func Load%sList(loader DataSourceLoader) ([]*%s, error) {\n' % (struct['name'], struct['name'])
    content += '\tbuf, err := loader.LoadDataByKey(%s)\n' % const_key_name(struct['name'])
    content += '\tif err != nil {\n'
    content += '\t    return nil, err\n'
    content += '\t}\n'
    content += '\tvar list []*%s\n' % struct['name']
    content += '\tvar r = csv.NewReader(buf)\n'
    content += '\tfor i := 0; ; i++ {\n'
    content += '\t row, err := r.Read()\n'
    content += '    if err == io.EOF {\n'
    content += '        break\n'
    content += '    }\n'
    content += '    if err != nil {\n'
    content += '        log.Errorf("%s: read csv %%v", err)\n'  % struct['name']
    content += '        return nil, err\n'
    content += '    }\n'
    content += '    var item %s\n' % struct['name']
    content += '    if err := item.ParseFromRow(row); err != nil {\n'
    content += '\t     log.Errorf("%s: parse row %%d, %%s, %%v", i+1, row, err)\n' % struct['name']
    content += '    return nil, err\n'
    content +=  '}\n'
    content +='\tlist = append(list, &item)\n'
    content += '}\n'
    content += 'return list, nil\n'
    content += '}\n\n'
    return content 
    

def gen_const_names(descriptors):
    content = 'const (\n'
    for struct in descriptors:
        content += '\t%s = "%s"\n' % (const_key_name(struct['name']), struct['name'].lower())
    content += ')\n\n'
    return content


def export_go_content(struct, params):
    content = ''
    content += gen_go_struct(struct)
    content += '}\n\n'
    content += gen_parse_method(struct)
    content += gen_load_method(struct)
    return content

    
# 执行导出
def run_export(info, params):
    descriptors = info['descriptors']
    content = '// This file is auto-generated by taxi at %s, DO NOT EDIT!\n\n' % exportutil.current_time_str()
    content += 'package %s\n' % params['pkg']
    content += gen_const_names(descriptors)
    
    data_only = False
    if 'data-only' in params:
        data_only = True
    
    for struct in descriptors:
        exportutil.setup_comment(struct)
        exportutil.setup_key_value_mode(struct)
        datarows = exportutil.validate_data_file(struct, params)
        struct['options']['datarows'] = datarows
        if not data_only:
            content += export_go_content(struct, params)

    if data_only:
        return
        
    outdir = '.'
    if 'outsrc-dir' in params:
        outdir = params['outsrc-dir']
    filename = outdir + '/stub.go'
    f = codecs.open(filename, 'w', 'utf-8')
    f.writelines(content)
    f.close()
    print('wrote to %s' % filename)

    # run go import
    os.system('goimports -w ' + filename)


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