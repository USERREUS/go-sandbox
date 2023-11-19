package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func modExp(base, exponent, modulus *big.Int) *big.Int {
	local_base := new(big.Int).Set(base)
	local_exponent := new(big.Int).Set(exponent)
	local_modulus := new(big.Int).Set(modulus)

	result := new(big.Int).SetInt64(1)
	for local_exponent.BitLen() > 0 {
		if local_exponent.Bit(0) == 1 {
			result = new(big.Int).Mod(new(big.Int).Mul(result, local_base), local_modulus)
		}
		local_base = new(big.Int).Mod(new(big.Int).Mul(local_base, local_base), local_modulus)
		local_exponent.Rsh(local_exponent, 1)
	}
	return result
}

func generateKeyPair() (*big.Int, *big.Int, *big.Int) {
	// Генерация простого числа p
	p, _ := rand.Prime(rand.Reader, 128)

	// Выбор генератора g
	g := new(big.Int).SetInt64(2) // В реальных системах нужно выбирать более сложные значения

	// Генерация закрытого ключа x
	x, _ := rand.Int(rand.Reader, p)

	// Вычисление открытого ключа y
	y := modExp(g, x, p)

	return p, g, y
}

func encrypt(p, g, y, message *big.Int) (*big.Int, *big.Int) {
	// Генерация случайного секретного числа k
	k, _ := rand.Int(rand.Reader, p)
	// Вычисление a = g^k mod p
	a := modExp(g, k, p)

	// Вычисление b = (y^k * message) mod p
	b := new(big.Int).Mod(new(big.Int).Mul(modExp(y, k, p), message), p)

	return a, b
}

func decrypt(p, x, a, b *big.Int) *big.Int {
	// Вычисление s = a^x mod p
	s := modExp(a, x, p)

	// Вычисление s^(-1) mod p
	sInverse := new(big.Int).ModInverse(s, p)

	// Вычисление расшифрованного сообщения message = (b * s^(-1)) mod p
	message := new(big.Int).Mod(new(big.Int).Mul(b, sInverse), p)

	return message
}

func main() {
	p, g, y := generateKeyPair()

	message := big.NewInt(10)

	a, b := encrypt(p, g, y, message)
	fmt.Println("Encrypted:", a, b)

	decryptedMessage := decrypt(p, y, a, b)
	fmt.Println("Decrypted:", decryptedMessage)
}
