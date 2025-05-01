package verifier

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"os"
	"testing"

	"github.com/stretchr/testify/assert")

func loadCertificate() *x509.Certificate{
	var f = "apko_0.20.1_linux_amd64.tar.gz.crt"
	r, _ := os.ReadFile(f)
	block, _ := pem.Decode(r)
	cert, _ := x509.ParseCertificate(block.Bytes)

	return cert
}

func loadCorrectSignature() *Signature {
	var f = "apko_0.20.1_linux_amd64.tar.gz.sig"
	bytes, _ := os.ReadFile(f)
	var s = Signature(bytes)
	return &s
}

func loadWrongSignature() *Signature {
	base64Encoded := "MEUCIAVLmr6q4iAilBhM2n+/IEC07TKzpGLVMTTl1kjTBrBMAiEAtcDe7ZSl4zwxvJjMLUkOoS4Hx9xSxdSktwMmomD8x6A="
	sig, _ := base64.StdEncoding.DecodeString(base64Encoded)
	var s = Signature(sig)
	return &s
}

func TestVerify(t *testing.T) {
	var success = VerificationResult(true)
	var failure = VerificationResult(false)

	validSha := "442d8baafc0c3a873b21a3add32f5c65f538fb5cbcf4a4a69ba098a2b730c5d2"
	invalidSha := "02907c168a1e9a440743efee9e70cce3f3fb04c5d5394b06fd924ba7c416c661"

	tests := []struct {
		cert    *x509.Certificate
		sig   *Signature
		sha string
		verificationResult *VerificationResult
		error bool
		description string
	}{
		{ cert: loadCertificate(), sig: loadCorrectSignature(), sha: validSha, verificationResult: &success, description: "Sha singed with the correct PrivateKey should pass verification with the PublicKey" },
		{ cert: loadCertificate(), sig: loadWrongSignature(), sha: validSha, verificationResult: &failure, description: "Invalid Signature should not pass verification" },
		{ cert: loadCertificate(), sig: loadCorrectSignature(), sha: invalidSha, verificationResult: &failure, description: "Invalid Sha should not pass verification" },	
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			v := CertificateVerifier {
				cert: tt.cert,
				sig: tt.sig,
				sha: tt.sha,
			}
			verificationResult, err := v.do()

			if (!assert.EqualValues(t, tt.verificationResult, verificationResult)){
				t.Fail()
			}

			isError := err!=nil

			if isError != tt.error {
				t.Fail()
			}

			if isError && bool(*verificationResult) {
				t.Fail()
			}

		})
	}
}

func TestVerifierIntegration(t *testing.T) {
	r := VerifyFile("apko_0.20.1_linux_amd64.tar.gz")
	assert.True(t, r)
}