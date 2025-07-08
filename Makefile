.PHONY: all flash build version install_tinygo_edit flash_ens160_example build_ens160_example

vi:
	tinygo-edit --target pico --editor nvim --wait

flash:
	tinygo flash -target=pico -monitor ./cmd/pico_co2/

build:
	tinygo build -target=pico -o main.uf2 ./cmd/pico_co2/

version:
	go version
	tinygo version
	tinygo-edit --version

install_tinygo_edit:
	go install github.com/sago35/tinygo-edit@latest

flash_ens160_example:
	tinygo flash -target=pico -monitor ./pkg/ens160/example/

build_ens160_example:
	tinygo build -target=pico -o main.uf2 ./pkg/ens160/example/
