// todo: esp32 outer firmware
#include <Arduino.h>
#include <AccelStepper.h>
#include "config_new.h"

AccelStepper ringStepper(AccelStepper::DRIVER, RING_STEPPER_STEP, RING_STEPPER_DIR);
AccelStepper zStepper(AccelStepper::DRIVER, Z_STEPPER_STEP, Z_STEPPER_DIR);
AccelStepper yawStepper(AccelStepper::DRIVER, YAW_STEPPER_STEP, YAW_STEPPER_DIR);
AccelStepper pitchStepper(AccelStepper::DRIVER, PITCH_STEPPER_STEP, PITCH_STEPPER_DIR);
AccelStepper pipetteStepper(AccelStepper::DRIVER, PIPETTE_STEPPER_STEP, PIPETTE_STEPPER_DIR);

const float SPIN_SPEED = 100.0; // steps/sec
const unsigned long DIRECTION_FLIP_INTERVAL_MS = 5000;

unsigned long lastDirectionFlip = 0;
unsigned long lastLimitPrint = 0;
int spinDirection = 1;

void setSpinDirection(int direction) {
  ringStepper.setSpeed(SPIN_SPEED * direction);
  ringStepper.setAcceleration(100);
  zStepper.setSpeed(SPIN_SPEED * direction);
  yawStepper.setSpeed(SPIN_SPEED * direction);
  pitchStepper.setSpeed(SPIN_SPEED * direction);
  pipetteStepper.setSpeed(SPIN_SPEED * direction);
}

void setup() {
  Serial.begin(9600);
  pinMode(PITCH_LIMIT_SWITCH, INPUT);
  pinMode(YAW_LIMIT_SWITCH, INPUT);
  pinMode(Z_LIMIT_SWITCH, INPUT);
  pinMode(RING_LIMIT_SWITCH, INPUT);
  pinMode(PIPETTE_LIMIT_SWITCH, INPUT);

  pinMode(STEPPER_SLEEP, OUTPUT);
  digitalWrite(STEPPER_SLEEP, HIGH);

  setSpinDirection(spinDirection);
}

void loop() {
  unsigned long now = millis();

  if (now - lastDirectionFlip >= DIRECTION_FLIP_INTERVAL_MS) {
    lastDirectionFlip = now;
    spinDirection = -spinDirection;
    setSpinDirection(spinDirection);
  }

  ringStepper.runSpeed();
  zStepper.runSpeed();
  yawStepper.runSpeed();
  pitchStepper.runSpeed();
  pipetteStepper.runSpeed();

  if (now - lastLimitPrint >= 2000) {
    lastLimitPrint = now;

    bool pitchLimitTriggered = digitalRead(PITCH_LIMIT_SWITCH);
    bool yawLimitTriggered = digitalRead(YAW_LIMIT_SWITCH);
    bool zLimitTriggered = digitalRead(Z_LIMIT_SWITCH);
    bool ringLimitTriggered = digitalRead(RING_LIMIT_SWITCH);
    bool pipetteLimitTriggered = digitalRead(PIPETTE_LIMIT_SWITCH);

    Serial.println(pitchLimitTriggered ? "Pitch limit switch: TRIGGERED" : "Pitch limit switch: clear");
    Serial.println(yawLimitTriggered ? "Yaw limit switch: TRIGGERED" : "Yaw limit switch: clear");
    Serial.println(zLimitTriggered ? "Z limit switch: TRIGGERED" : "Z limit switch: clear");
    Serial.println(ringLimitTriggered ? "Ring limit switch: TRIGGERED" : "Ring limit switch: clear");
    Serial.println(pipetteLimitTriggered ? "Pipette limit switch: TRIGGERED" : "Pipette limit switch: clear");
  }
}
