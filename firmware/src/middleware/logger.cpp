#include "logger.h"
#include "mqtt.h"
#include "../extras/topics_firmware/topics_firmware.h"

namespace Logger {

	namespace {
		enum Level level = DEBUG;
	}

	void SetLevel(enum Level _level) {
		level = _level;
	}

	void Crit(String str) {
		if (level >= ERROR) {
			Mqtt::Publish(TOPIC_LOGS_CRIT, str);
		}
	}

	void Error(String str) {
		if (level >= ERROR) {
			Mqtt::Publish(TOPIC_LOGS_ERROR, str);
		}
	}

	void Warn(String str) {
		if (level >= WARN) {
			Mqtt::Publish(TOPIC_LOGS_WARN, str);
		}
	}

	void Info(String str) {
		if (level >= INFO) {
			Mqtt::Publish(TOPIC_LOGS_INFO, str);
		}
	}

	void Debug(String str) {
		if (level >= DEBUG) {
			Mqtt::Publish(TOPIC_LOGS_DEBUG, str);
		}
	}
}
