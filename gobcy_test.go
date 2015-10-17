//Tests for the BlockCypher Go SDK. Test functions
//try to mirror file names where possible.
package gobcy

import (
	"log"
	"os"
	"testing"
)

var keys1, keys2 AddrKeychain
var txhash1, txhash2 string
var bcy API

func TestMain(m *testing.M) {
	//Set Coin/Chain to BlockCypher testnet
	bcy.Coin = "bcy"
	bcy.Chain = "test"
	//Set Token
	bcy.Token = "test-token"
	//Create/fund the test addresses
	var err error
	keys1, err = bcy.GenAddrKeychain()
	keys2, err = bcy.GenAddrKeychain()
	if err != nil {
		log.Fatal("Error generating test addresses: ", err)
	}
	txhash1, err = bcy.Faucet(keys1, 1e5)
	txhash2, err = bcy.Faucet(keys2, 2e5)
	if err != nil {
		log.Fatal("Error funding test addresses: ", err)
	}
	os.Exit(m.Run())
}

//TestsGetTXConf runs first, to test
//Confidence factor
func TestGetTXConf(t *testing.T) {
	conf, err := bcy.GetTXConf(txhash2)
	if err != nil {
		t.Error("Error encountered: ", err)
	}
	t.Logf("%+v\n", conf)
	return
}

func TestBlockchain(t *testing.T) {
	ch, err := bcy.GetChain()
	if err != nil {
		t.Error("GetChain error encountered: ", err)
	}
	t.Logf("%+v\n", ch)
	_, err = bcy.GetBlock(187621, "")
	if err != nil {
		t.Error("GetBlock via height error encountered: ", err)
	}
	bl, err := bcy.GetBlock(0, "0000ffeb0031885f2292475eac7f9c6f7bf5057e3b0017a09cd1994e71b431a4")
	if err != nil {
		t.Error("GetBlock via hash error encountered: ", err)
	}
	t.Logf("%+v\n", bl)
	_, err = bcy.GetBlock(187621, "0000ffeb0031885f2292475eac7f9c6f7bf5057e3b0017a09cd1994e71b431a4")
	if err == nil {
		t.Error("Expected error when querying both height and hash in GetBlock, did not receive one")
	}
	err = nil
	bl, err = bcy.GetBlockPage(0, "0000cb69e3c85ec1a4a17d8a66634c1cf136acc9dca9a5a71664a593f92bc46e", 0, 1)
	if err != nil {
		t.Error("GetBlockPage error encountered: ", err)
	}
	t.Logf("%+v\n", bl)
	bl2, err := bcy.GetBlockNextTXs(bl)
	if err != nil {
		t.Error("GetBlockNextTXs error encountered: ", err)
	}
	t.Logf("%+v\n", bl2)
	return
}

func TestAddress(t *testing.T) {
	addr, err := bcy.GetAddrBal(keys1.Address)
	if err != nil {
		t.Error("GetAddrBal error encountered: ", err)
	}
	t.Logf("%+v\n", addr)
	addr, err = bcy.GetAddr(keys1.Address)
	if err != nil {
		t.Error("GetAddr error encountered: ", err)
	}
	t.Logf("%+v\n", addr)
	addr, err = bcy.GetAddrFull(keys2.Address)
	if err != nil {
		t.Error("GetAddrFull error encountered: ", err)
	}
	t.Logf("%+v\n", addr)
	return
}

func TestGenAddrMultisig(t *testing.T) {
	pubkeys := []string{
		"02c716d071a76cbf0d29c29cacfec76e0ef8116b37389fb7a3e76d6d32cf59f4d3",
		"033ef4d5165637d99b673bcdbb7ead359cee6afd7aaf78d3da9d2392ee4102c8ea",
		"022b8934cc41e76cb4286b9f3ed57e2d27798395b04dd23711981a77dc216df8ca",
	}
	response, err := bcy.GenAddrMultisig(AddrKeychain{PubKeys: pubkeys, ScriptType: "multisig-2-of-3"})
	if err != nil {
		t.Error("Error encountered: ", err)
	}
	if response.Address != "De2gwq9GvNgvKgHCYRMKnPqss3pzWGSHiH" {
		t.Error("Response does not match expected address")
	}
	t.Logf("%+v\n", response)
	return
}

func TestWallet(t *testing.T) {
	wal, err := bcy.CreateWallet(Wallet{Name: "testwallet",
		Addresses: []string{keys1.Address}})
	if err != nil {
		t.Error("CreateWallet error encountered: ", err)
	}
	t.Logf("%+v\n", wal)
	list, err := bcy.ListWallets()
	if err != nil {
		t.Error("ListWallet error encountered: ", err)
	}
	t.Logf("%+v\n", list)
	wal, err = bcy.AddAddrWallet("testwallet", []string{keys2.Address}, false)
	if err != nil {
		t.Error("AddAddrWallet error encountered: ", err)
	}
	t.Logf("%+v\n", wal)
	err = bcy.DeleteAddrWallet("testwallet", []string{keys1.Address})
	if err != nil {
		t.Error("DeleteAddrWallet error encountered ", err)
	}
	addrs, err := bcy.GetAddrWallet("testwallet", false, false, false, false)
	if err != nil {
		t.Error("GetAddrWallet error encountered: ", err)
	}
	if addrs[0] != keys2.Address {
		t.Error("GetAddrWallet response does not match expected addresses")
	}
	wal, newAddrKeys, err := bcy.GenAddrWallet("testwallet")
	if err != nil {
		t.Error("GenAddrWallet error encountered: ", err)
	}
	t.Logf("%+v\n%+v\n", wal, newAddrKeys)
	err = bcy.DeleteWallet("testwallet")
	if err != nil {
		t.Error("DeleteWallet error encountered: ", err)
	}
	return
}

