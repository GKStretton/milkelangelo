#include <Arduino.h>
#include <Preferences.h>
#include "config.h"
#include "calibration.h"
#include "common/util.h"
#include "common/mathutil.h"
#include "common/ik_algorithm.h"
#include "middleware/logger.h"
#include "middleware/sleep.h"
#include "middleware/mqtt.h"
#include "app/state.h"
#include "app/navigation.h"
#include "app/controller.h"
#include "app/state_report.h"
#include "extras/topics_firmware/topics_firmware.h"

// Manual-jog axis speed topics. Not yet part of the shared asol-protos
// topics_firmware contract (see extras/topics_firmware/topics_firmware.h),
// so they're kept local here rather than hand-edited into that generated
// file. Follow the same "mega/req/..." legacy namespace so they're covered
// by the mega/req/# subscription in middleware/mqtt.cpp.
const char *TOPIC_MANUAL_RING_SPEED_SET = "mega/req/manual/ring-speed";
const char *TOPIC_MANUAL_Z_SPEED_SET = "mega/req/manual/z-speed";
const char *TOPIC_MANUAL_YAW_SPEED_SET = "mega/req/manual/yaw-speed";
const char *TOPIC_MANUAL_PITCH_SPEED_SET = "mega/req/manual/pitch-speed";
const char *TOPIC_MANUAL_PIPETTE_SPEED_SET = "mega/req/manual/pipette-speed";

State s = CreateStateObject();

Controller controller;

int updatesInLastSecond;
unsigned long lastUpdatesPerSecondTime = millis();

void initSteppers();
void runSteppers(State *s);
void topicHandler(String topic, String payload);

// incrementStartupCounter reads the startup counter from NVS, increments it,
// writes and prints it.
void incrementStartupCounter() {
	Preferences preferences;
	preferences.begin("outer", false);
	uint32_t counter = preferences.getUInt("startups", 0) + 1;
	preferences.putUInt("startups", counter);
	preferences.end();
	Logger::Info("Startup counter incremented to " + String(counter));
	s.startup_counter = counter;
}

void sleepHandler(Sleep::SleepStatus sleepStatus) {
	digitalWrite(STEPPER_SLEEP, LOW);

	if (Sleep::IsEStopActive()) {
		StateReport_SetStatus(machine_Status_E_STOP_ACTIVE);
		StateReport_Update(&s);
	}
}

void wakeHandler(Sleep::SleepStatus lastSleepStatus) {
	s.ClearState();
	incrementStartupCounter();
}

void setup()
{
	Serial.begin(9600);

	pinMode(E_STOP_PIN, INPUT);

	pinMode(PITCH_LIMIT_SWITCH, INPUT);
	pinMode(YAW_LIMIT_SWITCH, INPUT);
	pinMode(Z_LIMIT_SWITCH, INPUT);
	pinMode(RING_LIMIT_SWITCH, INPUT);
	pinMode(PIPETTE_LIMIT_SWITCH, INPUT);

	// make steppers sleep on start
	InitPin(STEPPER_SLEEP, LOW);

	initSteppers();

	// register callback, then connect. Logger/StateReport publish through
	// Mqtt, so nothing below this may call them until Mqtt::Init() has run.
	Mqtt::SetTopicHandler(topicHandler);
	Mqtt::Init();

	Logger::SetLevel(Logger::DEBUG);
	Logger::Info("setup start");

	Sleep::SetOnWakeHandler(wakeHandler);
	Sleep::SetOnSleepHandler(sleepHandler);

	// disabling so it doesn't always start after flash / reconnect
	// Sleep::Wake();

	controller.Init(&s);

	Logger::Info("Sending first state report");
	StateReport_SetStatus(machine_Status_SLEEPING);
	StateReport_Update(&s);
	Logger::Info("setup complete");
}

