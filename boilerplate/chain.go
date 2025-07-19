package boilerplate

import (
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/subtrahend-labs/gobt/client"
)

type BaseChainSubscriber struct {
	callbacks           []func(types.Header)
	onSubscriptionError func(err error)
	stopChan            chan bool
	restartChan         chan bool
}

func (b *BaseChainSubscriber) AddBlockCallback(f func(types.Header)) {
	b.callbacks = append(b.callbacks, f)
}

func (b *BaseChainSubscriber) SetOnSubscriptionError(f func(e error)) {
	b.onSubscriptionError = f
}

func NewChainSubscriber() *BaseChainSubscriber {
	return &BaseChainSubscriber{stopChan: make(chan bool, 1)}
}

func (b *BaseChainSubscriber) Stop() {
	close(b.stopChan)
}

func (b *BaseChainSubscriber) Restart() {
	b.stopChan <- true
}

func (b *BaseChainSubscriber) Start(c *client.Client) error {
	for {
		sub, err := c.Api.RPC.Chain.SubscribeFinalizedHeads()
		if err != nil {
			return err
		}
		for {
			select {
			case <-b.restartChan:
				sub.Unsubscribe()
				sub, err = c.Api.RPC.Chain.SubscribeFinalizedHeads()
				if err != nil {
					return err
				}
			case <-b.stopChan:
				return nil
			case head := <-sub.Chan():
				for _, exec := range b.callbacks {
					exec(head)
				}
			case err = <-sub.Err():
				b.onSubscriptionError(err)
				sub.Unsubscribe()
				sub, err = c.Api.RPC.Chain.SubscribeFinalizedHeads()
				if err != nil {
					return err
				}
			}
		}
	}
}
