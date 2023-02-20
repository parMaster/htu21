package htu21

import (
	"fmt"
	"sync"
	"time"

	"github.com/sigurn/crc8"
	"periph.io/x/conn/v3"
	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/physic"
)

// NewI2C returns an object that communicates over I²C to HTU21
// environmental sensor.
//
// HTU21 I²C device address is 0x40, according to the datasheet
// https://www.te.com/usa-en/models/4/00/028/790/CAT-HSC0004.html
func NewI2C(b i2c.Bus, addr uint16) (*Dev, error) {
	d := &Dev{d: &i2c.Dev{Bus: b, Addr: addr}}
	if err := d.Reset(); err != nil {
		return nil, err
	}
	return d, nil
}

// Dev is a handle to an initialized HTU21 device.
type Dev struct {
	d  conn.Conn
	mx sync.Mutex
}

// Performs soft reset of the device
// 0xFE - code for soft reset
// 15 milliseconds delay follows
func (d *Dev) Reset() error {
	d.mx.Lock()
	defer d.mx.Unlock()

	err := d.writeCommands([]byte{0xFE})
	if err != nil {
		return d.wrap(err)
	}
	doSleep(time.Millisecond * 15)
	return nil
}

// Sense requests a one time measurement as °C and % of relative humidity.
// 0xF3, delay 50 ms - temperature
// 0xF5, delay 16 ms - humidity
func (d *Dev) Sense(e *physic.Env) error {
	d.mx.Lock()
	defer d.mx.Unlock()

	// Temperature | 0xF3, delay 50 ms
	err := d.writeCommands([]byte{0xF3})
	if err != nil {
		return d.wrap(err)
	}
	doSleep(50 * time.Millisecond)

	rawData, err := d.read()
	if err != nil {
		return d.wrap(err)
	}

	data, err := d.validate(rawData)
	if err != nil {
		return d.wrap(err)
	}

	temp := -46.85 + (175.72 * (float32(data) / 65536))
	e.Temperature.Set(fmt.Sprintf("%fC", temp))

	// Humidity | 0xF5, delay 16 ms
	err = d.writeCommands([]byte{0xF5})
	if err != nil {
		return d.wrap(err)
	}
	doSleep(16 * time.Millisecond)

	rawData, err = d.read()
	if err != nil {
		return d.wrap(err)
	}

	data, err = d.validate(rawData)
	if err != nil {
		return d.wrap(err)
	}

	var humidity = float32(data)
	humidity = -6 + (125 * (humidity / 65536))

	e.Humidity.Set(fmt.Sprintf("%f%%", humidity))

	return nil
}

var crc8Table = crc8.MakeTable(crc8.Params{
	Poly:   0x31,
	Init:   0x0,
	RefIn:  false,
	RefOut: false,
	XorOut: 0x0,
	Check:  0x0,
})

// CRC Checksum using the polynomial given in the datasheet
func validateCrc(data []byte, expectedCrc byte) error {
	checksum := crc8.Checksum(data, crc8Table)

	if expectedCrc != checksum {
		return fmt.Errorf("CRC Checksum error. Expected %v, actual %v", expectedCrc, checksum)
	}
	return nil
}

func (d *Dev) validate(rawData []byte) (uint16, error) {

	if e := validateCrc(rawData[:2], rawData[2]); e != nil {
		return 0, e
	}

	data := uint16(rawData[0])
	data <<= 8
	data |= uint16(rawData[1])

	return data, nil
}

// writeCommands writes a command to the device.
// Warning: b may be modified!
func (d *Dev) writeCommands(b []byte) error {
	if err := d.d.Tx(b, nil); err != nil {
		return d.wrap(err)
	}
	return nil
}

// read from the device
func (d *Dev) read() ([]byte, error) {
	b := make([]byte, 3)
	if err := d.d.Tx(nil, b); err != nil {
		return nil, d.wrap(err)
	}
	return b, nil
}

func (d *Dev) wrap(err error) error {
	return fmt.Errorf("%e", err)
}

var doSleep = time.Sleep
