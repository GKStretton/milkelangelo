#include "controller.h"
#include "../config.h"
#include "state_report.h"

// Manual mode has no local input of its own: axis speeds are set directly
// by the mqtt manual/<axis>-speed topics (see main.cpp's topicHandler), and
// persist via each stepper's setSpeed() until changed or overridden by
// autoUpdate's moveTo() calls. This just keeps steppers awake and reports mode.
void Controller::manualUpdate(State *s) {
	StateReport_SetMode(machine_Mode_MANUAL);
	digitalWrite(STEPPER_SLEEP, HIGH);
}
