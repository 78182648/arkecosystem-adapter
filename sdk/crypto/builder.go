// This file is part of Ark Go Crypto.
//
// (c) Ark Ecosystem <info@ark.io>
//
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.

package crypto

func buildSignedTransaction(transaction *Transaction, passphrase string, secondPassphrase string) *Transaction {
	transaction.Sign(passphrase)

	if len(secondPassphrase) > 0 {
		transaction.SecondSign(secondPassphrase)
	}

	transaction.Id = transaction.GetId()

	return transaction
}

func buildMultiSignedTransaction(transaction *Transaction, signerIndex int, passphrase string) *Transaction {
	transaction.SignMulti(signerIndex, passphrase)

	transaction.Id = transaction.GetId()

	return transaction
}

func BuildTransferMySelf(recipient string, amount FlexToshi, senderpk string, senderid string, nonce uint64) *Transaction {
	transaction := &Transaction{
		SenderPublicKey: senderpk,
		SenderId:        senderid,
		Nonce:           nonce + 1,
		Amount: amount,
		RecipientId: recipient,
	}

	setCommonFields(transaction, TRANSACTION_TYPES.Transfer)

	transaction.Asset = &TransactionAsset{}
	transaction.Timestamp = GetTime()

	return transaction
}

func setCommonFields(transaction *Transaction, transactionType uint16) {
	if transaction.Fee == 0 {
		transaction.Fee = GetFee(transactionType)
	}

	if transaction.Network == 0 {
		transaction.Network = GetNetwork().Version
	}

	transaction.SecondSenderPublicKey = ""
	transaction.SecondSignature = ""

	if transaction.Timestamp == 0 {
		transaction.Timestamp = GetTime()
	}

	transaction.Type = transactionType
	transaction.TypeGroup = TRANSACTION_TYPE_GROUPS.Core
	transaction.Version = 2
}

/** Set all fields and sign a TransactionTypes.Transfer transaction.
 * Members of the supplied transaction that must be set when calling this function:
 *   Amount
 *   Expiration - optional, could be 0 to designate no expiration
 *   Fee - optional, if 0, then it will be set to a default fee
 *   Network - optional, if 0, then it will be set to the configured network
 *   Nonce
 *   RecipientId
 *   Timestamp - optional, if 0, then it will be set to the present time
 *   VendorField - optional */
func BuildTransfer(transaction *Transaction, passphrase string, secondPassphrase string) *Transaction {
	setCommonFields(transaction, TRANSACTION_TYPES.Transfer)

	transaction.Asset = &TransactionAsset{}

	return buildSignedTransaction(transaction, passphrase, secondPassphrase)
}

/** Set all fields and sign a multi signature TransactionTypes.Transfer transaction.
 * Members of the supplied transaction that must be set when calling this function:
 *   Amount
 *   Expiration - optional, could be 0 to designate no expiration
 *   Fee - optional, if 0, then it will be set to a default fee
 *   Network - optional, if 0, then it will be set to the configured network
 *   Nonce
 *   RecipientId
 *   Signatures - must be an array (could be empty); a new signature will be appended to it
 *   Timestamp - optional, if 0, then it will be set to the present time
 *   VendorField - optional */
func BuildTransferMultiSignature(transaction *Transaction, signerIndex int, passphrase string) *Transaction {
	setCommonFields(transaction, TRANSACTION_TYPES.Transfer)

	transaction.Asset = &TransactionAsset{}

	return buildMultiSignedTransaction(transaction, signerIndex, passphrase)
}

/** Set all fields and sign a TransactionTypes.SecondSignatureRegistration transaction.
 * Members of the supplied transaction that must be set when calling this function:
 *   Expiration - optional, could be 0 to designate no expiration
 *   Fee - optional, if 0, then it will be set to a default fee
 *   Network - optional, if 0, then it will be set to the configured network
 *   Nonce
 *   Timestamp - optional, if 0, then it will be set to the present time
 *   VendorField - optional */
func BuildSecondSignatureRegistration(transaction *Transaction, passphrase string, secondPassphrase string) *Transaction {
	setCommonFields(transaction, TRANSACTION_TYPES.SecondSignatureRegistration)

	secondPublicKey, _ := PublicKeyFromPassphrase(secondPassphrase)

	transaction.Amount = 0
	transaction.Asset = &TransactionAsset{
		Signature: &SecondSignatureRegistrationAsset{
			PublicKey: HexEncode(secondPublicKey.Serialize()),
		},
	}

	return buildSignedTransaction(transaction, passphrase, secondPassphrase)
}

/** Set all fields and sign a TransactionTypes.DelegateRegistration transaction.
 * Members of the supplied transaction that must be set when calling this function:
 *   Asset.Delegate.Username
 *   Expiration - optional, could be 0 to designate no expiration
 *   Fee - optional, if 0, then it will be set to a default fee
 *   Network - optional, if 0, then it will be set to the configured network
 *   Nonce
 *   Timestamp - optional, if 0, then it will be set to the present time
 *   VendorField - optional */
func BuildDelegateRegistration(transaction *Transaction, passphrase string, secondPassphrase string) *Transaction {
	setCommonFields(transaction, TRANSACTION_TYPES.DelegateRegistration)

	return buildSignedTransaction(transaction, passphrase, secondPassphrase)
}

/** Set all fields and sign a TransactionTypes.Vote transaction.
 * Members of the supplied transaction that must be set when calling this function:
 *   Asset.Votes
 *   Expiration - optional, could be 0 to designate no expiration
 *   Fee - optional, if 0, then it will be set to a default fee
 *   Network - optional, if 0, then it will be set to the configured network
 *   Nonce
 *   Timestamp - optional, if 0, then it will be set to the present time
 *   VendorField - optional */
