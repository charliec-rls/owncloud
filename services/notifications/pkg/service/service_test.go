package service_test

import (
	"context"
	"time"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/cs3org/reva/v2/pkg/utils"
	cs3mocks "github.com/cs3org/reva/v2/tests/cs3mocks/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	ogrpc "github.com/owncloud/ocis/v2/ocis-pkg/service/grpc"
	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
	settingssvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/settings/v0"
	"github.com/owncloud/ocis/v2/services/graph/pkg/config/defaults"
	"github.com/owncloud/ocis/v2/services/notifications/pkg/channels"
	"github.com/owncloud/ocis/v2/services/notifications/pkg/service"
	"github.com/stretchr/testify/mock"
	"go-micro.dev/v4/client"
	"google.golang.org/grpc"
)

var _ = Describe("Notifications", func() {
	var (
		gatewayClient   *cs3mocks.GatewayAPIClient
		gatewaySelector pool.Selectable[gateway.GatewayAPIClient]
		vs              *settingssvc.MockValueService
		sharer          = &user.User{
			Id: &user.UserId{
				OpaqueId: "sharer",
			},
			Mail:        "sharer@owncloud.com",
			DisplayName: "Dr. S. Harer",
		}
		sharee = &user.User{
			Id: &user.UserId{
				OpaqueId: "sharee",
			},
			Mail:        "sharee@owncloud.com",
			DisplayName: "Eric Expireling",
		}
		resourceid = &provider.ResourceId{
			StorageId: "storageid",
			SpaceId:   "spaceid",
			OpaqueId:  "itemid",
		}
	)

	BeforeEach(func() {
		pool.RemoveSelector("GatewaySelector" + "com.owncloud.api.gateway")
		gatewayClient = &cs3mocks.GatewayAPIClient{}
		gatewaySelector = pool.GetSelector[gateway.GatewayAPIClient](
			"GatewaySelector",
			"com.owncloud.api.gateway",
			func(cc *grpc.ClientConn) gateway.GatewayAPIClient {
				return gatewayClient
			},
		)

		gatewayClient.On("GetUser", mock.Anything, mock.Anything).Return(&user.GetUserResponse{Status: &rpc.Status{Code: rpc.Code_CODE_OK}, User: sharer}, nil).Once()
		gatewayClient.On("GetUser", mock.Anything, mock.Anything).Return(&user.GetUserResponse{Status: &rpc.Status{Code: rpc.Code_CODE_OK}, User: sharee}, nil).Once()
		gatewayClient.On("Authenticate", mock.Anything, mock.Anything).Return(&gateway.AuthenticateResponse{Status: &rpc.Status{Code: rpc.Code_CODE_OK}, User: sharer}, nil)
		gatewayClient.On("Stat", mock.Anything, mock.Anything).Return(&provider.StatResponse{Status: &rpc.Status{Code: rpc.Code_CODE_OK}, Info: &provider.ResourceInfo{Name: "secrets of the board", Space: &provider.StorageSpace{Name: "secret space"}}}, nil)
		vs = &settingssvc.MockValueService{}
		vs.GetValueByUniqueIdentifiersFunc = func(ctx context.Context, req *settingssvc.GetValueByUniqueIdentifiersRequest, opts ...client.CallOption) (*settingssvc.GetValueResponse, error) {
			return nil, nil
		}
	})

	DescribeTable("Sending notifications",
		func(tc testChannel, ev events.Event) {
			cfg := defaults.FullDefaultConfig()
			cfg.GRPCClientTLS = &shared.GRPCClientTLS{}
			_ = ogrpc.Configure(ogrpc.GetClientOptions(cfg.GRPCClientTLS)...)
			ch := make(chan events.Event)
			evts := service.NewEventsNotifier(ch, tc, log.NewLogger(), gatewaySelector, vs, "", "", "")
			go evts.Run()

			ch <- ev
			select {
			case <-tc.done:
				// finished
			case <-time.Tick(3 * time.Second):
				Fail("timeout waiting for notification")
			}
		},

		Entry("Share Created", testChannel{
			expectedReceipients: []string{sharee.GetMail()},
			expectedSubject:     "Dr. S. Harer shared 'secrets of the board' with you",
			expectedTextBody: `Hello Eric Expireling

Dr. S. Harer has shared "secrets of the board" with you.

Click here to view it: files/shares/with-me


---
ownCloud - Store. Share. Work.
https://owncloud.com
`,
			expectedSender: sharer.GetDisplayName(),
			done:           make(chan struct{}),
		}, events.Event{
			Event: events.ShareCreated{
				Sharer:        sharer.GetId(),
				GranteeUserID: sharee.GetId(),
				CTime:         utils.TimeToTS(time.Date(2023, 4, 17, 16, 42, 0, 0, time.UTC)),
				ItemID:        resourceid,
			},
		}),

		Entry("Share Expired", testChannel{
			expectedReceipients: []string{sharee.GetMail()},
			expectedSubject:     "Share to 'secrets of the board' expired at 2023-04-17 16:42:00",
			expectedTextBody: `Hello Eric Expireling,

Your share to secrets of the board has expired at 2023-04-17 16:42:00

Even though this share has been revoked you still might have access through other shares and/or space memberships.


---
ownCloud - Store. Share. Work.
https://owncloud.com
`,
			expectedSender: sharer.GetDisplayName(),
			done:           make(chan struct{}),
		}, events.Event{
			Event: events.ShareExpired{
				ShareOwner:    sharer.GetId(),
				GranteeUserID: sharee.GetId(),
				ExpiredAt:     time.Date(2023, 4, 17, 16, 42, 0, 0, time.UTC),
				ItemID:        resourceid,
			},
		}),

		Entry("Added to Space", testChannel{
			expectedReceipients: []string{sharee.GetMail()},
			expectedSubject:     "Dr. S. Harer invited you to join secret space",
			expectedTextBody: `Hello Eric Expireling,

Dr. S. Harer has invited you to join "secret space".

Click here to view it: f/spaceid


---
ownCloud - Store. Share. Work.
https://owncloud.com
`,
			expectedSender: sharer.GetDisplayName(),
			done:           make(chan struct{}),
		}, events.Event{
			Event: events.SpaceShared{
				Executant:     sharer.GetId(),
				Creator:       sharer.GetId(),
				GranteeUserID: sharee.GetId(),
				ID:            &provider.StorageSpaceId{OpaqueId: "spaceid"},
			},
		}),

		Entry("Removed from Space", testChannel{
			expectedReceipients: []string{sharee.GetMail()},
			expectedSubject:     "Dr. S. Harer removed you from secret space",
			expectedTextBody: `Hello Eric Expireling,

Dr. S. Harer has removed you from "secret space".

You might still have access through your other groups or direct membership.

Click here to check it: f/spaceid


---
ownCloud - Store. Share. Work.
https://owncloud.com
`,
			expectedSender: sharer.GetDisplayName(),
			done:           make(chan struct{}),
		}, events.Event{
			Event: events.SpaceUnshared{
				Executant:     sharer.GetId(),
				GranteeUserID: sharee.GetId(),
				ID:            &provider.StorageSpaceId{OpaqueId: "spaceid"},
			},
		}),

		Entry("Space Expired", testChannel{
			expectedReceipients: []string{sharee.GetMail()},
			expectedSubject:     "Membership of 'secret space' expired at 2023-04-17 16:42:00",
			expectedTextBody: `Hello Eric Expireling,

Your membership of space secret space has expired at 2023-04-17 16:42:00

Even though this membership has expired you still might have access through other shares and/or space memberships


---
ownCloud - Store. Share. Work.
https://owncloud.com
`,
			expectedSender: sharer.GetDisplayName(),
			done:           make(chan struct{}),
		}, events.Event{
			Event: events.SpaceMembershipExpired{
				SpaceOwner:    sharer.GetId(),
				GranteeUserID: sharee.GetId(),
				SpaceID:       &provider.StorageSpaceId{OpaqueId: "spaceid"},
				SpaceName:     "secret space",
				ExpiredAt:     time.Date(2023, 4, 17, 16, 42, 0, 0, time.UTC),
			},
		}),
	)
})

