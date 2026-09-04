#pragma once

#define CONTROLLER_UPDATE_PERIOD_MS 100

// Autonomous configuration
#define DO_DYE_COLLECTION true
#define ENABLE_IK_ACTUATION 1

// 0 disables auto-sleep-on-idle-timeout
#define SLEEP_TIME_MINUTES 0
#define SLEEP_PRINT_INTERVAL 5000

#define E_STOP_PIN 14

#define STEPPER_SLEEP 23

#define RING_STEPPER_STEP 22
#define RING_STEPPER_DIR 21
#define Z_STEPPER_STEP 19
#define Z_STEPPER_DIR 18
#define YAW_STEPPER_STEP 5
#define YAW_STEPPER_DIR 17
#define PITCH_STEPPER_STEP 16
#define PITCH_STEPPER_DIR 4
#define PIPETTE_STEPPER_STEP 0
#define PIPETTE_STEPPER_DIR 2

// Limit switches
#define PITCH_LIMIT_SWITCH 33
#define YAW_LIMIT_SWITCH 32
#define Z_LIMIT_SWITCH 35
#define RING_LIMIT_SWITCH 34
#define PIPETTE_LIMIT_SWITCH 25
