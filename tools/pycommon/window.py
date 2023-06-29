import cv2
import time


class Window:
    def __init__(self, name="window"):
        self.name = name
        self.exiting = False
        cv2.namedWindow(self.name, cv2.WINDOW_NORMAL)
        cv2.setMouseCallback(self.name, self.mouse_handler)

        self.placeholder_image = cv2.imread("../resources/static_img/fallback.png", cv2.IMREAD_UNCHANGED)

        print("Window created")

    def exit(self):
        self.exiting = True

    def mouse_handler(self, event, x, y, flags, param):
        pass

    def keyboard_handler(self, key):
        if key == 27:
            self.exiting = True

    # update returns the frame to be drawn
    def update(self) -> any:
        time.sleep(0.032)
        return self.placeholder_image

    def __draw(self):
        # get frame
        frame = self.update()

        # show frame
        cv2.imshow(self.name, frame)

        # wait / process keyboard input
        key = cv2.waitKey(1)
        self.keyboard_handler(key)

    def loop(self):
        while not self.exiting:
            self.__draw()
        print("Exiting...")
