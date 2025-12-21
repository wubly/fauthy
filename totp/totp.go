package totp

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"strings"
	"time"
)

func Generate(secret string, t time.Time) (string, int) {
	secret = strings.ToUpper(strings.ReplaceAll(secret, " ", ""))
	
	// Add padding if needed
	switch len(secret) % 8 {
	case 2, 4, 5, 7:
		secret += strings.Repeat("=", 8-(len(secret)%8))
	}
	
	key, err := base32.StdEncoding.DecodeString(secret)
	if err != nil {
		return "------", 0
	}

	counter := uint64(t.Unix() / 30)

	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], counter)

	h := hmac.New(sha1.New, key)
	h.Write(buf[:])
	sum := h.Sum(nil)

	offset := sum[len(sum)-1] & 0x0f
	code := (int(sum[offset])&0x7f)<<24 |
		(int(sum[offset+1])&0xff)<<16 |
		(int(sum[offset+2])&0xff)<<8 |
		(int(sum[offset+3]) & 0xff)

	code %= 1000000
	remaining := 30 - int(t.Unix()%30)

	return fmt.Sprintf("%06d", code), remaining
}
