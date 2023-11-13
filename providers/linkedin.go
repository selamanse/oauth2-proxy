package providers

import (
	"context"
	"errors"
	"net/http"
	"net/url"

	"github.com/oauth2-proxy/oauth2-proxy/v7/pkg/apis/sessions"
	"github.com/oauth2-proxy/oauth2-proxy/v7/pkg/requests"
)

// LinkedInProvider represents an LinkedIn based OpenID Connect Provider
// https://learn.microsoft.com/en-us/linkedin/consumer/integrations/self-serve/sign-in-with-linkedin-v2
type LinkedInProvider struct {
	*ProviderData
}

var _ Provider = (*LinkedInProvider)(nil)

const (
	linkedinProviderName = "LinkedIn"
	linkedinDefaultScope = "openid email profile"
)

var (
	// Default Login URL for LinkedIn.
	// Pre-parsed URL of https://www.linkedin.com/uas/oauth2/authorization.
	linkedinDefaultLoginURL = &url.URL{
		Scheme: "https",
		Host:   "www.linkedin.com",
		Path:   "/oauth/v2/authorization",
	}

	// Default Redeem URL for LinkedIn.
	// Pre-parsed URL of https://www.linkedin.com/uas/oauth2/accessToken.
	linkedinDefaultRedeemURL = &url.URL{
		Scheme: "https",
		Host:   "www.linkedin.com",
		Path:   "/oauth2/v2/accessToken",
	}

	// Default Profile URL for LinkedIn.
	// Pre-parsed URL of https://www.linkedin.com/v1/people/~/email-address.
	linkedinDefaultProfileURL = &url.URL{
		Scheme: "https",
		Host:   "api.linkedin.com",
		Path:   "/v2/userinfo",
	}

	// Default Validate URL for LinkedIn.
	linkedinDefaultValidateURL = &url.URL{
		Scheme: "https",
		Host:   "api.linkedin.com",
		Path:   "/v2/userinfo",
	}
)

// NewLinkedInProvider initiates a new LinkedInProvider
func NewLinkedInProvider(p *ProviderData) *LinkedInProvider {
	p.setProviderDefaults(providerDefaults{
		name:        linkedinProviderName,
		loginURL:    linkedinDefaultLoginURL,
		redeemURL:   linkedinDefaultRedeemURL,
		profileURL:  linkedinDefaultProfileURL,
		validateURL: linkedinDefaultValidateURL,
		scope:       linkedinDefaultScope,
	})
	p.getAuthorizationHeaderFunc = makeLinkedInHeader

	return &LinkedInProvider{ProviderData: p}
}

func makeLinkedInHeader(accessToken string) http.Header {
	// extra headers required by the LinkedIn API when making authenticated requests
	extraHeaders := map[string]string{
		acceptHeader:  acceptApplicationJSON,
		"x-li-format": "json",
	}
	return makeAuthorizationHeader(tokenTypeBearer, accessToken, extraHeaders)
}

// GetEmailAddress returns the Account email address
func (p *LinkedInProvider) GetEmailAddress(ctx context.Context, s *sessions.SessionState) (string, error) {
	if s.AccessToken == "" {
		return "", errors.New("missing access token")
	}

	requestURL := p.ProfileURL.String()
	json, err := requests.New(requestURL).
		WithContext(ctx).
		WithHeaders(makeLinkedInHeader(s.AccessToken)).
		Do().
		UnmarshalSimpleJSON()
	if err != nil {
		return "", err
	}
	email, err := json.Get("email").String()
	if err != nil {
		return "", err
	}
	return email, nil
}

// ValidateSession validates the AccessToken
func (p *LinkedInProvider) ValidateSession(ctx context.Context, s *sessions.SessionState) bool {
	return validateToken(ctx, p, s.AccessToken, makeLinkedInHeader(s.AccessToken))
}
