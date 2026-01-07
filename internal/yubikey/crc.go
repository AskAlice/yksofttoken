package yubikey

// CRC16 calculates the Yubikey CRC16
// The polynomial used is x^16 + x^15 + x^2 + 1 (0x8005 reversed = 0xa001)
func CRC16(data []byte) uint16 {
	var crc uint16 = 0xffff
	for _, b := range data {
		crc ^= uint16(b)
		for i := 0; i < 8; i++ {
			if crc&1 != 0 {
				crc = (crc >> 1) ^ 0x8408
			} else {
				crc >>= 1
			}
		}
	}
	return crc
}

// CRC16Valid checks if the CRC is valid (should be 0xf0b8 for valid data)
const CRCOKResidual uint16 = 0xf0b8
