#pragma once
#include <Arduino.h>
#include "../extras/nanopb/pb_encode.h"

namespace Mqtt {
	// Connects to wifi and the MQTT broker, and starts the network task.
	// Blocks until wifi is connected.
	void Init();

	// Registers the handler called (from the main loop, via Update()) for
	// every message received on a subscribed topic.
	void SetTopicHandler(void (*f)(String topic, String payload));

	// Drains received messages and calls the topic handler for each one.
	// Call every loop() iteration so message handling always happens
	// synchronously with the rest of the control loop.
	void Update();

	void Publish(const char *topic, const String &payload);
	void PublishProto(const char *topic, const pb_msgdesc_t *fields, const void *src_struct);
};
