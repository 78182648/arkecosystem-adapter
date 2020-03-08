package arkecosystem_addrec

import (
	"github.com/blocktree/arkecosystem-adapter/sdk/crypto2"
	"github.com/blocktree/openwallet/log"
)

var (
	Default = AddressDecoderV2{}
)

//AddressDecoderV2
type AddressDecoderV2 struct {
	IsTestNet bool
}

// GetAddressFromPublicKey takes a Lisk public key and returns the associated address
func GetAddressFromPublicKey(publicKey []byte) string {
	//publicKeyHash := sha256.Sum256(publicKey)

	pk, err := crypto2.PublicKeyFromBytes(publicKey)
	if err != nil {
		log.Error(err)
	}
	return pk.ToAddress()

}

//AddressEncode 地址编码
func (dec *AddressDecoderV2) AddressEncode(hash []byte, opts ...interface{}) (string, error) {
	address := GetAddressFromPublicKey(hash)
	return address, nil
}
