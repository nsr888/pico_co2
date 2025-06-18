run:
	mpremote run main.py
push_main:
	mpremote fs cp main.py :main.py
push_RTC:
	mpremote fs cp RTC_DS3231.py :RTC_DS3231.py
pull_main:
	mpremote fs cp :main.py main.py
ls:
	mpremote fs ls
pull_all:
	mpremote fs cp :main.py main.py
	mpremote fs cp :AHT21.py AHT21.py
	mpremote fs cp :ENS160.py ENS160.py
	mpremote fs cp :RTC_DS3231.py RTC_DS3231.py
	mpremote fs cp :ens160_state.dat ens160_state.dat
	mpremote fs cp :last_readings.json last_readings.json
	mpremote fs cp :picozero.py picozero.py
	mpremote fs cp :settings.py settings.py
	mpremote fs cp :smile_faces.py smile_faces.py
	mpremote fs cp :ssd1306.py ssd1306.py

