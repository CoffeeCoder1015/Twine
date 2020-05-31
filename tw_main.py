from scan_py import scan
import timer

SCAN = scan(Name='', matchMode='', searchType='')
SCAN.scan_mc(path=".")
rlsc = SCAN.rls
print(SCAN.path)
print(len(rlsc))
