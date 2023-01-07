# This script takes 
# - input directory or file; batch mode / individual mode
# - output directory
# 
# If crop configs are present it will crop the input
# It will also do post-processing on the images required to get them to final
# standard.

from skimage import io, exposure
import matplotlib.pyplot as plt
import cv2 as cv
import numpy as np
import time
import argparse
import yaml

def process_image(raw, crop_config=None):
	# post = exposure.equalize_hist(raw)
	# post = exposure.equalize_adapthist(raw, clip_limit=0.01)

	post = raw
	print(crop_config)
	# crop
	# equalise
	# any other post
	return post

def handle_directory(in_path, out_path):
	# for each image in directory:
	# 	handle_image(path)
	return

def handle_image(in_file, out_dir):
	raw = load_image(in_file)
	crop_config = load_crop_config(args.input)

	post = process_image(raw, crop_config=crop_config)

	preview_image(raw, post)
	save_image(out_dir, post)

	return

def preview_image(raw, post):
	# print(raw.shape)
	# print(post.shape)
	# print(raw[0,0,0])
	# print(post[0,0,0])
	res = cv.vconcat([raw, post])
	cv.namedWindow("win", cv.WINDOW_NORMAL)
	cv.imshow("win", res)

	while cv.waitKey() != ord('x'):
		time.sleep(0.01)
	cv.destroyAllWindows()

def load_image(path):
	img = cv.imread(path)
	return img / 255.0

def save_image(path, image):
	return None

def load_crop_config(path):
	config = None
	try:
		with open(path + ".yml", 'r') as f:
			config = yaml.load(f)
	except FileNotFoundError as err:
		pass

	return config

if __name__ == "__main__":
	parser = argparse.ArgumentParser()

	parser.add_argument("-i", "--input", action="store", help="input file or directory to process")
	parser.add_argument("-o", "--output", action="store", help="output directory for processed images")
	parser.add_argument("-v", "--verbose", action="store_true", help="verbose output, print debug information")

	args = parser.parse_args()

	if False: #is_directory(args.input):
		handle_directory(args.input)
	else:
		handle_image(args.input, args.output)
	
	print("done")