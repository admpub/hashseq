// Package hashseq This package provides ID masking using hash ids.
//
// Hash IDs are a string represnetation of numerical incrementing IDs,
// obfuscating the integer value. For more information see
// http://hashids.org/go/
package hashseq

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	hashid "github.com/admpub/go-hashids"
)

var (
	hashData = &hashid.HashIDData{
		Alphabet:   hashid.DefaultAlphabet,
		MinLength:  4,
		Salt:       "",
		Uint64Mode: true,
	}
	hashID *hashid.HashID
)

// HashID .
func HashID() *hashid.HashID {
	if hashID == nil {
		var err error
		hashID, err = hashid.NewWithData(hashData)
		if err != nil {
			panic(err)
		}
	}
	return hashID
}

// SetSalt Set the salt to use for ID obfuscation
func SetSalt(salt string) {
	if hashData.Salt == salt {
		return
	}
	hashData.Salt = salt
	var err error
	hashID, err = hashid.NewWithData(hashData)
	if err != nil {
		panic(err)
	}
}

// ID uint64
type ID uint64

// Return the hashid as an obfuscated string
func (id *ID) String() string {
	str, err := HashID().EncodeUint64([]uint64{id.Uint64()})
	if err != nil {
		return ""
	}
	return str
}

func (id ID) Uint64() uint64 {
	return uint64(id)
}

// MarshalJSON Returns the hashid as a string fulfilling the json.Marshaller interface
func (id ID) MarshalJSON() ([]byte, error) {
	str, err := HashID().EncodeUint64([]uint64{id.Uint64()})
	if err != nil {
		return nil, err
	}
	return json.Marshal(str)
}

// UnmarshalJSON Unmarshal a string and decode into an integer
func (id *ID) UnmarshalJSON(data []byte) error {
	decoded, err := Decode(data)
	if err != nil {
		return err
	}
	*id = ID(decoded)
	return nil
}

// Decode a hashid byte into an Id, setting its integer
func Decode(hashid []byte) (id uint64, err error) {
	return DecodeString(string(hashid))
}

// MustDecodeString .
func MustDecodeString(hashid string) uint64 {
	i, err := DecodeString(hashid)
	if err != nil {
		panic(err.Error())
	}
	return uint64(i)
}

// DecodeString Decode a hashid string into an Id, setting its integer
func DecodeString(h string) (id uint64, err error) {
	var ids []uint64
	ids, err = HashID().DecodeUint64WithError(h)
	if err != nil {
		return
	}
	return ids[0], nil
}

// Scan Database scanning
func (id *ID) Scan(value interface{}) (err error) {
	var data uint64

	// If the first four bytes of this are 0000
	switch v := value.(type) {
	// Same as []byte
	case uint64:
		data = v
	case nil:
		return
	default:
		return fmt.Errorf("Invalid format: can't convert %T into id.Id", value)
	}

	*id = ID(data)
	return nil
}

// Value This is called when saving the ID to a database
func (id ID) Value() (driver.Value, error) {
	return id.Uint64(), nil
}
