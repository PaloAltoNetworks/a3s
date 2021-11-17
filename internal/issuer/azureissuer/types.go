package azureissuer

import "fmt"

// ErrAzure represents an error that happened
// during operation related to Azure.
type ErrAzure struct {
	Err error
}

func (e ErrAzure) Error() string {
	return fmt.Sprintf("azure error: %s", e.Err)
}

func (e ErrAzure) Unwrap() error {
	return e.Err
}

type azureJWT struct {
	AIO      string `json:"aio"`
	AppID    string `json:"appid"`
	AppIDAcr string `json:"appidacr"`
	IDP      string `json:"idp"`
	OID      string `json:"oid"`
	RH       string `json:"rh"`
	TID      string `json:"tid"`
	UTI      string `json:"uti"`
	XmsMIRID string `json:"xms_mirid"`
}
