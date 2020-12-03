package tdam

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"gopkg.in/headzoo/surf.v1"
)

var tokenEndpoint string = "https://api.tdameritrade.com/v1/oauth2/token"

var tdamToken *TokenResponse

// TODO: this sets the token globally!
// eventually, this will need to be stored in a user session
// so multiple people can use the app with their own authentication
func (c *Client) SetToken(token *TokenResponse) error {
	// store the access token (bearer token) in memory,
	// write the refresh token to disk to use on subsequent runs
	tdamToken = token

	configDir, err := getConfigDir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(configDir, os.ModeDir&os.ModePerm); err != nil {
		return err
	}

	f, err := os.Create(path.Join(configDir, "tdam_refresh"))
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write([]byte(token.RefreshToken))
	return err
}

var authMutex sync.Mutex

// gets the active access token if available
// if not, looks for refresh token on disk and attempts to refresh
// if not, push through oauth flow
func (c *Client) TDAMToken() (string, error) {
	if c == nil {
		return "", fmt.Errorf("can't get a token with a nil client!")
	}
	// synchronizing this so we don't infinitely spawn browser windows
	authMutex.Lock()
	defer authMutex.Unlock()
	if tdamToken != nil && tdamToken.AccessExpiry.After(time.Now()) {
		return tdamToken.AccessToken, nil
	}
	//fmt.Printf("Access token missing or expired, checking refresh\n")
	refresh, err := getStoredRefreshToken()
	if err != nil || refresh == "" {
		fmt.Printf("Error fetching refresh token\n")

		fmt.Println("1")
		err := c.slogThroughOauthFlow()
		fmt.Println("2")
		if err != nil {
			fmt.Println(err)
			return "", err
		}
	}
	tdamToken, err = c.refreshToken(refresh)
	if err != nil {
		fmt.Printf("Error refreshing token: %s\n", err)
		fmt.Printf("To obtain a new token, visit %s\n", c.TdamAuthURL())
		//exec.Command("open", tdamAuthURL()).Run()
		//time.Sleep(10 * time.Second)
		return "", err
	}

	return tdamToken.AccessToken, nil
}

func (c *Client) TdamAuthURL() string {
	return fmt.Sprintf("https://auth.tdameritrade.com/oauth?client_id=%s@AMER.OAUTHAP&response_type=code&redirect_uri=%s", c.ConsumerKey, "https://localhost:8443/auth")
}

func getStoredRefreshToken() (string, error) {
	envToken := os.Getenv("REFRESH_TOKEN")
	if envToken != "" {
		return envToken, nil
	}
	configDir, err := getConfigDir()
	if err != nil {
		return "", err
	}
	f, err := os.Open(path.Join(configDir, "tdam_refresh"))
	if err != nil {
		return "", err
	}
	b, err := ioutil.ReadAll(f)
	return string(b), err
}

func getConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return path.Join(homeDir, ".tdam"), nil
}

type TokenResponse struct {
	AccessToken   string    `json:"access_token"`
	AccessExpiry  time.Time `json:"-"`
	RefreshToken  string    `json:"refresh_token"`
	RefreshExpiry time.Time `json:"-"`
	TokenType     string    `json:"token_type"`
	Error         string    `json:"error,omitempty"`

	ExpiresIn             int `json:"expires_in"`
	RefreshTokenExpiresIn int `json:"refresh_token_expires_in"`
}

func (c *Client) AuthHandler(w http.ResponseWriter, req *http.Request) {
	code := req.URL.Query().Get("code")

	token, err := c.getToken(code)
	if err != nil {
		log.Printf("Error getting token: %v\n", err)
	}

	w.Header().Set("Content-Type", "text/plain")
	if err := c.SetToken(token); err != nil {
		fmt.Printf("error setting token: %v\n", err)
		fmt.Fprintf(w, "error setting token: %v\n", err)
	} else {
		fmt.Fprintf(w, "token acquired! have fun:  %v", token)
	}
}

func (c *Client) getToken(code string) (*TokenResponse, error) {
	transport := &http.Transport{TLSClientConfig: &tls.Config{}}
	client := &http.Client{Transport: transport}

	form := url.Values{
		"grant_type":   []string{"authorization_code"},
		"access_type":  []string{"offline"},
		"code":         []string{code},
		"client_id":    []string{c.ConsumerKey},
		"redirect_uri": []string{"https://localhost:8443/auth"},
	}

	req, err := http.NewRequest("POST", tokenEndpoint, bytes.NewBuffer([]byte(form.Encode())))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	var token TokenResponse
	if err = json.NewDecoder(resp.Body).Decode(&token); err != nil {
		fmt.Printf("error decoding token response: %v\n", err)
		return nil, err
	}

	token.AccessExpiry = time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)
	token.RefreshExpiry = time.Now().Add(time.Duration(token.RefreshTokenExpiresIn) * time.Second)

	return &token, err
}

