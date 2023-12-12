package random

import (
	"math/rand"
	"strings"
)

const (
	// base
	lowercaseAlphabet = "abcdefghijklmnopqrstuvwxyz"
	uppercaseAlphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	nums              = "0123456789"
	special           = "-_"

	// combinations
	alphabet = lowercaseAlphabet + uppercaseAlphabet
	alphanum = alphabet + nums
	chars    = alphanum + special
)

type CharsetConfig int

const (
	Lowercase CharsetConfig = 1 << iota
	Uppercase
	Nums
	Special
	Alphabet = Lowercase | Uppercase
	Alphanum = Alphabet | Nums
	All      = Alphanum | Special
)

func GetCharset(config CharsetConfig) string {
	if config == 0 {
		config = All
	}

	var charset string

	if config&Lowercase != 0 {
		charset += lowercaseAlphabet
	}

	if config&Uppercase != 0 {
		charset += uppercaseAlphabet
	}

	if config&Nums != 0 {
		charset += nums
	}

	if config&Special != 0 {
		charset += special
	}

	return charset
}

func Int() int {
	return rand.Int()
}

// Intr returns a random int between min and max
func Intr(min, max int) int {
	return min + rand.Intn(max-min+1)
}

// Int32 returns a random int32 between min and max
func Int32(min, max int32) int32 {
	return min + rand.Int31n(max-min+1)
}

func StringFrom(n int, charset string) string {
	var sb strings.Builder
	k := len(charset)

	for i := 0; i < n; i++ {
		c := charset[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func StringWith(n int, config CharsetConfig) string {
	charset := GetCharset(config)
	return StringFrom(n, charset)
}

// String returns a random string (alphanumeric) with length n
func String(n int) string {
	var sb strings.Builder
	k := len(alphanum)

	for i := 0; i < n; i++ {
		c := alphanum[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// StringL returns a random string with length between min and max
func StringL(min, max int) string {
	n := Intr(min, max)
	return String(n)
}

// StringN returns a random string with length n
func Username() string {
	return String(8)
}

var emailDomains = []string{"gmail.com", "yahoo.com", "example.com", "hotmail.com", "outlook.com", "mail.com"}

// EmailDomain returns a random email domain
func EmailDomain() string {
	n := len(emailDomains)
	return emailDomains[rand.Intn(n)]
}

// Email returns a random email address
func Email() string {
	return String(8) + "@" + EmailDomain()
}
