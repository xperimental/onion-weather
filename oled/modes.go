package oled

import "errors"

func (o *oled) SetDisplayMode(inverted bool) error {
	cmd := OLED_EXP_NORMAL_DISPLAY
	if inverted {
		cmd = OLED_EXP_INVERT_DISPLAY
	}

	return sendCommand(o.dev, cmd)
}

func (o *oled) SetDisplayPower(powered bool) error {
	// int oledSetDisplayPower	(int bPowerOn);
	return errors.New("unimplemented: SetDisplayPower")
}

func (o *oled) SetBrightness(brightness uint8) error {
	// int oledSetBrightness (int brightness);
	return errors.New("unimplemented: SetBrightness")
}

func (o *oled) SetDim(dim bool) error {
	// int oledSetDim (int dim);
	return errors.New("unimplemented: SetDim")
}
