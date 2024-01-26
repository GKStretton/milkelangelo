#include "cover_servo.h"
#include "../common/util.h"
#include "../config.h"
#include "../middleware/logger.h"
#include <Servo.h>
#include "../drivers/tof10120.h"

Servo coverServo;

static void attach() {
	if (!coverServo.attached()) {
		coverServo.attach(COVER_SERVO_PIN);
	}
}

static void detach() {
	coverServo.detach();
	digitalWrite(COVER_SERVO_PIN, HIGH);
}

void CoverServo_Init() {
	// servo jolts on power up if this is left alone, or LOW
	InitPin(COVER_SERVO_PIN, HIGH);
}

void CoverServo_Open() {
	attach();

	Logger::Info("Moving cover from closed to open");
	for (int i = COVER_SERVO_CLOSED_US; i <= COVER_SERVO_OPEN_US; i+=COVER_SERVO_RESOLUTION_US) {
		coverServo.writeMicroseconds(i);
		delay(20);
	}
	// extra time to ensure open, if it was stuck at all
	delay(1500);
	
	detach();
}

void CoverServo_Close() {
	attach();

	Logger::Info("Moving cover from open to closed");
	for (int i = COVER_SERVO_OPEN_US; i >= COVER_SERVO_CLOSED_US; i-=COVER_SERVO_RESOLUTION_US) {
		coverServo.writeMicroseconds(i);
		delay(20);
	}
	delay(500);

	detach();
}

void CoverServo_SetMicroseconds(int us) {
	attach();

	coverServo.writeMicroseconds(us);
	Logger::Info("Set cover servo to " + String(us) + "us.");
	delay(500);

	detach();
}

bool CoverServo_IsOpen() {
	SetDualRelay(TOF_POWER_PIN, true);
	delay(500);
	// ensure cover is open
	float dist = TOF_GetDistance();
	// turn off
	SetDualRelay(TOF_POWER_PIN, false);

	if (dist < 10) {
		// error reading cover position
		Logger::Error("cover tof invalid reading " + String(dist) + "mm.");
		return false;
	}
	if (dist > 25) {
		Logger::Warn("cover tof reading too high at " + String(dist) + "mm.");
		return false;
	}
	Logger::Info("cover tof valid open reading at " + String(dist) + "mm.");

	return true;
}