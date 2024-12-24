package infraHelper

import (
	"math/rand"
)

type DummyValueGenerator struct {
	generatedUsername string
}

func (helper *DummyValueGenerator) GenPass(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+"
	charsetLen := len(charset)

	pass := make([]byte, length)
	for i := 0; i < length; i++ {
		randomIndex := rand.Intn(charsetLen)
		pass[i] = charset[randomIndex]
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
