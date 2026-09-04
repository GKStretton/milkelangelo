#include "util.h"

void InitPin(uint8_t pin, byte v) {
	pinMode(pin, OUTPUT);
	digitalWrite(pin, v);
}
