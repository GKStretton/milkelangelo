import numpy as np
import cv2 as cv

TOP_MASK = "resources/static_img/top-mask.png"
FRONT_MASK = "resources/static_img/front-mask.png"
top_mask = None
front_mask = None

# add_overlay adds the mask overlay with dims magXmag to an image at x,y
# is mag is unspecified, it will default to the width of the input image
def add_overlay(img: np.ndarray, x=0, y=0, mag=0):
    global top_mask
    if top_mask is None:
        top_mask = cv.imread(TOP_MASK)

    if mag == 0:
        mag = img.shape[0]

    resized_mask = cv.resize(top_mask, (mag, mag))
    overlay_image_alpha(img, np.zeros((mag, mag, 3)), x, y, 1 - resized_mask[:,:,0] / 255.0)

def add_feather(img: np.ndarray):
    global front_mask
    if front_mask is None:
        front_mask = cv.imread(FRONT_MASK)
    
    w, h = img.shape[1], img.shape[0]
    resized_mask = cv.resize(front_mask, (w, h))
    overlay_image_alpha(img, np.zeros((h, w, 3)), 0, 0, 1 - resized_mask[:,:,0] / 255.0)

# https://stackoverflow.com/a/45118011
def overlay_image_alpha(img, img_overlay, x, y, alpha_mask):
    """Overlay `img_overlay` onto `img` at (x, y) and blend using `alpha_mask`.

    `alpha_mask` must have same HxW as `img_overlay` and values in range [0, 1].
    """
    # Image ranges
    y1, y2 = max(0, y), min(img.shape[0], y + img_overlay.shape[0])
    x1, x2 = max(0, x), min(img.shape[1], x + img_overlay.shape[1])

    # Overlay ranges
    y1o, y2o = max(0, -y), min(img_overlay.shape[0], img.shape[0] - y)
    x1o, x2o = max(0, -x), min(img_overlay.shape[1], img.shape[1] - x)

    # Exit if nothing to do
    if y1 >= y2 or x1 >= x2 or y1o >= y2o or x1o >= x2o:
        return

    # Blend overlay within the determined ranges
    img_crop = img[y1:y2, x1:x2]
    img_overlay_crop = img_overlay[y1o:y2o, x1o:x2o]
    alpha = alpha_mask[y1o:y2o, x1o:x2o, np.newaxis]
    alpha_inv = 1.0 - alpha

    img_crop[:] = alpha * img_overlay_crop + alpha_inv * img_crop