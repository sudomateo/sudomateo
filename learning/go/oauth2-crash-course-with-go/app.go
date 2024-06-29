package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

// UserInfo represents information retrieved from GitHub's `/user` API.
type UserInfo struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
}

// App is our example web application that can speak OAuth2.
type App struct {
	OAuthConfig *oauth2.Config
	Logger      *slog.Logger
	Template    *template.Template

	// In-memory session storage. A production application would store this in
	// a persistent location such as a database.
	AccessToken  string
	RefreshToken string
	UserInfo     *UserInfo
}

// Root renders the home page that's user to sign in or show user information.
func (a *App) Root(w http.ResponseWriter, r *http.Request) {
	// We don't have an access token for this user. Render the sign-in page.
	if a.AccessToken == "" {
		w.WriteHeader(http.StatusOK)
		if err := a.Template.Execute(w, a); err != nil {
			a.Logger.Error("failed executing template", "error", err)
		}
		return
	}

    // At this point we have an access token for the user so we can retrieve
    // the user's information to personalize their experience.

	req, err := http.NewRequest(http.MethodGet, "https://api.github.com/user", nil)
	if err != nil {
		a.Logger.Error("failed creating request", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.AccessToken))
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		a.Logger.Error("failed retrieving user details", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var userInfo UserInfo

	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		a.Logger.Error("failed decoding user details", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Store the user's information to in-memory session storage.
	a.UserInfo = &userInfo

	w.WriteHeader(http.StatusOK)
	if err := a.Template.Execute(w, a); err != nil {
		a.Logger.Error("failed executing template", "error", err)
		return
	}
}

// OAuthCallback handles OAuth2 callback requests and exchanges given
// information for an access token and refresh token.
func (a *App) OAuthCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	token, err := a.OAuthConfig.Exchange(r.Context(), code)
	if err != nil {
		a.Logger.Error("failed oauth exchange", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	a.Logger.Info("completed oauth exchange",
		"token_type", token.Type(),
		"expiration", token.Expiry,
	)

	// Store the access token and refresh token in in-memory session storage.
	a.AccessToken = token.AccessToken
	a.RefreshToken = token.RefreshToken

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func LoggerMiddleware(logger *slog.Logger, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		receivedTime := time.Now()

		logger.Info("request received",
			"method", r.Method,
			"path", r.URL.Path,
		)

		handler(w, r)

		logger.Info("request complete",
			"method", r.Method,
			"path", r.URL.Path,
			"duration_ms", time.Since(receivedTime).Milliseconds(),
		)
	}
}
