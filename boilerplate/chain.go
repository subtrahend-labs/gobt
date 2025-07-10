package boilerplate

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/subtrahend-labs/gobt/client"
)

type BaseChainSubscriber struct {
	callbacks                   []func(types.Header)
	mainfunc                    func(i <-chan (bool), o chan<- (bool))
	onSubscriptionError         func(err error)
	onSubscriptionCreationError func(err error)
	startup                     func()
	NetUID                      int
}

func (b *BaseChainSubscriber) AddBlockCallback(f func(types.Header)) {
	b.callbacks = append(b.callbacks, f)
}

func (b *BaseChainSubscriber) SetStartupFunc(f func()) {
	b.startup = f
}

func (b *BaseChainSubscriber) SetMainFunc(f func(i <-chan (bool), o chan<- (bool))) {
	b.mainfunc = f
}

func (b *BaseChainSubscriber) SetOnSubscriptionError(f func(e error)) {
	b.onSubscriptionError = f
}

func (b *BaseChainSubscriber) SetOnSubscriptionCreationError(f func(e error)) {
	b.onSubscriptionCreationError = f
}

func NewChainSubscriber(n int) *BaseChainSubscriber {
	return &BaseChainSubscriber{NetUID: n}
}

func (b *BaseChainSubscriber) Start(c *client.Client) {
	// Handle graceful exits
	sigChan := make(chan os.Signal, 1)
	startCleanup := make(chan bool, 1)
	exitReady := make(chan bool, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go b.mainfunc(startCleanup, exitReady)
	for {
		sub, err := c.Api.RPC.Chain.SubscribeFinalizedHeads()
		if err != nil {
			b.onSubscriptionCreationError(err)
		}
		for {
			select {
			// Exit cleanly after system finishes current block
			case <-sigChan:
				// Send done signal
				startCleanup <- true
				// Wait for mainfunc to respond that it is done
				<-exitReady
				os.Exit(0)
			case head := <-sub.Chan():
				for _, exec := range b.callbacks {
					exec(head)
				}
			case err = <-sub.Err():
				b.onSubscriptionError(err)
				sub, err = c.Api.RPC.Chain.SubscribeFinalizedHeads()
				if err != nil {
					b.onSubscriptionCreationError(err)
				}
			}
		}
	}
}
