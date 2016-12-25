package oled

import "golang.org/x/exp/io/i2c"

func sendCommand(dev *i2c.Device, command uint8) error {
	buf := []byte{
		OLED_EXP_REG_COMMAND,
		command,
	}
	return dev.Write(buf)
}

func sendData(dev *i2c.Device, data uint8) error {
	buf := []byte{
		OLED_EXP_REG_DATA,
		data,
	}
	return dev.Write(buf)
}
