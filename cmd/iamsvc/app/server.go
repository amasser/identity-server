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
	"github.com/tierklinik-dobersberg/identity-server/repos/inmem"
	"github.com/tierklinik-dobersberg/identity-server/services/group"
	"github.com/tierklinik-dobersberg/identity-server/services/policy"
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
	//logger = logrusLogger.NewLogrusLogger(logrus.New())
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	dbPath, err := cmd.Flags().GetString("database")
	if err != nil {
		return err
	}

	var db *bbolt.Database
	if dbPath != ":memory:" {
		db, err = bbolt.OpenWithLogger(dbPath, log.With(logger, "component", "bbolt"))
		if err != nil {
			return err
		}
	}

	var users iam.UserRepository
	{
		if db == nil {
			users = inmem.NewUserRepository()
		} else {
			users = db.UserRepo()
		}
	}

	var groups iam.GroupRepository
	{
		if db == nil {
			groups = inmem.NewGroupRepository()
		} else {
			groups = db.GroupRepo()
		}
	}

	var members iam.MembershipRepository
	{
		if db == nil {
			members = inmem.NewMembershipRepository()
		} else {
			members = db.MembershipRepo()
		}
	}

	var policies iam.PolicyRepository
	{
		if db == nil {
			policies = inmem.NewPolicyRepository()
		} else {
			policies = db.PolicyRepo()
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

	jwtTokenExtractor := func(token string) (subject string, err error) {
		return "", nil
	}

	// User management service
	var us user.Service
	{
		us = user.NewService(users, as)
		us = user.NewLoggingService(log.With(logger, "component", "user"), us)
	}

	//  Group management service
	var gs group.Service
	{
		groupLogger := log.With(logger, "component", "group")
		gs = group.NewService(us, groups, members, groupLogger)
		gs = group.NewLoggingService(gs, groupLogger)
	}

	// Policy management service
	var ps policy.Service
	{
		ps = policy.NewService(policies)
		ps = policy.NewLoggingService(log.With(logger, "component", "policy"), ps)
	}

	// Setup HTTP server handlers
	mux := http.NewServeMux()
	httpLogger := log.With(logger, "component", "http")
	{
		mux.Handle("/v1/users/", user.MakeHandler(us, jwtTokenExtractor, httpLogger))
		mux.Handle("/v1/groups/", group.MakeHandler(gs, jwtTokenExtractor, httpLogger))
		mux.Handle("/v1/policies/", policy.MakeHandler(ps, jwtTokenExtractor, httpLogger))
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
