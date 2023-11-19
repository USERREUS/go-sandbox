package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// Генерация случайного простого числа
func generatePrime(bits int) (*big.Int, error) {
	return rand.Prime(rand.Reader, bits)
}

// Генерация примитивного корня в конечном поле
func generatePrimitiveRoot(p *big.Int) *big.Int {
	// Простой способ - найти первый примитивный корень
	g := new(big.Int).SetInt64(2)
	return g
}

// Генерация закрытого и открытого ключей
func generateKeys(p, g *big.Int) (*big.Int, *big.Int) {
	privateKey, _ := rand.Int(rand.Reader, new(big.Int).Sub(p, big.NewInt(2)))
	publicKey := new(big.Int).Exp(g, privateKey, p)
	return privateKey, publicKey
}

// Вычисление общего секрета
func computeSecretKey(publicKey, privateKey, p *big.Int) *big.Int {
	secretKey := new(big.Int).Exp(publicKey, privateKey, p)
	return secretKey
}

func main() {
	// Генерация случайного простого числа
	p, _ := generatePrime(256)

	// Генерация примитивного корня в конечном поле
	g := generatePrimitiveRoot(p)

	// Генерация ключей для Alice и Bob
	privateKeyAlice, publicKeyAlice := generateKeys(p, g)
	privateKeyBob, publicKeyBob := generateKeys(p, g)

	// Обмен открытыми ключами
	// Эти значения могут передаваться по открытому каналу
	sharedKeyAlice := computeSecretKey(publicKeyBob, privateKeyAlice, p)
	sharedKeyBob := computeSecretKey(publicKeyAlice, privateKeyBob, p)

	// Проверка, что обе стороны вычислили общий секретный ключ одинаково
	fmt.Println("Shared key for Alice:", sharedKeyAlice)
	fmt.Println("Shared key for Bob:", sharedKeyBob)
}
