package infrastructure

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/adaptive-ai-learn/backend/internal/auth_engine/domain"
)

type GoogleOAuth struct {
	clientID     string
	clientSecret string
	redirectURL  string
}

func NewGoogleOAuth(clientID, clientSecret, redirectURL string) *GoogleOAuth {
	return &GoogleOAuth{
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURL:  redirectURL,
	}
}

func (g *GoogleOAuth) GetLoginURL(state string) string {
	oauthURL := "https://accounts.google.com/o/oauth2/v2/auth"

	params := url.Values{}
	params.Add("client_id", g.clientID)
	params.Add("redirect_uri", g.redirectURL)
	params.Add("response_type", "code")
	params.Add("scope", "email profile")
	params.Add("state", state) // Used for CSRF protection
	params.Add("access_type", "offline")

	return fmt.Sprintf("%s?%s", oauthURL, params.Encode())
}

func (g *GoogleOAuth) ExchangeCodeForUser(ctx context.Context, code string) (*domain.GoogleUserData, error) {
	// 1. Exchange auth code for Google Access Token
	tokenURL := "https://oauth2.googleapis.com/token"

	data := url.Values{}
	data.Add("client_id", g.clientID)
	data.Add("client_secret", g.clientSecret)
	data.Add("code", code)
	data.Add("grant_type", "authorization_code")
	data.Add("redirect_uri", g.redirectURL)

	resp, err := http.PostForm(tokenURL, data)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to authenticate with google")
	}

	var tokenRes struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenRes); err != nil {
		return nil, err
	}

	// 2. Fetch User Profile
	userInfoURL := "https://www.googleapis.com/oauth2/v2/userinfo"
	req, _ := http.NewRequestWithContext(ctx, "GET", userInfoURL, nil)
	req.Header.Add("Authorization", "Bearer "+tokenRes.AccessToken)

	client := &http.Client{}
	userResp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user info: %w", err)
	}
	defer userResp.Body.Close()

	if userResp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch user info from google")
	}

	var googleUser struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}
	if err := json.NewDecoder(userResp.Body).Decode(&googleUser); err != nil {
		return nil, err
	}

	return &domain.GoogleUserData{
		GoogleID:  googleUser.ID,
		Email:     googleUser.Email,
		Name:      googleUser.Name,
		AvatarURL: googleUser.Picture,
	}, nil
}
