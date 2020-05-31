import os
import threading
from concurrent.futures import ThreadPoolExecutor

from director import cd


class scan:
    def __init__(self, Name='', matchMode='', searchType=''):
        self.Name = Name
        self.matchMode = matchMode

        if matchMode == "l":
            self.Name = self.Name.lower()
        self.searchType = searchType.lower()
        self.rls = []

        # set execution cmd
        # below (↓) is the execution command for scanning
        self.ex_txt = ""
        # ALL
        if self.searchType == '':
            self.ex_txt = ("""
fst = threading.Thread(target=fss,args=(self,))
dst = threading.Thread(target=dss,args=(self,))
threads.append(fst)
threads.append(dst)
fst.start()
dst.start()
""")
        # file
        if self.searchType == 'f':
            self.ex_txt = ("""
fst = threading.Thread(target=fss,args=(self,))
threads.append(fst)
fst.start()""")

        # folder
        if self.searchType == 'd':
            self.ex_txt = ("""
dst = threading.Thread(target=dss,args=(self,))
threads.append(dst)
dst.start()""")

        # scan filter command
        self.scf = ""
        # file
        if Name == "":
            self.scf = ("""
ap_obj = [self.CurDir, F_nme, "file"]
self.rls.append(ap_obj)
            """)
        else:
            self.scf = ("""
F_nme = self.mm_switch(F_nme)
if self.Name in F_nme or self.Name == "":
    ap_obj = [self.CurDir, F_nme, "file"]
    self.rls.append(ap_obj)""")

        self.scd = ""
        # folder
        if Name == "":
            self.scd = ("""
ap_obj = [self.CurDir, D_nme, "folder"]
self.rls.append(ap_obj)
            """)
        else:
            self.scd = ("""
D_nme = self.mm_switch(D_nme)
if self.Name in D_nme or self.Name == "":
    ap_obj = [self.CurDir, D_nme, "folder"]
    self.rls.append(ap_obj)
""")

    def mm_switch(self, string):
        if self.matchMode == "l":
            return string.lower()
        if self.matchMode == "":
            return string

    def scan_mc(self, path="."):
        if path == ".":
            path = os.getcwd()

        CD = cd(path=path)
        path = os.getcwd()
        self.path = path

        # search field setup
        ags_lst = os.listdir()
        files = []
        for s in ags_lst:
            if os.path.isdir(s) == False:
                files.append(s)

        for i in files:
            ags_lst.remove(i)

        # un-touched folder&file appending
        # ALL
        if self.searchType == "":
            for i in range(0, len(ags_lst)):
                if self.Name in self.mm_switch(ags_lst[i]):
                    ap_obj = [os.getcwd(), ags_lst[i], "folder"]
                    self.rls.append(ap_obj)

            for i in range(0, len(files)):
                if self.Name in self.mm_switch(files[i]):
                    ap_obj = [os.getcwd(), files[i], "file"]
                    self.rls.append(ap_obj)
        # Directories
        if self.searchType == 'd':
            for i in range(0, len(ags_lst)):
                if self.Name in self.mm_switch(ags_lst[i]):
                    ap_obj = [os.getcwd(), ags_lst[i], "folder"]
                    self.rls.append(ap_obj)
        # Files
        if self.searchType == 'f':
            for i in range(0, len(files)):
                if self.Name in self.mm_switch(files[i]):
                    ap_obj = [os.getcwd(), files[i], "file"]
                    self.rls.append(ap_obj)

        # scanner execution & processing organization
        length = len(ags_lst)
        ptx = ThreadPoolExecutor(max_workers=1000000)
        for s in range(0, length):
            ptx.submit(self.scanner, ags_lst[s])

        # wait for completion
        ptx.shutdown(wait=True)

        # return
        CD.ret_start()

    # internal function
    def scanner(self, path=None):
        self.curfPath = os.getcwd()
        self.working = True
        if path == "" or path == None or path == 0:
            path = self.curfPath
        else:
            path = self.curfPath+"\\"+path
        path = path.replace("\\", "/")

        def dss(self):
            for D_nme in Dir:
                exec(self.scd)

        def fss(self):
            for F_nme in Files:
                exec(self.scf)

        threads = []
        for self.CurDir, Dir, Files in os.walk(path):
            self.CurDir = self.CurDir.replace("\\", "/")
            exec(self.ex_txt)

<<<<<<< HEAD
        #concurrent wait for completetion
        def jthrd (thread):
            thread.join()

        tjpx = ThreadPoolExecutor(max_workers=100e+100)
=======
>>>>>>> 07eb7094fd23620d8cd5b4918016daa11eab1094
        for t in threads:
            tjpx.submit(jthrd,t)
