#pragma once

int getRingAndYawFromXY(float x, float y, float lastRing, float *ring, float *yaw, float minRingUnit, float maxRingUnit);
float calculateZCalibrationOffset(float yaw, float ring);
