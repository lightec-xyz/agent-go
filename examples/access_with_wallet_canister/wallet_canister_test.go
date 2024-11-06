package access_with_wallet_canister

import (
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/aviate-labs/agent-go/ic/wallet"

	"github.com/aviate-labs/agent-go"
	"github.com/aviate-labs/agent-go/candid/idl"
	"github.com/aviate-labs/agent-go/identity"
	"github.com/aviate-labs/agent-go/principal"
)

/* 1. 使用如下命令，得到 identity的 pem。
	dfx identity list
	dfx identity use xxx // one of the listed idnetity
	dfx identity export xxx

   2. 使用如下命令得到 identity的 wallet
   dfx identity get-wallet
*/

//https://a4gq6-oaaaa-aaaab-qaa4q-cai.raw.icp0.io/?id=x5kf3-raaaa-aaaah-aq2xq-cai

// replace with "dfx identity export xxx" result
var pem = []byte(`
-----BEGIN EC PRIVATE KEY-----
//identity's pem
-----END EC PRIVATE KEY-----`)

const BUILD_WITNESS_FEE = 10_000_000

func Test_BuildWitness_Mainnet(t *testing.T) {
	id, err := identity.NewSecp256k1IdentityFromPEMWithoutParameters(pem)
	if err != nil {
		panic(err)
	}

	host, err := url.Parse("https://icp0.io/")
	if err != nil {
		panic(err)
	}
	cfg := agent.Config{
		Identity:                       id,
		ClientConfig:                   &agent.ClientConfig{Host: host},
		FetchRootKey:                   true,
		PollTimeout:                    30 * time.Second,
		DisableSignedQueryVerification: false, //MUST BE TRUE TO ACCESS LOCAL REPLICA
	}
	a, err := wallet.NewAgent(principal.MustDecode("identity'w wallet"), cfg) //这里的canister id是 dfx identity get-wallet --network ic 结果
	if err != nil {
		panic(err)
	}

	balance, err := a.WalletBalance()
	if err != nil {
		panic(err)
	}
	fmt.Printf("balance:%v\n", balance)

	canisterId := principal.MustDecode("x5kf3-raaaa-aaaah-aq2xq-cai")

	var s1 string
	err = a.Query(canisterId, "genesis_sc_root", []any{}, []any{&s1})
	if err != nil {
		panic(err)
	}
	fmt.Printf("s1:%v\n", s1)

	//step4: build witness
	input4, err := idl.Marshal([]any{"79d104d3c3a5ef481febf4c53931b79ad90c55f58c826bdc9883c6e0949fe469", []string{"a27fde426e69c0b9b646c224a3c938424c5522f859ea22a8d78a3699b6526489", "2deeb285a6143c34c1695fb2871da698f2d93749484ba1dc9535aff3b48c5e98"}})
	if err != nil {
		panic(err)
	}
	fmt.Printf("input4:%v\n", input4)

	arg4 := struct {
		Canister   principal.Principal `ic:"canister" json:"canister"`
		MethodName string              `ic:"method_name" json:"method_name"`
		Args       []byte              `ic:"args" json:"args"`
		Cycles     uint64              `ic:"cycles" json:"cycles"`
	}{
		Canister:   canisterId,
		MethodName: "build_witness",
		Args:       []byte(input4),
		Cycles:     BUILD_WITNESS_FEE,
	}
	result, err := a.WalletCall(arg4)
	if err != nil {
		panic(err)
	}
	fmt.Printf("result:%v\n", result.Ok.Return)

	// type Result struct {
	// 	ok bool
	// 	res []byte
	// }

	// var r []byte
	// err = idl.Unmarshal(result.Ok.Return, []any{&r})
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("s4:%v\n", hex.EncodeToString(r))
}
