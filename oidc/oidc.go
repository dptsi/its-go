package oidc

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"bitbucket.org/dptsi/go-framework/contracts"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

var (
	ErrNoEndSessionEndpoint = errors.New("this_oidc_provider_does_not_support_end_session_endpoint")
	ErrInvalidState         = errors.New("invalid_state")
	ErrInvalidNonce         = errors.New("invalid_nonce")
	ErrInvalidCodeChallenge = errors.New("invalid_code_challenge")
	ErrInvalidIdToken       = errors.New("invalid_id_token")
	ErrRetrieveUserInfo     = errors.New("error_retrieve_user_info")
)

const (
	stateKey        = "oidc.state"
	idTokenKey      = "oidc.id_token"
	nonceKey        = "oidc.nonce"
	codeVerifierKey = "oidc.code_verifier"
)

type Client struct {
	provider    *oidc.Provider
	oauthConfig oauth2.Config
	storage     contracts.SessionStorage

	needToVerifyState bool
	needToVerifyNonce bool
	isPKCEEnabled     bool
}

func NewClient(
	ctx context.Context,
	storage contracts.SessionStorage,
	providerUrl string,
	clientID string,
	clientSecret string,
	redirectURL string,
	scopes []string,
) (*Client, error) {
	provider, err := oidc.NewProvider(ctx, providerUrl)
	if err != nil {
		return nil, err
	}
	cfg := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,

		// Discovery returns the OAuth2 endpoints.
		Endpoint: provider.Endpoint(),

		Scopes: scopes,
	}

	return &Client{provider, cfg, storage, true, true, true}, nil
}

func (c *Client) UserInfoEndpoint() string {
	return c.provider.UserInfoEndpoint()
}

func (c *Client) SetNeedToVerifyState(needToVerifyState bool) {
	c.needToVerifyState = needToVerifyState
}

func (c *Client) SetNeedToVerifyNonce(needToVerifyNonce bool) {
	c.needToVerifyNonce = needToVerifyNonce
}

func (c *Client) SetPKCEEnabled(isPKCEEnabled bool) {
	c.isPKCEEnabled = isPKCEEnabled
}

func (c *Client) RedirectURL(ctx context.Context, sess contracts.SessionData) (string, error) {
	authCodeOptions := make([]oauth2.AuthCodeOption, 0)

	if c.isPKCEEnabled {
		codeVerifier := oauth2.GenerateVerifier()
		sess.Set(codeVerifierKey, codeVerifier)
		authCodeOptions = append(authCodeOptions, oauth2.S256ChallengeOption(codeVerifier))
	}

	if c.needToVerifyNonce {
		nonce := uuid.NewString()
		sess.Set(nonceKey, nonce)
		authCodeOptions = append(authCodeOptions, oauth2.SetAuthURLParam("nonce", nonce))
	}

	state := ""
	if c.needToVerifyState {
		state = uuid.NewString()
		sess.Set(stateKey, state)
	}

	isNeedToSaveSession := c.needToVerifyState || c.needToVerifyNonce || c.isPKCEEnabled
	if err := c.storage.Save(ctx, sess); isNeedToSaveSession && err != nil {
		return "", fmt.Errorf("RedirectURL: %w", err)
	}

	return c.oauthConfig.AuthCodeURL(
		state,
		authCodeOptions...,
	), nil
}

func (c *Client) ExchangeCodeForToken(ctx context.Context, sess contracts.SessionData, code string, state string) (*oauth2.Token, *oidc.IDToken, error) {
	if err := c.verifyState(ctx, sess, state); c.needToVerifyState && err != nil {
		return nil, nil, fmt.Errorf("unable to exchange code: %w", err)
	}

	authCodeOptions := make([]oauth2.AuthCodeOption, 0)
	codeVerifier, err := c.GetCodeVerifierAndRemoveFromSession(ctx, sess)
	if c.isPKCEEnabled && err != nil {
		return nil, nil, fmt.Errorf("unable to exchange code: %w", err)
	}
	if c.isPKCEEnabled {
		authCodeOptions = append(authCodeOptions, oauth2.VerifierOption(codeVerifier))
	}

	token, err := c.oauthConfig.Exchange(ctx, code, authCodeOptions...)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to exchange code: %w", err)
	}
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, nil, errors.New("no_id_token_in_payload")
	}

	IDToken, err := c.parseAndVerifyIDToken(ctx, rawIDToken)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to exchange code: %w", err)
	}

	if c.needToVerifyNonce {
		c.verifyNonce(ctx, sess, IDToken.Nonce)
	}

	sess.Set(idTokenKey, rawIDToken)
	if err := c.storage.Save(ctx, sess); err != nil {
		return nil, nil, fmt.Errorf("unable to exchange code: %w", err)
	}

	return token, IDToken, nil
}

