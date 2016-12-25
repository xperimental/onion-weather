package oled

// Display can be used to access the OLED expansion on an Onion Omega.
type Display interface {
	// Init initializes and clears the display.
	Init() error

	// Close closes the device connection.
	Close() error

	// Clear clears the display.
	Clear() error

	// Write draws text on the display.
	Write(text string) error
	// WriteChar draws a single character on the display.
	// There are only a few ASCII chars available, others will be blank.
	WriteChar(char rune) error

	// SetDisplaPower toggles the display on or off.
	SetDisplayPower(powered bool) error
	// SetDisplayMode toggles the display between normal and inverted mode.
	SetDisplayMode(inverted bool) error
	// SetBrightness sets the brightness of the display (default is 207).
	SetBrightness(brightness uint8) error
	// SetDim toggles the display between dim (0) and default (207) brightness.
	SetDim(dim bool) error

	// SetColumnAddressing sets the display column addressing.
	SetColumnAddressing(startPixel, endPixel uint8) error
	// SetTextColumns sets the display column addressing to text mode.
	SetTextColumns() error
	// SetImageColumns sets the display column addressing to full width.
	SetImageColumns() error

	// SetCursor positions the cursor at the specified row and column.
	SetCursor(row, column uint8) error

	/* C functions:
	implemented:
	int oledDriverInit ();
	int oledSetDisplayMode (int bInvert);
	int oledSetColumnAddressing (int startPixel, int endPixel);
	int oledSetTextColumns ();
	int oledSetImageColumns ();
	int oledSetCursor (int row, int column);
	int oledClear ();
	int oledWriteChar (char c);
	int oledWrite (char *msg);

	existing:
	int oledSetDisplayPower	(int bPowerOn);
	int oledSetBrightness (int brightness);
	int oledSetDim (int dim);

	unimplemented:
	int oledSetMemoryMode (int mode);
	int oledSetCursorByPixel (int row, int pixel);
	int oledWriteByte (int byte);
	int oledReadLcdFile	(char* file, uint8_t *buffer);
	int oledDraw (uint8_t *buffer, int bytes);
	int oledScroll 	(int direction, int scrollSpeed, int startPage, int stopPage);
	int oledScrollDiagonal 	(int direction, int scrollSpeed, int fixedRows, int scrollRows, int verticalOffset, int startPage, int stopPage);
	int oledScrollStop ();
	*/
}
