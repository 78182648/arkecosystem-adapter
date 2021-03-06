// This file is part of Ark Go Crypto.
//
// (c) Ark Ecosystem <info@ark.io>
//
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.

package crypto

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"

	"github.com/fatih/structs"
)

func (transaction *Transaction) GetId() string {
	return fmt.Sprintf("%x", sha256.Sum256(transaction.Serialize(true, true, true)))
}

func (transaction *Transaction) Sign(passphrase string) {
	privateKey, _ := PrivateKeyFromPassphrase(passphrase)

	transaction.SenderPublicKey = HexEncode(privateKey.PublicKey.Serialize())

	hash := sha256.Sum256(transaction.Serialize(false, false, false))

	signature, err := privateKey.Sign(hash[:])
	if err == nil {
		transaction.Signature = HexEncode(signature)
	}
}

func (transaction *Transaction) SignMulti(signerIndex int, passphrase string) {
	privateKey, _ := PrivateKeyFromPassphrase(passphrase)

	hash := sha256.Sum256(transaction.Serialize(false, false, false))

	signature, err := privateKey.Sign(hash[:])
	if err == nil {
		var signatureWithIndex []byte
		signatureWithIndex = append(signatureWithIndex, byte(signerIndex))
		signatureWithIndex = append(signatureWithIndex, signature...)
		transaction.Signatures = append(transaction.Signatures, HexEncode(signatureWithIndex))
	}
}

func (transaction *Transaction) SecondSign(passphrase string) {
	privateKey, _ := PrivateKeyFromPassphrase(passphrase)

	hash := sha256.Sum256(transaction.Serialize(true, false, false))

	signature, err := privateKey.Sign(hash[:])
	if err == nil {
		transaction.SecondSignature = HexEncode(signature)
	}
}

func (transaction *Transaction) VerifyMultiSignature(multiSignatureAsset *MultiSignatureRegistrationAsset) (bool, error) {
	hash := sha256.Sum256(transaction.Serialize(false, false, false))

	publicKeyIndexes := make(map[int]bool)
	numVerified := 0

	for i := 0; i < len(transaction.Signatures); i++ {
		publicKeyIndex := int(HexDecode(transaction.Signatures[i][:2])[0])
		signature := HexDecode(transaction.Signatures[i][2:])

		if publicKeyIndexes[publicKeyIndex] {
			return false, fmt.Errorf("VerifyMultiSignature: duplicate signer index: %d", publicKeyIndex)
		}

		if publicKeyIndex >= len(multiSignatureAsset.PublicKeys) {
			return false, fmt.Errorf(
				"VerifyMultiSignature: signer index too large: %d, total of %d " +
				"signers have been registered",
				publicKeyIndex, len(multiSignatureAsset.PublicKeys))
		}

		publicKeyIndexes[publicKeyIndex] = true

		publicKey, err := PublicKeyFromBytes(HexDecode(multiSignatureAsset.PublicKeys[publicKeyIndex]))
		if err != nil {
			return false, err
		}

		verified, err := publicKey.Verify(signature, hash[:])

		if verified && err == nil {
			numVerified++
		}

		if numVerified >= int(multiSignatureAsset.Min) {
			return true, nil
		}

		if len(transaction.Signatures) - (i + 1 - numVerified) < int(multiSignatureAsset.Min) {
			return false, fmt.Errorf(
				"VerifyMultiSignature: less than the minimum %d signatures verified successfully",
				multiSignatureAsset.Min)
		}
	}

	return false, fmt.Errorf(
		"VerifyMultiSignature: less than the minimum %d signatures verified successfully (checked all)",
		multiSignatureAsset.Min)
}

func (transaction *Transaction) Verify(multiSignatureAsset ...*MultiSignatureRegistrationAsset) (bool, error) {
	if len(multiSignatureAsset) == 1 && multiSignatureAsset[0].Min > 0 {
		return transaction.VerifyMultiSignature(multiSignatureAsset[0])
	}

	publicKey, err := PublicKeyFromBytes(HexDecode(transaction.SenderPublicKey))

	if err != nil {
		return false, err
	}

	hash := sha256.Sum256(transaction.Serialize(false, false, true))

	return publicKey.Verify(HexDecode(transaction.Signature), hash[:])
}

