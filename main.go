package main

import (
	"context"
	"fmt"
	"log"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	rpcclienthttp "github.com/tendermint/tendermint/rpc/client/http"
)

func main() {
	client, err := rpcclienthttp.New("https://noble-testnet-rpc.polkachu.com:443", "/websocket")
	if err != nil {
		log.Fatal(err)
	}

	registry := codectypes.NewInterfaceRegistry()
	authtypes.RegisterInterfaces(registry)
	cryptocodec.RegisterInterfaces(registry)
	cdc := codec.NewProtoCodec(registry)

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

	response.Account.ProtoMessage()

	account := authtypes.BaseAccount{}
	if err := cdc.UnpackAny(response.Account, &account); err != nil {
		log.Fatal(err)
	}
}
