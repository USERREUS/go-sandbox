package main

import (
	"fmt"
	"io/ioutil"
)

// Функция ksa выполняет Key Scheduling Algorithm для инициализации массива состояния.
// Ключевой алгоритм позволяет переставлять элементы массива в соответствии с ключом.
func ksa(key []byte) []byte {
	s := make([]byte, 256)
	for i := 0; i < 256; i++ {
		s[i] = byte(i)
	}

	j := 0
	for i := 0; i < 256; i++ {
		j = (j + int(s[i]) + int(key[i%len(key)])) % 256
		s[i], s[j] = s[j], s[i]
	}

	return s
}

// Функция prga выполняет Pseudo-Random Generation Algorithm для генерации псевдослучайной последовательности.
// Алгоритм использует массив состояния, который был инициализирован в ksa.
func prga(s []byte, length int) []byte {
	i, j := 0, 0
	keyStream := make([]byte, length)

	for k := 0; k < length; k++ {
		i = (i + 1) % 256
		j = (j + int(s[i])) % 256
		s[i], s[j] = s[j], s[i]
		keyStream[k] = s[(int(s[i])+int(s[j]))%256]
	}

	return keyStream
}

// Функция rc4 выполняет шифрование или расшифрование данных с использованием алгоритма RC4.
func rc4(key, data []byte) []byte {
	// Инициализация массива состояния
	s := ksa(key)
	// Генерация псевдослучайной последовательности
	keyStream := prga(s, len(data))

	// Применение XOR между данными и ключевой последовательностью
	result := make([]byte, len(data))
	for i := 0; i < len(data); i++ {
		result[i] = data[i] ^ keyStream[i]
	}

	return result
}

// Функция encryptFile выполняет шифрование содержимого файла и сохраняет результат в другом файле.
func encryptFile(inputPath, outputPath string, key []byte) error {
	// Чтение содержимого входного файла
	plaintext, err := ioutil.ReadFile(inputPath)
	if err != nil {
		return err
	}

	// Шифрование содержимого файла
	ciphertext := rc4(key, plaintext)

	// Запись зашифрованного содержимого в выходной файл
	err = ioutil.WriteFile(outputPath, ciphertext, 0644)
	if err != nil {
		return err
	}

	fmt.Println("Шифрование завершено.")
	return nil
}

// Функция decryptFile выполняет расшифрование содержимого файла и сохраняет результат в другом файле.
func decryptFile(inputPath, outputPath string, key []byte) error {
	// Чтение зашифрованного содержимого входного файла
	ciphertext, err := ioutil.ReadFile(inputPath)
	if err != nil {
		return err
	}

	// Расшифрование содержимого файла
	plaintext := rc4(key, ciphertext)

	// Запись расшифрованного содержимого в выходной файл
	err = ioutil.WriteFile(outputPath, plaintext, 0644)
	if err != nil {
		return err
	}

	fmt.Println("Расшифрование завершено.")
	return nil
}

func main() {
	// Задание ключа и путей к файлам
	key := []byte("secretkey")
	inputPath := "input.jpg"
	encryptedPath := "encrypted.jpg"
	decryptedPath := "decrypted.jpg"

	// Шифрование файла
	err := encryptFile(inputPath, encryptedPath, key)
	if err != nil {
		fmt.Println("Ошибка шифрования файла:", err)
		return
	}

	// Расшифрование файла
	err = decryptFile(encryptedPath, decryptedPath, key)
	if err != nil {
		fmt.Println("Ошибка расшифрования файла:", err)
		return
	}
}
