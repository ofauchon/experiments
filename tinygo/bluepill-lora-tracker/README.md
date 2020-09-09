# Project Bluepill Lora Tracker

This is a demo application for the lora/sx127x, gps tinygo-drivers on bluepill STM32.


# Bluepill, RFM69, 

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

BLUEPILL  <=>     GPS       

PA2  (UART2 TX)
PA3  (UART2 RX)  => GPS Tx
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
make term   => Start terminal
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

# References

Sx1276 module
https://www.tindie.com/products/blkbox/sx1276-lora-module-for-arduino/
