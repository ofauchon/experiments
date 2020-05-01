
BM_DEV=/dev/ttyACM0

build:
	tinygo build -size short -target=bluepill -o main.hex main.go	

flash:
	arm-none-eabi-gdb -nx --batch -ex "target extended-remote $(BM_DEV)" -ex "monitor swdp_scan" -ex "attach 1" -ex load -ex compare-sections -ex kill ./main.hex

debug:
	arm-none-eabi-gdb -nx -ex "target extended-remote $(BM_DEV)"  -ex "monitor swdp_scan" -ex "attach 1" -ex "set mem inaccessible-by-default off" -ex "load main.hex" -ex "compare-sections"

term:
	picocom -b 9600 /dev/ttyACM1
