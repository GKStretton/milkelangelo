import serial

def handleInput(line: str):
    if line[0] == '>':
        data = line[1:].strip("\r\n\t ").split('\t')
        if len(data) == 2:
            print(data[0], ":", data[1])
        else:
            print("malformed data, expected 2 parameters after splitting on string:", line)
    elif line == "" or line == "\r\n" or line == "\n":
        return
    else:
        print(line)

if __name__ == "__main__":
    print("Attempting to open serial interface...")
    with serial.Serial('/dev/ttyACM0', 1000000, timeout=1) as ser:
        print("Opened.")

        while True:
            if ser.inWaiting() > 0:
                handleInput(ser.readline().decode("utf-8"))