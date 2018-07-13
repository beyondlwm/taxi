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


def gen_go_struct(struct):
    content = '// %s\n' % struct['comment']
    content += 'type %s struct {\n' % struct['camel_case_name']
    for field in struct['fields']:
        typename = type_mapping[field['type_name']]
        assert typename != "", field['type_name']
        content += '\t%s %s `json:"%s"`// %s\n' % (field['camel_case_name'], typename, field['name'], field['comment'])
    content += '}\n'
    return content


# 生成select语句
def gen_select_stmt_method(struct):
    content = 'func (p *%s) SelectStmt() string {\n' % struct['camel_case_name']
    content += '\treturn "SELECT '
    for i, field in enumerate(struct['fields']):
        content += '`%s`' % field['name']
        if i + 1 < len(struct['fields']):
            content += ', '
    content += ' FROM `%s`' % struct['name']

    primary_keys = get_key_list(struct['options'], 'primary_keys')
    if len(primary_keys) > 0:
        content += ' WHERE `%s`=?' % primary_keys[0]
    else:
        unique_keys = get_key_list(struct['options'], 'unique_keys')
        print('unique keys: ', unique_keys)
        if len(unique_keys) > 0:
            content += "WHERE "
            for i, key in enumerate(unique_keys):
                content += '`%s`=?' % key
                if i + 1 < len(unique_keys):
                    content += ' AND '
    content += '"\n}\n'
    return content

# 生成insert语句
def gen_insert_stmt_method(struct):
    content = 'func (p *%s)InsertStatment() string {\n' % struct['camel_case_name']
    content += '\t return fmt.Sprintf('
    for field in struct['fields']:
        content += '`%s`' % field['name']

# 执行导出
def run_export(info, params):
    descriptors = info['descriptors']
    content = 'package %s\n\n' % params['pkg']
    for descriptor in descriptors:
        content += gen_go_struct(descriptor)
        content += '\n'
        content += gen_select_stmt_method(descriptor)
        content += '\n'

    filename = 'stub.go'
    f = codecs.open(filename, 'w', 'utf8')
    f.writelines(content)
    f.close()
    print('wrote to %s' % filename)


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
