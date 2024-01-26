#pragma once
#include <Arduino.h>

namespace Logger {

	enum Level {
		NONE,
		CRIT,
		ERROR,
		WARN,
		INFO,
		DEBUG,
	};

	void SetLevel(enum Level level);

	// Critical errors require maintenance attention, maintainer will be emailed
	// by goo.
	void Crit(String str);
	void Error(String str);
	void Warn(String str);
	void Info(String str);
	void Debug(String str);

	void TestLogger();
};