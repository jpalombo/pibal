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

#define DIM 3
short accel[DIM]; // [x, y, z]           accel vector
short gyro[DIM];  // [x, y, z]           gyro vector
int gyrotrim = 0; //200;
int acceltrim = 0; //5000;

int mpu_open() {

    unsigned char devStatus; // return status after each device operation

    //monitor->watch(&gyrotrim, "gyrotrim");
    //monitor->watch(&acceltrim, "acceltrim");

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
    int gyrosum = 0;
    int count = 0;
    printf("Calibrating...\n");
    while (count < 300)
    {
        if (mpu_get_gyro_reg(gyro, &timestamp) == 0) {
            gyrosum += gyro[0];
            count++;
        }
        usleep(1000);
    }
    gyrotrim = gyrosum/count;
    printf("Gyro offset : %i\n", gyrotrim);


    printf("Done.\n");
    return 0;
}

int sensorAngle(int i){
    if (i < 0 || i > 2)
      return 0;
    if (mpu_get_accel_reg(accel, &timestamp) == 0)
        return accel[i] - acceltrim;
    else
        return 0;
}

void updateAngleTrim(int offset) {
    acceltrim += offset;
}

int sensorGyro(int i) {
    if (i < 0 || i > 2)
        return 0;
    if (mpu_get_gyro_reg(gyro, &timestamp) == 0)
        return gyro[i] - gyrotrim;
    else
        return 0;
}

void updateGyroTrim(int offset) {
    gyrotrim += offset;
}

int wheelcount()
{

}

int ms_close() {
    return 0;
}
