package boilerplate

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	sr25519 "github.com/ChainSafe/go-schnorrkel"
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/google/uuid"
	"github.com/vedhavyas/go-subkey/v2"
)

func sha256Hash(str []byte) string {
	h := sha256.New()
	h.Write(str)
	sum := h.Sum(nil)
	return hex.EncodeToString(sum)
}

// Takes sender ss58, sender public, sender private, receiver ss58 and body
func GetEpistulaHeaders(kp signature.KeyringPair, rSS58 string, body []byte) (map[string]string, error) {
	timestamp := time.Now().UnixMilli()
	uuid := uuid.New().String()
	timestampInterval := int64(math.Ceil(float64(timestamp) / 1e4))
	bodyHash := sha256Hash(body)
	message := fmt.Sprintf("%s.%s.%d.%s", bodyHash, uuid, timestamp, rSS58)
	requestSignature, err := signature.Sign([]byte(message), kp.URI)
	if err != nil {
		return nil, err
	}

	s1, _ := signature.Sign(fmt.Appendf([]byte{}, "%d.%s", timestampInterval-1, kp.Address), kp.URI)
	s2, _ := signature.Sign(fmt.Appendf([]byte{}, "%d.%s", timestampInterval, kp.Address), kp.URI)
	s3, _ := signature.Sign(fmt.Appendf([]byte{}, "%d.%s", timestampInterval+1, kp.Address), kp.URI)

	headers := map[string]string{
		"Epistula-Version":            "2",
		"Epistula-Timestamp":          fmt.Sprintf("%d", timestamp),
		"Epistula-Uuid":               uuid,
		"Epistula-Signed-By":          kp.Address,
		"Epistula-Signed-For":         rSS58,
		"Epistula-Request-Signature":  types.NewSignature(requestSignature).Hex(),
		"Epistula-Secret-Signature-0": types.NewSignature(s1).Hex(),
		"Epistula-Secret-Signature-1": types.NewSignature(s2).Hex(),
		"Epistula-Secret-Signature-2": types.NewSignature(s3).Hex(),
		"Content-Type":                "application/json",
		"Connection":                  "keep-alive",
	}

	return headers, nil
}

func ss58ToPublicKey(address string) ([]byte, error) {
	_, pubKey, err := subkey.SS58Decode(address)
	if err != nil {
		return nil, fmt.Errorf("decode address error: %w", err)
	}

	return pubKey, nil
}

func VerfySignature(signed_by string, msg []byte, sig []byte) bool {
	signat := new(sr25519.Signature)
	if err := signat.Decode([64]byte(sig)); err != nil {
		return false
	}
	pubk, err := ss58ToPublicKey(signed_by)
	if err != nil {
		fmt.Println(err)
		return false
	}
	pk, err := sr25519.NewPublicKey([32]byte(pubk))
	if err != nil {
		fmt.Println(err)
		return false
	}

	ok, err := pk.Verify(signat, sr25519.NewSigningContext([]byte("substrate"), msg))
	if err != nil || !ok {
		return false
	}

	return true
}

func VerifyEpistulaHeaders(self_ss58 string, sig string, body []byte, timestamp, uuid, signed_for, signed_by string) error {
	i, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		panic(err)
	}
	tm := time.Unix(i, 0)
	ALLOWED_DELTA := 8000 * time.Millisecond
	if time.Since(tm) > ALLOWED_DELTA {
		return errors.New("reques is too old")
	}
	bodyHash := sha256Hash(body)

	sig, ok := strings.CutPrefix(sig, "0x")
	if !ok {
		return errors.New("sig is not hex string")
	}
	bytes, err := hex.DecodeString(sig)
	if err != nil {
		return errors.New("failed decoding hex string")
	}
	message := fmt.Sprintf("%s.%s.%s.%s", bodyHash, uuid, timestamp, self_ss58)

	ok = VerfySignature(signed_by, []byte(message), bytes)
	if !ok {
		return errors.New("signature mismatch")
	}
	return nil
}
