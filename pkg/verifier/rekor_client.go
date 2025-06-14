package verifier

import (
	"bytes"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/kasia-kittel/rekor-verifier/pkg/log"
)


const (
    RekorInstanceURLV1 = "https://rekor.sigstore.dev/api/v1"
    retrieveUUIDEndpoint = "/index/retrieve"
    logEntryEndpoint = "/log/entries/"
)


type RekorClient struct {
    BaseURL    string
    HTTPClient *http.Client
}

func NewRekorClient(baseURL string) *RekorClient {
    return &RekorClient{
        BaseURL: baseURL,
        HTTPClient: &http.Client{
            Timeout: time.Minute,
        },
    }
}

type RetrieveUUIDRequest struct {
	Hash string `json:"hash"`
}

type RetrieveUUIDResponse []string

type UUID string

func (c *RekorClient) retrieveUUID(sha string) (*UUID, error) {

    requestPayload, err := json.Marshal(RetrieveUUIDRequest{Hash: "sha256:" + sha})

    if err != nil {
		log.StdOutLogger.Println(err.Error())
		return nil, err
	}

    request := bytes.NewBuffer(requestPayload)
    retrieveUrl := c.BaseURL + retrieveUUIDEndpoint
	rawResponse, err := c.HTTPClient.Post(retrieveUrl, "application/json", request)

    if err != nil {
		log.StdOutLogger.Println(err.Error())
        return nil, err
	}

	if rawResponse.StatusCode != 200 {
		log.StdOutLogger.Println(rawResponse.Status)
        return nil, err
	}

	body, err := io.ReadAll(rawResponse.Body)
	
    if err != nil {
		log.StdOutLogger.Println(err.Error())
        return nil, err
	}
    
    var response RetrieveUUIDResponse

	err = json.Unmarshal([]byte(body), &response)

	if err != nil {
		log.StdOutLogger.Println(err.Error())
        return nil, err
	}

	resp := UUID(response[0])

    return &resp, nil
}

type LogEntryResponse struct {
	BodyBase64 string `json:"body"`
}


type DecodedBody struct {
	Spec struct {
		Signature struct {
			Content string	
			PublicKey struct {
				Content string
			}
		}
	}
}


// For now, as this is a PoF, this will only work for Entries containing 
// Signature and Public Key
// TODO: log if the kind of entry is different
func (c *RekorClient) logEntry(uuid UUID) (*DecodedBody, error) {
	
    rawResponse, err := c.HTTPClient.Get(c.BaseURL + logEntryEndpoint + string(uuid))

	if err != nil {
		log.StdOutLogger.Println(err.Error())
        return nil, err
	}

	if rawResponse.StatusCode != 200 {
		log.StdOutLogger.Println(rawResponse.Status)
        return nil, err
	}

    body, err := io.ReadAll(rawResponse.Body)

	if err != nil {
		log.StdOutLogger.Println(err.Error())
        return nil, err
	}

	logEntryResponse := make(map[UUID]*json.RawMessage)

	err = json.Unmarshal([]byte(body), &logEntryResponse)
	if err != nil {
		log.StdOutLogger.Println(err.Error())	
        return nil, err
	}

	var logEntry LogEntryResponse 

	v := logEntryResponse[uuid]

	err = json.Unmarshal(*v, &logEntry)
	if err != nil {
		log.StdOutLogger.Println(err.Error())	
        return nil, err
	}

	decodedBodyString := logEntry.BodyBase64

	data, err := base64.StdEncoding.DecodeString(decodedBodyString)
    if err != nil {
		log.StdOutLogger.Println(err.Error())
        return nil, err
	}

	var decodedBody DecodedBody

	err = json.Unmarshal(data, &decodedBody)
	if err != nil {
		log.StdOutLogger.Println(err.Error())
        return nil, err	
	}

	return &decodedBody, nil
}


// helper functions
func extractCertificate(b *DecodedBody) (*x509.Certificate, error) {
 
    // the certificate is base64 encoded string
    decodedBody, err := base64.StdEncoding.DecodeString(b.Spec.Signature.PublicKey.Content)
    
    if err != nil {
        log.StdOutLogger.Println(err.Error())	
        return nil, err
    }

    pemBlock, _ := pem.Decode(decodedBody)

    if (pemBlock == nil) {
        err = errors.New("can't decode certificate data")
        log.StdOutLogger.Println(err.Error())	
        return nil, err
    }

	cert, err := x509.ParseCertificate(pemBlock.Bytes)

	if err != nil {
		log.StdOutLogger.Println(err.Error())
        return nil, err	
	}

    return cert, nil
}

type Signature []byte

func extractSignature(b *DecodedBody) (*Signature, error) {
    
    var sig Signature 
    
    sig, err := base64.StdEncoding.DecodeString(b.Spec.Signature.Content)

    if err != nil {
		log.StdOutLogger.Println(err.Error())
        return nil, err	
	}

    return &sig, nil
}
