# The automated tester should stay for testing. The ultimate goal is 
# human-centric control, so this should not become a major part of the project.

import pycommon.mqtt_client as mc
import time
import signal


def simulate_piece():
	print("Simulating piece...")
	print("Sending Wake")
	mc.wake()
	time.sleep(5)

	coords = [(-0.5, 0.5), (0.5, 0.5), (-0.5, -0.5), (0.5, -0.5)]
	dispense_amount = 50

	for vial in range(1, 7, 2):
		print("collecting from vial", vial)
		mc.collect(vial, len(coords) * dispense_amount)

		time.sleep(10)

		for i, coord in enumerate(coords):
			print("going to coord {}: ({}, {})".format(i, *coord))
			mc.goto_xy(*coord)
			time.sleep(3)
			print("requested dispense")
			mc.dispense(dispense_amount)
			time.sleep(10)
		
		print("sleeping before next vial")
		time.sleep(10)
	
	print("finished vials. waiting before sleep")
	time.sleep(10)
	print("shutting down")
	mc.shutdown()
	print("waiting after shutdown...")
	time.sleep(20)


def term_handler(signum, frame):
	print("term received, sending shutdown...")
	mc.shutdown()
	time.sleep(3)
	exit()

if __name__ == "__main__":
	signal.signal(signal.SIGTERM, term_handler)
	signal.signal(signal.SIGINT, term_handler)

	mc.connect()
	time.sleep(1)

	while True:
		simulate_piece()
