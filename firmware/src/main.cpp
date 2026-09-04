// todo: esp32 outer firmware
#include <Arduino.h>
#include <AccelStepper.h>
#include <PubSubClient.h>
#include <WiFi.h>
#include "config_new.h"

const char *WIFI_SSID = "TP-Link_8A2C";
const char *WIFI_PASSWORD = "85118010";

const char *MQTT_BROKER = "192.168.1.102";
const int MQTT_PORT = 1883;
const char *MQTT_CLIENT_ID = "milkelangelo-outer";
const char *TOPIC_STATUS = "milkelangelo/outer/status";              // publish
const char *TOPIC_SPIN_SPEED_SET = "milkelangelo/outer/spin-speed/set";  // subscribe: steps/sec, unsigned

const uint32_t STATUS_PUBLISH_INTERVAL_MS = 1000;

AccelStepper ringStepper(AccelStepper::DRIVER, RING_STEPPER_STEP, RING_STEPPER_DIR);
AccelStepper zStepper(AccelStepper::DRIVER, Z_STEPPER_STEP, Z_STEPPER_DIR);
AccelStepper yawStepper(AccelStepper::DRIVER, YAW_STEPPER_STEP, YAW_STEPPER_DIR);
AccelStepper pitchStepper(AccelStepper::DRIVER, PITCH_STEPPER_STEP, PITCH_STEPPER_DIR);
AccelStepper pipetteStepper(AccelStepper::DRIVER, PIPETTE_STEPPER_STEP, PIPETTE_STEPPER_DIR);

// Bit-banged from loop() on this core (the default Arduino task, pinned to
// core 1). The MQTT client runs in its own task pinned to core 0 so handling
// network traffic can never stall or jitter the step pulses.
WiFiClient wifiClient;
PubSubClient mqttClient(wifiClient);

float spinSpeed = 20.0; // steps/sec
const unsigned long DIRECTION_FLIP_INTERVAL_MS = 5000;

unsigned long lastDirectionFlip = 0;
unsigned long lastLimitPrint = 0;
int spinDirection = 1;

void setSpinDirection(int direction) {
  ringStepper.setSpeed(spinSpeed * direction);
  zStepper.setSpeed(spinSpeed * direction);
  yawStepper.setSpeed(spinSpeed * direction);
  pitchStepper.setSpeed(spinSpeed * direction);
  pipetteStepper.setSpeed(spinSpeed * direction);
}

void setSpinSpeed(float speed) {
  spinSpeed = speed;

  ringStepper.setMaxSpeed(spinSpeed);
  zStepper.setMaxSpeed(spinSpeed);
  yawStepper.setMaxSpeed(spinSpeed);
  pitchStepper.setMaxSpeed(spinSpeed);
  pipetteStepper.setMaxSpeed(spinSpeed);

  setSpinDirection(spinDirection);
  Serial.println("set spin speed to " + String(spinSpeed) + " steps/sec");
}

void mqttCallback(char *topic, byte *payload, unsigned int length) {
  if (strcmp(topic, TOPIC_SPIN_SPEED_SET) != 0) {
    return;
  }
  String value((char *)payload, length);
  setSpinSpeed(value.toFloat());
}

void publishStatus() {
  String body;
  body += "spinSpeed: " + String(spinSpeed) + " steps/sec\n";
  body += "spinDirection: " + String(spinDirection) + "\n";
  body += "estop: " + String(digitalRead(E_STOP_PIN)) + "\n";
  mqttClient.publish(TOPIC_STATUS, body.c_str());
}

void mqttReconnect() {
  while (!mqttClient.connected()) {
    Serial.print("connecting to MQTT broker...");
    if (mqttClient.connect(MQTT_CLIENT_ID)) {
      Serial.println("connected");
      mqttClient.subscribe(TOPIC_SPIN_SPEED_SET);
    } else {
      Serial.println("failed, rc=" + String(mqttClient.state()) + ", retrying in 2s");
      delay(2000);
    }
  }
}

void mqttTask(void *) {
  uint32_t lastStatusPublish = 0;
  for (;;) {
    if (!mqttClient.connected()) {
      mqttReconnect();
    }
    mqttClient.loop();

    uint32_t now = millis();
    if (now - lastStatusPublish >= STATUS_PUBLISH_INTERVAL_MS) {
      lastStatusPublish = now;
      publishStatus();
    }
    vTaskDelay(1);  // yield so WiFi/idle tasks on this core get scheduled
  }
}

void setup() {
  Serial.begin(9600);
  pinMode(PITCH_LIMIT_SWITCH, INPUT);
  pinMode(YAW_LIMIT_SWITCH, INPUT);
  pinMode(Z_LIMIT_SWITCH, INPUT);
  pinMode(RING_LIMIT_SWITCH, INPUT);
  pinMode(PIPETTE_LIMIT_SWITCH, INPUT);

  pinMode(E_STOP_PIN, INPUT);
  pinMode(STEPPER_SLEEP, OUTPUT);
  digitalWrite(STEPPER_SLEEP, HIGH);

  ringStepper.setMaxSpeed(spinSpeed);
  zStepper.setMaxSpeed(spinSpeed);
  yawStepper.setMaxSpeed(spinSpeed);
  pitchStepper.setMaxSpeed(spinSpeed);
  pipetteStepper.setMaxSpeed(spinSpeed);

  setSpinDirection(spinDirection);

  WiFi.mode(WIFI_STA);
  WiFi.setHostname("milkelangelo-outer");
  WiFi.begin(WIFI_SSID, WIFI_PASSWORD);
  Serial.print("connecting to " + String(WIFI_SSID));
  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
    Serial.print(".");
  }
  Serial.println();
  Serial.print("connected, IP address: ");
  Serial.println(WiFi.localIP());
  Serial.print("MAC address: ");
  Serial.println(WiFi.macAddress());

  mqttClient.setServer(MQTT_BROKER, MQTT_PORT);
  mqttClient.setCallback(mqttCallback);

  xTaskCreatePinnedToCore(mqttTask, "mqtt", 4096, nullptr, 1, nullptr, 0);
}

void loop() {
  unsigned long now = millis();

  if (now - lastDirectionFlip >= DIRECTION_FLIP_INTERVAL_MS) {
    lastDirectionFlip = now;
    spinDirection = -spinDirection;
    setSpinDirection(spinDirection);
  }

  if (digitalRead(E_STOP_PIN) == HIGH) {
    ringStepper.runSpeed();
    zStepper.runSpeed();
    yawStepper.runSpeed();
    pitchStepper.runSpeed();
    pipetteStepper.runSpeed();
  }

  /*
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
  */
}
