# Copyright (C) 2018-present ichenq@outlook.com. All rights reserved.
# Distributed under the terms and conditions of the Apache License.
# See accompanying files LICENSE.


from __future__ import print_function
from __future__ import unicode_literals
from __future__ import division

import re
import sys
import csv
import time
import codecs
import random
import string
import shutil
import tempfile
import filecmp
import os.path


def is_integer_type(typ):
     table =['int8', 'uint8', 'int16', 'uint16','int', 'int32', 'uint32', 'int64', 'uint64', ]
     return typ in table    


def is_floating_type(typ):
    return typ == 'float' or typ == 'float32' or typ == 'float64'
    
    
def current_time_str():
    return time.strftime('%Y-%m-%d', time.localtime())


def array_element_type(typ):
    return typ[6:-1]

# 最长串的大小
def max_field_length(table, key, f):
    max_len = 0
    for v in table:
        n = len(v[key])
        if f is not None:
            n = len(f(v[key]))
        if n > max_len:
            max_len = n
    return max_len
    
    
# 空格对齐
def pad_spaces(text, min_len):
    if len(text) < min_len:
        for n in range(min_len - len(text)):
            text += ' '
    return text  

    
def random_word(length):
   letters = string.ascii_lowercase
   return ''.join(random.choice(letters) for i in range(length))
   

# 转换为驼峰风格
def camel_case(name, compare=True):
    assert name != '_'

    if name.find('_') < 0:
        return name

    words = [x.title() for x in name.split('_')]
    for i in range(len(words)):
        word = words[i].upper()
        if word in commonInitialisms and compare:
            words[i] = word
    return ''.join(words)


def setup_key_value_mode(struct):
    struct['options']['parse-kv-mode'] = False
    if 'key-value-column' in struct['options']:
        kv = struct['options']['key-value-column'].split(',')
        assert len(kv) == 2
        struct['options']['parse-kv-mode'] = True
        struct['options']['key-column'] = int(kv[0])
        struct['options']['value-column'] = int(kv[1])


def setup_comment(struct):
    if struct['comment'].strip() == "":
        if 'class-comment' in struct['options']:
            struct['comment'] = struct['options']['class-comment']  


def parse_argv_kv():
    params = {}
    if len(sys.argv[2]) > 0:
        kvlist = sys.argv[2].split(',')
        for item in kvlist:
            kv = item.split('=')
            assert len(kv) == 2, item
            key = kv[0].strip()
            value = kv[1].strip()
            params[key] = value
    return params


def get_field_by_column_index(struct, column_idx):
    assert column_idx > 0
    idx = 0
    for field in struct['fields']:
        if field['column_index'] == column_idx:
            return idx, field
        idx += 1
    assert False
    
#    
def construct_struct_kv_fields(struct):
    rows = struct['options']['datarows']
    keycol = struct['options']['key-column']
    valcol = struct['options']['value-column']
    typcol = int(struct['options']['value-type-column'])
    assert keycol > 0 and valcol > 0 and typcol > 0
    comment_idx = -1
    if 'comment-column' in struct['options']:
        commentcol = int(struct['options']['comment-column'])
        assert commentcol > 0
        comment_field = {}
        comment_idx, comment_field = get_field_by_column_index(struct, commentcol)
    fields = []
    
    key_idx, key_field = get_field_by_column_index(struct, keycol)
    # print('get_field_by_column_index', key_idx, key_field)
    value_idx, value_field = get_field_by_column_index(struct, valcol)
    type_idx, type_field = get_field_by_column_index(struct, typcol)
    
    for i in range(len(rows)):
        # print(rows[i])
        name = rows[i][key_idx].strip()
        typename = rows[i][type_idx].strip()
        assert len(name) > 0, (rows[i], key_idx)
        comment = ''
        if comment_idx >= 0:
            comment = rows[i][comment_idx].strip()
        field = {
            'name': name,
            'camel_case_name': camel_case(name),
            'type_name': typename,
            'original_type_name': typename,
            'comment': comment,
        }
        fields.append(field)
    
    return fields

    
# 将excel里的浮点数四舍五入为整数
def validate_data_file(struct, params):
    datadir = './'
    datafile = struct['options']['datafile']
    if 'outdata-dir' in params:
        datadir = params['outdata-dir']
    filename = '%s/%s.csv' % (datadir, struct['name'].lower())

    rows = []
    with codecs.open(datafile, 'r', 'utf-8') as f:
        rows = [row for row in csv.reader(f)]

    new_rows = []
    fields = struct['fields']
    print(struct['name'], 'total rows', len(rows), 'total fields', len(fields))
    for row in rows:
        assert len(row) >= len(fields), (len(fields), len(row), row)
        for j in range(len(row)):
            if j >= len(fields):
                continue
            typename = fields[j]['type_name']
            if is_integer_type(typename) and len(row[j]) > 0:
                f = float(row[j])  # test if ok
                if row[j].find('.') >= 0:
                    print('round interger', row[j], '-->', round(f))
                    row[j] = str(round(float(row[j])))
            else:
                if is_floating_type(typename) and len(row[j]) > 0:
                    f = float(row[j]) # test if ok
                    
        # skip all empty row
        is_all_empty = True
        for text in row:
            if len(text.strip()) > 0:
                is_all_empty = False
                break
        if not is_all_empty:
            new_rows.append(row)

    assert len(new_rows) <= len(rows)
    f = codecs.open(filename, 'w','utf-8')
    w = csv.writer(f)
    w.writerows(new_rows) #写新文件

    print(datafile, '-->', filename)
    #shutil.copyfile(datafile, filename)
    return new_rows


#写入内容到文件，如果内容相同则不覆盖
def write_content_to_file(filename, content, enc):
    # write to a temp file
    # print(tempfile.gettempdir())
    tmp_filename = '%s/taxi_tmp_%s' % (tempfile.gettempdir(), random_word(10))
    f = codecs.open(tmp_filename, 'w', enc)
    f.writelines(content)
    f.close()
        
    if os.path.isfile(filename) and filecmp.cmp(tmp_filename, filename):
        print('file content not modified', filename)
    else:
        shutil.move(tmp_filename, filename)


commonInitialisms = {
    "API": True,
    "ASCII": True,
    "CPU": True,
    "CSS": True,
    "DNS": True,
    "EOF": True,
    "GUID": True,
    "HTML": True,
    "HTTP": True,
    "HTTPS": True,
    "ID": True,
    "IP": True,
    "JSON": True,
    "LHS": True,
    "QPS": True,
    "RAM": True,
    "RHS": True,
    "RPC": True,
    "SLA": True,
    "SMTP": True,
    "SQL": True,
    "SSH": True,
    "TCP": True,
    "TLS": True,
    "TTL": True,
    "UDP": True,
    "UI": True,
    "UUID": True,
    "URI": True,
    "URL": True,
    "UTF8": True,
    "VM": True,
    "XML": True,
    "XSRF": True,
    "XSS": True,
    "NPC": True,
    "VIP": True,
}