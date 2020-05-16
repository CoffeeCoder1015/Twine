commands = {

}

while True:
    cmd_in = input("Twine-CLI >>")

    err_han_lst = list(cmd_in)
    if len(err_han_lst) != 0:
        try:
            print(commands[cmd_in])
            exec(commands[cmd_in])
        except:
            if err_han_lst.count(" ") != len(err_han_lst):
                print("COMMAND NOT FOUND")