package wlib

import (
	"encoding/base64"
	"encoding/json"
	"github.com/filecoin-project/go-state-types/abi"
	builtin5 "github.com/filecoin-project/specs-actors/v5/actors/builtin"
	"github.com/smartystreets/assertions"
	"log"
	"testing"
)

type ret struct {
	Err   string `json:"err"`
	Param string `json:"param"`
}

func TestGenConstructorParam(t *testing.T) {
	data := `{"signers": ["f3xaczqsnxryrhirf4ptfsjb72nv3ogr5uhzsl6qd7l2zahkiaqqkw4fyeim2msfsjdi4sirimpitkc27wgv6q"], "threshold": 2, "unlock_duration": 120}`
	result := GenConstructorParamV3(data)
	var r ret
	err := json.Unmarshal([]byte(result), &r)
	if err != nil {
		log.Fatalf("返回的数据不正确：%v", err)
	}
	
	// {"param":"hIACGHgA"}
	chk := assertions.ShouldEqual(r.Err, "")
	if chk != "" {
		t.Fatal(chk)
	}
	t.Logf("result %s", result)
}

func TestGenProposeForSendParamV3(t *testing.T) {
	result := GenProposeForSendParamV3(
		"f3xaczqsnxryrhirf4ptfsjb72nv3ogr5uhzsl6qd7l2zahkiaqqkw4fyeim2msfsjdi4sirimpitkc27wgv6q",
		"2187000000000000000")
	
	var r ret
	err := json.Unmarshal([]byte(result), &r)
	if err != nil {
		log.Fatalf("返回的数据不正确：%v", err)
	}
	
	chk := assertions.ShouldEqual(r.Err, "")
	if chk != "" {
		t.Fatal(chk)
	}
	t.Logf("result %s", result)
	// Output: {"param":"hFgxA7gFmEm3jiJ0RLx8yySH+m1240e0PmS/QH9esgOpAIQVbhcEQzTJFkkaOSRFDHomoUkAHlnI6avHgAAAQA=="}
}

func TestGenProposalForWithdrawBalanceV3(t *testing.T) {
	result := GenProposalForWithdrawBalanceV3(
		"f02438",
		"2187000000000000000")
	
	var r ret
	err := json.Unmarshal([]byte(result), &r)
	if err != nil {
		log.Fatalf("返回的数据不正确：%v", err)
	}
	
	chk := assertions.ShouldEqual(r.Err, "")
	if chk != "" {
		t.Fatal(chk)
	}
	t.Logf("result %s", result)
	// Output: {"param":"hEMAhhNJAB5ZyOmrx4AAEEuBSQAeWcjpq8eAAA=="}
}

func TestGenProposalForChangeOwnerV3(t *testing.T) {
	result := GenProposalForChangeOwnerV3("f3xaczqsnxryrhirf4ptfsjb72nv3ogr5uhzsl6qd7l2zahkiaqqkw4fyeim2msfsjdi4sirimpitkc27wgv6q",
		"f02438",
		"2187000000000000000")
	
	var r ret
	err := json.Unmarshal([]byte(result), &r)
	if err != nil {
		log.Fatalf("返回的数据不正确：%v", err)
	}
	
	chk := assertions.ShouldEqual(r.Err, "")
	if chk != "" {
		t.Fatal(chk)
	}
	t.Logf("result %s", result)
	
	// Output: {"param":"hEMAhhNJAB5ZyOmrx4AAF1gzWDEDuAWYSbeOInREvHzLJIf6bXbjR7Q+ZL9Af16yA6kAhBVuFwRDNMkWSRo5JEUMeiah"}
}

func TestGenApprovalV3(t *testing.T) {
	enc, _ := SerializeParams(&WithdrawBalanceParams{
		AmountRequested: abi.NewTokenAmount(100000000000000000),
	})
	trans := &TransactionInput{
		TxID:      45,
		Requester: "f3xaczqsnxryrhirf4ptfsjb72nv3ogr5uhzsl6qd7l2zahkiaqqkw4fyeim2msfsjdi4sirimpitkc27wgv6q",
		To:        "f02438",
		Value:     "0",
		Method:    uint64(builtin5.MethodsMiner.WithdrawBalance),
		Params:    base64.StdEncoding.EncodeToString(enc),
	}
	
	payload, err := json.Marshal(trans)
	if err != nil {
		t.Fatalf("序列化json失败: %v", err)
	}
	result := GenApprovalV3(string(payload))
	
	var r ret
	err = json.Unmarshal([]byte(result), &r)
	if err != nil {
		log.Fatalf("返回的数据不正确：%v", err)
	}
	
	chk := assertions.ShouldEqual(r.Err, "")
	if chk != "" {
		t.Fatal(chk)
	}
	t.Logf("result %s", result)
	// Output: {"param":"ghgtWCCeUPHlmW2dfl7PfSCxH8chhwjk8WgXgbZZnmSLRHUHew=="}
}