func (c *Client) parseAndVerifyIDToken(ctx context.Context, rawIDToken string) (*oidc.IDToken, error) {
	// Parse and verify ID Token payload.
	var verifier = c.provider.Verifier(&oidc.Config{ClientID: c.oauthConfig.ClientID})
	parsed, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, ErrInvalidIdToken
	}

	return parsed, nil
}

func (c *Client) UserInfo(ctx context.Context, t *oauth2.Token) (*oidc.UserInfo, error) {
	userInfo, err := c.provider.UserInfo(ctx, oauth2.StaticTokenSource(t))
	if err != nil {
		return nil, ErrRetrieveUserInfo
	}

	return userInfo, nil
}

func (c *Client) RPInitiatedLogout(ctx context.Context, sess contracts.SessionData, postLogoutRedirectURI string) (string, error) {
	var claims struct {
		EndSessionEndpoint string `json:"end_session_endpoint"`
	}
	if err := c.provider.Claims(&claims); err != nil {
		return "", fmt.Errorf("rp initiated logout: get provider claims: %w", err)
	}
	endSessionEndpoint := claims.EndSessionEndpoint
	if endSessionEndpoint == "" {
		return "", ErrNoEndSessionEndpoint
	}
	req, err := http.NewRequest("GET", endSessionEndpoint, nil)
	if err != nil {
		return "", fmt.Errorf("rp initiated logout: make request: %w", err)
	}
	q := req.URL.Query()

	idTokenHintItf, exists := sess.Get(idTokenKey)
	if idTokenHint, ok := idTokenHintItf.(string); exists && ok && idTokenHint != "" {
		q.Add("id_token_hint", idTokenHint)
	}
	sess.Delete(idTokenKey)
	c.storage.Save(ctx, sess)

	if postLogoutRedirectURI != "" {
		q.Add("post_logout_redirect_uri", postLogoutRedirectURI)
	}

	req.URL.RawQuery = q.Encode()
	return req.URL.String(), nil
}

func (c *Client) verifyState(ctx context.Context, sess contracts.SessionData, state string) error {
	cookieState, ok := sess.Get(stateKey)
	if !ok {
		cookieState = ""
	}
	sess.Delete(stateKey)
	c.storage.Save(ctx, sess)

	if state == "" || state != cookieState {
		return ErrInvalidState
	}

	return nil
}

func (c *Client) GetCodeVerifierAndRemoveFromSession(ctx context.Context, sess contracts.SessionData) (string, error) {
	codeVerifierIf, ok := sess.Get(codeVerifierKey)
	sess.Delete(codeVerifierKey)
	if err := c.storage.Save(ctx, sess); err != nil {
		return "", err
	}
	codeVerifier := ""
	if ok {
		codeVerifier, ok = codeVerifierIf.(string)
		if !ok {
			codeVerifier = ""
		}
	}

	if codeVerifier == "" {
		return "", ErrInvalidCodeChallenge
	}

	return codeVerifier, nil
}

func (c *Client) verifyNonce(ctx context.Context, sess contracts.SessionData, nonce string) error {
	cookieNonce, ok := sess.Get(nonceKey)
	if !ok {
		cookieNonce = ""
	}
	sess.Delete(nonceKey)
	c.storage.Save(ctx, sess)

	if nonce != cookieNonce {
		return ErrInvalidNonce
	}

	return nil
}
