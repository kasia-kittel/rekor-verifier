package verifier

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func NewMockClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func TestRetrieveUUID(t *testing.T) {
	js, _ := json.Marshal([]string{"test-uuid"})

	client := NewMockClient(func(req *http.Request) *http.Response {
		
		assert.Equal(t, req.URL.String(), "https://rekor.sigstore.dev/api/v1/index/retrieve")
		
		return &http.Response{
			StatusCode: 200,
			Body:      io.NopCloser(bytes.NewReader(js)),
			Header:    make(http.Header),
		}
	})

	rekorClient := RekorClient{"https://rekor.sigstore.dev/api/v1", client}
	resp, _ := rekorClient.retrieveUUID("sha")
	assert.EqualValues(t, UUID("test-uuid"), *resp)
}

func createDecodedBody(publicKeyContent string, signatureContent string) *DecodedBody {
	// better way to init nested structs?
	decodedBody := &DecodedBody {}
	decodedBody.Spec.Signature.Content = signatureContent
	decodedBody.Spec.Signature.PublicKey.Content = publicKeyContent

	return decodedBody
}

const (
	validPublicKeyContent = "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUd5akNDQmxDZ0F3SUJBZ0lVUURFWmMzZVJ4WW1WZ2dGbTBBNS9Dd3R4QUVRd0NnWUlLb1pJemowRUF3TXcKTnpFVk1CTUdBMVVFQ2hNTWMybG5jM1J2Y21VdVpHVjJNUjR3SEFZRFZRUURFeFZ6YVdkemRHOXlaUzFwYm5SbApjbTFsWkdsaGRHVXdIaGNOTWpReE1USXdNVGMxTmpFNFdoY05NalF4TVRJd01UZ3dOakU0V2pBQU1Ga3dFd1lICktvWkl6ajBDQVFZSUtvWkl6ajBEQVFjRFFnQUV4ZFgxU3I2am1kWVBadXp6ZjRuT3EyRklnemltUGE4WXVsQ24KY1FpNEhSSW8weU1ibE8zMmZWSkV0RUhQYWpFNnZIc0ZveFlZMFU4d0M1M0preUc4ZWFPQ0JXOHdnZ1ZyTUE0RwpBMVVkRHdFQi93UUVBd0lIZ0RBVEJnTlZIU1VFRERBS0JnZ3JCZ0VGQlFjREF6QWRCZ05WSFE0RUZnUVVXdjA1CkE3azJSaldHOHp3VmhHVnY3OUxaeVUwd0h3WURWUjBqQkJnd0ZvQVUzOVBwejFZa0VaYjVxTmpwS0ZXaXhpNFkKWkQ4d1pRWURWUjBSQVFIL0JGc3dXWVpYYUhSMGNITTZMeTluYVhSb2RXSXVZMjl0TDJOb1lXbHVaM1ZoY21RdApaR1YyTDJGd2EyOHZMbWRwZEdoMVlpOTNiM0pyWm14dmQzTXZjbVZzWldGelpTNTVZVzFzUUhKbFpuTXZkR0ZuCmN5OTJNQzR5TUM0eE1Ea0dDaXNHQVFRQmc3OHdBUUVFSzJoMGRIQnpPaTh2ZEc5clpXNHVZV04wYVc5dWN5NW4KYVhSb2RXSjFjMlZ5WTI5dWRHVnVkQzVqYjIwd0VnWUtLd1lCQkFHRHZ6QUJBZ1FFY0hWemFEQTJCZ29yQmdFRQpBWU8vTUFFREJDaGpaV1V6TjJNM05HSXhNV1l6TXpabE9UazJOVFppTlRVM056aGhaV1ZsWkRnek16WXdNV000Ck1Cd0dDaXNHQVFRQmc3OHdBUVFFRGtOeVpXRjBaU0JTWld4bFlYTmxNQ0VHQ2lzR0FRUUJnNzh3QVFVRUUyTm8KWVdsdVozVmhjbVF0WkdWMkwyRndhMjh3SHdZS0t3WUJCQUdEdnpBQkJnUVJjbVZtY3k5MFlXZHpMM1l3TGpJdwpMakV3T3dZS0t3WUJCQUdEdnpBQkNBUXREQ3RvZEhSd2N6b3ZMM1J2YTJWdUxtRmpkR2x2Ym5NdVoybDBhSFZpCmRYTmxjbU52Ym5SbGJuUXVZMjl0TUdjR0Npc0dBUVFCZzc4d0FRa0VXUXhYYUhSMGNITTZMeTluYVhSb2RXSXUKWTI5dEwyTm9ZV2x1WjNWaGNtUXRaR1YyTDJGd2EyOHZMbWRwZEdoMVlpOTNiM0pyWm14dmQzTXZjbVZzWldGegpaUzU1WVcxc1FISmxabk12ZEdGbmN5OTJNQzR5TUM0eE1EZ0dDaXNHQVFRQmc3OHdBUW9FS2d3b1kyVmxNemRqCk56UmlNVEZtTXpNMlpUazVOalUyWWpVMU56YzRZV1ZsWldRNE16TTJNREZqT0RBZEJnb3JCZ0VFQVlPL01BRUwKQkE4TURXZHBkR2gxWWkxb2IzTjBaV1F3TmdZS0t3WUJCQUdEdnpBQkRBUW9EQ1pvZEhSd2N6b3ZMMmRwZEdoMQpZaTVqYjIwdlkyaGhhVzVuZFdGeVpDMWtaWFl2WVhCcmJ6QTRCZ29yQmdFRUFZTy9NQUVOQkNvTUtHTmxaVE0zCll6YzBZakV4WmpNek5tVTVPVFkxTm1JMU5UYzNPR0ZsWldWa09ETXpOakF4WXpnd0lRWUtLd1lCQkFHRHZ6QUIKRGdRVERCRnlaV1p6TDNSaFozTXZkakF1TWpBdU1UQVpCZ29yQmdFRUFZTy9NQUVQQkFzTUNUUTFOekF5TmpVMQpPVEF4QmdvckJnRUVBWU8vTUFFUUJDTU1JV2gwZEhCek9pOHZaMmwwYUhWaUxtTnZiUzlqYUdGcGJtZDFZWEprCkxXUmxkakFZQmdvckJnRUVBWU8vTUFFUkJBb01DRGczTkRNMk5qazVNR2NHQ2lzR0FRUUJnNzh3QVJJRVdReFgKYUhSMGNITTZMeTluYVhSb2RXSXVZMjl0TDJOb1lXbHVaM1ZoY21RdFpHVjJMMkZ3YTI4dkxtZHBkR2gxWWk5MwpiM0pyWm14dmQzTXZjbVZzWldGelpTNTVZVzFzUUhKbFpuTXZkR0ZuY3k5Mk1DNHlNQzR4TURnR0Npc0dBUVFCCmc3OHdBUk1FS2d3b1kyVmxNemRqTnpSaU1URm1Nek0yWlRrNU5qVTJZalUxTnpjNFlXVmxaV1E0TXpNMk1ERmoKT0RBVUJnb3JCZ0VFQVlPL01BRVVCQVlNQkhCMWMyZ3dXZ1lLS3dZQkJBR0R2ekFCRlFSTURFcG9kSFJ3Y3pvdgpMMmRwZEdoMVlpNWpiMjB2WTJoaGFXNW5kV0Z5WkMxa1pYWXZZWEJyYnk5aFkzUnBiMjV6TDNKMWJuTXZNVEU1Ck16ZzVPRGc0TmpndllYUjBaVzF3ZEhNdk1UQVdCZ29yQmdFRUFZTy9NQUVXQkFnTUJuQjFZbXhwWXpDQmlnWUsKS3dZQkJBSFdlUUlFQWdSOEJIb0FlQUIyQU4wOU1Hckd4eEV5WXhrZUhKbG5Od0tpU2w2NDNqeXQvNGVLY29BdgpLZTZPQUFBQmswcTN4c29BQUFRREFFY3dSUUlnWXlNdVBlMjlhMFNxemFHaGZ1ZXdhdGVlVCtNem5ZMjBNR0Y3ClZ5eW1GeHNDSVFDdzhBRndOM2tPcks0VHJiVDljbHVXWGVOdFVmNDFPbzF1cmJLbjFhaHZmakFLQmdncWhrak8KUFFRREF3Tm9BREJsQWpBa1VKcXB4R0FHcGxvVFVLK1NKOU9vZU5OWWVReUxPc2piM2YxdjRhUUp2a1ZEZjJubQo3cGpBYk1IMG1QTHpHMUlDTVFDWEpVSGd1TEdFZTBDZHFvVTljZmhGSTNEdFA2dkM0VnNscXpaWWJSUXpZWjVyCkZyY3dYeE1aSDkxV2JLL0NXSUk9Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K"
	validSignatureContent = "MEQCIGjwFf32DzBCrbSKYeJF+Ojss67EP4KE4T4kItxXhTzwAiALAxR6F7kd053HH7vXJVoWvw4fhqfr0HZZdX2hTOLHVg=="
)

