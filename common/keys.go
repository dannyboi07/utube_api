package common

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
)

var PrivateKey *rsa.PrivateKey
var PublicKey *rsa.PublicKey

func InitKeys() error {
	var err error
	PrivateKey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println("Failed to generate private key, err:", err)
		return err
	}

	PublicKey = &PrivateKey.PublicKey

	return nil
}

// Malformed token
// func InitKeys() error {
// 	var (
// 		privKeyBytes []byte
// 		pubKeyBytes  []byte
// 		err          error
// 	)
// 	privKeyBytes, err = ioutil.ReadFile("private.pem")
// 	if err != nil {
// 		if !os.IsNotExist(err) {
// 			fmt.Println("Failed to read private.pem file, err:", err)
// 			return err
// 		}

// 		fmt.Println("Creating private & public keys...")
// 		var privateKey *rsa.PrivateKey
// 		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
// 		if err != nil {
// 			fmt.Println("Failed to generate private key, err:", err)
// 			return err
// 		}

// 		var publicKey rsa.PublicKey = privateKey.PublicKey
// 		pubKeyBytes = x509.MarshalPKCS1PublicKey(&publicKey)

// 		privKeyBytes = x509.MarshalPKCS1PrivateKey(privateKey)

// 		var privateKeyBlock *pem.Block = &pem.Block{
// 			Type:  "RSA PRIVATE KEY",
// 			Bytes: privKeyBytes,
// 		}
// 		privKeyFile, err := os.Create("private.pem")
// 		if err != nil {
// 			fmt.Println("Failed to create private pem file, err:", err)
// 			return err
// 		}

// 		if err = pem.Encode(privKeyFile, privateKeyBlock); err != nil {
// 			fmt.Println("Failed to encode private key block to pem file, err:", err)
// 			return err
// 		}
// 		fmt.Println("Private key written to file...")

// 		var publicKeyBlock *pem.Block = &pem.Block{
// 			Type:  "PUBLIC KEY",
// 			Bytes: pubKeyBytes,
// 		}
// 		pubKeyFile, err := os.Create("public.pem")
// 		if err != nil {
// 			fmt.Println("Failed to create public pem file, err:", err)
// 			return err
// 		}

// 		if err = pem.Encode(pubKeyFile, publicKeyBlock); err != nil {
// 			fmt.Println("Failed to encode public key block to pem file, err:", err)
// 			return err
// 		}
// 		fmt.Println("Public key written to file...")
// 	}

// 	if privKeyPem, _ := pem.Decode(privKeyBytes); privKeyPem != nil {
// 		PrivateKey, err = x509.ParsePKCS1PrivateKey(privKeyPem.Bytes)
// 		if err != nil {
// 			fmt.Println("Failed to parse private key, err:", err)
// 			return err
// 		}
// 	} else {
// 		fmt.Println("Failed to decode private key pem")
// 		return err
// 	}

// 	PublicKey = &PrivateKey.PublicKey

// 	return nil
// }

// key is not the correct one
// func InitKeys() error {
// 	var (
// 		privKeyBytes []byte
// 		pubKeyBytes  []byte
// 		err          error
// 	)
// 	// PrivateKey, err = rsa.GenerateKey(rand.Reader, 2048)
// 	// if err != nil {
// 	// 	fmt.Println("Failed to generate private key, err:", err)
// 	// 	return err
// 	// }

// 	// PublicKey = PrivateKey.PublicKey

// 	privKeyBytes, err = ioutil.ReadFile("private.pem")
// 	if err != nil {
// 		if !os.IsNotExist(err) {
// 			fmt.Println("Failed to read private key file, err:", err)
// 			return err
// 		}

// 		fmt.Println("Creating private & public keys...")

// 		PrivateKey, err = rsa.GenerateKey(rand.Reader, 2048)
// 		if err != nil {
// 			fmt.Println("Failed to generate private key, err:", err)
// 			return err
// 		}

// 		PublicKey = PrivateKey.PublicKey
// 		privKeyBytes = x509.MarshalPKCS1PrivateKey(PrivateKey)
// 		var privateKeyBlock *pem.Block = &pem.Block{
// 			Type:  "RSA PRIVATE KEY",
// 			Bytes: privKeyBytes,
// 		}
// 		// fmt.Println("privKeyBytes", privKeyBytes[:10])

// 		var privateKeyFile *os.File
// 		privateKeyFile, err = os.Create("private.pem")
// 		if err != nil {
// 			fmt.Println("Failed to create private.pem file, err:", err)
// 			return err
// 		}

// 		if err = pem.Encode(privateKeyFile, privateKeyBlock); err != nil {
// 			fmt.Println("Failed to encode private PEM block to PEM file, err:", err)
// 			return err
// 		}
// 		privateKeyFile.Close()
// 		fmt.Println("Private key written to file")

// 		pubKeyBytes = x509.MarshalPKCS1PublicKey(&PublicKey)
// 		var publicKeyBlock *pem.Block = &pem.Block{
// 			Type:  "PUBLIC KEY",
// 			Bytes: pubKeyBytes,
// 		}

// 		var publicKeyFile *os.File
// 		publicKeyFile, err = os.Create("public.pem")
// 		if err != nil {
// 			fmt.Println("Failed to create public.pem file, err:", err)
// 			return err
// 		}

// 		if err = pem.Encode(publicKeyFile, publicKeyBlock); err != nil {
// 			fmt.Println("Failed to encode public PEM block to PEM file, err:", err)
// 			return err
// 		}
// 		publicKeyFile.Close()
// 		fmt.Println("Public key written to file")
// 	}

// 	if privateKeyPem, _ := pem.Decode(privKeyBytes); privateKeyPem != nil {
// 		PrivateKey, err = x509.ParsePKCS1PrivateKey(privateKeyPem.Bytes)
// 		if err != nil {
// 			fmt.Println("Failed to parse private key, err:", err)
// 			return err
// 		}
// 	} else {
// 		fmt.Println("Failed to decode private key PEM, empty bytes")
// 		return errors.New("Failed to decode private key PEM, empty bytes")
// 	}

// 	PublicKey = PrivateKey.PublicKey

// 	return nil
// }
