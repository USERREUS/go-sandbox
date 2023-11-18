package main

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
)

// Генерация ключей RSA
func generateRSAKeys(bits int) (N, E, D *big.Int, err error) {
	// Генерация простых чисел p и q
	p, err := rand.Prime(rand.Reader, bits/2)
	if err != nil {
		return nil, nil, nil, err
	}

	q, err := rand.Prime(rand.Reader, bits/2)
	if err != nil {
		return nil, nil, nil, err
	}

	// N = p * q
	N = new(big.Int).Mul(p, q)

	// (p-1)(q-1)
	phi := new(big.Int).Mul(new(big.Int).Sub(p, big.NewInt(1)), new(big.Int).Sub(q, big.NewInt(1)))

	// Выбор открытой экспоненты E, обычно 65537
	E = big.NewInt(65537)

	// Вычисление закрытой экспоненты D
	D = new(big.Int)
	D.ModInverse(E, phi)

	return N, E, D, nil
}

// Шифрование сообщения
func encrypt(message *big.Int, E, N *big.Int) *big.Int {
	// C = M^E mod N
	return new(big.Int).Exp(message, E, N)
}

// Расшифровка сообщения
func decrypt(ciphertext *big.Int, D, N *big.Int) *big.Int {
	// M = C^D mod N
	return new(big.Int).Exp(ciphertext, D, N)
}

// Генерация электронной подписи
func sign(message *big.Int, D, N *big.Int) *big.Int {
	// Вычисление хэша SHA-256 от сообщения
	hashedMessage := sha256.Sum256(message.Bytes())

	// Преобразование хэша в целое число
	hashedInt := new(big.Int).SetBytes(hashedMessage[:])

	// Создание электронной подписи: S = H(M)^D mod N
	return new(big.Int).Exp(hashedInt, D, N)
}

// Проверка электронной подписи
func verifySignature(message, signature, E, N *big.Int) bool {
	// Проверка подписи: S^E mod N
	verified := new(big.Int).Exp(signature, E, N)

	// Вычисление хэша SHA-256 от сообщения
	hashedMessage := sha256.Sum256(message.Bytes())

	// Преобразование хэша в целое число
	hashedInt := new(big.Int).SetBytes(hashedMessage[:])

	// Проверка совпадения с полученным хэшем
	return verified.Cmp(hashedInt) == 0
}

func main() {
	// Генерация ключей
	bits := 2048
	N, E, D, err := generateRSAKeys(bits)
	if err != nil {
		fmt.Println("Ошибка при генерации ключей:", err)
		return
	}

	// Пример сообщения для подписи
	message := big.NewInt(42)

	// Создание электронной подписи
	signature := sign(message, D, N)
	fmt.Println("Электронная подпись:", signature)

	// Проверка электронной подписи
	isValid := verifySignature(message, signature, E, N)
	fmt.Println("Подпись верна:", isValid)
}
