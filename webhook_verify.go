package conformvault

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// VerifyWebhookSignature verifies that a webhook payload was signed by ConformVault.
// The sigHeader is the value of the X-ConformVault-Signature header.
// The secret is the webhook signing secret returned when registering the endpoint.
func VerifyWebhookSignature(payload []byte, sigHeader string, secret string) bool {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	expectedSig := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(sigHeader), []byte(expectedSig))
}
