# -*- coding: UTF-8 -*-
# This is a sample Python script.
import utime
import driver

GPIO = driver.gpio()
leds=("led1", "led2", "led3", "led4", "led5")
for i in range(5):
    for led in leds:
        GPIO.open("/data/python/config/led.json", led)
        GPIO.write(0)
        utime.sleep_ms(200)
        # print("LED ON:",i)
        GPIO.write(1)
        utime.sleep_ms(200)
        # print("LED OFF:",i)
        GPIO.write(0)
        utime.sleep_ms(200)
        GPIO.close()