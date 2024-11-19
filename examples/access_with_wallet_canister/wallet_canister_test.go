package access_with_wallet_canister

import (
	"encoding/hex"
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

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
-----END EC PRIVATE KEY-----`)

const BUILD_WITNESS_FEE = 10_000_000
const SIGN_AND_VERIFY_FEE = 28_000_000_000

// dfx canister call x5kf3-raaaa-aaaah-aq2xq-cai build_witness '("98da4b13e7508b84712850e723b0f2f339e93df6250a0e93af7f231e85fdac8b", vec{"2cb50594fb592826b00b7e83d70ad96feef2a55a183222409883e1bfc6022662"})' --network ic --with-cycles 1000000000 --wallet 4ichq-aqaaa-aaaak-qcosa-cai
func Test_BuildWitness_Mainnet(t *testing.T) {
	id, err := identity.NewSecp256k1IdentityFromPEMWithoutParameters(pem)
	assert.NoError(t, err)

	host, err := url.Parse("https://icp0.io/")
	assert.NoError(t, err)

	cfg := agent.Config{
		Identity:                       id,
		ClientConfig:                   &agent.ClientConfig{Host: host},
		FetchRootKey:                   true,
		PollTimeout:                    30 * time.Second,
		DisableSignedQueryVerification: false, //MUST BE TRUE TO ACCESS LOCAL REPLICA
	}
	a, err := wallet.NewAgent(principal.MustDecode("4ichq-aqaaa-aaaak-qcosa-cai"), cfg) //这里的canister id是 dfx identity get-wallet --network ic 结果
	assert.NoError(t, err)

	balance, err := a.WalletBalance()
	assert.NoError(t, err)
	fmt.Printf("balance:%v\n", balance)

	canisterId := principal.MustDecode("x5kf3-raaaa-aaaah-aq2xq-cai")

	var s1 string
	err = a.Query(canisterId, "genesis_sc_root", []any{}, []any{&s1})
	assert.NoError(t, err)
	assert.Equal(t, "4d3646bb039cad3ec4347b1f6e81cac97eb579323758ef7777121459d1db965b", s1)

	//step4: build witness
	input4, err := idl.Marshal([]any{"98da4b13e7508b84712850e723b0f2f339e93df6250a0e93af7f231e85fdac8b", []string{"2cb50594fb592826b00b7e83d70ad96feef2a55a183222409883e1bfc6022662"}})
	assert.NoError(t, err)
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
	assert.NoError(t, err)
	fmt.Printf("result:%v\n", hex.EncodeToString(result.Ok.Return)) //hex.EncodeToString(result.Ok.Return)
}

func Test_VerifyAndSign_Mainnet_With_Cycles(t *testing.T) {
	id, err := identity.NewSecp256k1IdentityFromPEMWithoutParameters(pem)
	assert.NoError(t, err)

	host, err := url.Parse("https://icp0.io/")
	assert.NoError(t, err)

	cfg := agent.Config{
		Identity:                       id,
		ClientConfig:                   &agent.ClientConfig{Host: host},
		FetchRootKey:                   true,
		PollTimeout:                    30 * time.Second,
		DisableSignedQueryVerification: false, //MUST BE TRUE TO ACCESS LOCAL REPLICA
	}
	a, err := wallet.NewAgent(principal.MustDecode("4ichq-aqaaa-aaaak-qcosa-cai"), cfg) //这里的canister id是 dfx identity get-wallet --network ic 结果
	assert.NoError(t, err)

	balance, err := a.WalletBalance()
	assert.NoError(t, err)
	fmt.Printf("balance:%v\n", balance)

	canisterId := principal.MustDecode("xbsjj-qaaaa-aaaai-aqamq-cai")

	var s1 string
	err = a.Query(canisterId, "witness_builder_canister", []any{}, []any{&s1})
	assert.NoError(t, err)
	assert.Equal(t, "x5kf3-raaaa-aaaah-aq2xq-cai", s1)

	var s2 string
	err = a.Query(canisterId, "plonk_verifier_canister", []any{}, []any{&s2})
	assert.NoError(t, err)
	assert.Equal(t, "gbfw6-diaaa-aaaah-qpvia-cai", s2)

	// var s3 string
	// err = a.Query(canisterId, "public_key", []any{}, []any{&s3})
	// assert.NoError(t, err)
	// assert.Equal(t, "02971351ad0a4e80b4d61003a152c746bde6d7ac5cba52466727c611fdc8c20f5b", s3)

	//step4: build witness
	input4, err := idl.Marshal([]any{"98da4b13e7508b84712850e723b0f2f339e93df6250a0e93af7f231e85fdac8b", []string{"2cb50594fb592826b00b7e83d70ad96feef2a55a183222409883e1bfc6022662"}, "229baaee55571d2454fa2ba9f32e2b07135f0e10972009fd5e8a07b7e4025fe300ae9f755de6d73c11a82eb1b2cbc4fd2ba28c20fc215e51b7201ba31e262eb819709d1d8cfe03e7699b894d71052e6ccb5324c89825c2ecdf46d784f4a2e91510ac3e74bd63e57cd443ab4e8df6f8076742bf9a9e17cb61df35d08241eca8222c02e4390122e9c4b7bebb4a704fd78d96dba4cde80500e7240e15713301a4f8065fe4df227810ad02291f579fab794d24cba12c75c864fd47fcaa96343d2cfe05080197de901ad51b89e1035c6157b658f7b730a952e742b16ab99bb324307e12b6c16a4b4d0edd3af9d1eed0cc0ea001a8f1f3499f770ec6dd943119428db5100ac6efb2d4474c879247fd8d7132e33332c544fced798f85e24230bfd402890144705cf054a8784adf374e32da9729e541ee52963533a4a84abb04b83926aa03d7f95792708d08ae265db38423f35d5cf7b4203ccb28f0d330a7e0de694b6d0f324015bf06d3f7c5e8b5bc9f05018d98d8744f4ae99b16992adcc2b567704c1d7a69ec4b67334842aa1f105c00492fd92898a12b623ae942b0d27257c825191d2391cdc920925c37eddc3955bac7c562891b4f312a77b6cfd14beb54313bab0f92f3aa7b28e0756e51d2b6acee7f91c06f599ef8ac0fc84f3e9ca1a29f7dba05b1fc6dbb885994c78effe92941b231195cf990f4684ce80c8561840418de87134c1940edb3d7bfdc49c7202e28e795bd0547daee1e788764205c5e9b0995552896972dc9f8c8dc460f1b5a16db62ccf95fd27d02ec8972e69ee3e513bc43f01f2202ae9a09043697d1c5d49ed51813d81f51efc2aade8e480d8fff7014b6f41d09419fc87a071f116d4a04ffe99e50016b5b73bdd357c7fd72f0ac4bae6dbe25178b4dd99fe080c219741ccec1e2d713a29f9f4dd97ca323fbb9b49f572672022552b85f7e78be4dae71d3fd2adeef4549c4c584bc9787c51599a6f73c84382768aa27767f314f07cd3b8bcc4bb842daf2b3f43c9918e3513652ad8c9205071403995b6014a4dd5ef2fa03f1b497d0ec387be2a5f4c044f0f33f222a6e8ed5065ae326e0a46db583e874ce5567f0a215ed80ab39fa49159579529b5808f891012accdbb8424f962502df0a3b88b1abeca31104ae5094e83f4121b4774bd94621811743f001c3cb970bfc1f7b0125429963df1c8e4164183b38ad15cfacce942bd22ebc1146292580cee8469b4d870f4e3093987d2efbdb9e3e2147fbbfe0332835a812b42955a77a17c659c8c7b582e05874833abef31cdeb2c2c5b718de13"})
	assert.NoError(t, err)

	arg4 := struct {
		Canister   principal.Principal `ic:"canister" json:"canister"`
		MethodName string              `ic:"method_name" json:"method_name"`
		Args       []byte              `ic:"args" json:"args"`
		Cycles     uint64              `ic:"cycles" json:"cycles"`
	}{
		Canister:   canisterId,
		MethodName: "verify_and_sign",
		Args:       []byte(input4),
		Cycles:     SIGN_AND_VERIFY_FEE,
	}
	result, err := a.WalletCall(arg4)
	assert.NoError(t, err)
	fmt.Printf("result:%v\n", result.Ok.Return)
}

// dfx canister call xbsjj-qaaaa-aaaai-aqamq-cai verify_and_sign_free '("98da4b13e7508b84712850e723b0f2f339e93df6250a0e93af7f231e85fdac8b", vec {"2cb50594fb592826b00b7e83d70ad96feef2a55a183222409883e1bfc6022662"}, "229baaee55571d2454fa2ba9f32e2b07135f0e10972009fd5e8a07b7e4025fe300ae9f755de6d73c11a82eb1b2cbc4fd2ba28c20fc215e51b7201ba31e262eb819709d1d8cfe03e7699b894d71052e6ccb5324c89825c2ecdf46d784f4a2e91510ac3e74bd63e57cd443ab4e8df6f8076742bf9a9e17cb61df35d08241eca8222c02e4390122e9c4b7bebb4a704fd78d96dba4cde80500e7240e15713301a4f8065fe4df227810ad02291f579fab794d24cba12c75c864fd47fcaa96343d2cfe05080197de901ad51b89e1035c6157b658f7b730a952e742b16ab99bb324307e12b6c16a4b4d0edd3af9d1eed0cc0ea001a8f1f3499f770ec6dd943119428db5100ac6efb2d4474c879247fd8d7132e33332c544fced798f85e24230bfd402890144705cf054a8784adf374e32da9729e541ee52963533a4a84abb04b83926aa03d7f95792708d08ae265db38423f35d5cf7b4203ccb28f0d330a7e0de694b6d0f324015bf06d3f7c5e8b5bc9f05018d98d8744f4ae99b16992adcc2b567704c1d7a69ec4b67334842aa1f105c00492fd92898a12b623ae942b0d27257c825191d2391cdc920925c37eddc3955bac7c562891b4f312a77b6cfd14beb54313bab0f92f3aa7b28e0756e51d2b6acee7f91c06f599ef8ac0fc84f3e9ca1a29f7dba05b1fc6dbb885994c78effe92941b231195cf990f4684ce80c8561840418de87134c1940edb3d7bfdc49c7202e28e795bd0547daee1e788764205c5e9b0995552896972dc9f8c8dc460f1b5a16db62ccf95fd27d02ec8972e69ee3e513bc43f01f2202ae9a09043697d1c5d49ed51813d81f51efc2aade8e480d8fff7014b6f41d09419fc87a071f116d4a04ffe99e50016b5b73bdd357c7fd72f0ac4bae6dbe25178b4dd99fe080c219741ccec1e2d713a29f9f4dd97ca323fbb9b49f572672022552b85f7e78be4dae71d3fd2adeef4549c4c584bc9787c51599a6f73c84382768aa27767f314f07cd3b8bcc4bb842daf2b3f43c9918e3513652ad8c9205071403995b6014a4dd5ef2fa03f1b497d0ec387be2a5f4c044f0f33f222a6e8ed5065ae326e0a46db583e874ce5567f0a215ed80ab39fa49159579529b5808f891012accdbb8424f962502df0a3b88b1abeca31104ae5094e83f4121b4774bd94621811743f001c3cb970bfc1f7b0125429963df1c8e4164183b38ad15cfacce942bd22ebc1146292580cee8469b4d870f4e3093987d2efbdb9e3e2147fbbfe0332835a812b42955a77a17c659c8c7b582e05874833abef31cdeb2c2c5b718de13")' --with-cycles 28000000000 --wallet 4ichq-aqaaa-aaaak-qcosa-cai
// dfx canister call xbsjj-qaaaa-aaaai-aqamq-cai verify_and_sign '("98da4b13e7508b84712850e723b0f2f339e93df6250a0e93af7f231e85fdac8b", vec {"2cb50594fb592826b00b7e83d70ad96feef2a55a183222409883e1bfc6022662"}, "229baaee55571d2454fa2ba9f32e2b07135f0e10972009fd5e8a07b7e4025fe300ae9f755de6d73c11a82eb1b2cbc4fd2ba28c20fc215e51b7201ba31e262eb819709d1d8cfe03e7699b894d71052e6ccb5324c89825c2ecdf46d784f4a2e91510ac3e74bd63e57cd443ab4e8df6f8076742bf9a9e17cb61df35d08241eca8222c02e4390122e9c4b7bebb4a704fd78d96dba4cde80500e7240e15713301a4f8065fe4df227810ad02291f579fab794d24cba12c75c864fd47fcaa96343d2cfe05080197de901ad51b89e1035c6157b658f7b730a952e742b16ab99bb324307e12b6c16a4b4d0edd3af9d1eed0cc0ea001a8f1f3499f770ec6dd943119428db5100ac6efb2d4474c879247fd8d7132e33332c544fced798f85e24230bfd402890144705cf054a8784adf374e32da9729e541ee52963533a4a84abb04b83926aa03d7f95792708d08ae265db38423f35d5cf7b4203ccb28f0d330a7e0de694b6d0f324015bf06d3f7c5e8b5bc9f05018d98d8744f4ae99b16992adcc2b567704c1d7a69ec4b67334842aa1f105c00492fd92898a12b623ae942b0d27257c825191d2391cdc920925c37eddc3955bac7c562891b4f312a77b6cfd14beb54313bab0f92f3aa7b28e0756e51d2b6acee7f91c06f599ef8ac0fc84f3e9ca1a29f7dba05b1fc6dbb885994c78effe92941b231195cf990f4684ce80c8561840418de87134c1940edb3d7bfdc49c7202e28e795bd0547daee1e788764205c5e9b0995552896972dc9f8c8dc460f1b5a16db62ccf95fd27d02ec8972e69ee3e513bc43f01f2202ae9a09043697d1c5d49ed51813d81f51efc2aade8e480d8fff7014b6f41d09419fc87a071f116d4a04ffe99e50016b5b73bdd357c7fd72f0ac4bae6dbe25178b4dd99fe080c219741ccec1e2d713a29f9f4dd97ca323fbb9b49f572672022552b85f7e78be4dae71d3fd2adeef4549c4c584bc9787c51599a6f73c84382768aa27767f314f07cd3b8bcc4bb842daf2b3f43c9918e3513652ad8c9205071403995b6014a4dd5ef2fa03f1b497d0ec387be2a5f4c044f0f33f222a6e8ed5065ae326e0a46db583e874ce5567f0a215ed80ab39fa49159579529b5808f891012accdbb8424f962502df0a3b88b1abeca31104ae5094e83f4121b4774bd94621811743f001c3cb970bfc1f7b0125429963df1c8e4164183b38ad15cfacce942bd22ebc1146292580cee8469b4d870f4e3093987d2efbdb9e3e2147fbbfe0332835a812b42955a77a17c659c8c7b582e05874833abef31cdeb2c2c5b718de13")' --network ic --with-cycles 28000000000 --wallet 4ichq-aqaaa-aaaak-qcosa-cai
func Test_VerifyandSign_Mainnet_Free(t *testing.T) {
	id, err := identity.NewSecp256k1IdentityFromPEMWithoutParameters(pem)
	assert.NoError(t, err)

	host, err := url.Parse("https://icp0.io/")
	assert.NoError(t, err)

	cfg := agent.Config{
		Identity:                       id,
		ClientConfig:                   &agent.ClientConfig{Host: host},
		FetchRootKey:                   true,
		PollTimeout:                    30 * time.Second,
		DisableSignedQueryVerification: false, //MUST BE TRUE TO ACCESS LOCAL REPLICA
	}

	a, err := agent.New(cfg)
	assert.NoError(t, err)

	zkVerifierPrincipal := principal.MustDecode("xbsjj-qaaaa-aaaai-aqamq-cai")

	txIDHex := "98da4b13e7508b84712850e723b0f2f339e93df6250a0e93af7f231e85fdac8b"
	signHashesHex := []string{"2cb50594fb592826b00b7e83d70ad96feef2a55a183222409883e1bfc6022662"}
	proofHex := "229baaee55571d2454fa2ba9f32e2b07135f0e10972009fd5e8a07b7e4025fe300ae9f755de6d73c11a82eb1b2cbc4fd2ba28c20fc215e51b7201ba31e262eb819709d1d8cfe03e7699b894d71052e6ccb5324c89825c2ecdf46d784f4a2e91510ac3e74bd63e57cd443ab4e8df6f8076742bf9a9e17cb61df35d08241eca8222c02e4390122e9c4b7bebb4a704fd78d96dba4cde80500e7240e15713301a4f8065fe4df227810ad02291f579fab794d24cba12c75c864fd47fcaa96343d2cfe05080197de901ad51b89e1035c6157b658f7b730a952e742b16ab99bb324307e12b6c16a4b4d0edd3af9d1eed0cc0ea001a8f1f3499f770ec6dd943119428db5100ac6efb2d4474c879247fd8d7132e33332c544fced798f85e24230bfd402890144705cf054a8784adf374e32da9729e541ee52963533a4a84abb04b83926aa03d7f95792708d08ae265db38423f35d5cf7b4203ccb28f0d330a7e0de694b6d0f324015bf06d3f7c5e8b5bc9f05018d98d8744f4ae99b16992adcc2b567704c1d7a69ec4b67334842aa1f105c00492fd92898a12b623ae942b0d27257c825191d2391cdc920925c37eddc3955bac7c562891b4f312a77b6cfd14beb54313bab0f92f3aa7b28e0756e51d2b6acee7f91c06f599ef8ac0fc84f3e9ca1a29f7dba05b1fc6dbb885994c78effe92941b231195cf990f4684ce80c8561840418de87134c1940edb3d7bfdc49c7202e28e795bd0547daee1e788764205c5e9b0995552896972dc9f8c8dc460f1b5a16db62ccf95fd27d02ec8972e69ee3e513bc43f01f2202ae9a09043697d1c5d49ed51813d81f51efc2aade8e480d8fff7014b6f41d09419fc87a071f116d4a04ffe99e50016b5b73bdd357c7fd72f0ac4bae6dbe25178b4dd99fe080c219741ccec1e2d713a29f9f4dd97ca323fbb9b49f572672022552b85f7e78be4dae71d3fd2adeef4549c4c584bc9787c51599a6f73c84382768aa27767f314f07cd3b8bcc4bb842daf2b3f43c9918e3513652ad8c9205071403995b6014a4dd5ef2fa03f1b497d0ec387be2a5f4c044f0f33f222a6e8ed5065ae326e0a46db583e874ce5567f0a215ed80ab39fa49159579529b5808f891012accdbb8424f962502df0a3b88b1abeca31104ae5094e83f4121b4774bd94621811743f001c3cb970bfc1f7b0125429963df1c8e4164183b38ad15cfacce942bd22ebc1146292580cee8469b4d870f4e3093987d2efbdb9e3e2147fbbfe0332835a812b42955a77a17c659c8c7b582e05874833abef31cdeb2c2c5b718de13"

	var valid bool
	var texts []string

	err = a.Call(zkVerifierPrincipal, "verify_and_sign_free", []any{txIDHex, signHashesHex, proofHex}, []any{&valid, &texts})
	assert.NoError(t, err)
	assert.True(t, valid)

	for i, text := range texts {
		fmt.Printf("%v: text %v\n", i, text)
	}
}

func Test_VerifyandSign_Mainnet_With_Cycles_Debug(t *testing.T) {
	id, err := identity.NewSecp256k1IdentityFromPEMWithoutParameters(pem)
	assert.NoError(t, err)

	host, err := url.Parse("https://icp0.io/")
	assert.NoError(t, err)

	cfg := agent.Config{
		Identity:                       id,
		ClientConfig:                   &agent.ClientConfig{Host: host},
		FetchRootKey:                   true,
		PollTimeout:                    30 * time.Second,
		DisableSignedQueryVerification: false, //MUST BE TRUE TO ACCESS LOCAL REPLICA
	}
	a, err := wallet.NewAgent(principal.MustDecode("4ichq-aqaaa-aaaak-qcosa-cai"), cfg) //这里的canister id是 dfx identity get-wallet --network ic 结果
	assert.NoError(t, err)

	balance, err := a.WalletBalance()
	assert.NoError(t, err)
	fmt.Printf("balance:%v\n", balance)

	canisterId := principal.MustDecode("xbsjj-qaaaa-aaaai-aqamq-cai")

	var s1 string
	err = a.Query(canisterId, "witness_builder_canister", []any{}, []any{&s1})
	assert.NoError(t, err)
	assert.Equal(t, "x5kf3-raaaa-aaaah-aq2xq-cai", s1)

	var s2 string
	err = a.Query(canisterId, "plonk_verifier_canister", []any{}, []any{&s2})
	assert.NoError(t, err)
	assert.Equal(t, "gbfw6-diaaa-aaaah-qpvia-cai", s2)

	//step4: build witness
	txIDHex := "98da4b13e7508b84712850e723b0f2f339e93df6250a0e93af7f231e85fdac8b"
	signHashesHex := []string{"2cb50594fb592826b00b7e83d70ad96feef2a55a183222409883e1bfc6022662"}
	proofHex := "229baaee55571d2454fa2ba9f32e2b07135f0e10972009fd5e8a07b7e4025fe300ae9f755de6d73c11a82eb1b2cbc4fd2ba28c20fc215e51b7201ba31e262eb819709d1d8cfe03e7699b894d71052e6ccb5324c89825c2ecdf46d784f4a2e91510ac3e74bd63e57cd443ab4e8df6f8076742bf9a9e17cb61df35d08241eca8222c02e4390122e9c4b7bebb4a704fd78d96dba4cde80500e7240e15713301a4f8065fe4df227810ad02291f579fab794d24cba12c75c864fd47fcaa96343d2cfe05080197de901ad51b89e1035c6157b658f7b730a952e742b16ab99bb324307e12b6c16a4b4d0edd3af9d1eed0cc0ea001a8f1f3499f770ec6dd943119428db5100ac6efb2d4474c879247fd8d7132e33332c544fced798f85e24230bfd402890144705cf054a8784adf374e32da9729e541ee52963533a4a84abb04b83926aa03d7f95792708d08ae265db38423f35d5cf7b4203ccb28f0d330a7e0de694b6d0f324015bf06d3f7c5e8b5bc9f05018d98d8744f4ae99b16992adcc2b567704c1d7a69ec4b67334842aa1f105c00492fd92898a12b623ae942b0d27257c825191d2391cdc920925c37eddc3955bac7c562891b4f312a77b6cfd14beb54313bab0f92f3aa7b28e0756e51d2b6acee7f91c06f599ef8ac0fc84f3e9ca1a29f7dba05b1fc6dbb885994c78effe92941b231195cf990f4684ce80c8561840418de87134c1940edb3d7bfdc49c7202e28e795bd0547daee1e788764205c5e9b0995552896972dc9f8c8dc460f1b5a16db62ccf95fd27d02ec8972e69ee3e513bc43f01f2202ae9a09043697d1c5d49ed51813d81f51efc2aade8e480d8fff7014b6f41d09419fc87a071f116d4a04ffe99e50016b5b73bdd357c7fd72f0ac4bae6dbe25178b4dd99fe080c219741ccec1e2d713a29f9f4dd97ca323fbb9b49f572672022552b85f7e78be4dae71d3fd2adeef4549c4c584bc9787c51599a6f73c84382768aa27767f314f07cd3b8bcc4bb842daf2b3f43c9918e3513652ad8c9205071403995b6014a4dd5ef2fa03f1b497d0ec387be2a5f4c044f0f33f222a6e8ed5065ae326e0a46db583e874ce5567f0a215ed80ab39fa49159579529b5808f891012accdbb8424f962502df0a3b88b1abeca31104ae5094e83f4121b4774bd94621811743f001c3cb970bfc1f7b0125429963df1c8e4164183b38ad15cfacce942bd22ebc1146292580cee8469b4d870f4e3093987d2efbdb9e3e2147fbbfe0332835a812b42955a77a17c659c8c7b582e05874833abef31cdeb2c2c5b718de13"

	input4, err := idl.Marshal([]any{txIDHex, signHashesHex, proofHex})
	assert.NoError(t, err)

	arg4 := struct {
		Canister   principal.Principal `ic:"canister" json:"canister"`
		MethodName string              `ic:"method_name" json:"method_name"`
		Args       []byte              `ic:"args" json:"args"`
		Cycles     uint64              `ic:"cycles" json:"cycles"`
	}{
		Canister:   canisterId,
		MethodName: "verify_and_sign_free",
		Args:       []byte(input4),
		Cycles:     SIGN_AND_VERIFY_FEE,
	}
	result, err := a.WalletCall(arg4)
	assert.NoError(t, err)
	fmt.Printf("result:%v\n", result.Ok.Return)
}

func Test_VerifyandSign_Mainnet_redeem_0df910a0a0aba382efdda61a52043b2c4f2ea6b62cdefbb17f4e8ee13845c7d5(t *testing.T) {
	id, err := identity.NewSecp256k1IdentityFromPEMWithoutParameters(pem)
	assert.NoError(t, err)

	host, err := url.Parse("https://icp0.io/")
	assert.NoError(t, err)

	cfg := agent.Config{
		Identity:                       id,
		ClientConfig:                   &agent.ClientConfig{Host: host},
		FetchRootKey:                   true,
		PollTimeout:                    30 * time.Second,
		DisableSignedQueryVerification: false, //MUST BE TRUE TO ACCESS LOCAL REPLICA
	}
	a, err := wallet.NewAgent(principal.MustDecode("4ichq-aqaaa-aaaak-qcosa-cai"), cfg) //这里的canister id是 dfx identity get-wallet --network ic 结果
	assert.NoError(t, err)

	balance, err := a.WalletBalance()
	assert.NoError(t, err)
	fmt.Printf("balance:%v\n", balance)

	canisterId := principal.MustDecode("5jpjw-iaaaa-aaaak-qiuqa-cai")

	var s1 string
	err = a.Query(canisterId, "witness_builder_canister", []any{}, []any{&s1})
	assert.NoError(t, err)
	assert.Equal(t, "47er5-5qaaa-aaaak-qiuva-cai", s1)

	var s2 string
	err = a.Query(canisterId, "plonk_verifier_canister", []any{}, []any{&s2})
	assert.NoError(t, err)
	assert.Equal(t, "ysvvh-oqaaa-aaaak-qiupa-cai", s2)

	//step4: build witness
	txIDHex := "0F9D4C3FC125D6F3BB59B5323932283943F5D7C0F812DAC83A43F3914330184C"
	signHashesHex := []string{"37A53D761E4E9A9C0D7A27DB31720E15E0FDE21C3527EBE78EBA3E37A032F782", "D63F90B7C4151E4CF5763E5ACA488E0002331A7BFE8E5D3CE6C22CD946CAA04B", "259B02FED3B6874E8503769C4EE28153439994B8FA5845347046EEC88E7B11AD", "DF0DF8900DB9839EC023F5DE01F386F24F910875D0DEDC3E656929BABAE8A3A4"}
	proofHex := "20d5a9dd6fdff65a217259d6be1f44adb89e989a544eccb7ab7cffd8da29cb471d933e0c3d68694ad8ef3f75cfeaebf257e0739132487a393b47ab73cd968657281083a595c886a1f08c647587b529f6648d0685816d7becf743c77194ca8c0f2e068e395a51f4dd81408aef87256cd3a49b7dc9b154e9bfeeb5eae57abaff592168b757ca5ad1dfacc881a0f274b82d99450775269b16b0e8ca68740ca7eda2087544f3a0eb8dad09b789e6bf7015241ab1b22d07f8935e7c30a4b61da0ab8c044214fe0f4d327dc54bf7aa42c8dfcbcc7d6e6fdb60482ebb40d31592d213d60119a221a4ddeb8ec8046a0c97ab4a5dc9fd9c62a6d5fc6e7e0248f107c54cec26fd12aeb4a6fe5e104835eb840b5cbecfec71badaf34fb6bbfc55c1038c6e9f084d2a663cfc7b00222f704975436378f222bc9e7408c5c43b0a941f681839a82b6deaaf74eb2b3296132cbe3d933effac47ef7f8b8c999ce88bb355ec5b98e31bac52e9dc29c3c9b2b6405664cc504d78bcec66dd767328cf6c1c384991ffd50e4224fa492dc1afa95fd41a728b078c03448c39ef5402aa44c099d6103a00f908b2823c0905535e35caa7c121171f35984a988e435ece2b6e8aed27d50601b904d862bedd7e2bebfe88fcae9b8e5732dee24e6b6f104d55b597213f6b9ee2732d71932cb49ddff2537cf26624f39b6d65d22b96a5a341a1ffdd2e66d25ad4fc25f3bb290e5288a1e944ea0e31541db7fbb63394ab149bb35d4263e8d68fdc4e0f13b9ac4e6db0725196f34ae82043a92cf65b55b65661635cd31d80793dd3022d23e9deabd833085119f97fecf7011126b296952581f9c74379662b1136434300f3b578df057e986af4185bec6167f79a79d1b5e6da3c3aaf068d8610de42b8208f2b0b9c97770d230a7be74dc7b44256b41f6e3182f5d8f46a7043ab44994f1852315f915c6bd2c37013e84a76c43b85df3c662af4ef00698374059d7d13352251a4fc2f78bd5f219ffd54d3d4978264a8038745950fef9e854ab8e7d2c0331021dd1b4e66c428e00c83cb01571f2693b66bf1e6bb09b32d1caf31b6b925800123809da63ff127f53a642ddd90b4730916c6ed2f6b1277557de68619c8e1ff0e4b4aa08ab94044b916f014c0c38f557bfbf6488e489ebceadaa6b5785488a4017d739805c3f921a6957290e194e525aec78ba55793a40163dadc8db96bb7641f0090749cc5a89d6d28772ebd78d6344064a29dea7a72fff322c572f98250de0ad5ff15baed2d7f72147e134f408c901ab1c5dc78ec2a79afe9af51f8e576ec"

	input4, err := idl.Marshal([]any{txIDHex, signHashesHex, proofHex})
	assert.NoError(t, err)

	arg4 := struct {
		Canister   principal.Principal `ic:"canister" json:"canister"`
		MethodName string              `ic:"method_name" json:"method_name"`
		Args       []byte              `ic:"args" json:"args"`
		Cycles     uint64              `ic:"cycles" json:"cycles"`
	}{
		Canister:   canisterId,
		MethodName: "verify_and_sign",
		Args:       input4,
		Cycles:     SIGN_AND_VERIFY_FEE,
	}
	result, err := a.WalletCall(arg4)
	assert.NoError(t, err)
	fmt.Printf("result:%v\n", result.Ok.Return)
}
