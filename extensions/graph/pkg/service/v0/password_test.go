package svc_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	userv1beta1 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	revactx "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/rgrpc/status"
	"github.com/go-ldap/ldap/v3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	libregraph "github.com/owncloud/libre-graph-api-go"
	"github.com/owncloud/ocis/v2/extensions/graph/mocks"
	"github.com/owncloud/ocis/v2/extensions/graph/pkg/config"
	"github.com/owncloud/ocis/v2/extensions/graph/pkg/config/defaults"
	"github.com/owncloud/ocis/v2/extensions/graph/pkg/identity"
	service "github.com/owncloud/ocis/v2/extensions/graph/pkg/service/v0"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/stretchr/testify/mock"
)

type changePwTest struct {
	desc      string
	currentpw string
	newpw     string
	expected  int
}

var _ = Describe("Users changing their own password", func() {
	var (
		svc             service.Service
		gatewayClient   *mocks.GatewayClient
		httpClient      *mocks.HTTPClient
		ldapClient      *mocks.Client
		ldapConfig      config.LDAP
		identityBackend identity.Backend
		eventsPublisher mocks.Publisher
		ctx             context.Context
		cfg             *config.Config
		user            *userv1beta1.User
		err             error
	)

	JustBeforeEach(func() {
		ctx = context.Background()
		cfg = defaults.FullDefaultConfig()
		cfg.TokenManager.JWTSecret = "loremipsum"

		gatewayClient = &mocks.GatewayClient{}
		ldapClient = mockedLDAPClient()

		ldapConfig = config.LDAP{
			WriteEnabled:             true,
			UserDisplayNameAttribute: "displayName",
			UserNameAttribute:        "uid",
			UserEmailAttribute:       "mail",
			UserIDAttribute:          "ownclouduuid",
			UserSearchScope:          "sub",
			GroupNameAttribute:       "cn",
			GroupIDAttribute:         "ownclouduui",
			GroupSearchScope:         "sub",
		}
		loggger := log.NewLogger()
		identityBackend, err = identity.NewLDAPBackend(ldapClient, ldapConfig, &loggger)
		Expect(err).To(BeNil())

		httpClient = &mocks.HTTPClient{}
		eventsPublisher = mocks.Publisher{}
		svc = service.NewService(
			service.Config(cfg),
			service.WithGatewayClient(gatewayClient),
			service.WithIdentityBackend(identityBackend),
			service.WithHTTPClient(httpClient),
			service.EventsPublisher(&eventsPublisher),
		)
		user = &userv1beta1.User{
			Id: &userv1beta1.UserId{
				OpaqueId: "user",
			},
		}
		ctx = revactx.ContextSetUser(ctx, user)
	})

	It("fails if no user in context", func() {
		r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/me/changePassword", nil)
		rr := httptest.NewRecorder()
		svc.ChangeOwnPassword(rr, r)
		Expect(rr.Code).To(Equal(http.StatusInternalServerError))
	})

	DescribeTable("changing the password",
		func(current string, newpw string, authresult string, expected int) {
			switch authresult {
			case "error":
				gatewayClient.On("Authenticate", mock.Anything, mock.Anything).Return(nil, errors.New("fail"))
			case "deny":
				gatewayClient.On("Authenticate", mock.Anything, mock.Anything).Return(&gateway.AuthenticateResponse{
					Status: status.NewPermissionDenied(ctx, errors.New("wrong password"), "wrong password"),
					Token:  "authtoken",
				}, nil)
			default:
				gatewayClient.On("Authenticate", mock.Anything, mock.Anything).Return(&gateway.AuthenticateResponse{
					Status: status.NewOK(ctx),
					Token:  "authtoken",
				}, nil)
			}
			cpw := libregraph.NewPasswordChange()
			cpw.SetCurrentPassword(current)
			cpw.SetNewPassword(newpw)
			body, _ := json.Marshal(cpw)
			b := bytes.NewBuffer(body)
			r := httptest.NewRequest(http.MethodPost, "/graph/v1.0/me/changePassword", b).WithContext(ctx)
			rr := httptest.NewRecorder()
			svc.ChangeOwnPassword(rr, r)
			Expect(rr.Code).To(Equal(expected))
		},
		Entry("fails when current password is empty", "", "newpassword", "", http.StatusBadRequest),
		Entry("fails when new password is empty", "currentpassword", "", "", http.StatusBadRequest),
		Entry("fails when current and new password are equal", "password", "password", "", http.StatusBadRequest),
		Entry("fails authentication with current password errors", "currentpassword", "newpassword", "error", http.StatusInternalServerError),
		Entry("fails when current password is wrong", "currentpassword", "newpassword", "deny", http.StatusInternalServerError),
		Entry("succeeds when current password is correct", "currentpassword", "newpassword", "", http.StatusNoContent),
	)
})

func mockedLDAPClient() *mocks.Client {
	lm := &mocks.Client{}

	userEntry := ldap.NewEntry("uid=test", map[string][]string{
		"uid":         {"test"},
		"displayName": {"test"},
		"mail":        {"test@example.org"},
	})

	lm.On("Search", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(
			&ldap.SearchResult{Entries: []*ldap.Entry{userEntry}},
			nil)

	mr := ldap.NewModifyRequest("uid=test", nil)
	mr.Changes = []ldap.Change{
		{
			Operation: ldap.ReplaceAttribute,
			Modification: ldap.PartialAttribute{
				Type: "userPassword",
				Vals: []string{"newpassword"},
			},
		},
	}
	lm.On("Modify", mr).Return(nil)
	return lm
}
