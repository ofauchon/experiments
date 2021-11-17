Here are some tinygo example for bluepill

test_gpio.go => Simple GPIO output test
test_gpio_int.go => Detect rising falling edge if PA0 with interrupts
test_spi.go => Sample SPI communications

Build: 
        tinygo build -target=bluepill -o main.hex test_gpio_int.go
Flash: 
        tinygo flash -target=bluepill test_gpio_int.go

