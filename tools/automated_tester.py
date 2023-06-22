# The automated tester should stay for testing. The ultimate goal is
# human-centric control, so this should not become a major part of the project.

import pycommon.mqtt_client as mc
import time
import signal
import datetime
import argparse


def p(msg, *args):
    prefix = "{}:".format(datetime.datetime.now().time())
    print(prefix, msg, *args)


def simulate_piece(record=False):
    p("Simulating piece...")
    if record:
        p("Starting session")
        mc.begin_session()
    p("Sending Wake")
    mc.wake()
    time.sleep(30)

    coords = [(-0.5, 0.5), (0.5, 0.5), (-0.5, -0.5), (0.5, -0.5)]
    dispense_amount = 50

    for vial in range(1, 7, 2):
        if record:
            p("resuming session")
            mc.resume_session()

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
        time.sleep(5)
        if record:
            p("pausing session")
            mc.pause_session()
        time.sleep(10)

    p("finished vials. waiting before sleep")
    time.sleep(10)

    if record:
        p("resuming session")
        mc.resume_session()

    p("shutting down")
    mc.shutdown()
    time.sleep(2)
    p("waiting after shutdown...")
    time.sleep(20)
    if record:
        p("Ending session")
        mc.end_session()
        time.sleep(5)


def term_handler(signum, frame):
    p("term received, sending shutdown...")
    mc.shutdown()
    time.sleep(3)
    exit()


if __name__ == "__main__":
    signal.signal(signal.SIGTERM, term_handler)
    signal.signal(signal.SIGINT, term_handler)

    parser = argparse.ArgumentParser()
    parser.add_argument("-r", "--record", action="store_true",
                        help="if true, session request will be emitted to trigger recording")
    args = parser.parse_args()

    mc.connect()
    time.sleep(1)

    simulate_piece(record=args.record)

    print("Exiting.")
