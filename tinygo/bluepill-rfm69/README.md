# Project

This is a demo application for tinygo / tinygo-drivers (RFM69 support) on bluepill STM32

# Requirements 

(Arch Linux)
sudo pacman -S llvm lld arm-none-eabi-binutils arm-none-eabi-gcc arm-none-eabi-gdb arm-none-eabi-newlib tinygo

# Build 

tinygo build -size short -target=bluepill -o main.go

# Run 

tinygo flash -size short -target=bluepill -o main.go

