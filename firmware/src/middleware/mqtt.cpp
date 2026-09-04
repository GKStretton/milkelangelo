#include "mqtt.h"
#include <WiFi.h>
#include <PubSubClient.h>
#include "../extras/nanopb/pb_encode.h"

namespace Mqtt {
	namespace {
		const char *WIFI_SSID = "TP-Link_8A2C";
		const char *WIFI_PASSWORD = "85118010";

		const char *MQTT_BROKER = "192.168.1.102";
		const int MQTT_PORT = 1883;
		const char *MQTT_CLIENT_ID = "milkelangelo-outer";
		// "mega/req/#" matches the legacy topic namespace (see
		// ../extras/topics_firmware/topics_firmware.h) that the existing
		// backend/UI already expects, plus this firmware's own
		// mega/req/manual/* jog topics (main.cpp), which aren't yet part of
		// the shared asol-protos contract.
		const char *TOPIC_SUBSCRIBE_WILDCARD = "mega/req/#";

		const size_t TOPIC_MAX = 64;
		const size_t INBOUND_PAYLOAD_MAX = 192;
		const size_t OUTBOUND_PAYLOAD_MAX = 256;

		struct InboundMsg {
			char topic[TOPIC_MAX];
			char payload[INBOUND_PAYLOAD_MAX];
		};

		struct OutboundMsg {
			char topic[TOPIC_MAX];
			uint8_t payload[OUTBOUND_PAYLOAD_MAX];
			size_t length;
		};

		WiFiClient wifiClient;
		PubSubClient client(wifiClient);

		QueueHandle_t inboundQueue;
		QueueHandle_t outboundQueue;

		void (*topicHandler)(String topic, String payload) = nullptr;

		// Runs on the mqtt task. Only enqueues the message; state mutation
		// happens later on the main loop via Update(), so app state is never
		// touched concurrently from both cores.
		void onMessage(char *topic, byte *payload, unsigned int length) {
			InboundMsg msg;
			strncpy(msg.topic, topic, TOPIC_MAX - 1);
			msg.topic[TOPIC_MAX - 1] = '\0';

			size_t n = length < INBOUND_PAYLOAD_MAX - 1 ? length : INBOUND_PAYLOAD_MAX - 1;
			memcpy(msg.payload, payload, n);
			msg.payload[n] = '\0';

			xQueueSend(inboundQueue, &msg, 0);
		}

		void reconnect() {
			while (!client.connected()) {
				Serial.print("connecting to MQTT broker...");
				if (client.connect(MQTT_CLIENT_ID)) {
					Serial.println("connected");
					client.subscribe(TOPIC_SUBSCRIBE_WILDCARD);
				} else {
					Serial.println("failed, rc=" + String(client.state()) + ", retrying in 2s");
					delay(2000);
				}
			}
		}

		// All PubSubClient socket I/O (connect/loop/publish) happens
		// exclusively from this task, pinned to core 0, so it's never touched
		// concurrently from the stepper control loop on core 1.
		void task(void *) {
			for (;;) {
				if (!client.connected()) {
					reconnect();
				}
				client.loop();

				OutboundMsg out;
				while (xQueueReceive(outboundQueue, &out, 0) == pdTRUE) {
					client.publish(out.topic, out.payload, out.length);
				}

				vTaskDelay(1);  // yield so WiFi/idle tasks on this core get scheduled
			}
		}
	}

	void Init() {
		inboundQueue = xQueueCreate(16, sizeof(InboundMsg));
		outboundQueue = xQueueCreate(16, sizeof(OutboundMsg));

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

		client.setBufferSize(512);
		client.setServer(MQTT_BROKER, MQTT_PORT);
		client.setCallback(onMessage);

		xTaskCreatePinnedToCore(task, "mqtt", 4096, nullptr, 1, nullptr, 0);
	}

	void SetTopicHandler(void (*f)(String topic, String payload)) {
		topicHandler = f;
	}

	void Update() {
		if (topicHandler == nullptr) {
			return;
		}
		InboundMsg msg;
		while (xQueueReceive(inboundQueue, &msg, 0) == pdTRUE) {
			topicHandler(String(msg.topic), String(msg.payload));
		}
	}

	void Publish(const char *topic, const String &payload) {
		OutboundMsg out;
		strncpy(out.topic, topic, TOPIC_MAX - 1);
		out.topic[TOPIC_MAX - 1] = '\0';

		size_t n = payload.length() < OUTBOUND_PAYLOAD_MAX ? payload.length() : OUTBOUND_PAYLOAD_MAX;
		memcpy(out.payload, payload.c_str(), n);
		out.length = n;

		xQueueSend(outboundQueue, &out, 0);
	}

	void PublishProto(const char *topic, const pb_msgdesc_t *fields, const void *src_struct) {
		OutboundMsg out;
		strncpy(out.topic, topic, TOPIC_MAX - 1);
		out.topic[TOPIC_MAX - 1] = '\0';

		pb_ostream_t stream = pb_ostream_from_buffer(out.payload, OUTBOUND_PAYLOAD_MAX);
		if (!pb_encode(&stream, fields, src_struct)) {
			Serial.println("pb_encode failed in Mqtt::PublishProto");
			return;
		}
		out.length = stream.bytes_written;

		xQueueSend(outboundQueue, &out, 0);
	}
}
