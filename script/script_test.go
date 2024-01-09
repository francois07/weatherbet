package script_test

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"script"
	"testing"
)

func TestOPHashEqualSucces(t *testing.T) {
	message := "test"
	hash := script.OPHash(message)
	s := "message --- message OPDup OPHash " + hash + " OPEqualVerify"
	evalArgs := map[string]string{"message": message}
	evalRes := script.EvalScript(s, evalArgs)

	if !evalRes {
		t.Fatalf(`Got %v, expected true`, evalRes)
	}
}

func TestOPCheckSigSuccess(t *testing.T) {
	priKey := `-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQCQh/7uafm1IDob
nufvUMRoQEVAnmi0GIKheISbQUgqE++CryXihEjhE3rFYtdD+lDiLd2rxWPaT6XQ
ZTjzz0BB4hvcR2h77/DeUDeytNPiwy5Y/XmM9RKOKqKXN2QUYZ2QzEVI69baVDp3
+UknVHdL2RqV3hfE9sCW/E99r5i23d+F1O+mENeAykVzSX2hoaeimQVwXQQqGWTz
50pVsso4+uv4r/348osbD231fP1r6VXsC0uqRHrl49VD06a+aP+s4nrBh8tBl1e8
m7v2Pn96yWUcNYZchaHHQwqHS8h4ONsrDkd9g5cQL60lIqECaoKbyb9uoNUi+8cj
7nR0DXSJAgMBAAECggEAOkWKhCdYFS3MH8y+qF9BAONAyJ/ViR8EtAN6O3pjlQt6
mo7xUFWTaKPI8QI71l169OYXJKUt8BwCt3XIM4yQ1L9MClEEEEwaKhffjP6ZGykk
a56uviUk+Qq5iQw9HnaI55NkL9VeR6KD/FEWrHPCsWLz9A6aOMBfv8A0cHE2uis5
SeBP4ea2GWUgv0R2+koQIIy6eWBhI9LI92HF7aDTc3uiCEeeSQPMz/ceWnQdcZsV
XLJHTUnibjDv/Oi7JWr2hgA0WjltnbWgbsYx98IcaG78ENUhHsl0T1aFg7RcGAdU
ZSdlXSSlY3FPB38edrFAFiatheYQ+IU6DxTPZ4SczQKBgQDEi7V9l7qd8HRssJTd
y2RT74KuKtOhjaOvFtrzE/GypEJ890sKBKd8ENr5jDEwspj8QBGUQMtvSjHfnIJO
xdT59HmfJmW8KqBX71aOy3bBFbQFB6xfXDd6eaN/myMwOruh9A4oZooY8RHiooFJ
B3MiCucywtmVHqvtZN1z4y94RwKBgQC8QFbnrq+Xgy+Q4yGKWaAT/dXua92rnzjF
hzhVXRt8PeQawT2D4YQYx3OcdUk0RJTUQNBqDaSlfYf8tVSAZlki0zXw8hb3lYEe
FZROHkWHPhLcdCFtaxE7Kq2mjU8gv2/yM/NbTHl7kBtfwRc77JdgLyi3i/6GPklB
uvgs7S3krwKBgQCEokstn/C9mIDYwCkVq6XexqiHZDtAsFafV2sV0oWuqg58Sl2H
OzDTFoTPFn4zgLKgt5OlWjxus8EIR5PgGLzqmMJiVgUdgB6IeOkOn9tZ3Y2IP29h
QtflfKSK/mQ2rcvlNM9BEEFtJ3GMYWGhqLdXZ8gxhzBR40refsy64bstDQKBgDFE
pXn9PfdpXgmNaDnNOxgAVv0PPfSsty77NMMimw7pI8ncyTy6yNezW46XI5GKYWkr
jWA0MeMd93kr+/Ge17VFkdh9g4VIm4JEI4xOX+QFWupXemgonVne0ZPFZ/AqKiI5
dndujFzKWl+1KV+FjBigPwfKm9KGeqW5STp42IoBAoGBAMJ3CQGGbX3E57cLccb5
F9iLQgDaGMHIG/sRvFqJY1DMy1KM+de5gdfH3HnQXfy4fqDWzLkP38oKCzVPJaaf
myWweOr0UyDTI/M/YqrEqY9r5UBqPMxPjF+UV1u05I9MkD3jspKCUyK39ADHNsYS
Z3Cto9zj48IhwcytZO4i9P9A
-----END PRIVATE KEY-----`

	pubKey := `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAkIf+7mn5tSA6G57n71DE
aEBFQJ5otBiCoXiEm0FIKhPvgq8l4oRI4RN6xWLXQ/pQ4i3dq8Vj2k+l0GU4889A
QeIb3Edoe+/w3lA3srTT4sMuWP15jPUSjiqilzdkFGGdkMxFSOvW2lQ6d/lJJ1R3
S9kald4XxPbAlvxPfa+Ytt3fhdTvphDXgMpFc0l9oaGnopkFcF0EKhlk8+dKVbLK
OPrr+K/9+PKLGw9t9Xz9a+lV7AtLqkR65ePVQ9Omvmj/rOJ6wYfLQZdXvJu79j5/
esllHDWGXIWhx0MKh0vIeDjbKw5HfYOXEC+tJSKhAmqCm8m/bqDVIvvHI+50dA10
iQIDAQAB
-----END PUBLIC KEY-----`

	pubKey64 := base64.StdEncoding.EncodeToString([]byte(pubKey))
	pubKeyHash := sha256.Sum256([]byte(pubKey64))
	pubKeyHashHex := hex.EncodeToString(pubKeyHash[:])
	pemPriKey, _ := pem.Decode([]byte(priKey))

	rsaKey, _ := x509.ParsePKCS8PrivateKey(pemPriKey.Bytes)
	signed, _ := rsa.SignPKCS1v15(nil, rsaKey.(*rsa.PrivateKey), crypto.SHA256, pubKeyHash[:])
	signed64 := base64.StdEncoding.EncodeToString(signed)

	ok := script.OPCheckSig(pubKey64, pubKeyHashHex, signed64)

	if !ok {
		t.Fatalf("Got %v, expected true", ok)
	}
}

