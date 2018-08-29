package escpos

const (
	ESC byte = 27
	FS  byte = 28
	GS  byte = 29
	DLE byte = 16
	EOT byte = 4
	ENQ byte = 5
	SP  byte = 32
	HT  byte = 9
	LF  byte = 10
	CR  byte = 13
	FF  byte = 12
	SO  byte = 14
	CAN byte = 24
	AT  byte = 64
	I   byte = 105

	PC437    byte = 0
	KATAKANA byte = 1
	PC850    byte = 2
	PC860    byte = 3
	PC863    byte = 4
	PC865    byte = 5
	WPC1252  byte = 16
	PC866    byte = 17
	PC852    byte = 18
	PC858    byte = 19

	/**
	 * BarCode table
	 */
	UPC_A   byte = 0
	UPC_E   byte = 1
	EAN13   byte = 2
	EAN8    byte = 3
	CODE39  byte = 4
	ITF     byte = 5
	NW7     byte = 6
	CODE93  byte = 72
	CODE128 byte = 73
)