func (c *Client) refreshToken(code string) (*TokenResponse, error) {
	if c == nil {
		return nil, fmt.Errorf("can't get a token with a nil client!")
	}
	transport := &http.Transport{TLSClientConfig: &tls.Config{}}
	client := &http.Client{Transport: transport}

	form := url.Values{
		"grant_type":    []string{"refresh_token"},
		"access_type":   []string{"offline"},
		"refresh_token": []string{code},
		"client_id":     []string{c.ConsumerKey},
		"redirect_uri":  []string{"https://localhost:8443/auth"},
	}

	req, err := http.NewRequest("POST", tokenEndpoint, bytes.NewBuffer([]byte(form.Encode())))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	var token TokenResponse
	if err = json.NewDecoder(resp.Body).Decode(&token); err != nil {
		fmt.Printf("error decoding token response: %v\n", err)
		return nil, err
	}

	if token.Error != "" {
		return &token, fmt.Errorf("%s", token.Error)
	}

	token.AccessExpiry = time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)
	token.RefreshExpiry = time.Now().Add(time.Duration(token.RefreshTokenExpiresIn) * time.Second)

	return &token, err
}

func (c *Client) slogThroughOauthFlow() error {
	username := os.Getenv("TDAM_USERNAME")
	password := os.Getenv("TDAM_PASSWORD")
	if username == "" || password == "" {
		return fmt.Errorf("Must set TDAM_USERNAME and TDAM_PASSWORD for automatic oauth\n")
	}
	fmt.Printf("hang on, let's try to do this ourselves\n")

	bow := surf.NewBrowser()
	err := bow.Open(c.TdamAuthURL())
	if err != nil {
		return err
	}
	fmt.Println(bow.Title())

	title, err := bow.Find("div.title>h1>p").First().Html()
	if err != nil {
		return err
	}
	fmt.Printf("step 1 title: %#v\n", title)
	// at this point, title should be "Secure Log-in"
	// this title is wrapped in a <p> but 2fa pages aren't

	// Log in to the site.
	fm, err := bow.Form("#authform")
	if err != nil {
		return err
	}
	fm.Input("su_username", username)
	fm.Input("su_password", password)
	if err := fm.Submit(); err != nil {
		return err
	}
	_ = bow.Body() // the below doesn't seem to work unless we read the body.  Is there a better way?
	//fmt.Println(bow.Body())

	title, err = bow.Find("div.title>h1").First().Html()
	if err != nil {
		return err
	}
	fmt.Printf("step 2 title: %#v\n", title)
	// here, we either are at a 2fa page or the final page
	// 2fa would be "Get Code via Text Message"

	if title == "Get Code via Text Message" {
		fmt.Printf("lets try to use the security question\n")
		fm, err := bow.Form("#authform")
		if err != nil {
			return err
		}
		if err := fm.Click("init_secretquestion"); err != nil {
			return err
		}
		_ = bow.Body()
		//fmt.Println(bow.Body())
		title, err := bow.Find("div.title>h1").First().Html()
		if err != nil {
			return err
		}
		fmt.Printf("step 2b title: %#v\n", title)
		// here, we should have the title "Answer Security Question'"
		desc, err := bow.Find("div.description>p").Last().Html()
		if err != nil {
			return err
		}
		fmt.Printf("step 2b desc: %#v\n", desc)

		env := ""
		if strings.Contains(desc, "What is your paternal grandmother&#39;s first name?") {
			env = "TDAM_SQ_ANSWER_1"
		}
		if strings.Contains(desc, "What is your maternal grandmother&#39;s first name?") {
			env = "TDAM_SQ_ANSWER_2"
		}
		if strings.Contains(desc, "What was the name of the town your grandmother lived in?") {
			env = "TDAM_SQ_ANSWER_3"
		}
		if strings.Contains(desc, "What is your best friend&#39;s first name?") {
			env = "TDAM_SQ_ANSWER_4"
		}
		if env == "" {
			return fmt.Errorf("Unknown security question: %#v", desc)
		}
		sqAnswer := os.Getenv(env)
		if sqAnswer == "" {
			return fmt.Errorf("Please provide the answer to the following question in %s: %#v", env, desc)
		}

		// ideally we should check the security question but i've only gotten one so far
		fm, err = bow.Form("#authform")
		if err != nil {
			return err
		}
		fm.Input("su_secretquestion", sqAnswer)
		if err := fm.Submit(); err != nil {
			return err
		}
		_ = bow.Body()
		// fmt.Println(bow.Body())
	}

	// at this point, we should be at "TD Ameritrade Authorization"
	// which shows us the scope and asks us to accept
	title, err = bow.Find("div.title>h1").First().Html()
	if err != nil {
		return err
	}
	fmt.Printf("step 3 title: %#v\n", title)

	// set the transport to ignore https cert validation
	bow.SetTransport(&http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	})

	if err := fm.Submit(); err != nil {
		return err
	}
	fmt.Println(bow.Body()) // this should come from our callback server

	return nil
}
