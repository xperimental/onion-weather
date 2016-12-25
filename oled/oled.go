// Package oled can be used to access the OLED expansion on an Onion Omega.
package oled

import (
	"fmt"

	"golang.org/x/exp/io/i2c"
)

type oled struct {
	dev        *i2c.Device
	vccState   uint8
	memoryMode int

	buffer [int(OLED_EXP_WIDTH) * int(OLED_EXP_PAGES)]uint8
	cursor int

	cursorInRow        uint8
	bColumnsSetForText bool
}

// NewOled creates a new Oled struct connected to I2C.
func NewOled() (Display, error) {
	path := fmt.Sprintf(i2CDevPath, i2CDevNum)
	dev, err := i2c.Open(&i2c.Devfs{Dev: path}, i2COledAddr)
	if err != nil {
		return nil, err
	}

	return &oled{
		dev: dev,
	}, nil
}

func (o *oled) Close() error {
	return o.dev.Close()
}
