
BM_DEV=/dev/ttyACM0
OUT=main.elf
BOARD=lgt92
SRC=test_gpio.go


BM_TTY_GDB=$(shell grep -n "Black Magic GDB Server" /sys/class/tty/tty*/device/interface 2>&1 | cut -d/ -f5)
BM_TTY_USART=$(shell grep -n "Black Magic UART Port" /sys/class/tty/tty*/device/interface 2>&1 | cut -d/ -f5)

OPT=z


build:
	tinygo build -size=full -target=${BOARD} -opt=${OPT} -o $(OUT) ${SRC}

flash:
	@echo "******** BM GDB TTY:"$(BM_TTY_GDB)"************"
	arm-none-eabi-gdb -nx --batch -ex "target extended-remote /dev/$(BM_TTY_GDB)" -ex "monitor swdp_scan" -ex "attach 1" -ex load -ex compare-sections -ex kill $(OUT) 


debug:
	@echo "******** BM USART TTY:"$(BM_TTY_USART)"************"
	arm-none-eabi-gdb -nx -ex "target extended-remote /dev/$(BM_TTY_GDB)"  -ex "monitor swdp_scan" -ex "attach 1" -ex "set mem inaccessible-by-default off" -ex "file $(OUT)" -ex "load" 
	#arm-none-eabi-gdb -nx -ex "target extended-remote /dev/$(BM_TTY_GDB)"  -ex "monitor swdp_scan" -ex "attach 1" -ex "set mem inaccessible-by-default off" -ex "file $(OUT)" -ex "load $(OUT)" -ex "compare-sections"

term:
	@echo "******** BM USART TTY:"$(BM_TTY_USART)"************"
	picocom -b 9600 /dev/$(BM_TTY_USART) --imap lfcrlf
