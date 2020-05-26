import os

import ATM
from scan_py import scan

SCAN = scan(Name='',matchMode='',searchType='')
SCAN.scan_mc(path="..")

#cache init
cacher = ATM.cache_machine(name="Twine",start_pos=os.getcwd())
path=SCAN.path
rls =SCAN.rls

rls = list(str(rls))
del rls[0]
del rls[len(rls)-1]
rls="".join(rls).replace("'","").replace("], [","]\n[")

cacher.deposit(data=rls,name=path,cacheTarget="twine_cache")
print(cacher.withdraw(name=path,cacheTarget="twine_cache"))
cacher.clear(cacheTarget="twine_cache")