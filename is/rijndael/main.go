package main

import (
	"crypto/rand"
	"io"
	"io/ioutil"
	"log"
	"os"

	rijndael "is/rijndael/internal"
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
	//encryptedBuffer := encryptDataECB(c, mediaBuffer)
	iv, _ := generateRandomIV()
	//encryptedBuffer := encryptDataCBC(c, mediaBuffer, iv)
	encryptedBuffer := encryptDataCFB(c, mediaBuffer, iv)
	//encryptedBuffer := encryptDataOFB(c, mediaBuffer, iv)

	// Запись зашифрованных данных в файл
	encryptedFilePath := "encrypted.jpg"
	err = writeFile(encryptedFilePath, encryptedBuffer)
	if err != nil {
		log.Fatal(err)
	}

	// Расшифрование данных
	//decryptedBuffer := decryptDataECB(c, encryptedBuffer)
	//decryptedBuffer := decryptDataCBC(c, encryptedBuffer, iv)
	decryptedBuffer := decryptDataCFB(c, encryptedBuffer, iv)
	//decryptedBuffer := decryptDataOFB(c, encryptedBuffer, iv)

	// Запись расшифрованных данных в файл
	decryptedFilePath := "decrypted.jpg"
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

func generateRandomIV() ([32]byte, error) {
	var iv [32]byte
	_, err := rand.Read(iv[:])
	if err != nil {
		return [32]byte{}, err
	}
	return iv, nil
}

// func encryptDataCBC(c *rijndael.Cipher, data []byte) []byte {
// 	blockSize := 32
// 	numBlocks := (len(data) + blockSize - 1) / blockSize

// 	encryptedBuffer := make([]byte, len(data))

// 	for i := 0; i < numBlocks; i++ {
// 		start := i * blockSize
// 		end := (i + 1) * blockSize
// 		if end > len(data) {
// 			end = len(data)
// 		}

// 		var block [32]byte
// 		for i := start; i < end; i++ {
// 			block[i-start] = data[i] ^ c.Prev[i-start]
// 		}

// 		var encryptedBlock [32]byte

// 		// Зашифрование блока данных
// 		c.Encrypt(&encryptedBlock, &block)

// 		// Запись зашифрованного блока в буфер
// 		copy(encryptedBuffer[start:end], encryptedBlock[:])

// 		copy(c.Prev[:], encryptedBlock[:])
// 	}

// 	return encryptedBuffer
// }

func encryptDataCBC(c *rijndael.Cipher, data []byte, iv [32]byte) []byte {
	blockSize := 32
	numBlocks := (len(data) + blockSize - 1) / blockSize

	encryptedBuffer := make([]byte, len(data))
	prevBlock := iv

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

		// Применение XOR с предыдущим зашифрованным блоком
		for j := 0; j < blockSize; j++ {
			block[j] ^= prevBlock[j]
		}

		var encryptedBlock [32]byte

		// Зашифрование блока данных
		c.Encrypt(&encryptedBlock, &block)

		// Запись зашифрованного блока в буфер и обновление предыдущего блока
		copy(encryptedBuffer[start:end], encryptedBlock[:])
		prevBlock = encryptedBlock
	}

	return encryptedBuffer
}

