# The crop tool lets you configure the crop on the top webcam. Its output is
# used by the cropper in rtsp.

import pycommon.window as window
import cv2
import pycommon.image as image
import numpy as np
from pycommon.config_manager_client import *

HOST = "localhost"

TOP_CAM_CHOICE = "1"
FRONT_CAM_CHOICE = "2"

def name_from_choice(choice):
    if choice == TOP_CAM_CHOICE:
        return "crop_top-cam"
    if choice == FRONT_CAM_CHOICE:
        return "crop_front-cam"

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
        self.show_overlay = True

        current_yml = load_yaml(choice)
        if current_yml:
            self.load_config(current_yml)
        else:
            self.x1 = 0
            self.y1 = 0
            self.x2 = 100
            self.y2 = 100
        
        print("loaded (x1, y1); (x2, y2) as ({}, {}); ({}, {})".format(self.x1, self.y1, self.x2, self.y2))

        self.open_stream()
    
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

        # dimensions must be halvable for x264
        if (self.x2 - self.x1) % 2 == 1:
            self.x2-=1
        if (self.y2 - self.y1) % 2 == 1:
            self.y2-=1
    
    def keyboard_handler(self, key):
        super().keyboard_handler(key)
        if key == ord('o'):
            self.show_overlay = not self.show_overlay

    
    def update(self):
        ret, frame = self.stream.read()
        if ret == False:
            print("Frame empty, quitting")
            self.exit()

        self.frame_width=frame.shape[1]
        self.frame_height=frame.shape[0]
        mag = abs(self.x2 - self.x1)

        # draw crop location
        cv2.rectangle(frame,(self.x1,self.y1), (self.x2, self.y2), (0,0,255),2, cv2.LINE_AA)

        res = frame.copy()
        # mask with the vig
        if self.choice != FRONT_CAM_CHOICE and self.show_overlay:
            image.add_overlay(res, self.x1, self.y1, mag)

        return res

if __name__ == "__main__":
    print(TOP_CAM_CHOICE, "- top-cam")
    print(FRONT_CAM_CHOICE, "- front-cam")
    choice = input("Enter choice 1-2: ")
    if choice != TOP_CAM_CHOICE and choice != FRONT_CAM_CHOICE:
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