func TestOPCheckSigFailure(t *testing.T) {
	priKey := `-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQCQh/7uafm1IDob
nufvUMRoQEVAnmi0GIKheISbQUgqE++CryXihEjhE3rFYtdD+lDiLd2rxWPaT6XQ
ZTjzz0BB4hvcR2h77/DeUDeytNPiwy5Y/XmM9RKOKqKXN2QUYZ2QzEVI69baVDp3
+UknVHdL2RqV3hfE9sCW/E99r5i23d+F1O+mENeAykVzSX2hoaeimQVwXQQqGWTz
50pVsso4+uv4r/348osbD231fP1r6VXsC0uqRHrl49VD06a+aP+s4nrBh8tBl1e8
m7v2Pn96yWUcNYZchaHHQwqHS8h4ONsrDkd9g5cQL60lIqECaoKbyb9uoNUi+8cj
7nR0DXSJAgMBAAECggEAOkWKhCdYFS3MH8y+qF9BAONAyJ/ViR8EtAN6O3pjlQt6
mo7xUFWTaKPI8QI71l169OYXJKUt8BwCt3XIM4yQ1L9MClEEEEwaKhffjP6ZGykk
a56uviUk+Qq5iQw9HnaI55NkL9VeR6KD/FEWrHPCsWLz9A6aOMBfv8A0cHE2uis5
SeBP4ea2GWUgv0R2+koQIIy6eWBhI9LI92HF7aDTc3uiCEeeSQPMz/ceWnQdcZsV
XLJHTUnibjDv/Oi7JWr2hgA0WjltnbWgbsYx98IcaG78ENUhHsl0T1aFg7RcGAdU
ZSdlXSSlY3FPB38edrFAFiatheYQ+IU6DxTPZ4SczQKBgQDEi7V9l7qd8HRssJTd
y2RT74KuKtOhjaOvFtrzE/GypEJ890sKBKd8ENr5jDEwspj8QBGUQMtvSjHfnIJO
xdT59HmfJmW8KqBX71aOy3bBFbQFB6xfXDd6eaN/myMwOruh9A4oZooY8RHiooFJ
B3MiCucywtmVHqvtZN1z4y94RwKBgQC8QFbnrq+Xgy+Q4yGKWaAT/dXua92rnzjF
hzhVXRt8PeQawT2D4YQYx3OcdUk0RJTUQNBqDaSlfYf8tVSAZlki0zXw8hb3lYEe
FZROHkWHPhLcdCFtaxE7Kq2mjU8gv2/yM/NbTHl7kBtfwRc77JdgLyi3i/6GPklB
uvgs7S3krwKBgQCEokstn/C9mIDYwCkVq6XexqiHZDtAsFafV2sV0oWuqg58Sl2H
OzDTFoTPFn4zgLKgt5OlWjxus8EIR5PgGLzqmMJiVgUdgB6IeOkOn9tZ3Y2IP29h
QtflfKSK/mQ2rcvlNM9BEEFtJ3GMYWGhqLdXZ8gxhzBR40refsy64bstDQKBgDFE
pXn9PfdpXgmNaDnNOxgAVv0PPfSsty77NMMimw7pI8ncyTy6yNezW46XI5GKYWkr
jWA0MeMd93kr+/Ge17VFkdh9g4VIm4JEI4xOX+QFWupXemgonVne0ZPFZ/AqKiI5
dndujFzKWl+1KV+FjBigPwfKm9KGeqW5STp42IoBAoGBAMJ3CQGGbX3E57cLccb5
F9iLQgDaGMHIG/sRvFqJY1DMy1KM+de5gdfH3HnQXfy4fqDWzLkP38oKCzVPJaaf
myWweOr0UyDTI/M/YqrEqY9r5UBqPMxPjF+UV1u05I9MkD3jspKCUyK39ADHNsYS
Z3Cto9zj48IhwcytZO4i9P9A
-----END PRIVATE KEY-----`

	pubKey := `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAkIf+7mn5tSA6G57n71DE
aEBFQJ5otBiCoXiEm0FIKhPvgq8l4oRI4RN6xWLXQ/pQ4i3dq8Vj2k+l0GU4889A
QeIb3Edoe+/w3lA3srTT4sMuWP15jPUSjiqilzdkFGGdkMxFSOvW2lQ6d/lJJ1R3
S9kald4XxPbAlvxPfa+Ytt3fhdTvphDXgMpFc0l9oaGnopkFcF0EKhlk8+dKVbLK
OPrr+K/9+PKLGw9t9Xz9a+lV7AtLqkR65ePVQ9Omvmj/rOJ6wYfLQZdXvJu79j5/
esllHDWGXIWhx0MKh0vIeDjbKw5HfYOXEC+tJSKhAmqCm8m/bqDVIvvHI+50dA10
iQIDAQAB
-----END PUBLIC KEY-----`

	pubKey64 := base64.StdEncoding.EncodeToString([]byte(pubKey))
	pubKeyHash := sha256.Sum256([]byte(pubKey64))
	pemPriKey, _ := pem.Decode([]byte(priKey))

	rsaKey, _ := x509.ParsePKCS8PrivateKey(pemPriKey.Bytes)
	signed, _ := rsa.SignPKCS1v15(nil, rsaKey.(*rsa.PrivateKey), crypto.SHA256, pubKeyHash[:])
	signed64 := base64.StdEncoding.EncodeToString(signed)

	ok := script.OPCheckSig(pubKey64, "failed", signed64)

	if ok {
		t.Fatalf("Got %v, expected false", ok)
	}
}
