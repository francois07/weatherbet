package script

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"github.com/tidwall/gjson"
	"internal/stack"
	"io"
	"log"
	"net/http"
	"regexp"
)

func ParseScript(input string) ([]string, []string) {
	matchTokenDelimiter := regexp.MustCompile(`\s+|\n+]`)
	inputTokens := matchTokenDelimiter.Split(input, -1)
	var args []string
	var instructions []string

	matchArgsDelimiter := regexp.MustCompile(`^-+$`)
	for idx, token := range inputTokens {
		if matchArgsDelimiter.MatchString(token) {
			args = inputTokens[:idx]
			instructions = inputTokens[idx+1:]
			break
		}
	}

	return args, instructions
}

func EvalScript(script string, args map[string]string) bool {
	_, instructions := ParseScript(script)
	s := stack.Stack{}

	for _, instruction := range instructions {
		switch instruction {
		case "OPDup":
			s.Push(args[s.Pop()])
		case "OPHash":
			s.Push(OPHash(s.Pop()))
		case "OPEqualVerify":
			if !OPEqualVerify(s.Pop(), s.Pop()) {
				return false
			}
		case "OPCheckThirdParty":
			if !OPCheckThirdParty(s.Pop(), s.Pop(), s.Pop()) {
				return false
			}
		case "OPCheckSig":
			if !OPCheckSig(s.Pop(), s.Pop(), s.Pop()) {
				return false
			}
		default:
			s.Push(instruction)
		}
	}

	return true
}

func OPHash(input string) string {
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}

func OPEqualVerify(input1, input2 string) bool {
	return input1 == input2
}

func OPCheckSig(pubKey string, hash string, signature string) bool {
	pDec, errDec := base64.StdEncoding.DecodeString(pubKey)
	if errDec != nil {
		log.Fatal(errDec)
	}

	pemKey, _ := pem.Decode(pDec)

	rsaKey, errParse := x509.ParsePKIXPublicKey(pemKey.Bytes)
	if errParse != nil {
		log.Fatal(errParse)
	}

	sDec, errDecSig := base64.StdEncoding.DecodeString(signature)
	if errDecSig != nil {
		log.Fatal(errDecSig)
	}

	hexHash, errDecode := hex.DecodeString(hash)
	if errDecode != nil {
		return false
	}
	err := rsa.VerifyPKCS1v15(rsaKey.(*rsa.PublicKey), crypto.SHA256, hexHash, sDec)
	if err != nil {
		return false
	}
	return true
}

func OPCheckThirdParty(url string, finalField string, expectedValue string) bool {
	res, err := http.Get(url)
	if err != nil {
		return false
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return false
	}
	r := gjson.Get(string(body), finalField).String()

	return r == expectedValue
}
