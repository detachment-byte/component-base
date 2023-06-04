package auth

import "testing"

func TestEncrypt(t *testing.T) {
	plainText := "chenlt"
	key, _ := Encrypt(plainText)
	t.Logf(`chenlt -> %s`, key)
}

func TestCompare(t *testing.T) {
	plainText := "chenlt"
	key, _ := Encrypt(plainText)
	res := Compare(key, plainText)
	if res == nil {
		t.Failed()
	}
}
