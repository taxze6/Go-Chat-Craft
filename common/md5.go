package common

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
)

// The Md5encoder returns the lowercase value after encryption
func Md5encoder(code string) string {
	m := md5.New()
	_, _ = io.WriteString(m, code)
	return hex.EncodeToString(m.Sum(nil))
}

// The Md5StrToUpper returns the uppercase value after encryption.
func Md5StrToUpper(code string) string {
	return strings.ToUpper(Md5encoder(code))
}

// The SaltPassWord function adds salt to the password.
func SaltPassWord(pw string, salt string) string {
	saltPW := fmt.Sprintf("%s$%s", Md5encoder(pw), salt)
	return saltPW
}

// The CheckPassWord function verifies the password
func CheckPassWord(rpw, salt, pw string) bool {
	return pw == SaltPassWord(rpw, salt)
}
