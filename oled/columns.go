package oled

import "fmt"

func (o *oled) SetTextColumns() error {
	// set the column addressing for text mode
	if err := o.SetColumnAddressing(0, OLED_EXP_CHAR_COLUMN_PIXELS-1); err != nil {
		return err
	}
	o.bColumnsSetForText = true

	return nil
}

func (o *oled) SetImageColumns() error {
	// set the column addressing to full width
	if err := o.SetColumnAddressing(0, OLED_EXP_WIDTH-1); err != nil {
		return err
	}
	o.bColumnsSetForText = false

	return nil
}

func (o *oled) SetColumnAddressing(startPixel, endPixel uint8) error {
	// check the inputs
	if startPixel < 0 || startPixel >= OLED_EXP_WIDTH || startPixel >= endPixel {
		return fmt.Errorf("ERROR: Invalid start pixel (%d) for column address setup\n", startPixel)
	}
	if endPixel < 0 || endPixel >= OLED_EXP_WIDTH {
		return fmt.Errorf("ERROR: Invalid end pixel (%d) for column address setup\n", endPixel)
	}

	// set column addressing
	for _, cmd := range []uint8{
		OLED_EXP_COLUMN_ADDR,
		startPixel,
		endPixel,
	} {
		if err := sendCommand(o.dev, cmd); err != nil {
			return err
		}
	}

	return nil
}
