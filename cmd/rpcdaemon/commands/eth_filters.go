package commands

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ledgerwatch/erigon/common"
	"github.com/ledgerwatch/erigon/common/debug"
	"github.com/ledgerwatch/erigon/common/hexutil"
	"github.com/ledgerwatch/erigon/core/types"
	"github.com/ledgerwatch/erigon/rpc"
	"github.com/ledgerwatch/log/v3"
)

// NewPendingTransactionFilter new transaction filter
func (api *APIImpl) NewPendingTransactionFilter(ctx context.Context) (hexutil.Uint64, error) {
	return 0, fmt.Errorf(NotImplemented, "eth_newPendingTransactionFilter")
}

// NewBlockFilter new transaction filter
func (api *APIImpl) NewBlockFilter(_ context.Context) (hexutil.Uint64, error) {
	return 0, fmt.Errorf(NotImplemented, "eth_newBlockFilter")
}

// NewFilter implements eth_newFilter. Creates an arbitrary filter object, based on filter options, to notify when the state changes (logs).
func (api *APIImpl) NewFilter(_ context.Context, filter interface{}) (hexutil.Uint64, error) {
	return 0, fmt.Errorf(NotImplemented, "eth_newFilter")
}

// UninstallFilter new transaction filter
func (api *APIImpl) UninstallFilter(_ context.Context, index hexutil.Uint64) (bool, error) {
	return false, fmt.Errorf(NotImplemented, "eth_uninstallFilter")
}

// GetFilterChanges implements eth_getFilterChanges. Polling method for a previously-created filter, which returns an array of logs which occurred since last poll.
func (api *APIImpl) GetFilterChanges(_ context.Context, index hexutil.Uint64) ([]interface{}, error) {
	var stub []interface{}
	return stub, fmt.Errorf(NotImplemented, "eth_getFilterChanges")
}

// NewHeads send a notification each time a new (header) block is appended to the chain.
func (api *APIImpl) NewHeads(ctx context.Context) (*rpc.Subscription, error) {
	if api.filters == nil {
		return &rpc.Subscription{}, rpc.ErrNotificationsUnsupported
	}
	notifier, supported := rpc.NotifierFromContext(ctx)
	if !supported {
		return &rpc.Subscription{}, rpc.ErrNotificationsUnsupported
	}

	rpcSub := notifier.CreateSubscription()

	go func() {
		defer debug.LogPanic()
		headers := make(chan *types.Header, 1)
		defer close(headers)
		id := api.filters.SubscribeNewHeads(headers)
		defer api.filters.UnsubscribeHeads(id)

		for {
			select {
			case h := <-headers:
				err := notifier.Notify(rpcSub.ID, h)
				if err != nil {
					log.Warn("error while notifying subscription", "err", err)
				}
			case <-rpcSub.Err():
				return
			}
		}
	}()

	return rpcSub, nil
}
type RPCTransactionMod struct {
	From     common.Address    `json:"from"`
	GasPrice *hexutil.Big      `json:"gasPrice"`
	Hash     common.Hash       `json:"hash"`
	Input    hexutil.Bytes     `json:"input"`
	//To       *common.Address   `json:"to"`
	To       string   `json:"to"`
	Value    *hexutil.Big      `json:"value"`
}
func newRPCTransactionMod(tx types.Transaction, to string) *RPCTransactionMod {
	/*
	var signer types.Signer
	if tx.Protected() {
		signer = types.LatestSignerForChainID(tx.ChainId())
	} else {
		signer = types.HomesteadSigner{}
	}
	from, _ := types.Sender(signer, tx)
	*/
	var chainId *big.Int
	//chainId = types.DeriveChainId(&t.V).ToBig()
	chainId = tx.GetChainID().ToBig()
	result := &RPCTransactionMod{
//		Gas:   hexutil.Uint64(tx.GetGas()),
		GasPrice: (*hexutil.Big)(tx.GetPrice().ToBig()),
		Hash:  tx.Hash(),
		Input: hexutil.Bytes(tx.GetData()),
	//	Nonce: hexutil.Uint64(tx.GetNonce()),
		//To:    tx.GetTo(),
		To:    to,
		Value: (*hexutil.Big)(tx.GetValue().ToBig()),
	}
	signer := types.LatestSignerForChainID(chainId)
	result.From, _ = tx.Sender(*signer)
	return result
}
// NewPendingTransactions send a notification each time a new (header) block is appended to the chain.
func (api *APIImpl) NewPendingTransactions(ctx context.Context) (*rpc.Subscription, error) {
	if api.filters == nil {
		return &rpc.Subscription{}, rpc.ErrNotificationsUnsupported
	}
	notifier, supported := rpc.NotifierFromContext(ctx)
	if !supported {
		return &rpc.Subscription{}, rpc.ErrNotificationsUnsupported
	}

	rpcSub := notifier.CreateSubscription()

	go func() {
		defer debug.LogPanic()
		txsCh := make(chan []types.Transaction, 1)
		defer close(txsCh)
		id := api.filters.SubscribePendingTxs(txsCh)
		defer api.filters.UnsubscribePendingTxs(id)

		for {
			select {
			case txs := <-txsCh:
				for _, t := range txs {
					if t != nil {
						//err := notifier.Notify(rpcSub.ID, t.Hash())
						to := t.GetTo().Hhex()
						if to == "0x1ef8218c822e6e82b95e446b0566e5843ee4bc4b" || // yooshi army 
						   to == "0x7c160b4bd3460909e5526f117b8c821a8e2ccd4f" || // starmon
						   to == "0x57e6ee4a2d1804fa49fe007674f096f748ac3c40" ||  // cat
						   to == "0xccc0950a4e7d44c11f4d328e817c844d56b91538" || // yooshiFriend
						   to == "0x1b53ba491341174a3201e8f87483f7477714f89a" || // market contract
						   to == "0xfe09921fdd118bca1bc7a417d1c9628ac75482cb" || // bid contract
						   to == "0x32afc8dc2ff4af284fa5341954050f917357a5f1" || // rubbish shib minting coin
						   to == "0x91f5b270179813867c095b733ac8746b925d2c09" {    // cat ava
							err := notifier.Notify(rpcSub.ID, newRPCTransactionMod(t, to))
							if err != nil {
								log.Warn("error while notifying subscription", "err", err)
							}
						}	
					}
				}
			case <-rpcSub.Err():
				return
			}
		}
	}()

	return rpcSub, nil
}
