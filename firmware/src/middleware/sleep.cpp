#include "sleep.h"
#include <Arduino.h>
#include <Preferences.h>
#include "../config.h"
#include "../middleware/logger.h"
#include "../app/state_report.h"
#include "../extras/topics_firmware/topics_firmware.h"

namespace Sleep {
	namespace {
		const char *PREFERENCES_NAMESPACE = "outer";
		const char *SAFE_SHUTDOWN_KEY = "safeShutdown";

		unsigned long lastNod = millis();
		unsigned long lastPrint = millis() - SLEEP_PRINT_INTERVAL;
		bool sleeping = true;
		SleepStatus lastSleepStatus = SleepStatus::UNKNOWN;

		void (*externalSleepHandler)(SleepStatus sleepStatus) = NULL;
		void (*externalWakeHandler)(SleepStatus lastSleepStatus) = NULL;

		bool eStopActive() {
			return digitalRead(E_STOP_PIN) == LOW;
		}

		void onSleep(SleepStatus status) {
			Logger::Info("Going to sleep with status " + String(status));

			Preferences preferences;
			preferences.begin(PREFERENCES_NAMESPACE, false);
			preferences.putUChar(SAFE_SHUTDOWN_KEY, (uint8_t) status);
			preferences.end();

			if (externalSleepHandler != NULL) {
				Logger::Info("Calling externalSleepHandler");
				externalSleepHandler(status);
			}

			if (eStopActive()) {
				StateReport_SetStatus(machine_Status_E_STOP_ACTIVE);
			} else {
				StateReport_SetStatus(machine_Status_SLEEPING);
			}
			StateReport_ForceSend();
		}

		void onWake() {
			StateReport_SetStatus(machine_Status_WAKING_UP);
			StateReport_ForceSend();

			Logger::Info("Waking up");

			Preferences preferences;
			preferences.begin(PREFERENCES_NAMESPACE, false);
			uint8_t data = preferences.getUChar(SAFE_SHUTDOWN_KEY, (uint8_t) UNKNOWN);
			lastSleepStatus = (SleepStatus) data;
			// write back to unknown in case of sudden shutdown
			preferences.putUChar(SAFE_SHUTDOWN_KEY, (uint8_t) UNKNOWN);
			preferences.end();
			Logger::Info("read lastSleepStatus as " + String(lastSleepStatus));

			if (externalWakeHandler != NULL) {
				Logger::Info("Calling externalWakeHandler");
				externalWakeHandler(lastSleepStatus);
			}
		}

		// private version of isSleeping that does the actual checks
		bool isSleeping() {
			//! ordered by priority
			if (eStopActive()) {
				return true;
			}

			if (SLEEP_TIME_MINUTES > 0 && (millis() - lastNod) / 1000 / 60 >= SLEEP_TIME_MINUTES) {
				return true;
			}

			// persist current state by default
			return sleeping;
		}
	}

	void Update() {
		if (isSleeping()) {
			Sleep(UNKNOWN);
			if (eStopActive()) {
				StateReport_SetStatus(machine_Status_E_STOP_ACTIVE);
			} else {
				StateReport_SetStatus(machine_Status_SLEEPING);
			}
		} else {
			Wake();
		}

		if (millis() - lastPrint > SLEEP_PRINT_INTERVAL) {
			lastPrint = millis();
			if (eStopActive()) {
				Logger::Info("E_STOP ACTIVE");
			} else if (sleeping) {
				Logger::Info("sleeping...");
			} else {
				Logger::Info("awake");
			}
		}
	}

	void Wake() {
		if (eStopActive()) {
			return;
		}

		lastNod = millis();
		bool wasSleeping = sleeping;
		sleeping = false;
		if (wasSleeping) {
			onWake();
		}
	}

	void Sleep(SleepStatus status) {
		bool wasSleeping = sleeping;
		sleeping = true;
		if (!wasSleeping) {
			onSleep(status);
		}
	}

	bool IsSleeping() {
		return sleeping;
	}

	bool IsEStopActive() {
		return eStopActive();
	}

	SleepStatus GetLastSleepStatus() {
		return lastSleepStatus;
	}

	void OverrideLastSleepStatus(SleepStatus status) {
		lastSleepStatus = status;
	}

	void SetOnSleepHandler(void (*f)(SleepStatus)) {
		externalSleepHandler = f;
	}

	void SetOnWakeHandler(void (*f)(SleepStatus)) {
		externalWakeHandler = f;
	}
}
