package client

import "github.com/centrifuge/go-substrate-rpc-client/v4/signature"

type Option func(*Client)

func WithKeyring(keyring *signature.KeyringPair) Option {
	return func(c *Client) {
		c.keyring = keyring
	}
}

func WithNetwork(network uint16) Option {
	return func(c *Client) {
		c.network = network
	}
}
