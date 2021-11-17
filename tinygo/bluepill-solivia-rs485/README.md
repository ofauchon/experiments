#TinyGo Lora Gateway for Delta Solivia Inverters


Delta Solivia Inverters provide a RS845 communication bus to exchange informations with other devices.
Many informations like Power, Voltages, Current can be fetched by quering the bus. 

This project details both hardware interface design and application code for communicating with inverters.

The application implements : 

 * RS485 communications with Delta Solivia inverter
 * RS232 communications with serial console for debugging
 * Lora/Lorawan communications through SX127x driver


The hardware part consists of :

 * STM32F103 (bluepill board)
 * RFM95 Lora board
 * RS485-TTL interface

