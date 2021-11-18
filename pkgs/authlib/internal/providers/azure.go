package providers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

var azureServiceTokenURL = "http://169.254.169.254/metadata/identity/oauth2/token"

// AzureServiceIdentityToken will retrieve the service account token for
// the VM using the Metadata Identity Service of Azure.
func AzureServiceIdentityToken() (string, error) {
	body, err := issueRequest(azureServiceTokenURL)
	if err != nil {
		return "", err
	}

	token := struct {
		AccessToken string `json:"access_token"`
	}{}

	err = json.Unmarshal(body, &token)
	if err != nil {
		return "", fmt.Errorf("invalid token returned by metadata service: %s", err)
	}

	return token.AccessToken, nil
}

func issueRequest(baseuri string) ([]byte, error) {

	var endpoint *url.URL
	endpoint, err := url.Parse(baseuri)
	if err != nil {
		return nil, fmt.Errorf("unable to access the service account URL: %s", err)
	}

	parameters := url.Values{}
	parameters.Add("api-version", "2018-02-01")
	parameters.Add("resource", "https://management.azure.com")

	endpoint.RawQuery = parameters.Encode()
	req, err := http.NewRequest("GET", endpoint.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create HTTP request: %s", err)
	}
	req.Header.Add("Metadata", "true")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to issue request: %s", err)
	}

	defer resp.Body.Close() // nolint errcheck
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read data: %s", err)
	}

	return body, nil
}
