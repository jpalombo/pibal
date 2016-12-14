#ifndef SENSOR_H
#define SENSOR_H

#include <stdint.h>

int mpu_open();
int sensorAccel(int);
int sensorGyro(int);
int getLastError();

#endif
