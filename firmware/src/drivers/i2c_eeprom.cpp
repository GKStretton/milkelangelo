#include "i2c_eeprom.h"
#include "../config.h"
#include "../middleware/logger.h"
#include "../common/util.h"
#include <Wire.h>

typedef unsigned char byte;

static void hack(bool on) {
	// not sure but it seems like this needs to be on for the eeprom to
	// work, (the tof sensor power). some kind of I2C interference?
	if (on) {
		SetDualRelay(TOF_POWER_PIN, true);
		delay(200);
	} else {
		SetDualRelay(TOF_POWER_PIN, false);
	}
}

void I2C_EEPROM::WriteByte(int addr, uint8_t data) {
	hack(true);
	Logger::Info("Sending to eeprom for WRITE. (If hanging, check voltages)");
	Wire.beginTransmission(EEPROM_I2C_ADDRESS);
	Wire.write((int)(addr >> 8)); // MSB
	Wire.write((int)(addr & 0xFF)); // LSB
	Wire.write(data);
	Wire.endTransmission();
	hack(false);
}

uint8_t I2C_EEPROM::ReadByte(int addr) {
	hack(true);
	uint8_t readData = 0xFF;

	Logger::Info("Sending to eeprom for READ. (If hanging, check voltages)");
	Wire.beginTransmission(EEPROM_I2C_ADDRESS);
	Wire.write((int)(addr >> 8)); // MSB
	Wire.write((int)(addr & 0xFF)); // LSB
	Wire.endTransmission();

	Wire.requestFrom(EEPROM_I2C_ADDRESS, 1);
	if (Wire.available()) {
		readData = Wire.read();
	}
	return readData;
	hack(false);
}

void I2C_EEPROM::WriteFloat(int addr, float data) {
	hack(true);
	byte* bytes = (byte*) &data;

	Logger::Debug("Writing float to eeprom");
	Wire.beginTransmission(EEPROM_I2C_ADDRESS);
	Wire.write((int)(addr >> 8)); // MSB
	Wire.write((int)(addr & 0xFF)); // LSB
	for (int i = 0; i < 4; i++) {
		Wire.write(bytes[i]);
	}
	Wire.endTransmission();
	hack(false);
}

float I2C_EEPROM::ReadFloat(int addr) {
	hack(true);
	float result;
	byte* bytes = (byte*) &result;

	Logger::Debug("Reading float from eeprom");
	Wire.beginTransmission(EEPROM_I2C_ADDRESS);
	Wire.write((int)(addr >> 8)); // MSB
	Wire.write((int)(addr & 0xFF)); // LSB
	Wire.endTransmission();

	Wire.requestFrom(EEPROM_I2C_ADDRESS, 4);
	for (int i = 0; i < 4; i++) {
		if (Wire.available()) {
			bytes[i] = Wire.read();
		}
	}
	hack(false);
	return result;
}