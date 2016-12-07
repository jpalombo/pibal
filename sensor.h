#ifndef SENSOR_H
#define SENSOR_H

#include <stdint.h>

int mpu_open();
int sensorAngle(int);
int sensorGyro(int);

#endif
