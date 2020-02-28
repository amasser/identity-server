package authn

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/transport/http"
	"github.com/tierklinik-dobersberg/identity-server/pkg/enforcer"
)

type contextKey string

// ContextKeyJWT is used by NewAuthenticator to add the parsed
// and validated JWT token to the request context.
const ContextKeyJWT contextKey = "authn:jwt-token"

// SubjectExtractorFunc extracts and validates the JWT user subject from
// the token.
type SubjectExtractorFunc func(token string) (string, error)

// NewAuthenticator returns an endpoint.Middleware that extracts and
// validates an AuthN JWT access token. The user URN is added to the
// request context.
func NewAuthenticator(fn SubjectExtractorFunc) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			val := ctx.Value(http.ContextKeyRequestAuthorization)
			if val == nil {
				return nil, errors.New("not authorized") // TODO(ppacher): return approriate error
			}
			bearer := val.(string)

			if !strings.HasPrefix(bearer, "Bearer") {
				return nil, errors.New("invalid authorization key")
			}

			idToken := strings.Replace(bearer, "Bearer ", "", 1)
			accountID, err := fn(idToken)
			if err != nil {
				return nil, err
			}

			// We use ParseUnverified here because thze SubjectExtractorFunc is expected
			// to verify the token. We cannot do any verification here because we just
			// don't know enough about the token to parse.
			token, _, err := new(jwt.Parser).ParseUnverified(idToken, &jwt.StandardClaims{})
			if err != nil {
				return nil, err
			}

			// We add the whole JWT as well as an identity-server UserURN to
			// the request context.
			ctx = context.WithValue(ctx, ContextKeyJWT, token)
			ctx = enforcer.WithSubject(ctx, fmt.Sprintf("urn:iam::user/%s", accountID))

			return next(ctx, request)
		}
	}
}

// LogFields returns authn related fields that might be useful in
// log statements.
func LogFields(ctx context.Context) []interface{} {
	if val := ctx.Value(ContextKeyJWT); val != nil {
		token := val.(*jwt.Token)
		claims := token.Claims.(*jwt.StandardClaims)

		return []interface{}{
			"subject", claims.Subject,
			"issuer", claims.Issuer,
			"audience", claims.Audience,
			"expiresAt", claims.ExpiresAt,
		}
	}

	return nil
}
