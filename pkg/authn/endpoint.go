package authn

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/transport/http"
	"github.com/tierklinik-dobersberg/identity-server/pkg/enforcer"
)

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

			// Add the accountID as a UserURN to the request context.
			ctx = enforcer.WithSubject(ctx, fmt.Sprintf("urn:iam::user/%s", accountID))

			return next(ctx, request)
		}
	}
}
