package lib

import (
	"log"
	"testing"
)



func TestCrypto(t *testing.T) {
	log.Println(GenSaltPassword("admin", "1234567"))
}