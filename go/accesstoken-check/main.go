package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/oauth2/google"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

const ServiceAccount = "tg-headless-cms-develop@appspot.gserviceaccount.com"

func init() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := appengine.NewContext(r)
		ts := google.AppEngineTokenSource(ctx)
		token, err := ts.Token()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		log.Infof(ctx, "%#v", token)

		header, _ := json.Marshal(map[string]string{
			"typ": "JWT",
			"alg": "RS256",
		})
		now := time.Now().Unix()
		payload, _ := json.Marshal(map[string]interface{}{
			"iat":   now,
			"exp":   now + int64(1*time.Hour),
			"iss":   ServiceAccount,
			"sub":   ServiceAccount,
			"aud":   "echo.endpoints.sample.google.com",
			"email": ServiceAccount,
		})
		headerAndPayload := fmt.Sprintf(
			"%s.%s",
			base64.RawURLEncoding.EncodeToString(header),
			base64.RawURLEncoding.EncodeToString(payload),
		)

		_, sig, err := appengine.SignBytes(ctx, []byte(headerAndPayload))
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		jwt := fmt.Sprintf("%s.%s", headerAndPayload, base64.RawURLEncoding.EncodeToString(sig))
		io.WriteString(w, fmt.Sprintf("%s", jwt))
	})
}
