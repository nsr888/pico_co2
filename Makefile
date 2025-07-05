vi_pico:
	tinygo-edit --target pico --editor nvim --wait
vi_pico2w:
	tinygo-edit --target pico2-w --editor nvim --wait
vscode:
	tinygo-edit --target pico --editor code --wait
flash_pico:
	tinygo flash -target=pico -monitor ./cmd/pico_co2/
flash_pico2w:
	tinygo flash -target=pico2-w -monitor ./cmd/pico_co2/
build_pico:
	tinygo build -target=pico -o main.uf2 ./cmd/pico_co2/
build_pico2w:
	tinygo build -target=pico2-w -o main.uf2 ./cmd/pico_co2/
version:
	go version
	tinygo version
	tinygo-edit --version
install_tinygo_edit:
	go install github.com/sago35/tinygo-edit@latest
flash_example_pico:
	tinygo flash -target=pico -monitor ./pkg/ens160/example/
build_example_pico:
	tinygo build -target=pico -o main.uf2 ./pkg/ens160/example/
