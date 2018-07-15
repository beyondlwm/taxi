# taxi
自动代码生成工具

## Usage

参数说明：
* mode 导入模式，目前仅支持mysql和excel
* import-args 导入器参数
* export-args 导出脚本参数
* exporter-path 导出脚本路径

根据MySQL表导出结构体

导入器参数
* user  数据库账号
* passwd 数据库密码
* db  数据库名称

示例用法：

```bash
taxi --mode=mysql --import-args=user=root,passwd=holyshit,db=mydbname --export-args=pkg=proto --exporter-path=taxi\scripts\sql_exporter.py
```

Export Excel

导入器参数
* filename      excel文件名
* parse-mode    (active-only, all)是否导出所有sheet
* sheet         指定导出excel的某个sheet

示例用法：

```bash
taxi --mode=excel --import-args=filename=abc.xlsx --exporter-path=taxi\scripts\export_excel_cpp.py
```