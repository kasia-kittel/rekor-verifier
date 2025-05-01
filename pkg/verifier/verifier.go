package verifier

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/hex"
	"errors"
	"os"

	"github.com/kasia-kittel/rekor-verifier/internal/utils"
	"github.com/kasia-kittel/rekor-verifier/pkg/log"
)

type Verifier interface {
    do() (*VerificationResult, error)
}

type VerificationResult bool

type CertificateVerifier struct {
	cert *x509.Certificate
	sig *Signature
	sha string
}

func (v CertificateVerifier) do() (*VerificationResult, error){
	
	verificationResult := VerificationResult(false)

	switch v.cert.PublicKeyAlgorithm {
	// case x509.RSA :
	// case x509.DSA:
	// case x509.Ed25519:
	case x509.ECDSA:
		publicKey := v.cert.PublicKey.(*ecdsa.PublicKey)
		
		hexSha, err := hex.DecodeString(v.sha)

		if err != nil {
			log.StdLogger.Fatalln(err.Error())
			return nil, err	
		}

		// this is a bit simplified for now
		verificationResult = VerificationResult(ecdsa.VerifyASN1(publicKey, hexSha, *v.sig))
	default:
		return nil, errors.New("public key algorithm not supported")
	}

	return &verificationResult, nil
}

func VerifyFile(path string) bool {

	err := utils.CheckPathToFile(path)
	if err != nil {
		log.StdLogger.Fatalln(err.Error())
	}


	file, err := os.Open(path)
	if err != nil {
		log.StdLogger.Fatalln(err)
	}

	defer file.Close()

	// TODO improve here so there is no need for encoding to string
	sha, _ := utils.CalculateSHA256(file)
	return VerifySha(hex.EncodeToString(sha.Sum(nil)))
}


func VerifySha(sha string) bool {

	r := NewRekorClient(RekorInstanceURLV1)

	uuid, err := r.retrieveUUID(sha)

	if err != nil {
		log.StdLogger.Fatalln(err)
	}

	b, err := r.logEntry(*uuid)
	if err != nil {
		log.StdLogger.Fatalln(err)
	}

	cert, err := extractCertificate(b)
	if err != nil {
		log.StdLogger.Fatalln(err)
	}

	sig, err := extractSignature(b)
	if err != nil {
		log.StdLogger.Fatalln(err)
	}

	v := CertificateVerifier{
		cert: cert,
		sig: sig,
		sha: sha,
	}

	res, _ := v.do()

	return bool(*res)
}

