package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	rpcclienthttp "github.com/tendermint/tendermint/rpc/client/http"
)

func GetAccount(ctx context.Context, client *rpcclienthttp.HTTP, cdc codec.Codec) {
	request := authtypes.QueryAccountRequest{
		Address: "noble1fxyp2x52a7sg879skvlyft62sj7xqhjzjzntk2",
	}

	requestBytes, err := cdc.Marshal(&request)
	if err != nil {
		log.Fatal(err)
	}

	abciResponse, err := client.ABCIQueryWithOptions(
		context.Background(),
		"/cosmos.auth.v1beta1.Query/Account",
		requestBytes,
		rpcclient.ABCIQueryOptions{},
	)
	if err != nil {
		log.Fatal(err)
	}

	if abciResponse.Response.Code != 0 {
		log.Fatal(fmt.Errorf(
			"account query failed: %s, code: %d, log: %s",
			abciResponse.Response.Codespace,
			abciResponse.Response.Code,
			abciResponse.Response.Log,
		))
	}

	response := authtypes.QueryAccountResponse{}
	if err := cdc.Unmarshal(abciResponse.Response.Value, &response); err != nil {
		log.Fatal(err)
	}

	var account authtypes.AccountI
	if err := cdc.UnpackAny(response.Account, &account); err != nil {
		log.Fatal(err)
	}

	fmt.Println(account.GetAddress().String())
	fmt.Println(hex.EncodeToString(account.GetPubKey().Bytes()))
	fmt.Println(account.GetAccountNumber())
	fmt.Println(account.GetSequence())
}

func GetTx(ctx context.Context, client *rpcclienthttp.HTTP, cdc codec.Codec, txHash string) {
	txHashBytes, err := hex.DecodeString(txHash)
	if err != nil {
		log.Fatal(err)
	}

	result, err := client.Tx(ctx, txHashBytes, false)
	if err != nil {
		log.Fatal(err)
	} else if result.TxResult.Code != 0 {
		log.Fatal(result.TxResult.Log)
	}

	txResult, err := cdc.MarshalJSON(&result.TxResult)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(txResult))
}

func main() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount("noble", "noblepub")
	config.SetBech32PrefixForValidator("nobleval", "noblevalpub")
	config.SetBech32PrefixForConsensusNode("noblevalcons", "noblevalconspub")
	config.SetPurpose(44)
	config.SetCoinType(118)
	config.Seal()

	// client, err := rpcclienthttp.New("https://noble-testnet-rpc.polkachu.com:443", "/websocket")
	client, err := rpcclienthttp.New("https://noble-rpc.lavenderfive.com:443", "/websocket")
	if err != nil {
		log.Fatal(err)
	}

	registry := codectypes.NewInterfaceRegistry()
	std.RegisterInterfaces(registry)
	authtypes.RegisterInterfaces(registry)
	cdc := codec.NewProtoCodec(registry)

	// GetAccount(context.Background(), client, cdc)
	GetTx(context.Background(), client, cdc, "8AE9CC485850C0B310A80015CE153D43CAB2D20F2169AB9494EAECC4F925010D")
}
