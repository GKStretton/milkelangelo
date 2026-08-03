#include <AccelStepper.h>

// STEP/DIR driver (A4988, DRV8825, TMC2208, etc.) — change to match your wiring.
#define STEP_PIN 26
#define DIR_PIN 27
#define SLEEP_PIN 25

const float SPEED_STEPS_PER_SEC = 500.0;

AccelStepper stepper(AccelStepper::DRIVER, STEP_PIN, DIR_PIN);

unsigned long lastToggleMs = 0;
bool sleepState = HIGH;

void setup() {
  Serial.begin(9600);

  pinMode(SLEEP_PIN, OUTPUT);
  digitalWrite(SLEEP_PIN, sleepState);  // active-low SLEEP: HIGH keeps the driver awake
  Serial.println(sleepState ? "awake" : "asleep");

  stepper.setMaxSpeed(SPEED_STEPS_PER_SEC);
  stepper.setSpeed(SPEED_STEPS_PER_SEC);
}

void loop() {
  if (millis() - lastToggleMs >= 5000) {
    sleepState = !sleepState;
    digitalWrite(SLEEP_PIN, sleepState);
    Serial.println(sleepState ? "awake" : "asleep");
    lastToggleMs = millis();
  }

  stepper.runSpeed();
}
