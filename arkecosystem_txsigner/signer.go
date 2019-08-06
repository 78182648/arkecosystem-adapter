package arkecosystem_txsigner

import (
	"github.com/blocktree/arkecosystem-adapter/sdk/crypto"
)

var Default = &TransactionSigner{}

type TransactionSigner struct {
}

// SignTransactionHash 交易哈希签名算法
// required
func (singer *TransactionSigner) SignTransactionHash(msg []byte, privateKey []byte, eccType uint32) ([]byte, error) {
	//_, err := owcrypt.Signature(privateKey, nil, 0, msg, uint16(len(msg)), eccType)
	//if err != owcrypt.SUCCESS {
	//	return nil, fmt.Errorf("ECC sign hash failed")
	//}
	pk := crypto.PrivateKeyFromBytes(privateKey)
	signature2, _ := pk.Sign(msg)
	return signature2, nil
}



