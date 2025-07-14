# pico_co2
Raspberry Pico CO2 measurements

This project is designed to measure CO2 levels, air quality, and other environmental parameters using a Raspberry Pi Pico microcontroller. It integrates sensors and displays to provide real-time data visualization.

## Hardware installation

### Required components

* Raspberry Pico board
* SSD1306 display
* AHT20+ENS160 sensor

### Connection

![Image](https://github.com/user-attachments/assets/a3db4534-e092-4903-9f99-1fe5824c7bd3)

## Software installation

```bash
make flash
```

## Case

![Image](https://github.com/user-attachments/assets/db7151cf-6b90-4819-a6e3-5002a788fd2c)

## Development
### requirements for vim development
* go version go1.24.4 linux/amd64
* tinygo version 0.37.0 linux/amd64
* go install github.com/sago35/tinygo-edit@latest
* run as `tinygo-edit --target pico --editor nvim --wait`
