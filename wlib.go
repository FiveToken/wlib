package wlib

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/filecoin-project/go-address"
	crypto "github.com/filecoin-project/go-crypto"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/minio/blake2b-simd"
)

type Msg struct {
	Version    uint64 `json:"version"`
	To         string `json:"to"`
	From       string `json:"from"`
	Nonce      uint64 `json:"nonce"`
	Value      string `json:"value"`
	GasLimit   int64  `json:"gaslimit"`
	GasFeeCap  string `json:"gasfeecap"`
	GasPremium string `json:"gaspremium"`
	Method     uint64 `json:"method"`
	Params     string `json:"params"`
}

func GenAddress(pk, t string) string {
	if pk == "" {
		return ""
	}
	pkbytes, err := base64.StdEncoding.DecodeString(pk)

	if err != nil {
		return ""
	}
	var addr address.Address

	if t == "secp" {
		addr, err = address.NewSecp256k1Address(pkbytes)
	} else {
		addr, err = address.NewFromBytes(append([]byte{3}, pkbytes...))
	}

	if err != nil {
		return ""
	}
	return addr.String()
}

func AddressFromString(str string) string {
	addr, err := address.NewFromString(str)
	if err != nil {
		return ""
	}

	return addr.String()
}

func MessageCid(jsonstr string) (r string) {
	defer func() {
		if re := recover(); re != nil {
			r = fmt.Sprintf("%v", r)
		}
	}()
	var msg Msg
	err := json.Unmarshal([]byte(jsonstr), &msg)
	checkError(err)
	toAddr, err := address.NewFromString(msg.To)
	checkError(err)
	fromAddr, err := address.NewFromString(msg.From)
	checkError(err)
	v, err := big.FromString(msg.Value)
	checkError(err)
	gasfeecap, err := big.FromString(msg.GasFeeCap)
	checkError(err)
	gaspremium, err := big.FromString(msg.GasPremium)
	checkError(err)
	pbytes, err := base64.StdEncoding.DecodeString(msg.Params)
	checkError(err)

	tmsg := Message{
		Version:    msg.Version,
		To:         toAddr,
		From:       fromAddr,
		Nonce:      msg.Nonce,
		Value:      v,
		GasLimit:   msg.GasLimit,
		GasFeeCap:  abi.TokenAmount(gasfeecap),
		GasPremium: abi.TokenAmount(gaspremium),
		Method:     abi.MethodNum(msg.Method),
		Params:     pbytes,
	}

	return base64.StdEncoding.EncodeToString(tmsg.Cid().Bytes())
}
func GenCid(jsonstr string) (r string) {
	defer func() {
		if re := recover(); re != nil {
			r = fmt.Sprintf("%v", r)
		}
	}()
	var msg Msg
	err := json.Unmarshal([]byte(jsonstr), &msg)
	checkError(err)
	toAddr, err := address.NewFromString(msg.To)
	checkError(err)
	fromAddr, err := address.NewFromString(msg.From)
	checkError(err)
	v, err := big.FromString(msg.Value)
	checkError(err)
	gasfeecap, err := big.FromString(msg.GasFeeCap)
	checkError(err)
	gaspremium, err := big.FromString(msg.GasPremium)
	checkError(err)
	pbytes, err := base64.StdEncoding.DecodeString(msg.Params)
	checkError(err)

	tmsg := Message{
		Version:    msg.Version,
		To:         toAddr,
		From:       fromAddr,
		Nonce:      msg.Nonce,
		Value:      v,
		GasLimit:   msg.GasLimit,
		GasFeeCap:  abi.TokenAmount(gasfeecap),
		GasPremium: abi.TokenAmount(gaspremium),
		Method:     abi.MethodNum(msg.Method),
		Params:     pbytes,
	}

	return tmsg.Cid().String()
}

func SecpPrivateToPublic(ck string) string {
	if ck == "" {
		return ""
	}
	ckbytes, err := base64.StdEncoding.DecodeString(ck)

	if err != nil {
		return ""
	}

	pk := crypto.PublicKey(ckbytes)

	return base64.StdEncoding.EncodeToString(pk)
}

func SecpSign(ck string, msg string) string {
	if ck == "" {
		return ""
	}
	if msg == "" {
		return ""
	}

	ckbytes, err := base64.StdEncoding.DecodeString(ck)
	if err != nil {
		return ""
	}

	msgbytes, err := base64.StdEncoding.DecodeString(msg)
	if err != nil {
		return ""
	}

	b2sum := blake2b.Sum256(msgbytes)
	sig, err := crypto.Sign(ckbytes, b2sum[:])
	if err != nil {
		return ""
	}

	return base64.StdEncoding.EncodeToString(sig)
}

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}
