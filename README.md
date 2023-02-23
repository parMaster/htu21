# HTU21 (SHT21, SI7021) golang package

This package allows you to connect and read environmental data from HTU21 (SHT21, SI7021) sensor. I guess it could be qualified as a Periph.io driver, because it tries to comply with periph.io sensors interface like [bmxx80](https://periph.io/device/bmxx80/) in terms of the ways to access I²C, function calls and data returned.

## Example
You can find an example in cmd folder, it works on RPi 4 with HTU21 connected to 4-th I²C bus, sensor address is 0x40. List of I²C buses can be obtained with 

`i2cdetect -l`

Once you know the bus number, connected devices can be enumerated with

`i2cdetect -y 4`

## Additional links

Some links you can find useful working with this sensor, RPi GPIO and I²C:

- Sensor datasheet is here - https://www.te.com/usa-en/models/4/00/028/790/CAT-HSC0004.html
- Wonderful Raspberry Pi pinout guide - https://pinout.xyz/
- How to define additional/alternative I²C buses on RPi GPIO - https://www.instructables.com/Raspberry-PI-Multiple-I2c-Devices/
- *"Access GPIO pins without root. No access to /dev/gpiomem. Try running as root!"* problem solution - https://raspberrypi.stackexchange.com/questions/40105/access-gpio-pins-without-root-no-access-to-dev-mem-try-running-as-root

## Why does this package exist?

This package is a heavily rewritten [idahoakl/HTU21D-sensor](https://github.com/idahoakl/HTU21D-sensor/) library. Main differences are:

-  [periph.io i2c](http://periph.io/x/conn/v3/i2c) package to access I²C instead of proprietary one, idahoakl used
- [periph.io physic](http://periph.io/x/conn/v3/physic) package for environmental data (temperature and relative humidity)
- one function *Sense()* instead of separate functions for each measurement 
- no proprietary logging package
- no write/read user registry functions **for now**
  
I'd really like to get rid of proprietary CRC checking dependency as well.
