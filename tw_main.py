from scan_py import scan

SCAN = scan(Name='',matchMode='',searchType='')
SCAN.scan_mc(path="..")
rls =SCAN.rls
print(SCAN.path)
print(len(rls))