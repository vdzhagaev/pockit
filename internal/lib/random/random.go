package random

import (
	"crypto/rand"
	"fmt"
)

type Generator struct {
	alphabet string
	length   int
}

func (g *Generator) Generate() (string, error) {
	bytes := make([]byte, g.length)
	alphabetBytes := []byte(g.alphabet)
	alphabetLen := byte(len(alphabetBytes))

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	for i := range bytes {
		bytes[i] = alphabetBytes[bytes[i]%alphabetLen]
	}

	return string(bytes), nil
}

func NewGenerator(alphabet string, length int) *Generator {
	if alphabet == "" {
		alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	}
	return &Generator{
		alphabet: alphabet,
		length:   length,
	}
}

func NewRandomString(size int) (string, error) {
	const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-"
	generator := NewGenerator(alphabet, size)

	str, err := generator.Generate()

	if err != nil {
		return "", fmt.Errorf("failed to generate random string: %w", err)
	}

	return str, nil
}