func BuildVote(transaction *Transaction, passphrase string, secondPassphrase string) *Transaction {
	setCommonFields(transaction, TRANSACTION_TYPES.Vote)

	transaction.RecipientId, _ = AddressFromPassphrase(passphrase)

	return buildSignedTransaction(transaction, passphrase, secondPassphrase)
}

/** Set all fields and sign a TransactionTypes.MultiSignatureRegistration transaction.
 * Members of the supplied transaction that must be set when calling this function:
 *   Asset.MultiSignature
 *   Expiration - optional, could be 0 to designate no expiration
 *   Fee - optional, if 0, then it will be set to a default fee
 *   Network - optional, if 0, then it will be set to the configured network
 *   Nonce
 *   Timestamp - optional, if 0, then it will be set to the present time
 *   VendorField - optional */
func BuildMultiSignatureRegistration(transaction *Transaction, passphrase string, secondPassphrase string) *Transaction {
	setCommonFields(transaction, TRANSACTION_TYPES.MultiSignatureRegistration)

	return buildSignedTransaction(transaction, passphrase, secondPassphrase)
}

/** Set all fields and sign a TransactionTypes.Ipfs transaction.
 * Members of the supplied transaction that must be set when calling this function:
 *   Asset.Ipfs
 *   Expiration - optional, could be 0 to designate no expiration
 *   Fee - optional, if 0, then it will be set to a default fee
 *   Network - optional, if 0, then it will be set to the configured network
 *   Nonce
 *   Timestamp - optional, if 0, then it will be set to the present time
 *   VendorField - optional */
func BuildIpfs(transaction *Transaction, passphrase string, secondPassphrase string) *Transaction {
	setCommonFields(transaction, TRANSACTION_TYPES.Ipfs)

	return buildSignedTransaction(transaction, passphrase, secondPassphrase)
}

/** Set all fields and sign a TransactionTypes.MultiPayment transaction.
 * Members of the supplied transaction that must be set when calling this function:
 *   Asset.Payments
 *   Expiration - optional, could be 0 to designate no expiration
 *   Fee - optional, if 0, then it will be set to a default fee
 *   Network - optional, if 0, then it will be set to the configured network
 *   Nonce
 *   Timestamp - optional, if 0, then it will be set to the present time
 *   VendorField - optional */
func BuildMultiPayment(transaction *Transaction, passphrase string, secondPassphrase string) *Transaction {
	setCommonFields(transaction, TRANSACTION_TYPES.MultiPayment)

	return buildSignedTransaction(transaction, passphrase, secondPassphrase)
}

/** Set all fields and sign a TransactionTypes.DelegateResignation transaction.
 * Members of the supplied transaction that must be set when calling this function:
 *   Expiration - optional, could be 0 to designate no expiration
 *   Fee - optional, if 0, then it will be set to a default fee
 *   Network - optional, if 0, then it will be set to the configured network
 *   Nonce
 *   Timestamp - optional, if 0, then it will be set to the present time
 *   VendorField - optional */
func BuildDelegateResignation(transaction *Transaction, passphrase string, secondPassphrase string) *Transaction {
	setCommonFields(transaction, TRANSACTION_TYPES.DelegateResignation)

	return buildSignedTransaction(transaction, passphrase, secondPassphrase)
}

/** Set all fields and sign a TransactionTypes.HtlcLock transaction.
 * Members of the supplied transaction that must be set when calling this function:
 *   Asset.Lock
 *   Expiration - optional, could be 0 to designate no expiration
 *   Fee - optional, if 0, then it will be set to a default fee
 *   Network - optional, if 0, then it will be set to the configured network
 *   Nonce
 *   Timestamp - optional, if 0, then it will be set to the present time
 *   VendorField - optional */
func BuildHtlcLock(transaction *Transaction, passphrase string, secondPassphrase string) *Transaction {
	setCommonFields(transaction, TRANSACTION_TYPES.HtlcLock)

	return buildSignedTransaction(transaction, passphrase, secondPassphrase)
}

/** Set all fields and sign a TransactionTypes.HtlcClaim transaction.
 * Members of the supplied transaction that must be set when calling this function:
 *   Asset.Claim
 *   Expiration - optional, could be 0 to designate no expiration
 *   Fee - optional, if 0, then it will be set to a default fee
 *   Network - optional, if 0, then it will be set to the configured network
 *   Nonce
 *   Timestamp - optional, if 0, then it will be set to the present time
 *   VendorField - optional */
func BuildHtlcClaim(transaction *Transaction, passphrase string, secondPassphrase string) *Transaction {
	setCommonFields(transaction, TRANSACTION_TYPES.HtlcClaim)

	return buildSignedTransaction(transaction, passphrase, secondPassphrase)
}

/** Set all fields and sign a TransactionTypes.HtlcRefund transaction.
 * Members of the supplied transaction that must be set when calling this function:
 *   Asset.Refund
 *   Expiration - optional, could be 0 to designate no expiration
 *   Fee - optional, if 0, then it will be set to a default fee
 *   Network - optional, if 0, then it will be set to the configured network
 *   Nonce
 *   Timestamp - optional, if 0, then it will be set to the present time
 *   VendorField - optional */
func BuildHtlcRefund(transaction *Transaction, passphrase string, secondPassphrase string) *Transaction {
	setCommonFields(transaction, TRANSACTION_TYPES.HtlcRefund)

	return buildSignedTransaction(transaction, passphrase, secondPassphrase)
}
