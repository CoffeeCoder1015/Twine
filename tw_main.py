import os

from scan_py import scan

SCAN = scan(Name='',matchMode='',searchType='d')
SCAN.scan_mc(path="..")

path=SCAN.path
rls =SCAN.rls

print(rls)