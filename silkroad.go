package silkroad

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

var (
	userAgent                string
	allowedEnvironments      []string
	allowedEndpoints         []string
	allowedJTWSigningMethods []string
)

// init defines constants that will be used later
func init() {
	userAgent = fmt.Sprintf("go-silkroad/%s", Version)
	allowedEnvironments = []string{"production", "staging", "current", "next", "qa", "integration"}
	allowedEndpoints = []string{"iam", "oauth", "assets", "resources"}
	allowedJTWSigningMethods = []string{"HS256", "RSA"}
}

// Client is the struct that manages communication with the Silkroad APIs.
type Client struct {
	// client is the HTTP client to communicate with the API.
	httpClient *http.Client

	// Environment is used to define the target environment to speak with.
	Environment string

	// ClientName is the name that match the clientID
	// (Optional) The required information is the clientID
	ClientName string

	// ClientID is the application defined client on Silkroad.
	ClientID string

	// ClientSecret is the application secret hash that match with clientID.
	ClientSecret string

	// ClientScopes are those scopes the client will ask for to the platform when building the client connection
	// ClientScopes is a string with the scopes delimited by spaces.
	ClientScopes string

	// ClientDomain is the SR domain where to make the operations using the provided credentials.
	// (Optional) Every clientID only maps to one SR domain.
	ClientDomain string

	// ClientJWTSigningMethod defines the signing method configured for the client.
	// Must match with the one configured on the platform since it will understand only that one.
	// Only allowed signing methods at the moment are: HS256 and RSA
	ClientJWTSigningMethod string

	// TokenExpirationTime define the amount of time in seconds that a token must be valid.
	// It must be lower than 3600 seconds, since is the imposed requisite from the platform.
	TokenExpirationTime time.Duration

	// UserAgent defines the UserAgent to send in the Headers for every request to the platform.
	UserAgent string

	// IAM endpoint struct
	IAM *IAMService
}

// URLFor returns the formated url of the API using the actual url scheme
func (c *Client) URLFor(endpoint, uri string) (url string) {
	switch c.Environment {
	case "production":
		url = fmt.Sprintf("https://%s.bqws.io%s", endpoint, uri)
	default:
		url = fmt.Sprintf("https://%s-%s.bqws.io%s", endpoint, c.Environment, uri)
	}
	return
}

// NewRequest creates an API request.
// method is the HTTP method to use
// endpoint is the endpoint of SR to speak with
// url is the url to query. it must be preceded by a slash.
// body is, if specified, the value JSON encoded to be used as request body.
// mediaType is the desired media type to set in the request.
func (c *Client) NewRequest(method, endpoint, urlStr, mediaType string, body interface{}) (*http.Request, error) {
	u, _ := url.Parse(c.URLFor(endpoint, urlStr))

	buf := new(bytes.Buffer)
	if body != nil {
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", mediaType)
	req.Header.Add("Accept", mediaType)
	req.Header.Add("User-Agent", userAgent)
	// req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authBearer))
	return req, nil
}

type iamOauthTokenResponse struct {
	AccessToken  string `json:"accessToken,omitempty"`
	ExpiresAt    int    `json:"expiresAt,omitempty"`
	RefreshToken string `json:"refreshToken,omitempty"`
}

// // GetToken returns
// func (c *Client) GetToken() (string, error) {
// 	signingMethod := jwt.GetSigningMethod(c.ClientJWTSigningMethod)
// 	token := jwt.New(signingMethod)
// 	// Set some claims
// 	token.Claims["aud"] = "http://iam.bqws.io"
// 	token.Claims["exp"] = time.Now().Add(time.Second * c.TokenExpirationTime).Unix()
// 	token.Claims["iss"] = c.ClientID
// 	token.Claims["scope"] = c.ClientScopes
// 	token.Claims["domain"] = c.ClientDomain
// 	token.Claims["name"] = c.ClientName
//
// 	// Sign and get the complete encoded token as a string
// 	tokenString, err := token.SignedString([]byte(c.ClientSecret))
// 	if err != nil {
// 		return "", errJWTEncodingError
// 	}
//
// 	values := url.Values{}
// 	values.Set("grant_type", grantType)
// 	values.Set("assertion", tokenString)
//
// 	req, err := http.NewRequest("POST", fmt.Sprintf("%s", c.URLFor("iam", "/v1.0/oauth/token")), bytes.NewBufferString(values.Encode()))
// 	if err != nil {
// 		return "", err
// 	}
// 	req.Header.Add("User-Agent", userAgent)
// 	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
//
// 	res, err := c.httpClient.Do(req)
// 	if err != nil {
// 		return "", errClientNotAuthorized
// 	}
//
// 	defer res.Body.Close()
// 	contents, err := ioutil.ReadAll(res.Body)
// 	if err != nil {
// 		return "", errResponseError
// 	}
//
// 	var iamResponse iamOauthTokenResponse
// 	err = json.Unmarshal(contents, &iamResponse)
// 	if err != nil {
// 		return "", errJSONUnmarshalError
// 	}
//
// 	return iamResponse.AccessToken, nil
// }

// NewClient returns a new Silkroad API client.
// If a nil httpClient is provided, it will return a http.DefaultClient.
// If a empty environment is provided, it will use production as environment.
func NewClient(httpClient *http.Client, environment, clientID, clientName, clientSecret, clientScopes, clientDomain, clientJWTSigningMethod string, tokenExpirationTime time.Duration) (*Client, error) {

	var thisClient *Client

	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	if environment == "" {
		environment = "production"
	}

	// allowedEnvironments?
	if stringInSlice(allowedEnvironments, environment) == false {
		return nil, errInvalidEnvironment
	}

	// allowedJTWSigningMethods?
	if stringInSlice(allowedJTWSigningMethods, clientJWTSigningMethod) == false {
		return nil, errInvalidJWTSigningMethod
	}

	// incorrect Token Expiration Time?
	if tokenExpirationTime > 3600 || tokenExpirationTime == 0 {
		return nil, errInvalidTokenExpirationTime
	}

	// required parameters?
	if clientID == "" || clientSecret == "" {
		return nil, errMissingClientParams
	}

	thisClient = &Client{
		httpClient:             httpClient,
		Environment:            environment,
		ClientName:             clientName,
		ClientID:               clientID,
		ClientSecret:           clientSecret,
		ClientDomain:           clientDomain,
		ClientScopes:           clientScopes,
		ClientJWTSigningMethod: clientJWTSigningMethod,
		TokenExpirationTime:    tokenExpirationTime,
		UserAgent:              userAgent,
	}

	thisClient.IAM = &IAMService{client: thisClient}

	return thisClient, nil
}
