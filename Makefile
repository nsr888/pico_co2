.PHONY: all flash build version install_tinygo_edit flash_ens160_example build_ens160_example test-displays test-unit

vi:
	tinygo-edit --target pico --editor nvim --wait

flash:
	tinygo flash -size=short -target=pico -monitor ./cmd/pico_co2/

build:
	tinygo build -size=short -target=pico -o main.uf2 ./cmd/pico_co2/

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

test-displays:
	go run ./cmd/virtualdisplay/

test-unit:
	go test -v ./internal/display/...
