
from getch import getch
from termcolor import colored

def getTrueFalse(msg, default=False):
    defaultChoice = "T/f" if default else "t/F"
    print("{} [{}]: ".format(msg, defaultChoice), end='', flush=True)

    textChoice = getch().lower()
    choice = default
    if textChoice == 't':
        choice = True
    elif textChoice == 'f':
        choice = False

    print()
    return choice
    


functions = {}


def helpText(interface):
    print("Help:")
    for _, (key, [_, helptext]) in enumerate(functions.items()):
        print("\t{}: {}".format(key, helptext))


def camera(interface):
    choice = getTrueFalse("Camera on?")
    print("Setting camera state to", choice, "...")
    interface.SetCameraState(choice)
    print("Done.")

def water(interface):
    choice = getTrueFalse("Water flow?")
    print("Setting water flow to", choice, "...")
    interface.SetFlushFlow(choice)
    print("Done.")

def milk(interface):
    choice = getTrueFalse("Milk flow?")
    print("Setting milk flow to", choice, "...")
    interface.SetCanvasFlow(choice)
    print("Done.")

def air(interface):
    choice = getTrueFalse("Air flow?")
    print("Setting air flow to", choice, "...")
    interface.SetAirFlow(choice)
    print("Done.")

def registerFunctions():
    functions['h'] = [helpText, "Display help text"]
    functions['c'] = [camera, "Set the camera state"]
    functions['w'] = [water, "Set dispensing of " + colored("water", 'blue', attrs=['bold'])]
    functions['m'] = [milk, "Set dispensing of " + colored("milk", attrs=['bold'])]
    functions['a'] = [air, "Set dispensing of " + colored("air", 'grey', attrs=['bold'])]



if __name__ == "__main__":
    registerFunctions()
    print("Attempting to open serial interface...")
    with serial.Serial('/dev/ttyACM0', 1000000, timeout=1) as ser:
        print("Opened.")

        while True:
            print(colored("\n> ", 'red'), end='')
            key = getch()
            print(key)
            if key == 'q':
                break
            if key in functions:
                functions[key][0](interface)
            else:
                print("Invalid command '{}'. Printing help:".format(key))
                helpText(interface)


    print("done")