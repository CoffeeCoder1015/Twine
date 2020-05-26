import os
from concurrent.futures import ThreadPoolExecutor

from director import cd


class scan:
    def __init__ (self,Name='',matchMode='',searchType=''):
        self.Name = Name
        self.matchMode=matchMode
        
        if matchMode == "l":
            self.Name=self.Name.lower()
        self.searchType=searchType.lower()
        self.rls = []    

        #set execution cmd
        self.ex_txt=""
        if self.searchType == '':
            self.ex_txt="""tpx.submit(dss,self)\ntpx.submit(fss,self)"""

        if self.searchType == 'f':
            self.ex_txt="tpx.submit(fss,self)"

        if self.searchType == 'd':
            self.ex_txt="tpx.submit(dss,self)"

    def mm_switch (self,string):
        if self.matchMode == "l":
            return string.lower()
        if self.matchMode == "":
            return string

    def scan_mc(self,path="."):
        if path == ".":
            path = os.getcwd()
        
        CD = cd(path=path)
        path = os.getcwd()
        #outer retrieve
        self.path = path

        #search fireld setup
        ags_lst = os.listdir()
        files = []
        for s in ags_lst:
            if os.path.isdir(s) == False:
                files.append(s)
        
        for i in files:
            ags_lst.remove(i)

        #un-touched folder&file appending
        #ALL
        if self.searchType == "":
            for i in range(0,len(ags_lst)):
                if self.Name in self.mm_switch(ags_lst[i]):
                    ap_obj = [os.getcwd(),ags_lst[i],"folder"]
                    self.rls.append(ap_obj)

            for i in range(0,len(files)):
                if self.Name in self.mm_switch(files[i]):
                    ap_obj = [os.getcwd(),files[i],"file"]
                    self.rls.append(ap_obj)
        #Directories
        if self.searchType == 'd':
            for i in range(0,len(ags_lst)):
                if self.Name in self.mm_switch(ags_lst[i]):
                    ap_obj = [os.getcwd(),ags_lst[i],"folder"]
                    self.rls.append(ap_obj)
        #Files
        if self.searchType == 'f':
            for i in range(0,len(files)):
                if self.Name in self.mm_switch(files[i]):
                    ap_obj = [os.getcwd(),files[i],"file"]
                    self.rls.append(ap_obj)

        #scanner execution & processing organization
        ptpx = ThreadPoolExecutor(max_workers=8192)
        length = len(ags_lst)
        for s in range(0,length):
            ptpx.submit(self.scanner,path=ags_lst[s])

        #wait for completion
        ptpx.shutdown(wait=True)

        #return
        CD.ret_start()

    #internal function
    def scanner(self,path=None):
        self.curfPath = os.getcwd()
        self.working = True
        if path== "" or path==None or path==0:
            path = self.curfPath
        else:
            path = self.curfPath+"\\"+path

        def dss(self):
            for D_nme in Dir:
                D_nme = self.mm_switch(D_nme)
                if self.Name in D_nme or self.Name == "":
                    ap_obj = [CurDir,D_nme,"folder"]
                    self.rls.append(ap_obj)

        def fss(self):         
            for F_nme in Files:
                F_nme = self.mm_switch(F_nme)
                if self.Name in F_nme or self.Name == "":
                    ap_obj = [CurDir,F_nme,"file"]
                    self.rls.append(ap_obj)


        tpx = ThreadPoolExecutor(max_workers=8192)

        for CurDir,Dir,Files in os.walk(path):
            exec(self.ex_txt)

        tpx.shutdown(wait=True)