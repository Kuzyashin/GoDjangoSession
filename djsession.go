package DjangoSession

import (
	"bytes"
	"compress/zlib"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

const salt = "django.contrib.sessionsSessionStore"

var defaultSep = []byte{':'}

func djangoSignature(salt string, value []byte, secret string) []byte {
	key := make([]byte, 0, len(salt)+len(secret))
	key = append(key, salt...)
	key = append(key, secret...)
	mac := hmac.New(sha1.New, key)
	mac.Write(value)

	return []byte(hex.EncodeToString(mac.Sum(nil)))
}

func unsign(secret string, cookie []byte) ([]byte, error) {
	splitted := strings.SplitN(string(cookie), string(defaultSep), 2)
	sig := splitted[0]
	val := []byte(splitted[1])
	expectedSig := djangoSignature(salt, val, secret)
	if subtle.ConstantTimeCompare([]byte(sig), expectedSig) != 1 {
		return nil, fmt.Errorf("signature mismatch: '%s' != '%s'", sig, string(expectedSig))
	}
	return val, nil
}

func timestampUnsign(secret string, cookie []byte) ([]byte, error) {
	val, err := unsign(secret, cookie)
	if err != nil {
		return nil, fmt.Errorf("unsign('%s'): %s", string(cookie), err)
	}
	return val, nil
}

func signingLoads(secret, cookie string) (map[string]interface{}, error) {
	c := []byte(cookie) // XXX: does this escape?
	payload, err := timestampUnsign(secret, c)
	if err != nil {
		return nil, fmt.Errorf("timestampUnsign: %s", err)
	}
	decompress := false
	if payload[0] == '.' {
		decompress = true
		payload = payload[1:]
	}
	if decompress {
		r, err := zlib.NewReader(bytes.NewReader(payload))
		if err != nil {
			return nil, fmt.Errorf("zlib.NewReader: %s", err)
		}
		payload, err = ioutil.ReadAll(r)
		r.Close()
		if err != nil {
			return nil, fmt.Errorf("ReadAll(zlib): %s", err)
		}
	}
	o := make(map[string]interface{})
	json.Unmarshal(payload, &o)
	return o, nil
}

func Decode(secret, cookie string) (map[string]interface{}, error) {
	return signingLoads(secret, cookie)
}
