package oled

func (o *oled) SetDisplayInverted(inverted bool) error {
	cmd := OLED_EXP_NORMAL_DISPLAY
	if inverted {
		cmd = OLED_EXP_INVERT_DISPLAY
	}

	return sendCommand(o.dev, cmd)
}

func (o *oled) SetDisplayPower(powered bool) error {
	cmd := OLED_EXP_DISPLAY_OFF
	if powered {
		cmd = OLED_EXP_DISPLAY_ON
	}

	return sendCommand(o.dev, cmd)
}

func (o *oled) SetBrightness(brightness uint8) error {
	if brightness < OLED_EXP_CONTRAST_MIN {
		brightness = OLED_EXP_CONTRAST_MIN
	}
	if brightness > OLED_EXP_CONTRAST_MAX {
		brightness = OLED_EXP_CONTRAST_MAX
	}

	for _, cmd := range []uint8{
		OLED_EXP_SET_CONTRAST,
		brightness,
	} {
		if err := sendCommand(o.dev, cmd); err != nil {
			return err
		}
	}

	return nil
}

func (o *oled) SetDim(dim bool) error {
	if dim {
		return o.SetBrightness(0)
	}

	return o.SetBrightness(207)
}
