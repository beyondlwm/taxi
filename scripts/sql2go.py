# Copyright (C) 2018-present ichenq@outlook.com. All rights reserved.
# Distributed under the terms and conditions of the Apache License.
# See accompanying files LICENSE.

import codecs
import json
import os
import sys

type_mapping = {
    'bool': 'bool',
    'int8': 'int8',
    'uint8': 'uint8',
    'int16': 'int16',
    'uint16': 'uint16',
    'int': 'int',
    'int32': 'int32',
    'uint32': 'uint32',
    'int64': 'int64',
    'uint64': 'uint64',
    'float': 'float32',
    'float32': 'float32',
    'float64': 'float64',
    'string': 'string',
    'bytes': '[]byte',
    'datetime': 'time.Time',
}


def get_key_list(options, key):
    if key in options:
        value = options[key]
        if len(value) > 0:
            return value.split(',')
    return []


def gen_go_struct(struct, params):
    content = '// %s\n' % struct['comment']
    name = struct['camel_case_name']
    if 'prefix' in params:
        name = params['prefix'] + name
    content += 'type %s struct {\n' % name
    for field in struct['fields']:
        typename = type_mapping[field['type_name']]
        assert typename != "", field['type_name']
        content += '\t%s %s `json:"%s"`// %s\n' % (field['camel_case_name'], typename, field['name'], field['comment'])
    content += '}\n'
    return content



def gen_where_clause(struct):
    content = ''
    keys = []
    primary_keys = get_key_list(struct['options'], 'primary_keys')
    if len(primary_keys) > 0:
        content += ' WHERE `%s`=?' % primary_keys[0]
        keys.append(primary_keys[0])
    else:
        unique_keys = get_key_list(struct['options'], 'unique_keys')
        print('unique keys: ', unique_keys)
        if len(unique_keys) > 0:
            content += " WHERE "
            for i, key in enumerate(unique_keys):
                keys.append(key)
                content += '`%s`=?' % key
                if i + 1 < len(unique_keys):
                    content += ' AND '
    return content, keys


# 生成select语句
def gen_select_stmt_variable(struct, params):
    clause, keys = gen_where_clause(struct)
    name = struct['camel_case_name']
    if 'prefix' in params:
        name = params['prefix'] + name    
    content = '\tconst Sql%sStmt = "SELECT ' % name
    for i, field in enumerate(struct['fields']):
        content += '`%s`' % field['name']
        if i + 1 < len(struct['fields']):
            content += ', '
    content += ' FROM `%s`' % struct['name']
    content += clause
    content += '"\n\n'
    return content


# 生成insert语句
def gen_insert_stmt_method(struct, params):
    name = struct['camel_case_name']
    if 'prefix' in params:
        name = params['prefix'] + name       
    content = 'func (p *%s) InsertStmt() *storage.SqlOperation {\n' % name
    content += '\t return storage.NewSqlOperation("INSERT INTO `%s`(' % struct['name']
    mark = ''
    for i, field in enumerate(struct['fields']):
        mark += '?, '
        content += '`%s`' % field['name']
        if i + 1 < len(struct['fields']):
            content += ', '
    content += ') VALUES(%s)", ' % mark[:-2]
    clause = ''
    for field in struct['fields']:
        clause += 'p.%s, ' % field['camel_case_name']
    content += clause[:-2]
    content += ')\n}\n'
    return content


# 生成update语句
def gen_update_stmt_method(struct, params):
    name = struct['camel_case_name']
    if 'prefix' in params:
        name = params['prefix'] + name          
    clause, keys = gen_where_clause(struct)
    assert len(clause) > 0, struct
    content = 'func (p *%s) UpdateStmt() *storage.SqlOperation {\n' % name
    content += '\t return storage.NewSqlOperation("UPDATE `%s` SET ' % struct['name']
    for i, field in enumerate(struct['fields']):
        if field['name'] not in keys:
            content += '`%s`=?' % field['name']
            if i + 1 < len(struct['fields']):
                content += ', '
    content += clause
    content += '",'
    clause = ''
    for field in struct['fields']:
        if field['name'] not in keys:
            clause += 'p.%s, ' % field['camel_case_name']
    for field in struct['fields']:
        if field['name'] in keys:
            clause += 'p.%s, ' % field['camel_case_name']
    content += clause[:-2]
    content += ')\n}\n'
    return content    


# 生成delete语句
def gen_remove_stamt_method(struct, params):
    name = struct['camel_case_name']
    if 'prefix' in params:
        name = params['prefix'] + name      
    clause, keys = gen_where_clause(struct)
    assert len(clause) > 0, struct  
    content = 'func (p *%s) DeleteStmt() *storage.SqlOperation {\n' % name
    content += '\t return storage.NewSqlOperation("DELETE FROM `%s`' % struct['name']
    content += clause
    content += '", '
    clause = ''
    for field in struct['fields']:
        if field['name'] in keys:
            clause += 'p.%s, ' % field['camel_case_name']
    content += clause[:-2]
    content += ')\n}\n'
    return content


# 执行导出
def run_export(info, params):
    descriptors = info['descriptors']
    content = 'package %s\n\n' % params['pkg']
    content += 'import (\n'
    content += '\t"fatchoy/storage"'
    content += ')\n\n'
    for struct in descriptors:
        content += gen_go_struct(struct, params)
        content += '\n'
        content += gen_select_stmt_variable(struct, params)
        content += '\n'
        content += gen_insert_stmt_method(struct, params)
        content += '\n'
        content += gen_update_stmt_method(struct, params)
        content += '\n'
        content += gen_remove_stamt_method(struct, params)

    filename = 'stub.go'
    if 'out' in params:
        filename = params['out']
    f = codecs.open(filename, 'w', 'utf8')
    f.writelines(content)
    f.close()
    print('wrote to %s' % filename)
    os.system('goimports -w ' + filename)


def main():
    params = {}
    print(sys.argv)
    if len(sys.argv) < 2:
        print('no arguments specified')
        sys.exit(1)

    if len(sys.argv) >= 3:
        if len(sys.argv[2]) > 0:
            kvlist = sys.argv[2].split(',')
            for item in kvlist:
                kv = item.split('=')
                assert len(kv) == 2, item
                params[kv[0]] = kv[1]

    info = json.loads(sys.argv[1])
    f = codecs.open(info['filepath'], 'r', 'utf8')
    obj = json.loads(f.read())
    f.close()
    #print(obj)
    if info['format'] == 'json':
        run_export(obj, params)
    else:
        assert False, "unsupported format " + info['format']


if __name__ == '__main__':
    main()
