package wlib

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/exitcode"
	builtin5 "github.com/filecoin-project/specs-actors/v5/actors/builtin"
	"github.com/minio/blake2b-simd"
	"golang.org/x/xerrors"

	cbg "github.com/whyrusleeping/cbor-gen"
)

type Out struct {
	Err   string `json:"err,omitempty"`
	Param []byte `json:"param,omitempty"`
}

type GenConstructorParamInput struct {
	Signers        []string `json:"signers"`
	Threshold      uint64   `json:"threshold"`
	UnlockDuration int64    `json:"unlock_duration"`
	StartEpoch     int64    `json:"start_epoch"`
}
type GenChangeWorkerParamInput struct {
	NewWorker       string   `json:"new_worker"`
	NewControlAddrs []string `json:"new_control_addrs"`
}

func (g *GenChangeWorkerParamInput) TransferToSpec(miner string) string {
	receiver, err := address.NewFromString(miner)
	if err != nil {
		return genOut(nil,
			xerrors.Errorf("failed to parse receiver(%s): %v", miner, err))
	}
	worker, e := address.NewFromString(g.NewWorker)
	var controllers []address.Address
	if e != nil {
		return genOut(nil, e)
	}
	for _, c := range g.NewControlAddrs {
		addr, err := address.NewFromString(c)
		if err != nil {
			return genOut(nil, err)
		}
		controllers = append(controllers, addr)
	}
	spec := &ChangeWorkerAddressParams{
		NewControlAddrs: controllers,
		NewWorker:       worker,
	}
	encParams, err := SerializeParams(spec)
	if err != nil {
		return genOut(nil, err)
	}
	enc, err := SerializeParams(&ProposeParams{
		To:     receiver,
		Value:  abi.NewTokenAmount(0),
		Method: builtin5.MethodsMiner.ChangeWorkerAddress,
		Params: encParams,
	})
	return genOut(enc, nil)

}
func (g *GenConstructorParamInput) TransferToSpec() string {
	var signers []address.Address
	for _, s := range g.Signers {
		addr, err := address.NewFromString(s)
		if err != nil {
			return genOut(nil, err)
		}
		signers = append(signers, addr)
	}

	spec := &ConstructorParams{
		Signers:               signers,
		NumApprovalsThreshold: g.Threshold,
		UnlockDuration:        abi.ChainEpoch(g.UnlockDuration),
		StartEpoch:            abi.ChainEpoch(g.StartEpoch),
	}

	enc, err := SerializeParams(spec)
	if err != nil {
		return genOut(nil, err)
	}

	execParams := &ExecParams{
		CodeCID:           builtin5.MultisigActorCodeID,
		ConstructorParams: enc,
	}

	enc, err = SerializeParams(execParams)
	if err != nil {
		return genOut(nil, err)
	}

	return genOut(enc, nil)
}

type TransactionInput struct {
	TxID      int64  `json:"tx_id"`
	Requester string `json:"requester"` // Requester必须填提案发起人账号的ID 比如f01234
	To        string `json:"to"`
	Value     string `json:"value"`
	Method    uint64 `json:"method"`
	Params    string `json:"params,omitempty"`
}

func (d *TransactionInput) Transfer() ([]byte, error) {
	receiver, err := address.NewFromString(d.To)
	if err != nil {
		return nil, xerrors.Errorf("invalid receiver(%s): %v", d.To, err)
	}

	requester, err := address.NewFromString(d.Requester)
	if err != nil {
		return nil, xerrors.Errorf("invalid receiver(%s): %v", d.To, err)
	}

	amount, err := ParseFIL(fmt.Sprintf("%s afil", d.Value))
	if err != nil {
		return nil, xerrors.Errorf("invalid value(%s): %v", d.Value, err)
	}

	param, err := base64.StdEncoding.DecodeString(d.Params)
	if err != nil {
		return nil, xerrors.Errorf("invalid param(%s): %v", d.Params, err)
	}

	hashData := ProposalHashData{
		Requester: requester,
		To:        receiver,
		Value:     BigInt(amount),
		Method:    abi.MethodNum(d.Method),
		Params:    param,
	}

	data, err := hashData.Serialize()
	if err != nil {
		return nil, xerrors.Errorf("failed to serialize proposal hash data: %v", err)
	}

	hash := blake2b.Sum256(data)

	enc, err := SerializeParams(&TxnIDParams{
		ID:           TxnID(d.TxID),
		ProposalHash: hash[:],
	})

	if err != nil {
		return nil, xerrors.Errorf("failed to serialize TxnIDParams: %v", err)
	}

	return enc, nil
}

func genOut(param []byte, err error) string {
	out := &Out{}
	if err != nil {
		out.Err = err.Error()
	}

	// json会自动做base64
	if param != nil {
		out.Param = param
	}

	result, err := json.Marshal(out)
	if err != nil {
		return fmt.Sprintf(`{"err":  "%s"}`, err)
	}

	return string(result)
}

func SerializeParams(i cbg.CBORMarshaler) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := i.MarshalCBOR(buf); err != nil {
		return nil, xerrors.Errorf("failed to encode parameter(ExitCode: %d): %v", exitcode.ErrSerialization, err)
	}
	return buf.Bytes(), nil
}
