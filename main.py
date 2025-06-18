import json
# import network
import os
import sys
import time

import AHT21
import ENS160
import framebuf
import RTC_DS3231
import settings
# import urequests
from machine import I2C, Pin, WDT
from smile_faces import faces
from ssd1306 import SSD1306_I2C

i2c_dev = I2C(0, scl=Pin(5), sda=Pin(4), freq=200000)
i2c_addr = [hex(ii) for ii in i2c_dev.scan()]  # get I2C address in hex format
print("I2C Configuration: {}".format(i2c_dev))  # print I2C params
print("I2C Addresses    : {}".format(i2c_addr))  # I2C devices address

rtc = RTC_DS3231.RTC(i2c=i2c_dev)

# It is encoded like this: sec min hour week day month year.
# Uncomment the line below to set time. Do this only once otherwise time will be set everytime the code is executed.
# rtc.DS3231_SetTime(b'\x30\x15\x21\x02\x17\x06\x25')

t = rtc.DS3231_ReadTime(1)  # Read RTC and receive data in Mode 1 (see RTC_DS3231.py)
print("Time             : {}".format(t))

# i2c = I2C(1, scl=Pin(27), sda=Pin(26))
# oled = SSD1306_I2C(128, 64, i2c)

# fb = framebuf.FrameBuffer(faces["excellent"], 32, 32, framebuf.MONO_HLSB)

# oled.fill(0)
# oled.blit(fb, 48, 16)
# oled.show()
# time.sleep(5)


# setup display
pix_res_x = 128  # SSD1306 horizontal resolution
pix_res_y = 64  # SSD1306 vertical resolution

# смайлы для AQI
# Размер каждого: 32x32, формат framebuf.MONO_HLSB
smile_excellent = faces["excellent"]
smile_good = faces["good"]
smile_moderate = faces["moderate"]
smile_poor = faces["poor"]
smile_unhealthy = faces["unhealthy"]

oled = SSD1306_I2C(pix_res_x, pix_res_y, i2c_dev)  # oled controller

# boot pattern
led = Pin("LED", Pin.OUT)
print("Playing LED boot pattern...")
for _ in range(3):
    led.on()
    time.sleep(0.5)
    led.off()
    time.sleep(0.5)

# try display last save data
last_values = None
try:
    with open("last_readings.json", "r") as f:
        last_values = json.loads(f.read())
        print("Last readings loaded:", last_values)
except Exception as e:
    print("No last readings found:", e)

