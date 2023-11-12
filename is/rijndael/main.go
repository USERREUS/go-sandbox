package main

import (
	"crypto/rand"
	"io"
	"io/ioutil"
	"log"
	"os"

	rijndael "is/rijndael/internal" // Замените "your_package_path" на фактический путь к вашему пакету Rijndael
)

func main() {
	// Чтение медиафайла в буфер
	filePath := "input.jpg"
	mediaBuffer, err := readFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	// Инициализация шифра Rijndael
	key := generateRandomKey()
	c := rijndael.NewCipher(&key)

	// Шифрование данных
	encryptedBuffer := encryptData(c, mediaBuffer)

	// Запись зашифрованных данных в файл
	encryptedFilePath := "your_encrypted_media_file.jpg"
	err = writeFile(encryptedFilePath, encryptedBuffer)
	if err != nil {
		log.Fatal(err)
	}

	// Расшифрование данных
	decryptedBuffer := decryptData(c, encryptedBuffer)

	// Запись расшифрованных данных в файл
	decryptedFilePath := "your_decrypted_media_file.jpg"
	err = writeFile(decryptedFilePath, decryptedBuffer)
	if err != nil {
		log.Fatal(err)
	}
}

func readFile(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return ioutil.ReadAll(file)
}

func writeFile(filePath string, data []byte) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	return err
}

func encryptDataCBC(c *rijndael.Cipher, data []byte) []byte {
	blockSize := 32
	numBlocks := (len(data) + blockSize - 1) / blockSize

	encryptedBuffer := make([]byte, len(data))

	for i := 0; i < numBlocks; i++ {
		start := i * blockSize
		end := (i + 1) * blockSize
		if end > len(data) {
			end = len(data)
		}

		var block [32]byte
		for i := start; i < end; i++ {
			block[i-start] = data[i] ^ c.Prev[i-start]
		}

		var encryptedBlock [32]byte

		// Зашифрование блока данных
		c.Encrypt(&encryptedBlock, &block)

		// Запись зашифрованного блока в буфер
		copy(encryptedBuffer[start:end], encryptedBlock[:])

		copy(c.Prev[:], encryptedBlock[:])
	}

	return encryptedBuffer
}

func encryptData(c *rijndael.Cipher, data []byte) []byte {
	blockSize := 32
	numBlocks := (len(data) + blockSize - 1) / blockSize

	encryptedBuffer := make([]byte, len(data))

	for i := 0; i < numBlocks; i++ {
		start := i * blockSize
		end := (i + 1) * blockSize
		if end > len(data) {
			end = len(data)
		}

		var block [32]byte
		for i := start; i < end; i++ {
			block[i-start] = data[i]
		}
		var encryptedBlock [32]byte

		// Зашифрование блока данных
		c.Encrypt(&encryptedBlock, &block)

		// Запись зашифрованного блока в буфер
		copy(encryptedBuffer[start:end], encryptedBlock[:])
	}

	return encryptedBuffer
}

func decryptDataCBC(c *rijndael.Cipher, data []byte) []byte {
	blockSize := 32
	numBlocks := (len(data) + blockSize - 1) / blockSize

	decryptedBuffer := make([]byte, len(data))

	for i := 0; i < numBlocks; i++ {
		start := i * blockSize
		end := (i + 1) * blockSize
		if end > len(data) {
			end = len(data)
		}

		var temp [32]byte
		copy(temp[:], data[start:end])

		var block [32]byte
		copy(block[:], data[start:end])
		var decryptedBlock [32]byte
		// Расшифрование блока данных
		c.Decrypt(&decryptedBlock, &block)

		for i := start; i < end; i++ {
			decryptedBlock[i-start] ^= c.Prev[i-start]
		}

		// Запись расшифрованного блока в буфер
		copy(decryptedBuffer[start:end], decryptedBlock[:])

		copy(c.Prev[:], temp[:])
	}

	return decryptedBuffer
}

func decryptData(c *rijndael.Cipher, data []byte) []byte {
	blockSize := 32
	numBlocks := (len(data) + blockSize - 1) / blockSize

	decryptedBuffer := make([]byte, len(data))

	for i := 0; i < numBlocks; i++ {
		start := i * blockSize
		end := (i + 1) * blockSize
		if end > len(data) {
			end = len(data)
		}

		var block [32]byte
		for i := start; i < end; i++ {
			block[i-start] = data[i]
		}
		var decryptedBlock [32]byte

		// Расшифрование блока данных
		c.Decrypt(&decryptedBlock, &block)

		// Запись расшифрованного блока в буфер
		copy(decryptedBuffer[start:end], decryptedBlock[:])
	}

	return decryptedBuffer
}

func generateRandomKey() [32]byte {
	var key [32]byte
	if _, err := io.ReadFull(rand.Reader, key[:]); err != nil {
		panic(err)
	}
	return key
}
