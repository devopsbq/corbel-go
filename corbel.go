package corbel

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
)

var (
	userAgent                string
	allowedEnvironments      []string
	allowedEndpoints         []string
	allowedJTWSigningMethods []string
)

// init defines constants that will be used later
func init() {
	userAgent = fmt.Sprintf("corbel-go/%s", Version)
	allowedEndpoints = []string{"iam", "oauth", "assets", "resources"}
	allowedJTWSigningMethods = []string{"HS256", "RSA"}
}

// Client is the struct that manages communication with the Corbel APIs.
type Client struct {
	// client is the HTTP client to communicate with the API.
	httpClient *http.Client

	// Endpoint is a structure that stores the uri for each resource
	Endpoints map[string]string

	// ClientName is the name that match the clientID
	// (Optional) The required information is the clientID
	ClientName string

	// ClientID is the application defined client on Corbel.
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
	TokenExpirationTime uint64

	// UserAgent defines the UserAgent to send in the Headers for every request to the platform.
	UserAgent string

	// Token is the actual token to send as Authentication Bearer
	CurrentToken string

	// CurrentTokenExpiresAt is the unix time where the token will expire
	CurrentTokenExpiresAt int64

	// CurrentRefreshToken is the current refresh token received from the IAM service
	CurrentRefreshToken string

	// IAM endpoint struct
	IAM       *IAMService
	Resources *ResourcesService
	Assets    *AssetsService

	// Logger
	logger *logrus.Logger

	// LogLevel for the logger.
	LogLevel string
}

// URLFor returns the formated url of the API using the actual url scheme
func (c *Client) URLFor(endpoint, uri string) string {
	return fmt.Sprintf("%s%s", c.Endpoints[endpoint], uri)
}

// Token returns the token to use as bearer. If the token has already expired
// it refresh it.
func (c *Client) Token() string {
	// if CurrentToken == "" then return it as is
	if c.CurrentToken == "" {
		return c.CurrentToken
	}
	// if we have CurrentToken check if already expired
	if c.CurrentTokenExpiresAt <= time.Now().Unix()*1000 {
		c.logger.Debug("refreshing token")
		if c.CurrentRefreshToken != "" {
			_ = c.IAM.RefreshToken()
		} else {
			_ = c.IAM.OauthToken()
		}
	}
	return c.CurrentToken
}

// DefaultClient return a client with most of its values set to the default ones
func DefaultClient(endpoints map[string]string, clientID, clientName, clientSecret, clientScopes, clientDomain string) (*Client, error) {
	return NewClient(nil, endpoints, clientID, clientName, clientSecret, clientScopes, clientDomain, "HS256", 3600, "info")
}

// NewClient returns a new Corbel API client.
// If a nil httpClient is provided, it will return a http.DefaultClient.
// If a empty environment is provided, it will use production as environment.
func NewClient(httpClient *http.Client, endpoints map[string]string, clientID, clientName, clientSecret, clientScopes, clientDomain, clientJWTSigningMethod string, tokenExpirationTime uint64, logLevel string) (*Client, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	if len(endpoints) == 0 {
		endpoints = map[string]string{"iam": "https://iam.bqws.io", "resources": "https://resources.bqws.io"}
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

	thisClient := &Client{
		httpClient:             httpClient,
		Endpoints:              endpoints,
		ClientName:             clientName,
		ClientID:               clientID,
		ClientSecret:           clientSecret,
		ClientDomain:           clientDomain,
		ClientScopes:           clientScopes,
		ClientJWTSigningMethod: clientJWTSigningMethod,
		TokenExpirationTime:    tokenExpirationTime * 1000,
		UserAgent:              userAgent,
	}

	thisClient.IAM = &IAMService{client: thisClient}
	thisClient.Resources = &ResourcesService{client: thisClient}
	thisClient.Assets = &AssetsService{client: thisClient}

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return nil, errInvalidLogLevel
	}
	thisClient.logger = logrus.New()
	thisClient.logger.Level = level
	thisClient.LogLevel = logLevel

	return thisClient, nil
}
