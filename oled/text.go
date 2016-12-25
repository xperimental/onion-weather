package oled

func (o *oled) Write(text string) error {
	// set addressing mode to page
	//oledSetMemoryMode(OLED_EXP_MEM_PAGE_ADDR_MODE);	// want automatic newlines enabled

	// set column addressing to fit 126 characters that are 6 pixels wide
	if !o.bColumnsSetForText {
		if err := o.SetTextColumns(); err != nil {
			return err
		}
	}

	// write each character
	for _, char := range text {
		// check for newline character
		if char == '\n' {
			// move the cursor to the next row
			for i := o.cursorInRow; i <= OLED_EXP_CHAR_COLUMNS; i++ {
				if err := o.WriteChar(' '); err != nil {
					return err
				}
			}
		}

		if err := o.WriteChar(char); err != nil {
			return err
		}
	}

	// reset the column addressing
	return o.SetImageColumns()
}

func (o *oled) WriteChar(char rune) error {
	// ensure character is in the table
	buf := mapASCII(char)
	for _, data := range buf {
		if err := sendData(o.dev, data); err != nil {
			return err
		}

		// increment row cursor
		if o.cursorInRow == OLED_EXP_CHAR_COLUMNS-1 {
			o.cursorInRow = 0
		} else {
			o.cursorInRow++
		}
	}

	return nil
}
