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

    def mm_switch(self, string):
        if self.matchMode == "l":
            return string.lower()
        if self.matchMode == "":
            return string

    def scan_mc(self, path="."):
        if path == ".":
            path = os.getcwd()

        CD = cd(path=path)
        path = os.getcwd().replace("\\", "/")
        self.path = path.replace("\\", "/")

        # search field setup
        ags_lst = os.listdir()
        files = []
        for s in ags_lst:
            if os.path.isdir(s) == False:
                files.append(s)

        for i in files:
            ags_lst.remove(i)

        # un-touched folder&file appending
        if self.searchType == "" or self.searchType == 'd':
            for i in range(0, len(ags_lst)):
                if self.Name in self.mm_switch(ags_lst[i]):
                    ap_obj = [os.getcwd(), ags_lst[i], "folder"]
                    self.rls.append(ap_obj)

        if self.searchType == "" or self.searchType == 'f':
            for i in range(0, len(files)):
                if self.Name in self.mm_switch(files[i]):
                    ap_obj = [os.getcwd(), files[i], "file"]
                    self.rls.append(ap_obj)

        # scanner execution & processing organization
        length = len(ags_lst)
        threads = []

        for s in range(0, length):
            thrd = threading.Thread(target=self.scanner, args=(ags_lst[s],))
            thrd.start()
            threads.append(thrd)

        # concurrent wait for completetion
        def jthrd(thread):
            thread.join()

        tjpx = ThreadPoolExecutor(max_workers=10000)
        for t in threads:
            tjpx.submit(jthrd, t)

        tjpx.shutdown(True)

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

        ope = []

        def ap_func(self, type, lst_ID: int):
            for i in range(0, len(ope)):
                for j in range(0, len(ope[i][lst_ID])):
                    ap_obj = [ope[i][0], ope[i][lst_ID][j], type]
                    self.rls.append(ap_obj)

        # get files
        for CurDir, Dir, Files in os.walk(path):
            CurDir = CurDir.replace("\\", "/")
            f_struct = [CurDir, Dir, Files]
            ope.append(f_struct[:])

        # filter
        def filt(self, lst):
            del_lst = []

            for j in lst:
                if self.Name not in j:
                    del_lst.append(j)

            for l in del_lst:
                lst.remove(l)
            del_lst.clear()

        tpx = ThreadPoolExecutor()
        for i in range(0, len(ope)):
            tpx.submit(filt, self, ope[i][1])
            tpx.submit(filt, self, ope[i][2])

        tpx.shutdown(wait=True)

        if self.searchType == "" or self.searchType == "f":
            ap_func(self, "file", 2)
        if self.searchType == "" or self.searchType == "d":
            ap_func(self, "folder", 1)
