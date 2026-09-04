#include "controller.h"
#include "../config.h"
#include "state_report.h"

unsigned long lastControlUpdate = millis();

void Controller::Init(State *s) {
}

void Controller::Update(State *s) {
	if (millis() - lastControlUpdate > CONTROLLER_UPDATE_PERIOD_MS)
	{
		lastControlUpdate = millis();

		if (s->manualRequested) {
			manualUpdate(s);
		} else {
			autoUpdate(s);
		}

		StateReport_Update(s);
	}
}
