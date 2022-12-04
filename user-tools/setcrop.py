# The crop tool lets you configure the crop on the top webcam. Its output is
# used by the cropper in rtsp.

import pycommon.window as window
import cv2
from argparse import ArgumentParser
import pycommon.image as image
import yaml
import numpy as np
import sys
import os
from pycommon.config_manager_client import read_remote_crop_config
from pycommon.config_manager_client import write_remote_crop_config

TOP_MASK = "resources/static_img/top-mask.png"
HOST = "DEPTH"

def write_yaml(yml):
    local_file = ""#input("write to local yaml? enter path if so:")
    if local_file != "":
        with open(local_file, 'w') as f:
            yaml.dump(result, f)
    else:
        write_remote_crop_config(yml)

def load_yaml():
    local_file = ""#input("read from local yaml? enter path if so:")
    local = False
    if local_file != "":
        local = True
        print("local mode")

    # LOAD EXISTING CONFIG
    if local:
        print("loading from", local_file, "...")
        if os.path.isfile(local_file):
            with open(local_file, 'r') as f:
                yml = yaml.load(f)
                return yml
        else:
            print("file not found, proceeding with 0 values")
    else:
        return read_remote_crop_config()


class CropWindow(window.Window):
    def __init__(self):
        super().__init__()

        # load mask
        self.mask = cv2.imread(TOP_MASK)

        current_yml = load_yaml()
        if current_yml:
            self.load_config(current_yml)
        else:
            self.x1 = 0
            self.y1 = 0
            self.x2 = 100
            self.y2 = 100
        
        print("loaded (x1, y1); (x2, y2) as ({}, {}); ({}, {})".format(self.x1, self.y1, self.x2, self.y2))

        print("Opening stream...")
        self.stream = cv2.VideoCapture("rtsp://{}:8554/top-cam".format(HOST))#, cv2.CAP_GSTREAMER)
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
            self.y2 = y
        elif flags & cv2.EVENT_FLAG_LBUTTON:
            self.x1 = x
            self.y1 = y
    
    def keyboard_handler(self, key):
        super().keyboard_handler(key)

        if key != -1:
            print(key)
    
    def update(self):
        ret, frame = self.stream.read()
        if ret == False:
            print("Frame empty, quitting")
            self.exit()
        self.frame_width=frame.shape[1]
        self.frame_height=frame.shape[0]
        mag = abs(self.x2 - self.x1)

        # draw crop location
        cv2.rectangle(frame,(self.x1,self.y1), (self.x2, self.y1 + mag), (0,0,255),2, cv2.LINE_AA)

        # mask with the vig
        resized_mask = cv2.resize(self.mask, (mag, mag))
        res = frame.copy()
        image.overlay_image_alpha(res, np.zeros((mag, mag, 3)), self.x1, self.y1, 1 - resized_mask[:,:,0] / 255.0)

        return res
    

if __name__ == "__main__":
    win = CropWindow()
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

    write_yaml(result)