# dispaly last data
if last_values:
    oled.fill(0)
    oled.text(f"AQI: {int(last_values['aqi'])}", 5, 5)
    oled.text(f"eCO2: {int(last_values['eco2'])}", 5, 15)
    oled.text(f"TVOC: {int(last_values['tvoc'])}", 5, 25)
    oled.text(f"Humid: {round(last_values['humidity'], 1)}%", 5, 35)
    oled.text(f"Temp: {round(last_values['temperature'], 1)}C", 5, 45)
    oled.text(last_values["status"], ((128 - len(last_values["status"]) * 8) // 2), 55)
    oled.show()
    time.sleep(3)

# set up
print("Setting up ENS160 and AHT21 interface via I2C...")
aht = AHT21.AHT21(i2c_dev)

# calibration status
ens_calibrated: bool = False
startup_calibration_time = None
display_delay_status_calibrate = 10_000  # 10 sec

# save calibration status
ens_state_saved: bool = False
state_save_interval_ms = 1 * 60 * 1000  # 1 минута
next_state_save_time = time.ticks_ms() + state_save_interval_ms

ens = ENS160.ENS160(i2c_dev)
# Не вызывайте ens.reset() каждый раз при старте, иначе потеряете загруженное состояние.
# Вызывайте reset() только при проблемах (например, если датчик не отвечает или возвращает нули).
# ens.reset()

# Попытка восстановить сохранённое состояние ENS160
# Попытка восстановления только если не установлен флаг "начать с нуля"
load_previous_state = True  # ← можешь переключать вручную при переносе
state_restored = False

if load_previous_state:
    try:
        with open("ens160_state.dat", "rb") as f:
            state = f.read()
            ens.set_state(state)
            ens_calibrated = True
            # startup_calibration_time = time.ticks_ms()
            state_restored = True
            print("ENS160 state calibrate restored.")
    except:
        state = None
        print("No ENS160 state found or restore failed.")
else:
    print("Starting with fresh ENS160 calibration (state not loaded).")

time.sleep(0.5)
ens.operating_mode = 2
time.sleep(2.0)

# connect to wifi
"""
print("Preparing for wifi connection...")
wifi_con_attempt:int = 0
wlan = network.WLAN(network.STA_IF)
wlan.active(True)
while wlan.isconnected() == False:

    wifi_con_attempt = wifi_con_attempt + 1

    # blip light
    led.on()
    time.sleep(0.1)
    led.off()
    
    print("Attempt #" + str(wifi_con_attempt) + " to connect to wifi...")
    wlan.connect(settings.ssid, settings.password)
    time.sleep(3)
print("Connected to wifi after " + str(wifi_con_attempt) + " tries!")
my_ip:str = str(wlan.ifconfig()[0])
print("My IP Address: " + my_ip)
"""

# create watchdog timer
wdt = WDT(timeout=8388)  # 8,388 ms is the limit (8.388 seconds)
wdt.feed()
print("Watchdog timer now activated.")

# Enter infinite loop
samples_uploaded: int = 0
while True:
    # LED on while we are doing something
    led.on()

    # Чтение температуры и влажности один раз за цикл
    rht = aht.read()
    # print("RAW AHT:", rht)
    humidity = rht[0]
    temperature = rht[1]
    ens.set_envdata(temperature, humidity)

    # take reading from ENS160
    print("Taking ENS160 measurements... ")
    aqi: int = ens.AQI
    eco2: int = ens.CO2
    tvoc: int = ens.TVOC

    if aqi == 1:
        status = "- Excellent -"
    elif aqi == 2:
        status = "- Good -"
    elif aqi == 3:
        status = "- Moderate -"
    elif aqi == 4:
        status = "- Poor -"
    else:
        status = "- Unhealthy -"

    print(
        "AQI: "
        + str(aqi)
        + ", ECO2: "
        + str(eco2)
        + ", TVOC: "
        + str(tvoc)
        + ", Status: "
        + str(status)
    )

    wdt.feed()

    # Загружаем последние валидные значения, если еще не загружены
    if last_values is None:
        try:
            with open("last_readings.json", "r") as f:
                last_values = json.loads(f.read())
        except:
            last_values = {
                "aqi": 0,
                "eco2": 400,
                "tvoc": 0,
                "humidity": 0.0,
                "temperature": 0.0,
                "status": "-",
            }

    # Подготовка данных для отображения — если текущие значения 0, то берём из файла
    display_aqi = aqi if aqi != 0 else last_values["aqi"]
    display_eco2 = eco2 if eco2 > 400 else last_values["eco2"]
    display_tvoc = tvoc if tvoc != 0 else last_values["tvoc"]
    display_temp = temperature if temperature != 0 else last_values["temperature"]
    display_hum = humidity if humidity != 0 else last_values["humidity"]
    display_status = status if aqi != 0 else last_values["status"]

    smile_map = {
        1: smile_excellent,
        2: smile_good,
        3: smile_moderate,
        4: smile_poor,
        5: smile_unhealthy,
    }
    icon_buf = smile_map.get(display_aqi, smile_unhealthy)
    fb = framebuf.FrameBuffer(icon_buf, 32, 32, framebuf.MONO_HLSB)

    # если данные валидны — проверим, нужно ли сохранять
    if aqi != 0 and eco2 > 400 and tvoc != 0:
        try:
            body = {
                "aqi": aqi,
                "eco2": eco2,
                "tvoc": tvoc,
                "humidity": humidity,
                "temperature": temperature,
                "status": display_status,
            }
            print("Measurements taken! " + str(body))

            with open("last_readings.json", "w") as f:
                f.write(json.dumps(body))
                f.flush()
                os.sync()
            last_values = body  # обновим текущие валидные значения
            print("Last readings saved.")
        except Exception as e:
            print("Failed to save last readings:", e)

        # сохраним состояние ENS160 при валидных показаниях
        if not ens_calibrated:
            ens_calibrated = True
            if startup_calibration_time is None:
                startup_calibration_time = time.ticks_ms()

        now = time.ticks_ms()
        if not ens_state_saved or time.ticks_diff(now, next_state_save_time) >= 0:
            try:
                state = ens.get_state()
                with open("ens160_state.dat", "wb") as f:
                    f.write(state)
                    f.flush()
                    os.sync()
                ens_state_saved = True
                next_state_save_time = now + state_save_interval_ms
                print("ENS160 state saved.")
            except Exception as e:
                print("Failed to save ENS160 state:", e)

    # if there was an error getting legit values from the ENS160
    if (aqi == 0 or eco2 <= 400 or tvoc == 0) and not state_restored:
        ens.operating_mode = 1
        time.sleep(0.1)
        ens.operating_mode = 2

        print(
            "AQI, ECO2, or TVOC readings were not successful! Going into troubleshooting mode."
        )
        while (
            aqi == 0 or eco2 <= 400 or tvoc == 0
        ):  # go into troubleshooting mode, trying to recover ENS160 functionality, until solved for.

            # print msg
            print(
                "AQI/ECO2/TVOC reading unsuccessful at last attempt. Will try to reset again."
            )

            # flash quickly to show an issue
            print("Playing troubleshooting LED pattern...")
            for x in range(0, 5):
                led.on()
                time.sleep(0.05)
                led.off()
                time.sleep(0.05)
                wdt.feed()

            print("Reinitializing ENS160 (soft)...")
            ens.operating_mode = 1
            time.sleep(0.1)
            ens.operating_mode = 2
            wdt.feed()  # feed watchdog timer

            # take sample
            print("Sampling ENS160 after reset...")
            time.sleep(1)
            led.on()
            aqi = ens.AQI
            eco2 = ens.CO2
            tvoc = ens.TVOC
            wdt.feed()
            led.off()

    wdt.feed()

    # display
    # Clear the oled display in case it has junk on it.
    oled.fill(0)

    # Blit the image from the framebuffer to the oled display
    oled.blit(fb, 96, 0)

    # Display like integer values
    # oled.text(f'aqi: {int(aqi)}', 5, 5)
    # oled.text(f'eco2: {int(eco2)}', 5, 15)
    # oled.text(f'tvoc: {int(tvoc)}', 5, 25)
    # oled.text(f'humid: {int(humidity)}%', 5, 35)
    # led.text(f'temp: {int(temperature)}C', 5, 45)

    # Display with one decimal place.
    # oled.text(f'AQI: {int(aqi)}', 5, 5)
    # oled.text(f'eCO2: {int(eco2)}', 5, 15)
    # oled.text(f'TVOC: {int(tvoc)}', 5, 25)
    # oled.text(f'Humid: {round(humidity, 1)}%', 5, 35)
    # oled.text(f'Temp: {round(temperature, 1)}C', 5, 45)
    oled.text(f"AQI: {int(display_aqi)}", 5, 5)
    oled.text(f"eCO2: {int(display_eco2)}", 5, 15)
    oled.text(f"TVOC: {int(display_tvoc)}", 5, 25)
    oled.text(f"Humid: {round(display_hum, 1)}%", 5, 35)
    oled.text(f"Temp: {round(display_temp, 1)}C", 5, 45)

    # Display calibration status
    # Показывать Calibrated только в первые 10 секунд после старта
    if (
        ens_calibrated
        and time.ticks_diff(time.ticks_ms(), startup_calibration_time)
        < display_delay_status_calibrate
    ):
        oled.text("Calibrated", 5, 55)
    else:
        oled.text(display_status, ((128 - len(display_status) * 8) // 2), 55)

    # Finally update the oled display so the image & text is displayed
    oled.show()

    # make HTTP call
    """
    print("Making HTTP call...")
    wdt.feed()
    pr = urequests.post(settings.post_url, json=body)
    wdt.feed()
    pr.close()
    print("HTTP call made!")

    # if the status code of the HTTP response was not succesful (not in the 200 range), go into an infinite loop of
    # This is here so the program will stop if, for example, the receiving endpoint is no longer active or accepting.
    if str(pr.status_code)[0:1] != "2":
        while True: 
            led.on()
            wdt.feed()
            time.sleep(1)
            led.off()
            wdt.feed()
            time.sleep(1)
    else:
        print("Sample upload successfully accepted!")
    """
    # increment tracker
    samples_uploaded = samples_uploaded + 1

    # wait for time
    led.off()  # led off while doing nothing (just waiting)
    next_loop: int = time.ticks_ms() + (1000 * settings.sample_time_seconds)
    while time.ticks_ms() < next_loop:
        print(
            "Sampling #"
            + str(samples_uploaded + 1)
            + " next in "
            + str(round((next_loop - time.ticks_ms()) / 1000, 0))
            + " seconds..."
        )
        time.sleep(1)
        wdt.feed()
