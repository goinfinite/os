package infraHelper

import (
	"math/rand"
	"strings"
)

type DummyValueGenerator struct {
	generatedUsername string
}

func (helper *DummyValueGenerator) oneCharCharsetGuarantor(
	originalString []byte,
	charset string,
) []byte {
	if strings.ContainsAny(string(originalString), charset) {
		return originalString
	}

	randomStringIndex := rand.Intn(len(originalString))
	isFirstChar := randomStringIndex == 0
	if isFirstChar {
		randomStringIndex++
	}
	isLastChar := randomStringIndex == len(originalString)-1
	if isLastChar {
		randomStringIndex--
	}
	if randomStringIndex >= len(originalString) {
		randomStringIndex = len(originalString) - 1
	}

	randomCharsetIndex := rand.Intn(len(charset))
	originalString[randomStringIndex] = charset[randomCharsetIndex]

	return originalString
}

func (helper *DummyValueGenerator) GenPass(length int) string {
	lowercaseAlphabetCharset := "abcdefghijklmnopqrstuvwxyz"
	uppercaseAlphabetCharset := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numericCharset := "0123456789"
	alphanumericCharset := lowercaseAlphabetCharset + uppercaseAlphabetCharset + numericCharset
	symbolCharset := "!@#$%^&*()_+"

	pass := make([]byte, length)
	previousCharset := alphanumericCharset
	currentCharset := alphanumericCharset

	for i := 0; i < length; i++ {
		currentCharset = alphanumericCharset

		if previousCharset != symbolCharset {
			nextCharsetShouldUseSymbol := rand.Float32() < 0.1
			isTipChar := i == 0 || i == length-1
			if nextCharsetShouldUseSymbol && !isTipChar {
				currentCharset = symbolCharset
			}
		}
		currentCharsetLen := len(currentCharset)

		randomCharsetIndex := rand.Intn(currentCharsetLen)
		pass[i] = currentCharset[randomCharsetIndex]
		previousCharset = currentCharset
	}

	if length > 4 {
		pass = helper.oneCharCharsetGuarantor(pass, lowercaseAlphabetCharset)
		pass = helper.oneCharCharsetGuarantor(pass, uppercaseAlphabetCharset)
		pass = helper.oneCharCharsetGuarantor(pass, numericCharset)
		pass = helper.oneCharCharsetGuarantor(pass, "!@#$%^&*()_+")
	}

	return string(pass)
}

func (helper *DummyValueGenerator) GenUsername() string {
	dummyUsernames := []string{
		"yoda", "obi_wan", "anakin", "luke", "leia", "rey", "kylo",
	}
	helper.generatedUsername = dummyUsernames[rand.Intn(len(dummyUsernames))]

	return helper.generatedUsername
}

func (helper *DummyValueGenerator) GenMailAddress(mailUsername *string) string {
	dummyMailAddressDomains := []string{
		"@republic.gov", "@senate.gov", "@empire.gov",
	}
	randomMailAddressDomain := dummyMailAddressDomains[rand.Intn(len(dummyMailAddressDomains))]

	if mailUsername == nil {
		if helper.generatedUsername == "" {
			helper.GenUsername()
		}

		mailUsername = &helper.generatedUsername
	}

	return *mailUsername + randomMailAddressDomain
}
