#include <Arduino.h>
#include "tof10120.h";
#include "../config.h"
#include "../middleware/logger.h"
#include <Wire.h>

float TOF_GetDistance() {
	// https://surtrtech.com/2019/03/18/easy-use-of-tof-10120-laser-rangefinder-to-measure-distance-with-arduino-lcd/
	Wire.beginTransmission(TOF_I2C_ADDRESS);
	Wire.write(byte(TOF_I2C_DISTANCE_REGISTER));
	Wire.endTransmission();

	delay(1);

	// req 2 bytes
	Wire.requestFrom(TOF_I2C_ADDRESS, 2);
	if (Wire.available() >= 2) {
		uint8_t h = Wire.read();
		uint8_t l = Wire.read();
		return h << 8 | l;
	}

	// failed to read
	return -1.0f;
}