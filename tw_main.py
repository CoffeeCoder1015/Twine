import os
import time

from termcolor import colored

from director import cd
from scan_py import scan
import platform


class TWINE_interface:
    def __init__(self):
        self.console_txt = colored(os.getcwd(), "blue")+colored(" >>", "green")
        self.buffer_scan_lst = None

    def regen_console_txt(self):
        self.console_txt = colored(os.getcwd(), "blue")+colored(" >>", "green")

    def change_dir(self, dir):
        cd(path=dir)
        self.regen_console_txt()

    def list_dir(self):
        # file/folder selection setup
        ope = os.listdir(".")
        for i in range(0, len(ope)):
            if os.path.isdir(ope[i]) == True:
                ope[i] = [os.getcwd(), ope[i], "folder"]
            else:
                ope[i] = [os.getcwd(), ope[i], "file"]

        self.buffer_scan_lst = ope
        self.show()

    def show(self):
        if self.buffer_scan_lst != None:
            pad_l = len(str(len(self.buffer_scan_lst)))
            for i in range(0, len(self.buffer_scan_lst)):
                j = self.buffer_scan_lst[i]
                k_pad = format(i, "0%s" % (pad_l))
                ss = colored(str(k_pad)+" ", "blue") + \
                    colored(" | ", "green").join(j)
                print(ss)
        else:
            print((colored(
                "NOTHING HAS BEEN SCANNED\ntype scan/ls to load something into self.buffer_scan_lst", "red")))

    def Open(self, ID=None, inParentDir=None):
        sysVer = platform.system()
        openPath = False
        if ID == None:
            ID = int(input("ID of item:"))
        # type correction (ID)
        elif type(ID) == str:
            if not "\\" in str(ID) and not "/" in str(ID):
                ID = int(ID)
            else:
                openPath = True
        # type correction (inParentDir)
        if inParentDir == None or not "i" in str(inParentDir).lower():
            inParentDir = False
        elif "i" in str(inParentDir).lower():
            inParentDir = True

        fileExplorer = ""
        if sysVer == "Windows":
            fileExplorer = "explorer"
        elif sysVer == "Darwin":
            fileExplorer = "open"
        else:
            fileExplorer = "xdg-open"

        if openPath == False:
            #path string reformatting
            fps = ""
            if sysVer == "Windows":
                ope = self.buffer_scan_lst[ID]
                fps = list(ope[0].replace("/", "\\"))

                if fps[len(fps)-1] != "\\":
                    fps.extend("\\")
                if inParentDir == False:
                    fps.extend(ope[1])
                fps = "".join(fps)
            else:
                ope = self.buffer_scan_lst[ID]
                fps = list(ope[0])

                if fps[len(fps)-1] != "/":
                    fps.extend("/")
                if inParentDir == False:
                    fps.extend(ope[1])
                fps = "".join(fps)
            execc = f"{fileExplorer} {fps}".format()
        else:
            execc = f"{fileExplorer} {ID}".format()
        os.system(execc)

    def ar_input(self, ip_txt: str, bk_txt: str, err_msg: str, acpt_resp: list):
        # argument required input system
        while True:
            v_ip = input(ip_txt)
            checks = 0
            for i in acpt_resp:
                if v_ip == i:
                    checks += 1
            if checks > 0:
                print(bk_txt)
                return v_ip
            else:
                print(err_msg)

    def search(self):
        UNKNOWN_OPTION = colored("UNKNOWN OPTION", "red")
        UNKNOWN_DIRECTORY = colored("UNKNOWN DIRECTORY", "red")
        BREAKER = colored("-------------------------------", "green")

        # NAME
        print("<leave empty> = scan for all")
        Name = input("Search name:")
        print(BREAKER)

        # MATCH MODE
        print("c = contains the name, <leave empty> = exact match")
        matchMode = self.ar_input(
            "Match mode (c,<leave empty>):", BREAKER, UNKNOWN_OPTION, ["c", ""])

        # SEARCH TYPE
        print("f = file, d = folder, <leave empty> = file and folder")
        searchType = self.ar_input(
            "Search type (f,d,<leave empty>):", BREAKER, UNKNOWN_OPTION, ["f", "d", ""])

        # PATH
        print("<leave empty> = currenct directory")
        while True:
            path = input("Path:")
            if path == "":
                path = "."
            if os.path.isdir(path) == True:
                break
            else:
                print(UNKNOWN_DIRECTORY)
        print(BREAKER)

        SCAN = scan(Name=Name, matchMode=matchMode, searchType=searchType)
        SCAN.scan_mc(path=path)
        self.buffer_scan_lst = SCAN.rls
        self.show()

    def Help(self):
        htxt = (
            """{tmargin}
cd    {b}  change directory
ls    {b}  list currenct directory
show  {b}  show whats scanned
open  {b}  open selected file or folder
search{b}  search for file or folder
help  {b}  show this text
<KeyboardInterrupt>{b}  exit program
{tmargin}"""
        ).format(b=colored(" | ", "green"), tmargin=colored("----------------------------", "blue"))
        print(htxt)

    def shell_sesh(self):
        # mapting
        cmd_lst = {
            "cd": self.change_dir,
            "ls": self.list_dir,
            "show": self.show,
            "open": self.Open,
            "search": self.search,
            "help": self.Help

        }

        cmd_map = list(cmd_lst.keys())

        # execution
        while True:
            try:
                cmd_in = input(self.console_txt)
                cmd_in = [str(i).lower() for i in cmd_in.split(" ")]

                i = 0
                for _ in range(0, len(cmd_in)):
                    if cmd_in[i] == "":
                        del cmd_in[i]
                    else:
                        i += 1

                if cmd_in.count("") != len(cmd_in):
                    if cmd_in[0] in cmd_map:

                        try:

                            if cmd_in[0] == "cd":
                                cmd_lst[cmd_in[0]](" ".join(cmd_in[1:]))
                            elif cmd_in[0] == "open" and 2 <= len(cmd_in) <= 3:
                                cmd_lst[cmd_in[0]](*cmd_in[1:])
                            else:
                                cmd_lst[cmd_in[0]]()

                        except Exception as e:
                            print(colored("COMMAND ERROR", "red"))
                            print(e)

                    else:
                        print(colored("UNKNOWN COMMAND", "red"))

            except KeyboardInterrupt:
                exit()


intf = TWINE_interface()
intf.shell_sesh()