void initSteppers() {
	s.pitchStepper.setMaxSpeed(1250 * SPEED_MULT);
	s.pitchStepper.setAcceleration(1600 * SPEED_MULT);
	s.pitchStepper.setPinsInverted(true);
	s.pitchStepper.SetLimitSwitchPin(PITCH_LIMIT_SWITCH);
	s.pitchStepper.SetAtTargetUnitThreshold(0);

	s.yawStepper.setMaxSpeed(1250 * SPEED_MULT);
	s.yawStepper.setAcceleration(1600 * SPEED_MULT);
	s.yawStepper.setPinsInverted(true);
	s.yawStepper.SetLimitSwitchPin(YAW_LIMIT_SWITCH);
	s.yawStepper.SetAtTargetUnitThreshold(0);

	s.zStepper.setMaxSpeed(1250 * SPEED_MULT);
	s.zStepper.setAcceleration(800 * SPEED_MULT);
	s.zStepper.SetLimitSwitchPin(Z_LIMIT_SWITCH);
	s.zStepper.SetAtTargetUnitThreshold(0);

	s.ringStepper.setPinsInverted(true);
	s.ringStepper.setMaxSpeed(1250 * SPEED_MULT);
	s.ringStepper.setAcceleration(800 * SPEED_MULT);
	s.ringStepper.SetLimitSwitchPin(RING_LIMIT_SWITCH);

	s.pipetteStepper.setMaxSpeed(1250 * SPEED_MULT);
	s.pipetteStepper.setAcceleration(800 * SPEED_MULT);
	s.pipetteStepper.setPinsInverted(true);
	s.pipetteStepper.SetLimitSwitchPin(PIPETTE_LIMIT_SWITCH);
}

// unpackCommaSeparatedValues splits payload on ',' into up to n values.
void unpackCommaSeparatedValues(String payload, String values[], int n) {
	int value_index = 0;
	for (unsigned int i = 0; i < payload.length(); i++) {
		if (value_index >= n) return;

		if (payload[i] == ',') {
			values[++value_index] = "";
			continue;
		}
		values[value_index] += payload[i];
	}
}

