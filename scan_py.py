import json
from concurrent.futures import ThreadPoolExecutor
import os
import numpy as np
import ATM
from director import cd

class scan:
    def __init__ (self,Name='',matchMode='',searchType='',returnFormat=''):
        self.Name = Name
        self.matchMode=matchMode
        
        if matchMode == "l":
            self.Name=self.Name.lower()
            
        self.searchType=searchType.lower()

        self.returnFormat=returnFormat.lower()
        self.ret_lst_raw = np.empty((0,3))         

    def mm_switch (self,string):
        if self.matchMode == "l":
            return string.lower()
        else:
            return string

    def scan_mc(self,path="."):
        if path == ".":
            path = os.getcwd()
        
        #cache init
        cacher = ATM.cache_machine(name="Twine",Type="b",start_pos=os.getcwd())

        cd(path=path)
        path = os.getcwd()
        
        #search fireld setup
        ags_lst = os.listdir()
        files = []
        for s in ags_lst:
            if os.path.isdir(s) == False:
                files.append(s)
        
        for i in files:
            ags_lst.remove(i)

        #un-touched folder&file appending
        for i in range(0,len(ags_lst)):
            if self.Name in self.mm_switch(ags_lst[i]):
                ap_obj = [os.getcwd(),ags_lst[i],"folder"]
                self.ret_lst_raw = np.append(self.ret_lst_raw,[ap_obj],axis=0)

        for i in range(0,len(files)):
            if self.Name in self.mm_switch(files[i]):
                ap_obj = [os.getcwd(),files[i],"file"]
                self.ret_lst_raw = np.append(self.ret_lst_raw,[ap_obj],axis=0)

        #scanner execution & processing organization
        ptpx = ThreadPoolExecutor(max_workers=8192)
        length = len(ags_lst)
        self.scan_threads = length
        for s in range(0,length):
            ptpx.submit(self.scanner,path=ags_lst[s])

        #wait for completion
        while self.scan_threads > 0:
            print("Running Threads:",self.scan_threads,end="\r")

        print("\n-----")
        print(self.ret_lst_raw)
        #deposit cache
        #self.ret_lst_raw = str(self.ret_lst_raw.tolist()).replace("'","")
        #cacher.deposit(data=self.ret_lst_raw,name=path,cacheTarget="twine_cache")
        #print(cacher.withdraw(name=path,cacheTarget="twine_cache"))

    #internal function
    def scanner(self,path=None):
        self.curfPath = os.getcwd()
        self.working = True
        if path== "" or path==None or path==0:
            path = self.curfPath
        else:
            path = self.curfPath+"\\"+path

        def dss(self,ID):
            self.working_threads.append(["dss",ID])  

            for D_nme in Dir:
                D_nme = self.mm_switch(D_nme)
                if self.Name in D_nme:
                    ad_obj = np.array([CurDir,D_nme,"folder"])
                    self.ret_lst_raw = np.append(self.ret_lst_raw,[ad_obj],axis=0)

            self.comp_threads+=1
            self.working_threads.remove(["dss",ID])

        def fss(self,ID):         
            self.working_threads.append(["fss",ID]) 

            for F_nme in Files:
                F_nme = self.mm_switch(F_nme)
                if self.Name in F_nme:
                    ad_obj = np.array([CurDir,F_nme,"file"])
                    self.ret_lst_raw = np.append(self.ret_lst_raw,[ad_obj],axis=0)
            self.comp_threads+=1
            self.working_threads.remove(["fss",ID])


        tpx = ThreadPoolExecutor(max_workers=8192)
        d_ss_thrs = 0
        f_ss_thrs = 0
        self.comp_threads = 0
        self.working_threads = []

        if self.searchType == '':
            for CurDir,Dir,Files in os.walk(path):
                    d_ss_thrs+=1
                    tpx.submit(dss,self,d_ss_thrs)
                    f_ss_thrs+=1
                    tpx.submit(fss,self,f_ss_thrs)
            
        if self.searchType == 'f':
            for CurDir,Dir,Files in os.walk(path):
                    f_ss_thrs+=1
                    tpx.submit(fss,self,f_ss_thrs)
        
        if self.searchType == 'd':
            for CurDir,Dir,Files in os.walk(path):
                    d_ss_thrs+=1
                    tpx.submit(dss,self,d_ss_thrs)
                    
        self.scan_threads-=1