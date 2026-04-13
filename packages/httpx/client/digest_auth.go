package client

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// digestAuth holds the credentials for HTTP Digest authentication.
type digestAuth struct {
	username string
	password string
}

// digestChallenge holds the parsed WWW-Authenticate challenge parameters.
type digestChallenge struct {
	realm     string
	nonce     string
	opaque    string
	qop       string
	algorithm string
}

// parseDigestChallenge parses the WWW-Authenticate header value for Digest auth.
func parseDigestChallenge(header string) *digestChallenge {
	if !strings.HasPrefix(header, "Digest ") {
		return nil
	}

	challenge := &digestChallenge{}
	params := header[len("Digest "):]

	for _, param := range strings.Split(params, ",") {
		param = strings.TrimSpace(param)
		parts := strings.SplitN(param, "=", 2)

		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.Trim(strings.TrimSpace(parts[1]), `"`)

		switch key {
		case "realm":
			challenge.realm = value
		case "nonce":
			challenge.nonce = value
		case "opaque":
			challenge.opaque = value
		case "qop":
			challenge.qop = value
		case "algorithm":
			challenge.algorithm = value
		}
	}

	return challenge
}

// computeDigestResponse computes the Digest authentication response hash.
func computeDigestResponse(d *digestAuth, c *digestChallenge, method, uri, cnonce, nc string) string {
	ha1 := md5Hash(d.username + ":" + c.realm + ":" + d.password)

	ha2 := md5Hash(method + ":" + uri)

	if c.qop == "auth" || c.qop == "auth-int" {
		return md5Hash(ha1 + ":" + c.nonce + ":" + nc + ":" + cnonce + ":" + c.qop + ":" + ha2)
	}

	return md5Hash(ha1 + ":" + c.nonce + ":" + ha2)
}

// buildDigestHeader builds the Authorization header value for Digest auth.
func buildDigestHeader(d *digestAuth, c *digestChallenge, method, uri string) string {
	cnonce := generateCNonce()
	nc := "00000001"

	response := computeDigestResponse(d, c, method, uri, cnonce, nc)

	parts := []string{
		fmt.Sprintf(`username="%s"`, d.username),
		fmt.Sprintf(`realm="%s"`, c.realm),
		fmt.Sprintf(`nonce="%s"`, c.nonce),
		fmt.Sprintf(`uri="%s"`, uri),
		fmt.Sprintf(`response="%s"`, response),
	}

	if c.qop != "" {
		parts = append(parts,
			fmt.Sprintf(`qop=%s`, c.qop),
			fmt.Sprintf(`nc=%s`, nc),
			fmt.Sprintf(`cnonce="%s"`, cnonce),
		)
	}

	if c.opaque != "" {
		parts = append(parts, fmt.Sprintf(`opaque="%s"`, c.opaque))
	}

	if c.algorithm != "" {
		parts = append(parts, fmt.Sprintf(`algorithm=%s`, c.algorithm))
	}

	return "Digest " + strings.Join(parts, ", ")
}

// digestMiddleware returns a middleware that handles HTTP Digest authentication.
func digestMiddleware(d *digestAuth) Middleware {
	return func(req *http.Request, next RoundTripFunc) (*http.Response, error) {
		resp, err := next(req)

		if err != nil {
			return resp, err
		}

		if resp.StatusCode != http.StatusUnauthorized {
			return resp, nil
		}

		wwwAuth := resp.Header.Get("WWW-Authenticate")

		if !strings.HasPrefix(wwwAuth, "Digest ") {
			return resp, nil
		}

		challenge := parseDigestChallenge(wwwAuth)

		if challenge == nil {
			return resp, nil
		}

		// Close the initial response body.
		if resp.Body != nil {
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
		}

		uri := req.URL.RequestURI()
		authHeader := buildDigestHeader(d, challenge, req.Method, uri)

		req.Header.Set("Authorization", authHeader)

		return next(req)
	}
}

func md5Hash(text string) string {
	h := md5.New()
	h.Write([]byte(text))

	return hex.EncodeToString(h.Sum(nil))
}

func generateCNonce() string {
	b := make([]byte, 8)

	rand.Read(b)

	return hex.EncodeToString(b)
}
