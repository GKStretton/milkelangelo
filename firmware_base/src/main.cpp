#include <Arduino.h>
#include <ContinuousStepper.h>
#include <Preferences.h>
#include <PubSubClient.h>
#include <WiFi.h>

#define STEP_PIN 2
#define DIR_PIN 0
#define SLEEP_PIN 19
#define E_STOP_PIN 27

const char *WIFI_SSID = "TP-Link_8A2C";
const char *WIFI_PASSWORD = "85118010";

const char *MQTT_BROKER = "192.168.1.102";
const int MQTT_PORT = 1883;
const char *MQTT_CLIENT_ID = "milkelangelo-base";
const char *TOPIC_SPEED_SET = "milkelangelo/base/speed/set";  // subscribe: steps/sec, signed
const char *TOPIC_STATUS = "milkelangelo/base/status";        // publish

const float ACCELERATION_STEPS_PER_SEC2 = 800.0;
const uint32_t STATUS_PUBLISH_INTERVAL_MS = 1000;

// Bit-banged from loop() on this core (the default Arduino task, pinned to
// core 1). The MQTT client runs in its own task pinned to core 0 so handling
// network traffic can never stall or jitter the step pulses.
ContinuousStepper<StepperDriver> stepper;
WiFiClient wifiClient;
PubSubClient mqttClient(wifiClient);
Preferences preferences;
uint32_t startupCount;

void setSpeed(float speed) {
  digitalWrite(SLEEP_PIN, speed == 0 ? LOW : HIGH);  // active-low SLEEP: sleep when stopped
  stepper.spin(speed);
  Serial.println("set speed to " + String(speed) + " steps/sec");
}

void mqttCallback(char *topic, byte *payload, unsigned int length) {
  if (strcmp(topic, TOPIC_SPEED_SET) != 0) {
    return;
  }
  String value((char *)payload, length);
  setSpeed(value.toFloat());
}

void publishStatus() {
  String body;
  body += "speed: " + String(stepper.speed()) + " steps/sec\n";
  body += "estop: " + String(digitalRead(E_STOP_PIN)) + "\n";
  body += "startups: " + String(startupCount) + "\n";
  mqttClient.publish(TOPIC_STATUS, body.c_str());
}

void mqttReconnect() {
  while (!mqttClient.connected()) {
    Serial.print("connecting to MQTT broker...");
    if (mqttClient.connect(MQTT_CLIENT_ID)) {
      Serial.println("connected");
      mqttClient.subscribe(TOPIC_SPEED_SET);
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

  preferences.begin("driver", false);
  startupCount = preferences.getUInt("startups", 0) + 1;
  preferences.putUInt("startups", startupCount);
  preferences.end();
  Serial.print("startup count: ");
  Serial.println(startupCount);

  pinMode(E_STOP_PIN, INPUT);
  pinMode(SLEEP_PIN, OUTPUT);
  digitalWrite(SLEEP_PIN, LOW);  // active-low SLEEP: sleep by default until a nonzero speed is set

  stepper.begin(STEP_PIN, DIR_PIN);
  stepper.setAcceleration(ACCELERATION_STEPS_PER_SEC2);
  stepper.spin(0);  // stopped until the API sets a speed

  WiFi.mode(WIFI_STA);
  WiFi.setHostname("milkelangelo-base-driver");
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
  if (digitalRead(E_STOP_PIN) == HIGH) {
    stepper.loop();
  }
}
