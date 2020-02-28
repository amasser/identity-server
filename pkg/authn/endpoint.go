package authn

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/transport/http"
	"github.com/tierklinik-dobersberg/identity-server/pkg/enforcer"
	"gopkg.in/square/go-jose.v2/jwt"
)

type contextKey string

// ContextKeyJWTClaims is used by NewAuthenticator to add the parsed
// and validated JWT token to the request context.
const ContextKeyJWTClaims contextKey = "authn:jwt-token"

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

			token, err := jwt.ParseSigned(idToken)
			if err != nil {
				return nil, err
			}

			// We use UnsafeClaimsWithoutVerification here because the SubjectExtractorFunc is expected
			// to verify the token. We cannot do any verification here because we just
			// don't know enough about the token to parse.
			var claims jwt.Claims
			token.UnsafeClaimsWithoutVerification(&claims)

			// fn should verify the token here ...
			accountID, err := fn(idToken)
			if err != nil {
				return nil, fmt.Errorf("issuer=%q: %w", claims.Issuer, err)
			}

			// We add the whole JWT as well as an identity-server UserURN to
			// the request context.
			ctx = context.WithValue(ctx, ContextKeyJWTClaims, claims)
			ctx = enforcer.WithSubject(ctx, fmt.Sprintf("urn:iam::user/%s", accountID))

			return next(ctx, request)
		}
	}
}

// LogFields returns authn related fields that might be useful in
// log statements.
func LogFields(ctx context.Context) []interface{} {
	if val := ctx.Value(ContextKeyJWTClaims); val != nil {
		claims := val.(jwt.Claims)

		return []interface{}{
			"subject", claims.Subject,
			"issuer", claims.Issuer,
			"audience", claims.Audience,
			"expiresAt", claims.Expiry,
		}
	}

	return nil
}
