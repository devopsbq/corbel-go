package silkroad

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// OauthToken gets an access token
//
// API Docs: http://docs.silkroadiam.apiary.io/#reference/authorization/oauthtoken
func (i *IAMService) OauthToken() (string, error) {
	signingMethod := jwt.GetSigningMethod(i.client.ClientJWTSigningMethod)
	token := jwt.New(signingMethod)
	// Required JWT Claims for SR
	token.Claims["aud"] = "http://iam.bqws.io"
	token.Claims["exp"] = time.Now().Add(time.Second * i.client.TokenExpirationTime).Unix()
	token.Claims["iss"] = i.client.ClientID
	token.Claims["scope"] = i.client.ClientScopes
	token.Claims["domain"] = i.client.ClientDomain
	token.Claims["name"] = i.client.ClientName

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(i.client.ClientSecret))
	if err != nil {
		return "", errJWTEncodingError
	}

	values := url.Values{}
	values.Set("grant_type", grantType)
	values.Set("assertion", tokenString)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s", i.client.URLFor("iam", "/v1.0/oauth/token")), bytes.NewBufferString(values.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Add("User-Agent", userAgent)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := i.client.httpClient.Do(req)
	if err != nil {
		return "", errClientNotAuthorized
	}

	defer res.Body.Close()
	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", errResponseError
	}

	var iamResponse iamOauthTokenResponse
	err = json.Unmarshal(contents, &iamResponse)
	if err != nil {
		return "", errJSONUnmarshalError
	}

	return iamResponse.AccessToken, nil
}
