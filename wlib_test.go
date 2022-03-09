package wlib_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.forceup.in/FilecoinWallet/FilWallet/wlib"
)

func TestAddressFromPK(t *testing.T) {
	addr := wlib.GenAddress("122333", "bls")
	require.Equal(t, addr, "")

	a2 := wlib.GenAddress("jYpwp93MHABEeFfrcDnUo6Hpi0eqEuzARRHVhDRwfkSEd9lhm80EU3Sg+RI28PFu", "bls")
	require.Equal(t, a2, "t3rwfhbj65zqoaardyk7vxaoouuoq6tc2hvijozqcfchkyindqpzcii56zmgn42bctosqpserw6dyw5y2uf5aq")

	a3 := wlib.GenAddress("BENtLYLjXvhSHLCDdmMdg/cRHmfsfWgs/vWsMMrUJHAjBhCRolg+f6aBLThC+9xddiCoeb3f2xyuzPYVqwNsWQI=", "secp")
	require.Equal(t, a3, "t1lrgw6ss5nu5lbhqmmtthc7hmxg6hlt5r6txpy3i")
}

func TestAddressFromString(t *testing.T) {
	addr := wlib.AddressFromString("t020146")
	require.Equal(t, addr, "t020146")
}

func TestMessageCid(t *testing.T) {
	cidStr := wlib.MessageCid(`
	{
		"Version": 0,
		"To": "f125p5nhte6kwrigoxrcaxftwpinlgspfnqd2zaui",
		"From": "f153zbrv25wvfrqf2vrvlk2qmpietuu6wexiyerja",
		"Nonce": 0,
		"Value": "10000000000000000000",
		"GasLimit": 1000000000000,
		"GasFeeCap": "10000000",
		"GasPremium": "10000000",
		"Method": 0,
		"Params": ""
	}
	`)
	require.Equal(t, cidStr, "AXGg5AIgA7aUiB+WKlJZi77CrBo4OgwytRmXbBXj8ratzAtshGM=")
}

func TestSecpPrivateToPublic(t *testing.T) {
	privateKey := "p7ZGtfT3MyOdkVaEaE2LzT12fcl2N95jsiYuvBZZ1NA="
	publicKey := "BN4D/kPsYngi68E2wyEsJgqeaIyv/nqBK07s7TokD1CUUtGkVTtJmvMBCE2b0ygksRzVDXYJJzUweVRWDfn0SB0="
	pk := wlib.SecpPrivateToPublic(privateKey)

	require.Equal(t, pk, publicKey)

	sig := wlib.SecpSign(privateKey, "SGVsbG8gV29ybGQh")

	require.Equal(t, sig, "CjSBxOfeEIyWJuKgo7od+wrd+xGZJDbOkIBg+sIo9kZ3bWP0WhIYhkS9Pf9hIbtftszIuzHKFulT0hneCFBEMwE=")
}

func TestSecpSign(t *testing.T) {
	cidStr := wlib.MessageCid(`
	{
		"Version": 0,
		"To": "f125p5nhte6kwrigoxrcaxftwpinlgspfnqd2zaui",
		"From": "f153zbrv25wvfrqf2vrvlk2qmpietuu6wexiyerja",
		"Nonce": 0,
		"Value": "10000000000000000000",
		"GasLimit": 1000000000000,
		"GasFeeCap": "10000000",
		"GasPremium": "10000000",
		"Method": 0,
		"Params": ""
	}
	`)

	require.Equal(t, cidStr, "AXGg5AIgA7aUiB+WKlJZi77CrBo4OgwytRmXbBXj8ratzAtshGM=")

	sig := wlib.SecpSign("67WMRDA2ldmfcQ87DSHCy+ppKs3iSyNjxfBD7dR68Qw=", cidStr)

	require.Equal(t, sig, "jHF0ghnCwyl7XNEfgXx1+9sjbg3lJe09gEux/+m5pRFudpQEeFxxt9ZACHNDE//u31r3GBZ4aYixpV8xYp57HgA=")
}
