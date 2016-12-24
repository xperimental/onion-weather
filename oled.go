package main

import (
	"fmt"
	"time"

	"golang.org/x/exp/io/i2c"
)

const (
	i2CDevPath = "/dev/i2c-%d"
)

type Oled struct {
	dev        *i2c.Device
	vccState   uint8
	memoryMode int

	buffer [int(OLED_EXP_WIDTH) * int(OLED_EXP_PAGES)]uint8
	cursor int

	cursorInRow        uint8
	bColumnsSetForText bool
}

func NewOled() (*Oled, error) {
	path := fmt.Sprintf(i2CDevPath, i2CDevNum)
	dev, err := i2c.Open(&i2c.Devfs{Dev: path}, i2COledAddr)
	if err != nil {
		return nil, err
	}

	return &Oled{
		dev: dev,
	}, nil
}

func (o *Oled) Init() error {
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

	if err := o.SetDisplayMode(false); err != nil {
		return err
	}

	if err := o.Clear(); err != nil {
		return err
	}

	return nil
}

func (o *Oled) SetDisplayMode(inverted bool) error {
	cmd := OLED_EXP_NORMAL_DISPLAY
	if inverted {
		cmd = OLED_EXP_INVERT_DISPLAY
	}

	return sendCommand(o.dev, cmd)
}

func (o *Oled) Clear() error {
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

func (o *Oled) SetTextColumns() error {
	// set the column addressing for text mode
	if err := o.SetColumnAddressing(0, OLED_EXP_CHAR_COLUMN_PIXELS-1); err != nil {
		return err
	}
	o.bColumnsSetForText = true

	return nil
}

func (o *Oled) SetImageColumns() error {
	// set the column addressing to full width
	if err := o.SetColumnAddressing(0, OLED_EXP_WIDTH-1); err != nil {
		return err
	}
	o.bColumnsSetForText = false

	return nil
}

func (o *Oled) SetColumnAddressing(startPixel, endPixel uint8) error {
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

func (o *Oled) SetCursor(row, column uint8) error {
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

func (o *Oled) Write(text string) error {
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

func (o *Oled) WriteChar(char rune) error {
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

// Close closes the device connection.
func (o *Oled) Close() error {
	return o.dev.Close()
}

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
