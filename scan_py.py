import os
import threading
from concurrent.futures import ThreadPoolExecutor
from director import cd

class scan:
    def __init__(self, Name='', matchMode='', searchType=''):
        self.Name = Name
        self.matchMode = matchMode

        if matchMode == "c":
            self.Name = self.Name.lower()
        self.searchType = searchType.lower()
        self.rls = []

    def mm_switch(self, string):
        if self.matchMode == "c":
            return string.lower()
        elif self.matchMode == "":
            return string

    def scan_mc(self, path="."):
        if path == ".":
            path = os.getcwd()

        cd(path=path)
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

        CurDir = os.getcwd().replace("\\", "/").replace("//", "/")+"/"
        # un-touched folder&file appending
        if self.searchType == "" or self.searchType == 'd':
            for i in ags_lst:
                if self.Name in self.mm_switch(i):
                    self.rls.append([CurDir, i, "folder"])

        if self.searchType == "" or self.searchType == 'f':
            for i in files:
                if self.Name in self.mm_switch(i):
                    self.rls.append([CurDir, i, "file"])

        length = len(ags_lst)
        tpx = ThreadPoolExecutor(length)
        # scanner execution & processing organization
        self.current = -1
        for i, v in enumerate(ags_lst):
            tpx.submit(self.scanner, i, v)

        for i in range(0, length):
            self.current = i
            while self.current != -1:
                pass
            print("thread", i, "completed")

    # internal function
    def scanner(self, idn, path=None):
        self.curfPath = os.getcwd()
        self.working = True
        if path == "" or path == None or path == 0:
            path = self.curfPath
        else:
            path = self.curfPath+"\\"+path
        path = path.replace("\\", "/")

        ope = []

        # get files
        for CurDir, Dir, Files in os.walk(path):
            CurDir = CurDir.replace("\\", "/").replace("//", "/")+"/"
            if self.searchType == "" or self.searchType == "f":
                ope += [[CurDir, i, "file"]
                        for i in Files if self.Name in self.mm_switch(i)]
            if self.searchType == "" or self.searchType == "d":
                ope += [[CurDir, i, "folder"]
                        for i in Dir if self.Name in self.mm_switch(i)]

        while True:
            if idn == self.current:
                self.rls += ope
                self.current = -1
                break
