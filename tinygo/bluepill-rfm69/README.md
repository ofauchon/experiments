# Project

This is a demo application for rhe RFM69 tinygo-drivers on bluepill STM32.


# Bluepill to RFM69 wiring

```
BLUEPILL  <=>     RFM69       

3.3V      <=>     3.3V
GND       <=>     GND

PA0       <=>     RST (Reset)
PA1       <=>     NSS (Chip Select, Slave Select)

PB1       <=>     DIO0  (Interrupt on packet RX)
PB0       <=>     DIO2  (Data in continuous mode)

PA9  (UART1 TX)  => Serial adapter RX
PA10 (UART1 RX)  => Serial adapter TX
```

# Toolchain installation (on Linux)

(Arch Linux preparation)
```
sudo pacman -S llvm lld arm-none-eabi-binutils arm-none-eabi-gcc arm-none-eabi-gdb arm-none-eabi-newlib tinygo
```

# Build/Flash/Debug
a
```
make build  => Build the code
make flash  => Flash with blackmagick
make debug  => Start debug session with blackmagick
```

# Connect serial port 

```
$ picocom /dev/ttyUSB0  -b 9600
```

# Commands: 

```
reset, send")
get: temp|mode|freq|regs")
set: freq <433900>")
mode: <rx,tx,standby,sleep> ")
```

# example

```
>> reset
Reset done !
>> set freq 433900000
Freq set to  433900000
>> send aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
Scheduled data to send : aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
Will send bulk
Bulk TX DONE 
INTB1:  false
Pcket sent ok in    210 ms
>> mode rx
waitformode start
waitformode ok
Mode changed !