func encryptDataECB(c *rijndael.Cipher, data []byte) []byte {
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

func encryptDataCFB(c *rijndael.Cipher, data []byte, iv [32]byte) []byte {
	blockSize := 32
	numBlocks := (len(data) + blockSize - 1) / blockSize

	encryptedBuffer := make([]byte, len(data))
	previousCipherText := iv

	for i := 0; i < numBlocks; i++ {
		start := i * blockSize
		end := (i + 1) * blockSize
		if end > len(data) {
			end = len(data)
		}

		// Шифрование предыдущего шифротекста
		c.Encrypt(&previousCipherText, &previousCipherText)

		// XOR шифротекста с блоком данных
		for j := start; j < end; j++ {
			encryptedBuffer[j] = data[j] ^ previousCipherText[j-start]
		}
	}

	return encryptedBuffer
}

func encryptDataOFB(c *rijndael.Cipher, data []byte, iv [32]byte) []byte {
	blockSize := 32
	numBlocks := (len(data) + blockSize - 1) / blockSize

	encryptedBuffer := make([]byte, len(data))
	feedback := iv

	for i := 0; i < numBlocks; i++ {
		start := i * blockSize
		end := (i + 1) * blockSize
		if end > len(data) {
			end = len(data)
		}

		// Шифрование текущего feedback'а
		c.Encrypt(&feedback, &feedback)

		// XOR шифротекста с блоком данных
		for j := start; j < end; j++ {
			encryptedBuffer[j] = data[j] ^ feedback[j-start]
		}
	}

	return encryptedBuffer
}

// func decryptDataCBC(c *rijndael.Cipher, data []byte) []byte {
// 	blockSize := 32
// 	numBlocks := (len(data) + blockSize - 1) / blockSize

// 	decryptedBuffer := make([]byte, len(data))

// 	for i := 0; i < numBlocks; i++ {
// 		start := i * blockSize
// 		end := (i + 1) * blockSize
// 		if end > len(data) {
// 			end = len(data)
// 		}

// 		var temp [32]byte
// 		copy(temp[:], data[start:end])

// 		var block [32]byte
// 		copy(block[:], data[start:end])
// 		var decryptedBlock [32]byte
// 		// Расшифрование блока данных
// 		c.Decrypt(&decryptedBlock, &block)

// 		for i := start; i < end; i++ {
// 			decryptedBlock[i-start] ^= c.Prev[i-start]
// 		}

// 		// Запись расшифрованного блока в буфер
// 		copy(decryptedBuffer[start:end], decryptedBlock[:])

// 		copy(c.Prev[:], temp[:])
// 	}

// 	return decryptedBuffer
// }

func decryptDataCBC(c *rijndael.Cipher, data []byte, iv [32]byte) []byte {
	blockSize := 32
	numBlocks := (len(data) + blockSize - 1) / blockSize

	decryptedBuffer := make([]byte, len(data))
	prevBlock := iv

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

		// Применение XOR с предыдущим зашифрованным блоком для восстановления исходных данных
		for j := 0; j < blockSize; j++ {
			decryptedBlock[j] ^= prevBlock[j]
		}

		// Запись расшифрованного блока в буфер и обновление предыдущего блока
		copy(decryptedBuffer[start:end], decryptedBlock[:])
		prevBlock = block
	}

	return decryptedBuffer
}

func decryptDataECB(c *rijndael.Cipher, data []byte) []byte {
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

func decryptDataCFB(c *rijndael.Cipher, data []byte, iv [32]byte) []byte {
	blockSize := 32
	numBlocks := (len(data) + blockSize - 1) / blockSize

	decryptedBuffer := make([]byte, len(data))
	previousCipherText := iv

	for i := 0; i < numBlocks; i++ {
		start := i * blockSize
		end := (i + 1) * blockSize
		if end > len(data) {
			end = len(data)
		}

		// Шифрование предыдущего шифротекста
		c.Encrypt(&previousCipherText, &previousCipherText)

		// XOR шифротекста с блоком данных
		for j := start; j < end; j++ {
			decryptedBuffer[j] = data[j] ^ previousCipherText[j-start]
		}
	}

	return decryptedBuffer
}

func decryptDataOFB(c *rijndael.Cipher, data []byte, iv [32]byte) []byte {
	blockSize := 32
	numBlocks := (len(data) + blockSize - 1) / blockSize

	decryptedBuffer := make([]byte, len(data))
	feedback := iv

	for i := 0; i < numBlocks; i++ {
		start := i * blockSize
		end := (i + 1) * blockSize
		if end > len(data) {
			end = len(data)
		}

		// Шифрование текущего feedback'а
		c.Encrypt(&feedback, &feedback)

		// XOR шифротекста с блоком данных
		for j := start; j < end; j++ {
			decryptedBuffer[j] = data[j] ^ feedback[j-start]
		}
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
