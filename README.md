# pico_co2
Raspberry Pico CO2 measurements

This project is designed to measure CO2 levels, air quality, and other environmental parameters using a Raspberry Pi Pico microcontroller. It integrates sensors and displays to provide real-time data visualization.

## Hardware installation

### Required components

* Raspberry Pico board
* SSD1306 display
* AHT20+ENS160 sensor

### Connection

![pico_co2_](https://github.com/user-attachments/assets/616411a3-a43a-46e7-acf2-e7d9d982135e)

## Software installation

```bash
make flash
```

## Case

![20250714_065327_](https://github.com/user-attachments/assets/1119b4b1-9fd0-45b7-aa6f-d544de5fbf7b)


## Development
### requirements for vim development
* go version go1.24.4 linux/amd64
* tinygo version 0.37.0 linux/amd64
* go install github.com/sago35/tinygo-edit@latest
* run as `tinygo-edit --target pico --editor nvim --wait`
