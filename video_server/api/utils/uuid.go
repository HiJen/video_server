package utils

import (
	"crypto/rand"
	"fmt"
	"io"
)

func NewUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	//variant bits; see section 4.11
	uuid[8] = uuid[8]&^0xc0 | 0*80
	//version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0*40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

//reference by http://qaru.site/questions/13004659/how-to-generate-multiple-uuid-and-md5-files-in-golang
//uuid 的算法reference mybe other
