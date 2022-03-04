package main

// tar.xz the entire working directory
// and then delete the working directory
// and then encrypt the tarball with a password

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"

	"github.com/mholt/archiver"
	"golang.org/x/crypto/scrypt"
	tbget "i2pgit.org/idk/i2p.plugins.tor-manager/get"
)

func DecryptTarXZifThere(source, password string) error {
	if tbget.FileExists(source + ".tar.xz.crypt") {
		return DecryptTarXzip(source, password)
	}
	log.Println("DecryptTarXZifThere: no encrypted store at:", source+".tar.xz.crypt")
	return nil
}

func TarXzip(source string) error {
	target := fmt.Sprintf("%s.tar.xz", source)
	txz := archiver.NewTarXz()
	//	txz.CompressionLevel = 12
	err := txz.Archive([]string{source}, target)
	if err != nil {
		return fmt.Errorf("TarGzip: TarGz() failed: %s", err.Error())
	}
	return nil
}

func EncryptTarXZip(source, password string) error {
	log.Println("EncryptTarXZip:", source)
	err := TarXzip(source)
	if err != nil {
		return err
	}
	sourceBytes, err := ioutil.ReadFile(source + ".tar.xz")
	if err != nil {
		return err
	}
	encryptedSourceBytes, err := Encrypt([]byte(password), sourceBytes)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(source+".tar.xz.crypt", encryptedSourceBytes, 0644)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(source+".tar.xz", os.O_RDWR, 0)
	if err != nil {
		log.Println("Error:", err)
	}
	info, err := file.Stat()
	if err != nil {
		log.Println("Error:", err)
	}
	bytes := info.Size()
	file.Truncate(0)
	file.Write(make([]byte, bytes))
	file.Close()
	os.Remove(source + ".tar.xz")
	OverwriteDirectoryContents(source)
	os.RemoveAll(source)
	return nil
}

func UnTarXzip(source string) error {
	target := strings.Replace(source, ".tar.xz", "", 1)
	txz := archiver.NewTarXz()
	txz.Tar.OverwriteExisting = true
	txz.Tar.ContinueOnError = true
	err := txz.Unarchive(source, target)
	if err != nil {
		return fmt.Errorf("TarGzip: Unarchive() failed: %s", err.Error())
	}
	return nil
}

func DecryptTarXzip(source, password string) error {
	sourceBytes, err := ioutil.ReadFile(source + ".tar.xz.crypt")
	if err != nil {
		return err
	}
	decryptedSourceBytes, err := Decrypt([]byte(password), sourceBytes)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(source+".tar.xz", decryptedSourceBytes, 0644)
	if err != nil {
		return err
	}
	err = UnTarXzip(source + ".tar.xz")
	if err != nil {
		return err
	}
	os.Remove(source + ".tar.xz.crypt")
	return nil
}

// Borrowing **VERY** heavily from: https://bruinsslot.jp/post/golang-crypto/

func Encrypt(key, data []byte) ([]byte, error) {
	key, salt, err := DeriveKey(key, nil)
	if err != nil {
		return nil, err
	}

	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)

	ciphertext = append(ciphertext, salt...)

	return ciphertext, nil
}

func Decrypt(key, data []byte) ([]byte, error) {
	salt, data := data[len(data)-32:], data[:len(data)-32]

	key, _, err := DeriveKey(key, salt)
	if err != nil {
		return nil, err
	}

	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, err
	}

	nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func DeriveKey(password, salt []byte) ([]byte, []byte, error) {
	if salt == nil {
		salt = make([]byte, 32)
		if _, err := rand.Read(salt); err != nil {
			return nil, nil, err
		}
	}

	key, err := scrypt.Key(password, salt, 1048576, 8, 1, 32)
	if err != nil {
		return nil, nil, err
	}

	return key, salt, nil
}
