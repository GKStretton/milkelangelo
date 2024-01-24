#include "mathutil.h"
#include <math.h>

float AngleBetweenVectors(float x0, float y0, float x1, float y1) {
	float dot = x0*x1 + y0*y1;
	float det = x0*y1 - y0*x1;
	return atan2f(det, dot) * 180 / PI_F;
}

void boundXYToCircle(float *x, float *y, float radius) {
	float mag = (float)hypotf((float)*x, (float)*y);
	if (mag > radius)
	{
		*x = *x / mag * radius;
		*y = *y / mag * radius;
	}
}

void boundToSignedMaximum(float *n, float range) {
	if (*n > range) {
		*n = range;
	}
	if (*n < -range) {
		*n = -range;
	}
}

bool numInRange(float num, float min, float max) {
	return num >= min && num <= max;
}