# This script takes 
# - input directory or file; batch mode / individual mode
# - output directory
# 
# If crop configs are present it will crop the input
# It will also do post-processing on the images required to get them to final
# standard.

from skimage import io, exposure
import skimage.color as color
import matplotlib.pyplot as plt
import cv2 as cv
import pycommon.image as image
import numpy as np
import time
import argparse
import yaml


def process_image(raw, crop_config=None):
	cropped = raw.copy()
	# crop
	if crop_config is not None:
		x1 = crop_config['left_abs']
		x2 = crop_config['right_abs']
		y1 = crop_config['top_abs']
		y2 = crop_config['bottom_abs']
		cropped = cropped[y1:y2,x1:x2,:]

		image.add_overlay(cropped)

	cropped = cv.rotate(cropped, cv.ROTATE_180)
	post = cropped.copy()

	# post = exposure.adjust_gamma(post, 0.5)
	# equalise
	# post = exposure.equalize_hist(post)
	# post = exposure.equalize_adapthist(post, clip_limit=0.01)


	saturation_factor=1.3

	hsv = color.rgb2hsv(post)
	hsv[:,:,1] = hsv[:,:,1] * saturation_factor
	post = color.hsv2rgb(hsv)
	

	# any other post

	return cropped, post

def handle_directory(in_path, out_path):
	# for each image in directory:
	# 	handle_image(path)
	return

def handle_image(in_file, out_dir):
	raw = load_image(in_file)
	crop_config = load_crop_config(args.input)

	cropped, post = process_image(raw, crop_config=crop_config)

	preview_image(raw, cropped, post)
	save_image(out_dir, post)

	return

def preview_image(raw, cropped, post):
	# print(raw[0,0,0])
	# print(post[0,0,0])
	after = cv.hconcat([cropped, post])
	print(raw.shape)
	print(after.shape)
	if raw.shape[1] > after.shape[1]:
		after = np.pad(after, ((0,0),(0,raw.shape[1]-after.shape[1]), (0,0)))
	else:
		raw = np.pad(raw, ((0,0),(0,after.shape[1]-raw.shape[1]), (0,0)))
	print(raw.shape)
	print(after.shape)
	res = cv.vconcat([raw, after])
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