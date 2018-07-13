# taxi
Code generate toolchain

## Usage

Export MySQL

Arguments:
* user=
* passwd=
* db=


`taxi --mode=mysql --import-args=user=root,passwd=holyshit,db=mydbname --export-args=pkg=proto --export-template-dir=taxi\scripts`

Export Excel

Arguments:
* filename=
* parse-mode=active-only, all
* sheet=