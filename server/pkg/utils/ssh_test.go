package utils

import "testing"

func TestVerifySsh(t *testing.T) {
	err := VerifySsh("192.168.100.30", 22, "ttt", "red")
	if err != nil {
		t.Log(err)
		t.Fail()
	}
}
