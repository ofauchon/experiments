This code will turn a bluepill into UART bridge for testing Lora E5 external module. 

Steps: 

- Connect computer on BP UART1 (RX: PA10, TX: PA9, baudrate: 9600)
- Connect loraE5 on BP UART2 (RX: PA3, TX: PA2, baudrate: 9600) 


Run serial console and type "AT" command. 
The AT command will be sent to LoraE5, and you should see

```
AT
  +AT: OK
```

You may use "crlf" on MacOSX while using picocom:

```
picocom --omap crlf -b 9600 /dev/ttyUSB0 
```


