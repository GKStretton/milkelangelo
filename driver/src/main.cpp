#include <Arduino.h>
#include <ContinuousStepper.h>
#include <WebServer.h>
#include <WiFi.h>

#define STEP_PIN 2
#define DIR_PIN 0
#define SLEEP_PIN 19
#define ESTOP_PIN 5

const char *AP_SSID = "milkelangelo-driver";
const char *AP_PASSWORD = "driver123";  // WPA2 requires >= 8 chars
// http://192.168.4.1/speed?value=100

const float ACCELERATION_STEPS_PER_SEC2 = 800.0;

// Bit-banged from loop() on this core (the default Arduino task, pinned to
// core 1). The web server runs in its own task pinned to core 0 so handling
// HTTP requests can never stall or jitter the step pulses.
ContinuousStepper<StepperDriver> stepper;
WebServer server(80);

void handleSetSpeed() {
  if (!server.hasArg("value")) {
    server.send(400, "text/plain", "missing 'value' query param (steps/sec, signed)\n");
    return;
  }
  float speed = server.arg("value").toFloat();
  stepper.spin(speed);
  Serial.println("set speed to " + String(speed) + " steps/sec");
  server.send(200, "text/plain", "speed set to " + String(speed) + " steps/sec\n");
}

void handleStatus() {
  String body;
  body += "speed: " + String(stepper.speed()) + " steps/sec\n";
  body += "estop: " + String(digitalRead(ESTOP_PIN)) + "\n";
  server.send(200, "text/plain", body);
}

void webServerTask(void *) {
  for (;;) {
    server.handleClient();
    vTaskDelay(1);  // yield so WiFi/idle tasks on this core get scheduled
  }
}

void setup() {
  Serial.begin(9600);

  pinMode(ESTOP_PIN, INPUT);
  pinMode(SLEEP_PIN, OUTPUT);
  digitalWrite(SLEEP_PIN, HIGH);  // active-low SLEEP: keep the driver awake permanently

  stepper.begin(STEP_PIN, DIR_PIN);
  stepper.setAcceleration(ACCELERATION_STEPS_PER_SEC2);
  stepper.spin(0);  // stopped until the API sets a speed

  WiFi.softAP(AP_SSID, AP_PASSWORD);
  Serial.print("AP IP address: ");
  Serial.println(WiFi.softAPIP());

  server.on("/speed", handleSetSpeed);
  server.on("/status", handleStatus);
  server.begin();

  xTaskCreatePinnedToCore(webServerTask, "webServer", 4096, nullptr, 1, nullptr, 0);
}

void loop() {
  // if (digitalRead(ESTOP_PIN) == HIGH) {
  stepper.loop();
  // }
}
