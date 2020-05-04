import json
import os

class cache_machine:
    def __init__(self,name,Type,start_pos):
        self.start_pos = start_pos
        self.type = Type
        self.cfn = "_config_atm_cacher.py.json"
        self.set_template = json.dumps({
            "path":f".\\cache_atm_{name}".format()
        },indent=4)
        self.init_env()

    #internal function
    def set_error_checker(self):
        with open(self.start_pos+"\\"+self.cfn,'w') as fIO:
                fIO.write(self.set_template)
    
    #internal function
    def init_env(self):
        fconfig = os.listdir(self.start_pos)
        f_ex = fconfig.count(self.cfn)
        if f_ex == 0:
            self.set_error_checker()
        if f_ex == 1:
            with open(self.start_pos+"\\"+self.cfn,'r') as fIO:
                settings=json.loads(fIO.read())
            try:
                PATH = self.start_pos+"\\"+settings['path']
                self.PATH = PATH
                don = os.path.isdir(PATH)
                if don == False:
                    os.mkdir(PATH)
            except:
                self.set_error_checker()

    #internal function
    def cache_type_ckr(self):
        if self.type == "b" and self.cacheTarget == None or self.cacheTarget == "":
            raise Exception("""
-----------------------
supply cacheTarget argument with target file
when 'Type' is set to 'b'
-----------------------""")

    #internal function
    def che_init(self):
        with open(self.cache_file,'w') as fIO:
            fIO.write("{\n}")

    #internal function (error handeling)
    def type_b_cache_fixer(self,file):
        with open(file,"r+") as fIO:
            data = fIO.readlines()

            if data[0] == "{\n" and data[len(data)-1] == "}" and "".join(data).count("{")+"".join(data).count("}") == 2:
                return True
            data = list("".join(data))
            while data.count("}") != 0:
                del data[len(data)-1]
            data.extend("}")
            data = "".join(data)
            fIO.seek(0)
            fIO.truncate(0)
            fIO.write(data)
            

    def deposit(self,data,name,cacheTarget=None):
        r"""The function overwrites data, please fetch previous data
        and add new data onto it if you want to append data
        """
        self.cacheTarget = cacheTarget
        self.cache_type_ckr()
        data = list(bytes(data,'utf-8').hex())
        I= 1
        for i in range(0,len(data)):
            if i/I == 50:
                data.insert(i,"\n")
                I+=1
        
        data="".join(data)
        if self.type == "f":
            with open(self.PATH+"\\"+name+".atmc",'w')as fIO:
                fIO.write(data)
        
        if self.type == "b":
            self.cache_file = self.PATH+"\\"+cacheTarget+".json"

            #error handeling
            cft = os.path.isfile(self.cache_file)
            if cft == False:
                self.che_init()
            else:
                fs = os.path.getsize(self.cache_file)
                if fs == 0:
                    self.che_init()


            self.type_b_cache_fixer(self.cache_file)

            with open(self.cache_file,'r+') as fIO:
                raw_dat = fIO.read()
                cache_data = json.loads(raw_dat)
                cache_data[name] = data.replace("\n","")
                cache_data = json.dumps(cache_data,indent=4)
                fIO.seek(0)
                fIO.writelines(cache_data)

            self.type_b_cache_fixer(self.cache_file)

    #internal function
    def decode(self,raw_dat):
        byte_dat = bytes.fromhex(raw_dat)
        str_dat = str(byte_dat.decode('utf-8'))
        return str_dat

    def withdraw(self,name,cacheTarget=None):
        self.cache_type_ckr()
        if self.type == "f":
            with open(self.PATH+"\\"+name+".atmc",'r') as fIO:
                raw_dat = fIO.read()
                return self.decode(raw_dat)

        if self.type == "b":
            self.cache_file = self.PATH+"\\"+cacheTarget+".json"
            self.type_b_cache_fixer(self.cache_file)
            with open(self.cache_file,'r+') as fIO:
                raw_dat = fIO.read()
                cache_data = json.loads(raw_dat)
                return self.decode(cache_data[name])