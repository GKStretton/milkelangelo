# The automated tester should stay for testing. The ultimate goal is 
# human-centric control, so this should not become a major part of the project.

import pycommon.mqtt_client as mc
import time
import signal
import datetime

def p(msg, *args):
	prefix = "{}:".format(datetime.datetime.now().time())
	print(prefix, msg, *args)

def simulate_piece():
	p("Simulating piece...")
	p("Sending Wake")
	mc.wake()
	time.sleep(30)

	coords = [(-0.5, 0.5), (0.5, 0.5), (-0.5, -0.5), (0.5, -0.5)]
	dispense_amount = 50

	for vial in range(1, 7, 2):
		p("collecting from vial", vial)
		mc.collect(vial, len(coords) * dispense_amount)

		time.sleep(15)

		for i, coord in enumerate(coords):
			p("going to coord {}: ({}, {})".format(i, *coord))
			mc.goto_xy(*coord)
			time.sleep(10)
			p("requested dispense")
			mc.dispense(dispense_amount)
			time.sleep(5)
		
		p("sleeping before next vial")
		time.sleep(10)
	
	p("finished vials. waiting before sleep")
	time.sleep(10)
	p("shutting down")
	mc.shutdown()
	p("waiting after shutdown...")
	time.sleep(20)


def term_handler(signum, frame):
	p("term received, sending shutdown...")
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
