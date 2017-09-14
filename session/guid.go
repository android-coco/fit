package session

import (
	"crypto/rand"
	"encoding/hex"
)

// Used in string method conversion
const dash byte = '-'

func safeRandom(dest []byte) {
	if _, err := rand.Read(dest); err != nil {
		panic(err)
	}
}

// uuid representation compliant with specification
type uuid [16]byte

// setVersion sets version bits.
func (u *uuid) setVersion(v byte) {
	u[6] = (u[6] & 0x0f) | (v << 4)
}

// setVariant sets variant bits as described in RFC 4122.
func (u *uuid) setVariant() {
	u[8] = (u[8] & 0xbf) | 0x80
}

// returns random generated UUID.
func GUID() string {
	u := uuid{}
	safeRandom(u[:])
	u.setVersion(4)
	u.setVariant()
	return u.string()
}

// Returns canonical string representation of UUID:
// xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx.
func (u uuid) string() string {
	buf := make([]byte, 36)

	hex.Encode(buf[0:8], u[0:4])
	buf[8] = dash
	hex.Encode(buf[9:13], u[4:6])
	buf[13] = dash
	hex.Encode(buf[14:18], u[6:8])
	buf[18] = dash
	hex.Encode(buf[19:23], u[8:10])
	buf[23] = dash
	hex.Encode(buf[24:], u[10:])

	return string(buf)
}
