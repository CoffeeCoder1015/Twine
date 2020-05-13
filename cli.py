commands = {

}

while True:
    cmd_in = input("Twine-CLI >>")

    try:
        print(commands[cmd_in])
        exec(commands[cmd_in])
    except:
        if cmd_in != "":
            print("COMMAND NOT FOUND")