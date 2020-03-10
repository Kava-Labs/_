package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/keys"
	crkeys "github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkrest "github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authrest "github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/kava-labs/kava/app"
	"github.com/tendermint/go-amino"
)

func init() {
	config := sdk.GetConfig()
	app.SetBech32AddressPrefixes(config)
	app.SetBip44CoinType(config)
	config.Seal()
}

func main() {
	proposalContent := gov.ContentFromProposalType("A Test Title", "A test description on this proposal.", gov.ProposalTypeText)
	addr, err := sdk.AccAddressFromBech32("kava1ffv7nhd3z6sych2qpqkk03ec6hzkmufy0r2s4c") // validator
	if err != nil {
		panic(err)
	}

	// create a message to send to the blockchain
	msg := gov.NewMsgSubmitProposal(
		proposalContent,
		sdk.NewCoins(sdk.NewInt64Coin("stake", 1000)),
		addr,
	)

	// create a vote on a proposal to send to the blockchain
	vote := gov.NewMsgVote(addr, uint64(0), types.OptionYes)
	fmt.Printf("\nvote:%s\n", vote)

	// creating a transaction
	//	tx := authtypes.NewStdTx([]sdk.Msg{msg}, authtypes.StdFee{}, []authtypes.StdSignature{}, "a test memo")

	// helper methods for transactions
	cdc := app.MakeCodec() // make codec for the app
	// transaction builder
	// create a keybase

	// TODO - IMPORTANT - this needs to be set manually and does NOT work with tilde i.e. ~/ does NOT work
	keybase, err := keys.NewKeyBaseFromDir("/Users/john/.kvcli/")
	if err != nil {
		panic(err)
	}
	_, err = keybase.List()
	// fmt.Printf("Keys: %s\n\n", all)
	if err != nil {
		panic(err)
	}

	// the test address
	address := "kava1ffv7nhd3z6sych2qpqkk03ec6hzkmufy0r2s4c"

	accountNumber, sequenceNumber := getAccountNumberAndSequenceNumber(cdc, address)

	// cast to the generic msg type
	msgToSend := []sdk.Msg{msg}
	keyname := "vlad"      // TODO - IMPORTANT this must match the keys in the startchain.sh script
	password := "password" // TODO - IMPORTANT this must match the keys in the startchain.sh script

	sendMsgToBlockchain(cdc, accountNumber, sequenceNumber, keyname, password, msgToSend, keybase)

}

// sendMsgToBlockchain sends a message to the blockchain via the rest api
func sendMsgToBlockchain(cdc *amino.Codec, accountNumber uint64, sequenceNumber uint64, keyname string,
	password string, msg []sdk.Msg, keybase crkeys.Keybase) {
	txBldr := auth.NewTxBuilderFromCLI().
		WithTxEncoder(authclient.GetTxEncoder(cdc)).WithChainID("testing").
		WithKeybase(keybase).WithAccountNumber(accountNumber).
		WithSequence(sequenceNumber)

		// build and sign the transaction
	// this is the *Amino* encoded version of the transaction
	// fmt.Printf("%+v", txBldr.Keybase())
	txBytes, err := txBldr.BuildAndSign("vlad", "password", msg)
	if err != nil {
		panic(err)
	}
	// fmt.Printf("txBytes: %s", txBytes)

	// need to convert the Amino encoded version back to an actual go struct
	var tx auth.StdTx
	cdc.UnmarshalBinaryLengthPrefixed(txBytes, &tx) // might be unmarshal binary bare

	// now we re-marshall it again into json
	jsonBytes, err := cdc.MarshalJSON(
		authrest.BroadcastReq{
			Tx:   tx,
			Mode: "block",
		},
	)
	if err != nil {
		panic(err)
	}
	// fmt.Println("post body: ", string(jsonBytes))

	resp, err := http.Post("http://localhost:1317/txs", "application/json", bytes.NewBuffer(jsonBytes))
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// fmt.Println(string(body))

}

// getAccountNumberAndSequenceNumber gets an account number and sequence number from the blockchain
func getAccountNumberAndSequenceNumber(cdc *amino.Codec, address string) (accountNumber uint64, sequenceNumber uint64) {

	// we need to setup the account number and sequence in order to have a valid transaction
	resp, err := http.Get("http://localhost:1317/auth/accounts/" + address)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var bodyUnmarshalled sdkrest.ResponseWithHeight
	err = cdc.UnmarshalJSON(body, &bodyUnmarshalled)
	if err != nil {
		panic(err)
	}

	var account authexported.Account
	err = cdc.UnmarshalJSON(bodyUnmarshalled.Result, &account)
	if err != nil {
		panic(err)
	}
	// fmt.Printf("\n\naccount: %s\n\n", account)

	return account.GetAccountNumber(), account.GetSequence()

}
