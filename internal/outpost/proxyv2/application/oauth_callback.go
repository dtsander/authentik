package application

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"go.uber.org/zap"
	"goauthentik.io/internal/config"
	"goauthentik.io/internal/outpost/proxyv2/constants"
	"golang.org/x/oauth2"
)

func (a *Application) handleAuthCallback(rw http.ResponseWriter, r *http.Request) {
	state := a.stateFromRequest(r)
	if state == nil {
		a.log.Warn("invalid state")
		a.redirect(rw, r)
		return
	}
	claims, err := a.redeemCallback(r.URL, r.Context())
	if err != nil {
		a.log.Warn("failed to redeem code", zap.Error(err))
		a.redirect(rw, r)
		return
	}
	s, err := a.sessions.Get(r, a.SessionName())
	if err != nil {
		a.log.Debug("failed to get session", zap.Error(err), config.Trace())
	}
	s.Options.MaxAge = int(time.Until(time.Unix(int64(claims.Exp), 0)).Seconds())
	s.Values[constants.SessionClaims] = &claims
	err = s.Save(r, rw)
	if err != nil {
		a.log.Warn("failed to save session", zap.Error(err))
		rw.WriteHeader(400)
		return
	}
	a.redirect(rw, r)
}

func (a *Application) redeemCallback(u *url.URL, c context.Context) (*Claims, error) {
	code := u.Query().Get("code")
	if code == "" {
		return nil, fmt.Errorf("blank code")
	}

	ctx := context.WithValue(c, oauth2.HTTPClient, a.publicHostHTTPClient)
	// Verify state and errors.
	oauth2Token, err := a.oauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	jwt := oauth2Token.AccessToken
	a.log.Debug("access_token", config.Trace(), zap.String("jwt", jwt))

	// Parse and verify ID Token payload.
	idToken, err := a.tokenVerifier.Verify(ctx, jwt)
	if err != nil {
		return nil, err
	}

	// Extract custom claims
	var claims *Claims
	if err := idToken.Claims(&claims); err != nil {
		return nil, err
	}
	if claims.Proxy == nil {
		claims.Proxy = &ProxyClaims{}
	}
	claims.RawToken = jwt
	return claims, nil
}