func TestHDWallet(t *testing.T) {
	wal, err := bcy.CreateHDWallet(HDWallet{Name: "testhdwallet",
		ExtPubKey: "xpub661MyMwAqRbcFtXgS5sYJABqqG9YLmC4Q1Rdap9gSE8NqtwybGhePY2gZ29ESFjqJoCu1Rupje8YtGqsefD265TMg7usUDFdp6W1EGMcet8"})
	if err != nil {
		t.Error("CreateHDWallet error encountered: ", err)
	}
	t.Logf("%+v\n", wal)
	//Will reenable once ListHDWallet fixed
	/* list, err := bcy.ListHDWallets()
	if err != nil {
		t.Error("ListHDWallet error encountered: ", err)
	}
	t.Logf("%+v\n", list)*/
	addrs, err := bcy.GetAddrHDWallet("testhdwallet", false, false, false, false)
	if err != nil {
		t.Error("GetAddrHDWallet error encountered: ", err)
	}
	t.Logf("%+v\n", addrs)
	newhd, err := bcy.DeriveAddrHDWallet("testhdwallet", 1, false, 0)
	if err != nil {
		t.Error("DeriveAddrHDWallet error encountered: ", err)
	}
	t.Logf("%+v\n", newhd)
	wal, err = bcy.GetHDWallet("testhdwallet")
	if err != nil {
		t.Error("GetHDWallet error encountered: ", err)
	}
	t.Logf("%+v\n", wal)
	err = bcy.DeleteHDWallet("testhdwallet")
	if err != nil {
		t.Error("DeleteHDWallet error encountered: ", err)
	}
	return
}

func TestTX(t *testing.T) {
	txs, err := bcy.GetUnTX()
	if err != nil {
		t.Error("GetUnTX error encountered: ", err)
	}
	t.Logf("%+v\n", txs)
	tx, err := bcy.GetTX(txhash1)
	if err != nil {
		t.Error("GetTX error encountered: ", err)
	}
	t.Logf("%+v\n", tx)
	//Create New TXSkeleton
	temp := TempNewTX(keys2.Address, keys1.Address, 45000)
	skel, err := bcy.NewTX(temp, true)
	if err != nil {
		t.Error("NewTX error encountered: ", err)
	}
	t.Logf("%+v\n", skel)
	//Sign TXSkeleton
	err = skel.Sign([]string{keys2.Private})
	if err != nil {
		t.Error("*TXSkel.Sign error encountered: ", err)
	}
	//Send TXSkeleton
	skel, err = bcy.SendTX(skel)
	if err != nil {
		t.Error("SendTX error encountered: ", err)
	}
	t.Logf("%+v\n", skel)
	return
}

func TestMicro(t *testing.T) {
	result, err := bcy.SendMicro(MicroTX{Priv: keys2.Private, ToAddr: keys1.Address, Value: 25000})
	if err != nil {
		t.Error("Error encountered: ", err)
	}
	t.Logf("%+v\n", result)
	//Test public key signing method
	micpub, err := bcy.SendMicro(MicroTX{Pubkey: keys2.Public, ToAddr: keys1.Address, Value: 15000})
	if err != nil {
		t.Error("Error encountered: ", err)
	}
	t.Logf("First MicroTX Pubkey call: %+v\n", micpub)
	//Sign resultpub
	err = micpub.Sign(keys2.Private)
	if err != nil {
		t.Error("Error encountered: ", err)
	}
	t.Logf("Signed MicroTX: %+v\n", micpub)
	//Send signed resultpub
	resultpub, err := bcy.SendMicro(micpub)
	if err != nil {
		t.Error("Error encountered: ", err)
	}
	t.Logf("Final MicroTX Pubkey call: %+v\n", resultpub)
	return
}

func TestHook(t *testing.T) {
	hook, err := bcy.CreateHook(Hook{Event: "new-block", URL: "https://my.domain.com/api/callbacks/doublespend?secret=justbetweenus"})
	if err != nil {
		t.Error("PostHook error encountered: ", err)
	}
	t.Logf("%+v\n", hook)
	if err = bcy.DeleteHook(hook.ID); err != nil {
		t.Error("DeleteHook error encountered: ", err)
	}
	hooks, err := bcy.ListHooks()
	if err != nil {
		t.Error("ListHooks error encountered: ", err)
	}
	//Should be empty
	t.Logf("%+v\n", hooks)
	return
}

func TestPayFwd(t *testing.T) {
	pay, err := bcy.CreatePayFwd(PayFwd{Destination: keys1.Address})
	if err != nil {
		t.Error("CreatePayFwd error encountered: ", err)
	}
	t.Logf("%+v\n", pay)
	pay, err := bcy.GetPayFwd(pay.ID)
	if err != nil {
		t.Error("GetPayFwd error encountered: ", err)
	}
	t.Logf("%+v\n", pay)
	if err = bcy.DeletePayFwd(pay.ID); err != nil {
		t.Error("DeletePayFwd error encountered: ", err)
	}
	pays, err := bcy.ListPayFwds()
	if err != nil {
		t.Error("ListPayFwds error encountered: ", err)
	}
	//Should be empty
	t.Logf("%+v\n", pays)
	return
}