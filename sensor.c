//#include <unistd.h>
//#include <string.h>
//#include <math.h>

#include <wiringPi.h>

#include <stdio.h>
#include <stdlib.h>
#include "sensor.h"
#include "inv_mpu.h"
#include "inv_mpu_dmp_motion_driver.h"

// MPU control/status vars
unsigned long timestamp;

#define DIM    3
short accel[DIM];       // [x, y, z]           accel vector
short gyro[DIM];        // [x, y, z]           gyro vector
int gyrotrim[DIM] = { 0, 0, 0 };
int acceltrim[DIM] = { 0, 0, 260 };
int lasterror = 0;

int mpu_open()
{
	unsigned char devStatus; // return status after each device operation

	// initialize device
	printf("Initializing MPU...\n");
	if (mpu_init(NULL) != 0) {
		printf("MPU init failed!\n");
		return -1;
	}
	printf("Setting MPU sensors...\n");
	if (mpu_set_sensors(INV_XYZ_GYRO | INV_XYZ_ACCEL) != 0) {
		printf("Failed to set sensors!\n");
		return -1;
	}
	printf("Setting GYRO sensitivity...\n");
	if (mpu_set_gyro_fsr(250) != 0) {
		printf("Failed to set gyro sensitivity!\n");
		return -1;
	}
	printf("Setting ACCEL sensitivity...\n");
	if (mpu_set_accel_fsr(2) != 0) {
		printf("Failed to set accel sensitivity!\n");
		return -1;
	}
	if (mpu_set_compass_sample_rate(10) != 0) {
		printf("Failed to set compass sample rate!\n");
		return -1;
	}
	// verify connection
	printf("Powering up MPU...\n");
	mpu_get_power_state(&devStatus);
	printf(devStatus ? "MPU6050 connection successful\n" : "MPU6050 connection failed %u\n", devStatus);

	// calibrating
	int gyrosum[DIM] = { 0, 0, 0 };
	int count = 0;
	int i;
	printf("Calibrating...\n");
	while (count < 300) {
		if (mpu_get_gyro_reg(gyro, &timestamp) == 0) {
			for (i = 0; i < DIM; i++)
				gyrosum[i] += gyro[i];
			count++;
		}
		usleep(1000);
	}
	for (i = 0; i < DIM; i++)
		gyrotrim[i] = gyrosum[i] / count;
	printf("Gyro offset : %i, %i, %i\n", gyrotrim[0], gyrotrim[1], gyrotrim[2]);
	printf("Done.\n");
	return 0;
}

int getLastError()
{
	return lasterror;
}

int sensorAccel(int i)
{
	if (i < 0 || i > 2) {
		lasterror = 1;
		return 0;
	}
	if (mpu_get_accel_reg(accel, &timestamp) == 0) {
		lasterror = 0;
		return accel[i] - acceltrim[i];
	} else {
		lasterror = 2;
		return 0;
	}
}

void updateAngleTrim(int offsetx, int offsety, int offsetz)
{
	acceltrim[0] += offsetx;
	acceltrim[1] += offsety;
	acceltrim[2] += offsetz;
}

int sensorGyro(int i)
{
	if (i < 0 || i > 2) {
    lasterror = 1;
		return 0;
  }
	if (mpu_get_gyro_reg(gyro, &timestamp) == 0) {
    lasterror = 0;
		return gyro[i] - gyrotrim[i];
  }	else {
    lasterror = 2;
		return 0;
  }
}

void updateGyroTrim(int offsetx, int offsety, int offsetz)
{
	gyrotrim[0] += offsetx;
	gyrotrim[1] += offsety;
	gyrotrim[2] += offsetz;
}
