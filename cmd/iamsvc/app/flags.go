package app

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/tierklinik-dobersberg/identity-server/pkg/authn"
	"gopkg.in/square/go-jose.v2/jwt"
)

func addHTTPTransportFlags(flags *pflag.FlagSet) {
	flags.StringP("http.listen", "l", ":8080", "Address to listen for HTTP requests")
}

func addRepoFlags(cmd *cobra.Command) {
	flags := cmd.Flags()

	flags.String("database", "./iam.db", "Path to bbolt database")
}

func addAuthNFlags(cmd *cobra.Command) {
	flags := cmd.Flags()

	flags.StringP("authn.server", "a", "http://localhost:8090", "Address of the AuthN-server endpoint")
	flags.String("authn.user", "hello", "Username for private authn-server endpoints")
	flags.String("authn.password", "world", "Password for private authn-server endpoints")
	flags.String("authn.issuer", "", "Issuer for the authn-server endpoint. Defaults to the value of --authn.server")
	flags.String("authn.audience", "", "The audience for JWT access tokens")
	flags.Bool("disable-authorization", false, "Disable policy based authorization. Only use for bootstrapping or testing. DO NOT USE IN PRODUCTION.")

	cmd.MarkFlagRequired("authn.audience")
}

func getAuthnConfig(cmd *cobra.Command) (authn.Config, error) {
	f := cmd.Flags()

	var (
		audience, _ = f.GetString("authn.audience")
		server, _   = f.GetString("authn.server")
		password, _ = f.GetString("authn.password")
		user, _     = f.GetString("authn.user")
		issuer, _   = f.GetString("authn.issuer")
	)

	if issuer == "" {
		s, err := url.Parse(server)
		if err != nil {
			return authn.Config{}, err
		}

		issuer = fmt.Sprintf("%s://%s", s.Scheme, s.Host)
	}

	if audience == "" {
		httpListen, _ := f.GetString("http.listen")
		s, err := url.Parse(httpListen)
		if err != nil {
			return authn.Config{}, err
		}

		audience = s.Hostname()
	}

	return authn.Config{
		Audiences:          jwt.Audience{audience},
		PrivateBaseAddress: server,
		Password:           password,
		Username:           user,
		Issuer:             issuer,
	}, nil
}
