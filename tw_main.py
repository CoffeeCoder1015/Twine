from scan_py import scan

SCAN = scan(Name='', matchMode='', searchType='')
SCAN.scan_mc(path="..")
rlsc = SCAN.rls
print(len(rlsc))