void topicHandler(String topic, String payload)
{
	if (topic == TOPIC_WAKE)
	{
		Sleep::Wake();
		return;
	}
	else if (topic == TOPIC_STATE_REPORT_REQUEST)
	{
		StateReport_ForceSend();
	}
	if (Sleep::IsSleeping())
	{
		// if asleep, only listen for wake and state report
		return;
	}

	if (topic == TOPIC_SLEEP)
	{
		Sleep::Sleep(Sleep::UNKNOWN);
	}
	else if (topic == TOPIC_SHUTDOWN)
	{
		s.shutdownRequested = true;
	}
	else if (topic == TOPIC_UNCALIBRATE)
	{
		s.pitchStepper.MarkAsNotCalibrated();
		s.yawStepper.MarkAsNotCalibrated();
		s.zStepper.MarkAsNotCalibrated();
		s.ringStepper.MarkAsNotCalibrated();
		s.pipetteStepper.MarkAsNotCalibrated();
		s.calibrationCleared = true;
	}
	else if (topic == TOPIC_COLLECT) {
		String values[] = {"", ""};
		unpackCommaSeparatedValues(payload, values, 2);
		int vial = values[0].toInt();
		float ul = values[1].toFloat();

		if (!s.collectionRequest.requestCompleted) {
			Logger::Info("cannot collect because collection request " + String(s.collectionRequest.requestNumber) + " is still in progress");
		} else {
			s.collectionRequest.requestNumber++;
			s.collectionRequest.requestCompleted = false;
			s.collectionRequest.vialNumber = vial;
			s.collectionRequest.ulVolume = ul;
			Logger::Info("created collection request " + String(s.collectionRequest.requestNumber) + " for " + String(ul) + "ul of vial " + String(vial));
		}
	}
	else if (topic == TOPIC_DISPENSE) {
		String values[] = {""};
		unpackCommaSeparatedValues(payload, values, 1);
		float ul = values[0].toFloat();
		if (!s.pipetteState.spent) {
			s.pipetteState.dispenseRequested = true;
			s.pipetteState.ulVolumeHeldTarget -= ul;
			s.pipetteState.dispenseRequestNumber++;
			if (s.pipetteState.ulVolumeHeldTarget <= 0) {
				s.pipetteState.ulVolumeHeldTarget = 0;
			}
			Logger::Info("dispensed " + String(ul) + ", ulVolumeHeldTarget is now " + String(s.pipetteState.ulVolumeHeldTarget));
		} else {
			Logger::Info("Cannot dispense because already spent");
		}
	}
	else if (topic == TOPIC_GOTO_NODE)
	{
		long num = payload.toInt();
		s.forceIdleLocation = num == machine_Node_IDLE_LOCATION;
		s.SetGlobalNavigationTarget((machine_Node)num);
		Logger::Debug("Set globalTargetNode to " + String(num));
	}
	else if (topic == TOPIC_GOTO_XY) {
		String values[] = {"", ""};
		unpackCommaSeparatedValues(payload, values, 2);
		float target_x = values[0].toFloat();
		float target_y = values[1].toFloat();
		Logger::Info("recieved req for target_x, target_y to " + String(target_x) + ", " + String(target_y));

		boundXYToCircle(&target_x, &target_y, IK_TARGET_RADIUS_FRAC);
		Logger::Info("constrained target_x, target_y to " + String(target_x) + ", " + String(target_y));

		float ring, yaw;
		int code = getRingAndYawFromXY(target_x, target_y,
						s.ringStepper.PositionToUnit(s.ringStepper.currentPosition()),
						&ring, &yaw,
						s.ringStepper.GetMinUnit(), s.ringStepper.GetMaxUnit());

		if (code != 0) {
			Logger::Error("error code fromgetRingAndYawFromXY, aborting");
			return;
		}

		if (ring < s.ringStepper.GetMinUnit() || ring > s.ringStepper.GetMaxUnit()) {
			Logger::Error("Unexpected ring value " + String(ring) + " detected, aborting ik!");
			return;
		}
		boundToSignedMaximum(&yaw, MAX_BOWL_YAW);
		Logger::Info("Setting x,y, and target_ring=" + String(ring) + " and target_yaw=" + String(yaw));
		s.target_x = target_x;
		s.target_y = target_y;
		s.target_ring = ring;
		s.target_yaw = yaw;
	}
	else if (topic == TOPIC_TOGGLE_MANUAL)
	{
		s.manualRequested = !s.manualRequested;
		Logger::Info("Toggled manualRequested mode to " + String(s.manualRequested));
	}
	else if (topic == TOPIC_SET_IK_Z) {
		float z = payload.toFloat();
		if (z < MIN_BOWL_Z || z > s.zStepper.GetMaxUnit()) {
			Logger::Error("z level " + payload + " out of range.");
			return;
		}
		s.ik_target_z = z;
	}
	else if (topic == TOPIC_MARK_SAFE_TO_CALIBRATE) {
		s.overrideCalibrationBlock = true;
		Logger::Info("Set overrideCalibrationBlock true per mqtt request");
	}
	else if (topic == TOPIC_MAINTENANCE) {
		s.target_ring = MAINTENANCE_RING_ANGLE;
		s.forceIdleLocation = false;
		Navigation::SetGlobalNavigationTarget(&s, machine_Node_OUTER_HANDOVER);
	}
	else if (topic == TOPIC_GOTO_RING_IDLE_POS) {
		s.target_ring = IDLE_RING_ANGLE;
	}
	else if (topic == TOPIC_MANUAL_RING_SPEED_SET) {
		if (!s.manualRequested) {
			Logger::Warn("ignoring manual ring speed, not in manual mode");
			return;
		}
		s.ringStepper.setSpeed(payload.toFloat());
	}
	else if (topic == TOPIC_MANUAL_Z_SPEED_SET) {
		if (!s.manualRequested) {
			Logger::Warn("ignoring manual z speed, not in manual mode");
			return;
		}
		s.zStepper.setSpeed(payload.toFloat());
	}
	else if (topic == TOPIC_MANUAL_YAW_SPEED_SET) {
		if (!s.manualRequested) {
			Logger::Warn("ignoring manual yaw speed, not in manual mode");
			return;
		}
		s.yawStepper.setSpeed(payload.toFloat());
	}
	else if (topic == TOPIC_MANUAL_PITCH_SPEED_SET) {
		if (!s.manualRequested) {
			Logger::Warn("ignoring manual pitch speed, not in manual mode");
			return;
		}
		s.pitchStepper.setSpeed(payload.toFloat());
	}
	else if (topic == TOPIC_MANUAL_PIPETTE_SPEED_SET) {
		if (!s.manualRequested) {
			Logger::Warn("ignoring manual pipette speed, not in manual mode");
			return;
		}
		s.pipetteStepper.setSpeed(payload.toFloat());
	}
	else
	{
		Logger::Debug("no handler for " + topic + " (payload = " + payload + ")");
	}
}

void runSteppers(State *s)
{
	s->ringStepper.Update();
	s->pitchStepper.Update();
	s->yawStepper.Update();
	s->zStepper.Update();
	s->pipetteStepper.Update();
}

void loop()
{
	Mqtt::Update();

	Sleep::Update();
	if (Sleep::IsSleeping())
	{
		StateReport_Update(&s);
		delay(200);
		return;
	}

	controller.Update(&s);

	runSteppers(&s);

	updatesInLastSecond++;
	if (millis() - lastUpdatesPerSecondTime > 1000)
	{
		s.updatesPerSecond = updatesInLastSecond;
		updatesInLastSecond = 0;
		lastUpdatesPerSecondTime = millis();
	}
}
