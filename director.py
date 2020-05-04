import os

class cd:
    start_pos=os.getcwd()
    def __init__(self,path="\\N/A\\",top=False):
        self.start_pos = cd.start_pos
        if top == False and path != "\\N/A\\":
            os.chdir(path)
        if top == True:
            self.top()

    def top(self):
        SourceDir = []
        curDir = list(os.getcwd())
        for i in range(0,len(curDir)):
            if curDir[i] == "\\":
                SourceDir.extend(curDir[i])
                break
            SourceDir.extend(curDir[i])
        SourceDir="".join(SourceDir)
        os.chdir(SourceDir)

    def ret_start(self):
        os.chdir(self.start_pos)