var _ = Describe("Notifications X-Site Scripting", func() {
	var (
		gatewayClient   *cs3mocks.GatewayAPIClient
		gatewaySelector pool.Selectable[gateway.GatewayAPIClient]
		vs              *settingssvc.MockValueService
		sharer          = &user.User{
			Id: &user.UserId{
				OpaqueId: "sharer",
			},
			Mail:        "sharer@owncloud.com",
			DisplayName: "Dr. O'reilly",
		}
		sharee = &user.User{
			Id: &user.UserId{
				OpaqueId: "sharee",
			},
			Mail:        "sharee@owncloud.com",
			DisplayName: "<script>alert('Eric Expireling');</script>",
		}
		resourceid = &provider.ResourceId{
			StorageId: "storageid",
			SpaceId:   "spaceid",
			OpaqueId:  "itemid",
		}
	)

	BeforeEach(func() {
		pool.RemoveSelector("GatewaySelector" + "com.owncloud.api.gateway")
		gatewayClient = &cs3mocks.GatewayAPIClient{}
		gatewaySelector = pool.GetSelector[gateway.GatewayAPIClient](
			"GatewaySelector",
			"com.owncloud.api.gateway",
			func(cc *grpc.ClientConn) gateway.GatewayAPIClient {
				return gatewayClient
			},
		)

		gatewayClient.On("GetUser", mock.Anything, mock.Anything).Return(&user.GetUserResponse{Status: &rpc.Status{Code: rpc.Code_CODE_OK}, User: sharer}, nil).Once()
		gatewayClient.On("GetUser", mock.Anything, mock.Anything).Return(&user.GetUserResponse{Status: &rpc.Status{Code: rpc.Code_CODE_OK}, User: sharee}, nil).Once()
		gatewayClient.On("Authenticate", mock.Anything, mock.Anything).Return(&gateway.AuthenticateResponse{Status: &rpc.Status{Code: rpc.Code_CODE_OK}, User: sharer}, nil)
		gatewayClient.On("Stat", mock.Anything, mock.Anything).Return(&provider.StatResponse{
			Status: &rpc.Status{Code: rpc.Code_CODE_OK},
			Info: &provider.ResourceInfo{
				Name:  "<script>alert('secrets of the board');</script>",
				Space: &provider.StorageSpace{Name: "<script>alert('secret space');</script>"}},
		}, nil)
		vs = &settingssvc.MockValueService{}
		vs.GetValueByUniqueIdentifiersFunc = func(ctx context.Context, req *settingssvc.GetValueByUniqueIdentifiersRequest, opts ...client.CallOption) (*settingssvc.GetValueResponse, error) {
			return nil, nil
		}
	})

	DescribeTable("Sending notifications",
		func(tc testChannel, ev events.Event) {
			cfg := defaults.FullDefaultConfig()
			cfg.GRPCClientTLS = &shared.GRPCClientTLS{}
			_ = ogrpc.Configure(ogrpc.GetClientOptions(cfg.GRPCClientTLS)...)
			ch := make(chan events.Event)
			evts := service.NewEventsNotifier(ch, tc, log.NewLogger(), gatewaySelector, vs, "", "", "")
			go evts.Run()

			ch <- ev
			select {
			case <-tc.done:
				// finished
			case <-time.Tick(3 * time.Second):
				Fail("timeout waiting for notification")
			}
		},

		Entry("Share Created", testChannel{
			expectedReceipients: []string{sharee.GetMail()},
			expectedSubject:     "Dr. O'reilly shared '<script>alert('secrets of the board');</script>' with you",
			expectedTextBody: `Hello <script>alert('Eric Expireling');</script>

Dr. O'reilly has shared "<script>alert('secrets of the board');</script>" with you.

Click here to view it: files/shares/with-me


---
ownCloud - Store. Share. Work.
https://owncloud.com
`,
			expectedHTMLBody: `<!DOCTYPE html>
<html>
<body>
<table cellspacing="0" cellpadding="0" border="0" width="100%">
    <tr>
        <td>
            <table cellspacing="0" cellpadding="0" border="0" width="600px">
                <tr>
                    <td width="20px">&nbsp;</td>
                    <td style="font-weight:normal; font-size:0.8em; line-height:1.2em; font-family:verdana,'arial',sans;">
                        Hello &lt;script&gt;alert(&#39;Eric Expireling&#39;);&lt;/script&gt;
                        <br><br>
                        Dr. O&#39;reilly has shared "&lt;script&gt;alert(&#39;secrets of the board&#39;);&lt;/script&gt;" with you.
                        <br><br>
                        <a href="files/shares/with-me">Click here to view it: files/shares/with-me</a>
                    </td>
                </tr>
                <tr>
                    <td colspan="2">&nbsp;</td>
                </tr>
                <tr>
                    <td width="20px">&nbsp;</td>
                    <td style="font-weight:normal; font-size:0.8em; line-height:1.2em; font-family:verdana,'arial',sans;">
                        <footer>
                            <br>
                            <br>
                            --- <br>
                            ownCloud - Store. Share. Work.<br>
                            <a href="https://owncloud.com">https://owncloud.com</a>
                        </footer>
                    </td>
                </tr>
                <tr>
                    <td colspan="2">&nbsp;</td>
                </tr>
            </table>
        </td>
    </tr>
</table>
</body>
</html>
`,
			expectedSender: sharer.GetDisplayName(),

			done: make(chan struct{}),
		}, events.Event{
			Event: events.ShareCreated{
				Sharer:        sharer.GetId(),
				GranteeUserID: sharee.GetId(),
				CTime:         utils.TimeToTS(time.Date(2023, 4, 17, 16, 42, 0, 0, time.UTC)),
				ItemID:        resourceid,
			},
		}),

		Entry("Added to Space", testChannel{
			expectedReceipients: []string{sharee.GetMail()},
			expectedSubject:     "Dr. O'reilly invited you to join <script>alert('secret space');</script>",
			expectedTextBody: `Hello <script>alert('Eric Expireling');</script>,

Dr. O'reilly has invited you to join "<script>alert('secret space');</script>".

Click here to view it: f/spaceid


---
ownCloud - Store. Share. Work.
https://owncloud.com
`,
			expectedSender: sharer.GetDisplayName(),
			expectedHTMLBody: `<!DOCTYPE html>
<html>
<body>
<table cellspacing="0" cellpadding="0" border="0" width="100%">
    <tr>
        <td>
            <table cellspacing="0" cellpadding="0" border="0" width="600px">
                <tr>
                    <td width="20px">&nbsp;</td>
                    <td style="font-weight:normal; font-size:0.8em; line-height:1.2em; font-family:verdana,'arial',sans;">
                        Hello &lt;script&gt;alert(&#39;Eric Expireling&#39;);&lt;/script&gt;,
                        <br><br>
                        Dr. O&#39;reilly has invited you to join "&lt;script&gt;alert(&#39;secret space&#39;);&lt;/script&gt;".
                        <br><br>
                        <a href="f/spaceid">Click here to view it: f/spaceid</a>
                    </td>
                </tr>
                <tr>
                    <td colspan="2">&nbsp;</td>
                </tr>
                <tr>
                    <td width="20px">&nbsp;</td>
                    <td style="font-weight:normal; font-size:0.8em; line-height:1.2em; font-family:verdana,'arial',sans;">
                        <footer>
                            <br>
                            <br>
                            --- <br>
                            ownCloud - Store. Share. Work.<br>
                            <a href="https://owncloud.com">https://owncloud.com</a>
                        </footer>
                    </td>
                </tr>
                <tr>
                    <td colspan="2">&nbsp;</td>
                </tr>
            </table>
        </td>
    </tr>
</table>
</body>
</html>
`,
			done: make(chan struct{}),
		}, events.Event{
			Event: events.SpaceShared{
				Executant:     sharer.GetId(),
				Creator:       sharer.GetId(),
				GranteeUserID: sharee.GetId(),
				ID:            &provider.StorageSpaceId{OpaqueId: "spaceid"},
			},
		}),
	)
})

// NOTE: This is explictitly not testing the message itself. Should we?
type testChannel struct {
	expectedReceipients []string
	expectedSubject     string
	expectedTextBody    string
	expectedHTMLBody    string
	expectedSender      string
	done                chan struct{}
}

func (tc testChannel) SendMessage(ctx context.Context, m *channels.Message) error {
	defer GinkgoRecover()

	Expect(tc.expectedReceipients).To(Equal(m.Recipient))
	Expect(tc.expectedSubject).To(Equal(m.Subject))
	Expect(tc.expectedTextBody).To(Equal(m.TextBody))
	Expect(tc.expectedSender).To(Equal(m.Sender))
	if tc.expectedHTMLBody != "" {
		Expect(tc.expectedHTMLBody).To(Equal(m.HTMLBody))
	}
	tc.done <- struct{}{}
	return nil
}
