# The crop tool lets you configure the crop on the top webcam. Its output is
# used by the cropper in rtsp.

import pycommon.window as window
import cv2
import pycommon.image as image
import numpy as np
from pycommon.config_manager_client import *

HOST = "DEPTH"

# location to load dslr image from if doing that crop
# note, this uses the mount-md0 script to mount remote filesystem
DSLR_CAPTURE_LOCATION = "/mnt/md0/light-stores/dslr-preview.jpg"

TOP_CAM_CHOICE = "1"
FRONT_CAM_CHOICE = "2"
DSLR_CHOICE = "3"

def name_from_choice(choice):
    if choice == TOP_CAM_CHOICE:
        return "crop_top-cam"
    if choice == FRONT_CAM_CHOICE:
        return "crop_front-cam"
    if choice == DSLR_CHOICE:
        return "crop_dslr"

def write_yaml(choice, yml):
    name = name_from_choice(choice)
    write_remote_crop_config(name, yml)

def load_yaml(choice):
    name = name_from_choice(choice)
    return read_remote_crop_config(name)


class CropWindow(window.Window):
    def __init__(self, choice):
        super().__init__()

        self.choice = choice

        current_yml = load_yaml(choice)
        if current_yml:
            self.load_config(current_yml)
        else:
            self.x1 = 0
            self.y1 = 0
            self.x2 = 100
            self.y2 = 100
        
        print("loaded (x1, y1); (x2, y2) as ({}, {}); ({}, {})".format(self.x1, self.y1, self.x2, self.y2))

        if choice == TOP_CAM_CHOICE or choice == FRONT_CAM_CHOICE:
            self.open_stream()
        else:
            self.load_dslr_image()
    
    def load_dslr_image(self):
        # request dslr capture to be taken
        succ = trigger_dslr_capture()
        if not succ:
            print("failed to do dslr capture, reading anyway")
        self.dslr_capture = cv2.imread(DSLR_CAPTURE_LOCATION, cv2.IMREAD_UNCHANGED)
        self.dslr_capture = cv2.rotate(self.dslr_capture, cv2.ROTATE_180)
        # self.dslr_capture = np.zeros((1000, 1000, 3))
    
    def open_stream(self):
        print("Opening stream...")
        endpoint = "top-cam" if self.choice == TOP_CAM_CHOICE else "front-cam"
        self.stream = cv2.VideoCapture("rtsp://{}:8554/{}".format(HOST, endpoint))#, cv2.CAP_GSTREAMER)
        self.stream.set(cv2.CAP_PROP_BUFFERSIZE, 2)

        if not self.stream.isOpened():
            print("Error loading webcam stream, aborting.")
            self.exit()
        print("Stream opened")
    
    def load_config(self, current_yml):
        self.x1 = current_yml['left_abs']
        self.y1 = current_yml['top_abs']
        self.x2 = current_yml['right_abs']
        self.y2 = current_yml['bottom_abs']

        print("loaded values from yaml:", current_yml)
        
    def mouse_handler(self, event, x, y, flags, param):
        super().mouse_handler(event, x, y, flags, param)
        if flags & cv2.EVENT_FLAG_LBUTTON and flags & cv2.EVENT_FLAG_SHIFTKEY:
            self.x2 = x
            if self.choice != FRONT_CAM_CHOICE:
                mag = abs(self.x2 - self.x1)
                self.y2 = self.y1 + mag
            else:
                self.y2 = y
        elif flags & cv2.EVENT_FLAG_LBUTTON:
            self.x1 = x
            self.y1 = y

            if self.choice != FRONT_CAM_CHOICE:
                mag = abs(self.x2 - self.x1)
                self.y2 = self.y1 + mag
        
        # constrain to image
        self.x1 = max(0, self.x1)
        self.y1 = max(0, self.y1)
        self.x2 = min(self.frame_width, self.x2)
        self.y2 = min(self.frame_height, self.y2)
    
    def keyboard_handler(self, key):
        super().keyboard_handler(key)
    
    def update(self):
        if self.choice == TOP_CAM_CHOICE or self.choice == FRONT_CAM_CHOICE:
            ret, frame = self.stream.read()
            if ret == False:
                print("Frame empty, quitting")
                self.exit()
        else:
            frame = self.dslr_capture.copy()

        self.frame_width=frame.shape[1]
        self.frame_height=frame.shape[0]
        mag = abs(self.x2 - self.x1)

        # draw crop location
        cv2.rectangle(frame,(self.x1,self.y1), (self.x2, self.y2), (0,0,255),2, cv2.LINE_AA)

        res = frame.copy()
        # mask with the vig
        if self.choice != FRONT_CAM_CHOICE:
            image.add_overlay(res, self.x1, self.y1, mag)

        return res

if __name__ == "__main__":
    print(TOP_CAM_CHOICE, "- top-cam")
    print(FRONT_CAM_CHOICE, "- front-cam")
    print(DSLR_CHOICE, "- dslr")
    choice = input("Enter choice 1-3: ")
    if choice != TOP_CAM_CHOICE and choice != FRONT_CAM_CHOICE and choice != DSLR_CHOICE:
        print("invalid choice:", choice)
        exit()

    win = CropWindow(choice)
    win.loop()

    result = {
        "left_abs": win.x1,
        "right_abs": win.x2,
        "top_abs": win.y1,
        "bottom_abs": win.y2,
        "left_rel": win.x1,
        "right_rel": win.frame_width - win.x2,
        "top_rel": win.y1,
        "bottom_rel": win.frame_height - win.y2,
    }

    print(result)

    write_yaml(choice, result)