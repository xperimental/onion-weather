package oled

import "fmt"

func (o *oled) SetCursor(row, column uint8) error {
	// check the inputs
	if row < 0 || row >= OLED_EXP_CHAR_ROWS {
		return fmt.Errorf("ERROR: Attempting to set cursor to invalid row '%d'\n", row)
	}
	if column < 0 || column >= OLED_EXP_CHAR_COLUMNS {
		return fmt.Errorf("ERROR: Attempting to set cursor to invalid column '%d'\n", column)
	}

	for _, cmd := range []uint8{
		// set page address
		OLED_EXP_ADDR_BASE_PAGE_START + row,
		// set column lower address
		OLED_EXP_SET_LOW_COLUMN + (OLED_EXP_CHAR_LENGTH * column & 0x0F),
		// set column higher address
		OLED_EXP_SET_HIGH_COLUMN + ((OLED_EXP_CHAR_LENGTH * column >> 4) & 0x0F),
	} {
		if err := sendCommand(o.dev, cmd); err != nil {
			return err
		}
	}

	return nil
}
