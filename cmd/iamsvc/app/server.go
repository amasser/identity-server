package app

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/spf13/cobra"
	"github.com/tierklinik-dobersberg/identity-server/iam"
	"github.com/tierklinik-dobersberg/identity-server/pkg/authn"
	"github.com/tierklinik-dobersberg/identity-server/repos/bbolt"
	"github.com/tierklinik-dobersberg/identity-server/services/user"
)

// NewIAMCommand returns a new IAM server command
func NewIAMCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "iamsvc",
		Run: func(cmd *cobra.Command, args []string) {
			if err := runMain(cmd, args); err != nil {
				fmt.Println(err.Error())
				return
			}
		},
	}

	addHTTPTransportFlags(cmd.Flags())
	addAuthNFlags(cmd)
	addRepoFlags(cmd)

	return cmd
}

func runMain(cmd *cobra.Command, args []string) error {
	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	// Prepare user repository (bbolt)
	var users iam.UserRepository
	{
		path, err := cmd.Flags().GetString("database")
		if err != nil {
			return err
		}
		users, err = bbolt.Open(path)
		if err != nil {
			return err
		}
	}

	// Create authn client service
	var as authn.Service
	{
		cfg, err := getAuthnConfig(cmd)
		if err != nil {
			return err
		}
		as, err = authn.NewService(cfg)
		if err != nil {
			return err
		}
	}

	// Create user management service
	var us user.Service
	{
		us = user.NewService(users, as)
		us = user.NewLoggingService(log.With(logger, "component", "user"), us)
	}

	// Setup HTTP server handlers
	mux := http.NewServeMux()
	httpLogger := log.With(logger, "component", "http")
	{
		mux.Handle("/v1/users/", user.MakeHandler(us, httpLogger))
	}
	http.Handle("/", mux)

	httpAddr, _ := cmd.Flags().GetString("http.listen")
	errs := make(chan error, 2)
	go func() {
		logger.Log("transport", "http", "address", httpAddr, "msg", "listening")
		errs <- http.ListenAndServe(httpAddr, nil)
	}()

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	logger.Log("terminated", <-errs)
	return nil
}
