package wlib

import (
	init0 "github.com/filecoin-project/specs-actors/actors/builtin/init"
	miner0 "github.com/filecoin-project/specs-actors/actors/builtin/miner"
	multisig0 "github.com/filecoin-project/specs-actors/actors/builtin/multisig"
	 "github.com/filecoin-project/specs-actors/v5/actors/builtin/power"
	multisig2 "github.com/filecoin-project/specs-actors/v2/actors/builtin/multisig"
)

// type ConstructorParams struct {
// 	Signers               []addr.Address
// 	NumApprovalsThreshold uint64
// 	UnlockDuration        abi.ChainEpoch
// 	StartEpoch            abi.ChainEpoch
// }
type ConstructorParams = multisig2.ConstructorParams
type ExecParams = init0.ExecParams

//type ProposeParams struct {
//	To     addr.Address
//	Value  abi.TokenAmount
//	Method abi.MethodNum
//	Params []byte
//}
type ProposeParams = multisig0.ProposeParams
type WithdrawBalanceParams = miner0.WithdrawBalanceParams
type ChangeWorkerAddressParams = miner0.ChangeWorkerAddressParams
type ProposalHashData = multisig0.ProposalHashData
type Transaction = multisig0.Transaction
type TxnIDParams = multisig0.TxnIDParams
type TxnID = multisig0.TxnID
type CreateMinerParams= power.CreateMinerParams