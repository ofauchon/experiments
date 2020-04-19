
build:
	tinygo build -size short -target=bluepill -o main.hex main.go	

flash:
	arm-none-eabi-gdb -nx --batch -ex "target extended-remote /dev/ttyACM1" -ex "monitor swdp_scan" -ex "attach 1" -ex load -ex compare-sections -ex kill ./main.hex

term:
	picocom -b 9600 /dev/ttyACM2