func TestExtractCertificate(t *testing.T){
	tests := []struct {
		body    *DecodedBody
		error   bool
		description string
	}{
		{ body: createDecodedBody(validPublicKeyContent, ""), error: false, description: "Valid public key content should return valid certificate" },
		{ body: createDecodedBody("abc", ""), error: true, description: "Invalid base64 encoded string should return error " },
		{ body: createDecodedBody("YWJj", ""), error:true, description: "Valid base64 encoded string with invalid content should return error " },	
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			cert, err := extractCertificate(tt.body)
			isError := err!=nil

			if isError != tt.error {
				t.Fail()
			}

			if (!isError && cert == nil) {
				t.Fail()
			}

		})
	}
}

func TestExtractSignature(t *testing.T){	
	tests := []struct {
		body    *DecodedBody
		error   bool
		description string
	}{
		{ body: createDecodedBody("", validSignatureContent), error: false, description: "Valid public key content should return valid certificate" },
		{ body: createDecodedBody("", "abc"), error: true, description: "Invalid base64 encoded string should return error " },
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			
			sig, err := extractSignature(tt.body)

			isError := err!=nil

			if isError != tt.error {
				t.Fail()
			}

			if (!isError && sig == nil) {
				t.Fail()
			}
		})
	}
}

func TestRetrieveUUIDIntegration(t *testing.T) {
	
	if testing.Short() {
        t.Skip("skipping integration test")
    }
	
	rekorClient := NewRekorClient(RekorInstanceURLV1)
	uuid, _ := rekorClient.retrieveUUID("442d8baafc0c3a873b21a3add32f5c65f538fb5cbcf4a4a69ba098a2b730c5d2")
	assert.EqualValues(t, UUID("108e9186e8c5677a8d6736bdd79170adf94bd127aea751274d1d62504e88b058af7552d91dea0f26"), *uuid)
}

func TestLogEntryIntegration(t *testing.T) {
	
	if testing.Short() {
        t.Skip("skipping integration test")
    }
	
	rekorClient := NewRekorClient(RekorInstanceURLV1)
	body, _ := rekorClient.logEntry(UUID("108e9186e8c5677a8d6736bdd79170adf94bd127aea751274d1d62504e88b058af7552d91dea0f26"))

	expectedSignature := "MEQCIGjwFf32DzBCrbSKYeJF+Ojss67EP4KE4T4kItxXhTzwAiALAxR6F7kd053HH7vXJVoWvw4fhqfr0HZZdX2hTOLHVg=="
	assert.EqualValues(t, expectedSignature, body.Spec.Signature.Content)
}
