package wlib

import (
	"encoding/json"
	"fmt"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	big2 "github.com/filecoin-project/go-state-types/big"
	builtin5 "github.com/filecoin-project/specs-actors/v5/actors/builtin"
	"golang.org/x/xerrors"
	"strconv"

	//builtin4 "github.com/filecoin-project/specs-actors/v4/actors/builtin"
)

// GenConstructorParamV3 创建多签钱包，Spec Actor V3版本,从SpecActorV2版本开始延续
// 传入的参数 json格式
// {
//  "signers": ["f3xaczqsnxryrhirf4ptfsjb72nv3ogr5uhzsl6qd7l2zahkiaqqkw4fyeim2msfsjdi4sirimpitkc27wgv6q"],
//  "threshold": 2,
//  "unlock_duration": 120
// }
//
// 输出的结果 json格式
// {"param":"gtgqUwABVQAOZmlsLzMvbXVsdGlzaWdGhIACGHgA"}
// param 做base64编码
// 如果有错误，则返回
// {"err": "some error here"}
func GenConstructorParamV3(input string) string {
	var param GenConstructorParamInput
	err := json.Unmarshal([]byte(input), &param)
	if err != nil {
		return genOut(nil, xerrors.Errorf("invalid input json: %s, %v", input, err))
	}

	return param.TransferToSpec()
}

func GenProposeForSendParamV3(to, value string) string {
	receiver, amount, err := parseReceiverAndAmount(to, value)
	if err != nil {
		return genOut(nil,
			xerrors.Errorf("failed to parse receiver(%s) or value(%s): %v", to, value, err))
	}

	param, err := SerializeParams(&ProposeParams{
		To:     receiver,
		Value:  amount,
		Method: 0,
		Params: nil,
	})

	if err != nil {
		return genOut(nil, xerrors.Errorf("failed to serialize params: %v", err))
	}

	return genOut(param, nil)
}

func GenProposalForWithdrawBalanceV3(miner, value string) string {
	receiver, amount, err := parseReceiverAndAmount(miner, value)
	if err != nil {
		return genOut(nil,
			xerrors.Errorf("failed to parse receiver(%s) or value(%s): %v", miner, value, err))
	}

	enc, err := SerializeParams(&WithdrawBalanceParams{
		AmountRequested: amount,
	})
	if err != nil {
		return genOut(nil,
			xerrors.Errorf("failed to serialize WithdrawBalanceParams: %v", err))
	}

	enc, err = SerializeParams(&ProposeParams{
		To:     receiver,
		Value:  abi.NewTokenAmount(0),
		Method: builtin5.MethodsMiner.WithdrawBalance,
		Params: enc,
	})

	if err != nil {
		return genOut(nil,
			xerrors.Errorf("failed to serialize ProposeParams: %v", err))
	}

	return genOut(enc, nil)
}
func GenProposalForChangeWorkerAddress(miner, params string) string {
	var param GenChangeWorkerParamInput
	err := json.Unmarshal([]byte(params), &param)
	if err != nil {
		return genOut(nil, xerrors.Errorf("invalid input json: %s, %v", params, err))
	}
	return param.TransferToSpec(miner)
}

func parseReceiverAndAmount(to, value string) (address.Address, abi.TokenAmount, error) {
	receiver, err := address.NewFromString(to)
	if err != nil {
		return address.Undef, big2.Zero(), err
	}

	if value == "0" {
		return receiver, big2.Zero(), nil
	}

	amount, err := ParseFIL(fmt.Sprintf("%s afil", value))
	if err != nil {
		return address.Undef, big2.Zero(), err
	}

	return receiver, abi.TokenAmount(amount), nil
}

func GenProposalForChangeOwnerV3(self, miner, value string) string {
	if len(value) == 0 {
		value = "0"
	}

	newOwner, err := address.NewFromString(self)
	if err != nil {
		return genOut(nil,
			xerrors.Errorf("invalid newOwner address(%s): %v", self, err))
	}

	receiver, amount, err := parseReceiverAndAmount(miner, value)
	if err != nil {
		return genOut(nil,
			xerrors.Errorf("invalid miner address(%s): %v", miner, err))
	}

	enc, err := SerializeParams(&newOwner)
	if err != nil {
		return genOut(nil,
			xerrors.Errorf("failed to serialize newOwner params: %v", err))
	}

	enc, err = SerializeParams(&ProposeParams{
		To:     receiver,
		Value:  amount,
		Method: builtin5.MethodsMiner.ChangeOwnerAddress,
		Params: enc,
	})

	return genOut(enc, nil)
}

func GenApprovalV3(tx string) string {
	var txInput TransactionInput
	err := json.Unmarshal([]byte(tx), &txInput)
	if err != nil {
		return genOut(nil, xerrors.Errorf("failed to unmarshal json: %v", err))
	}

	enc, err := txInput.Transfer()
	if err != nil {
		return genOut(nil, err)
	}

	return genOut(enc, nil)
}
func GenConfirmUpdateWorkerKey(miner string) string {
	minerAddr, err := address.NewFromString(miner)
	if err != nil {
		return genOut(nil,
			xerrors.Errorf("invalid miner address(%s): %v", miner, err))
	}
	enc, err:= SerializeParams(&ProposeParams{
		To:     minerAddr,
		Value:  abi.NewTokenAmount(0),
		Method: builtin5.MethodsMiner.ConfirmUpdateWorkerKey,
		Params: nil,
	})
	return genOut(enc,nil)
}
func GenCreateMiner(ownerAddr, workerAddr, sealType string) string {
    t,e:=strconv.Atoi(sealType)
    if e!=nil{
		return genOut(nil,
			xerrors.Errorf("wrong sealType", e))
	}
	owner,_:=address.NewFromString(ownerAddr)
	worker,_:=address.NewFromString(workerAddr)
	sys,_:=address.NewFromString("f04")
	enc, err := SerializeParams(&CreateMinerParams{
		Multiaddrs: nil,
		Peer: nil,
		WindowPoStProofType: abi.RegisteredPoStProof(t),
		Worker: worker,
		Owner: owner,
	})
	if err != nil {
		return genOut(nil,
			xerrors.Errorf("failed to serialize CreateMiner: %v", err))
	}

	enc, err = SerializeParams(&ProposeParams{
		To:     sys,
		Value:  abi.NewTokenAmount(0),
		Method: builtin5.MethodsPower.CreateMiner,
		Params: enc,
	})

	if err != nil {
		return genOut(nil,
			xerrors.Errorf("failed to serialize ProposeParams: %v", err))
	}

	return genOut(enc, nil)
}