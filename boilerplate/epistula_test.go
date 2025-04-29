//go:build integration
// +build integration

package boilerplate_test

import (
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/stretchr/testify/assert"
	"github.com/subtrahend-labs/gobt/boilerplate"
)

func TestAdminUtilsModuleExtrinsics(t *testing.T) {
	t.Parallel()
	t.Run("VerifyEpistulaHeaders", func(t *testing.T) {
		t.Parallel()
		bob, _ := signature.KeyringPairFromSecret("//Bob", 42)
		alice, _ := signature.KeyringPairFromSecret("//Alice", 42)
		msg := []byte("Hello World")
		headers, _ := boilerplate.GetEpistulaHeaders(alice, bob.Address, msg)
		res := boilerplate.VerifyEpistulaHeaders(
			bob,
			headers["Epistula-Request-Signature"],
			msg,
			headers["Epistula-Timestamp"],
			headers["Epistula-Uuid"],
			headers["Epistula-Signed-For"],
			headers["Epistula-Signed-By"],
		)
		if res != nil {
			assert.FailNow(t, res.Error())
		}
		msg = []byte("Hello-World")
		res = boilerplate.VerifyEpistulaHeaders(
			bob,
			headers["Epistula-Request-Signature"],
			msg,
			headers["Epistula-Timestamp"],
			headers["Epistula-Uuid"],
			headers["Epistula-Signed-For"],
			headers["Epistula-Signed-By"],
		)
		if res == nil {
			assert.FailNow(t, "Signature should be mismatched, but returned true")
		}
	})
}
