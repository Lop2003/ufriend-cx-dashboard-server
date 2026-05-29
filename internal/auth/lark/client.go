package lark

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	defaultBaseURL     = "https://open.larksuite.com"
	defaultAccountsURL = "https://accounts.larksuite.com"
	// offline_access จำเป็นสำหรับ refresh_token (single-use, ~30 วัน)
	defaultOAuthScope = "offline_access"
)

// Client เรียก Lark Open API สำหรับ OAuth flow
type Client struct {
	baseURL     string // open.larksuite.com — token / user_info APIs
	accountsURL string // accounts.larksuite.com — หน้า authorize
	appID       string
	appSecret   string
	oauthScope  string
	httpClient  *http.Client
}

func NewClient(appID, appSecret, baseURL, accountsURL, oauthScope string) *Client {
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	if accountsURL == "" {
		accountsURL = defaultAccountsURL
	}
	if oauthScope == "" {
		oauthScope = defaultOAuthScope
	}
	return &Client{
		baseURL:     baseURL,
		accountsURL: accountsURL,
		appID:       appID,
		appSecret:   appSecret,
		oauthScope:  oauthScope,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// AuthorizeURL สร้าง URL ตาม Lark Web App OAuth spec
// https://accounts.larksuite.com/open-apis/authen/v1/authorize?client_id=...&redirect_uri=...&scope=...&state=...
func (c *Client) AuthorizeURL(redirectURI, state string) string {
	params := url.Values{}
	params.Set("client_id", c.appID)
	params.Set("redirect_uri", redirectURI)
	params.Set("state", state)
	if c.oauthScope != "" {
		params.Set("scope", c.oauthScope)
	}
	return fmt.Sprintf("%s/open-apis/authen/v1/authorize?%s", c.accountsURL, params.Encode())
}

type appAccessTokenResp struct {
	Code           int    `json:"code"`
	Msg            string `json:"msg"`
	AppAccessToken string `json:"app_access_token"`
	Expire         int    `json:"expire"`
}

// GetAppAccessToken ขั้นตอนที่ 1 — แลก app_id + app_secret เป็น app_access_token
func (c *Client) GetAppAccessToken() (string, error) {
	body := map[string]string{
		"app_id":     c.appID,
		"app_secret": c.appSecret,
	}
	var resp appAccessTokenResp
	if err := c.postJSON("/open-apis/auth/v3/app_access_token/internal", "", body, &resp); err != nil {
		return "", err
	}
	if resp.Code != 0 {
		return "", fmt.Errorf("lark app_access_token: %s (code=%d)", resp.Msg, resp.Code)
	}
	return resp.AppAccessToken, nil
}

// UserToken ผลลัพธ์จากการแลก authorization code
type UserToken struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int
	UserID       string
	Name         string
	AvatarURL    string
	OpenID       string
}

type accessTokenResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
		Name         string `json:"name"`
		AvatarURL    string `json:"avatar_url"`
		OpenID       string `json:"open_id"`
		UserID       string `json:"user_id"`
	} `json:"data"`
}

// ExchangeCode ขั้นตอนที่ 2 — แลก authorization code เป็น user access/refresh token
func (c *Client) ExchangeCode(appAccessToken, code string) (*UserToken, error) {
	body := map[string]string{
		"grant_type": "authorization_code",
		"code":       code,
	}
	var resp accessTokenResp
	if err := c.postJSON("/open-apis/authen/v1/access_token", appAccessToken, body, &resp); err != nil {
		return nil, err
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("lark access_token: %s (code=%d)", resp.Msg, resp.Code)
	}
	userID := resp.Data.UserID
	if userID == "" {
		userID = resp.Data.OpenID
	}
	return &UserToken{
		AccessToken:  resp.Data.AccessToken,
		RefreshToken: resp.Data.RefreshToken,
		ExpiresIn:    resp.Data.ExpiresIn,
		UserID:       userID,
		Name:         resp.Data.Name,
		AvatarURL:    resp.Data.AvatarURL,
		OpenID:       resp.Data.OpenID,
	}, nil
}

type refreshTokenResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
	} `json:"data"`
}

// RefreshAccessToken ใช้ refresh_token แลก token ชุดใหม่ (single-use — ต้องบันทึก refresh ใหม่ทันที)
func (c *Client) RefreshAccessToken(appAccessToken, refreshToken string) (*UserToken, error) {
	body := map[string]string{
		"grant_type":    "refresh_token",
		"refresh_token": refreshToken,
	}
	var resp refreshTokenResp
	if err := c.postJSON("/open-apis/authen/v1/refresh_access_token", appAccessToken, body, &resp); err != nil {
		return nil, err
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("lark refresh_access_token: %s (code=%d)", resp.Msg, resp.Code)
	}
	return &UserToken{
		AccessToken:  resp.Data.AccessToken,
		RefreshToken: resp.Data.RefreshToken,
		ExpiresIn:    resp.Data.ExpiresIn,
	}, nil
}

// UserProfile ข้อมูลผู้ใช้จาก Lark
type UserProfile struct {
	Name      string `json:"name"`
	EnName    string `json:"en_name"`
	AvatarURL string `json:"avatar_url"`
	OpenID    string `json:"open_id"`
	UnionID   string `json:"union_id"`
	Email     string `json:"email"`
	UserID    string `json:"user_id"`
	TenantKey string `json:"tenant_key"`
}

type userInfoResp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data UserProfile `json:"data"`
}

// GetUserInfo ดึง profile ผู้ใช้ด้วย user_access_token
func (c *Client) GetUserInfo(userAccessToken string) (*UserProfile, error) {
	req, err := http.NewRequest(http.MethodGet, c.baseURL+"/open-apis/authen/v1/user_info", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+userAccessToken)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	raw, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var resp userInfoResp
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, err
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("lark user_info: %s (code=%d)", resp.Msg, resp.Code)
	}
	return &resp.Data, nil
}

func (c *Client) postJSON(path, bearer string, body any, out any) error {
	payload, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, c.baseURL+path, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	if bearer != "" {
		req.Header.Set("Authorization", "Bearer "+bearer)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	raw, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(raw, out)
}