func (transaction *Transaction) SecondVerify(secondPublicKey *PublicKey) (bool, error) {
	hash := sha256.Sum256(transaction.Serialize(true, false, false))

	return secondPublicKey.Verify(HexDecode(transaction.SecondSignature), hash[:])
}

func isSchnorrSignature(length int) bool {
	// Logic copied from
	// https://github.com/ArkEcosystem/core/blob/0663b0f/packages/crypto/src/transactions/deserializer.ts#L173
	// length is in number of bytes (raw / binary)
	return (
	    length == 64 || // signature
	    length == 128 || // signature + secondSignature
		length % 65 == 0 || // `signatures` of a multi signature transaction, type != MultiSignatureRegistration (4)
		length % 65 == 64 || // type == MultiSignatureRegistration (4)
		length % 65 == 63) // type == MultiSignatureRegistration (4) + secondSignature
}

func ECDSASignatureLen(signature []byte) int {
	return int(signature[1] + 2)
}

func beginningMultiSignature(signature []byte) bool {
	return signature[0] == 0xFF
}

func (transaction *Transaction) ParseSignaturesECDSA(signatures []byte) *Transaction {
	signaturesLen := len(signatures)

	firstSignatureLen := ECDSASignatureLen(signatures)

	transaction.Signature = HexEncode(signatures[:firstSignatureLen])

	o := firstSignatureLen

	if o == signaturesLen {
		return transaction
	}

	if !beginningMultiSignature(signatures[o:]) {
		secondSignatureLen := ECDSASignatureLen(signatures[o:])
		transaction.SecondSignature = HexEncode(signatures[o:o + secondSignatureLen])
		o += secondSignatureLen
	}

	if o == signaturesLen {
		return transaction
	}

	if o != signaturesLen {
		log.Fatal("All signatures parsed, but ", signaturesLen - o,
			" bytes remain in the buffer: ", HexEncode(signatures))
	}

	return transaction
}

func (transaction *Transaction) ParseSignaturesSchnorr(signatures []byte) *Transaction {
	const schnorrSignatureLen = 64

	signaturesLen := len(signatures)
	o := 0

	canReadNonMultiSignature := func () bool {
		remaining := signaturesLen - o
		return remaining >= schnorrSignatureLen && remaining % 65 != 0
	}

	readSchnorrSignature := func () string {
		sig := HexEncode(signatures[o:o + schnorrSignatureLen])
		o += schnorrSignatureLen
		return sig
	}

	if canReadNonMultiSignature() {
		transaction.Signature = readSchnorrSignature()
	}

	if canReadNonMultiSignature() {
		transaction.SecondSignature = readSchnorrSignature()
	}

	if signaturesLen - o == 0 {
		return transaction
	}

	if (signaturesLen - o) % 65 != 0 {
		log.Fatalf("Cannot parse Schnorr signatures: remaining bytes not multiple of 65: %d", signaturesLen - o)
	}

	count := (signaturesLen - o) / 65

	for i := 0; i < count; i++ {
		signaturePlusPrefix := HexEncode(signatures[o:o + 1 + schnorrSignatureLen])
		o += 1 + schnorrSignatureLen

		transaction.Signatures = append(transaction.Signatures, signaturePlusPrefix)
	}

	return transaction
}

func (transaction *Transaction) ParseSignatures(sigOffset int) *Transaction {
	signatures := transaction.Serialized[sigOffset:]
	signaturesLen := len(signatures)

	if signaturesLen == 0 {
		transaction.Signature = ""
		return transaction
	}

	if isSchnorrSignature(signaturesLen) {
		return transaction.ParseSignaturesSchnorr(signatures)
	}

	return transaction.ParseSignaturesECDSA(signatures)
}

func (transaction *Transaction) ToMap() map[string]interface{} {
	return structs.Map(transaction)
}

func (transaction *Transaction) ToJson() (string, error) {
	jason, err := json.Marshal(transaction)

	if err != nil {
		return "", err
	}

	return string(jason), nil
}
