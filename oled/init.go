package oled

import "time"

func (o *oled) Init() error {
	o.vccState = OLED_EXP_SWITCH_CAP_VCC

	if err := sendCommand(o.dev, OLED_EXP_DISPLAY_OFF); err != nil {
		return err
	}
	time.Sleep(4500 * time.Microsecond)

	for _, c := range []uint8{
		OLED_EXP_SET_DISPLAY_CLOCK_DIV,
		0x80, // The suggested ratio 0x80
		OLED_EXP_SET_MULTIPLEX,
		0x3F,
		OLED_EXP_SET_DISPLAY_OFFSET,
		0x00, // no offset
		OLED_EXP_SET_START_LINE | 0x00, // Set start line to line #0
		OLED_EXP_CHARGE_PUMP,
	} {
		if err := sendCommand(o.dev, c); err != nil {
			return err
		}
	}
	if o.vccState == OLED_EXP_EXTERNAL_VCC {
		if err := sendCommand(o.dev, 0x10); err != nil {
			return err
		}
	} else {
		if err := sendCommand(o.dev, 0x14); err != nil {
			return err
		}
	}

	for _, c := range []uint8{
		OLED_EXP_MEMORY_MODE,
		OLED_EXP_MEM_HORIZONTAL_ADDR_MODE,
		OLED_EXP_SEG_REMAP | 0x01,
		OLED_EXP_COM_SCAN_DEC,
		OLED_EXP_SET_COM_PINS,
		0x12,
		OLED_EXP_SET_CONTRAST,
	} {
		if err := sendCommand(o.dev, c); err != nil {
			return err
		}
	}
	if o.vccState == OLED_EXP_EXTERNAL_VCC {
		if err := sendCommand(o.dev, OLED_EXP_DEF_CONTRAST_EXTERNAL_VCC); err != nil {
			return err
		}
	} else {
		if err := sendCommand(o.dev, OLED_EXP_DEF_CONTRAST_SWITCH_CAP_VCC); err != nil {
			return err
		}
	}
	if err := sendCommand(o.dev, OLED_EXP_SET_PRECHARGE); err != nil {
		return err
	}
	if o.vccState == OLED_EXP_EXTERNAL_VCC {
		if err := sendCommand(o.dev, 0x22); err != nil {
			return err
		}
	} else {
		if err := sendCommand(o.dev, 0xF1); err != nil {
			return err
		}
	}
	for _, c := range []uint8{
		OLED_EXP_SET_VCOM_DETECT,
		0x40,
		OLED_EXP_DISPLAY_ALL_ON_RESUME,
		OLED_EXP_NORMAL_DISPLAY,
		OLED_EXP_SEG_REMAP,
		OLED_EXP_COM_SCAN_INC,
		OLED_EXP_DISPLAY_ON,
	} {
		if err := sendCommand(o.dev, c); err != nil {
			return err
		}
	}
	time.Sleep(4500 * time.Microsecond)

	if err := o.SetDisplayInverted(false); err != nil {
		return err
	}

	if err := o.Clear(); err != nil {
		return err
	}

	return nil
}

func (o *oled) Clear() error {
	// set the column addressing for the full width
	if err := o.SetImageColumns(); err != nil {
		return err
	}

	// display off
	if err := sendCommand(o.dev, OLED_EXP_DISPLAY_OFF); err != nil {
		return err
	}

	// write a blank space to each character
	for row := uint8(0); row < OLED_EXP_CHAR_ROWS; row++ {
		if err := o.SetCursor(row, 0); err != nil {
			return err
		}

		for col := uint8(0); col < OLED_EXP_WIDTH; col++ {
			if err := sendData(o.dev, 0x00); err != nil {
				return err
			}
		}
	}

	// display on
	if err := sendCommand(o.dev, OLED_EXP_DISPLAY_ON); err != nil {
		return err
	}

	// reset the cursor to (0, 0)
	return o.SetCursor(0, 0)
}
