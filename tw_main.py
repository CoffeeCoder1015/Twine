import os

from scan_py import scan

SCAN = scan(Name='',matchMode='',searchType='')
#SCAN.scan_mc(path="..\\..")
SCAN.scanner(path=".")
rls =SCAN.rls
for i in rls:
    print(i)