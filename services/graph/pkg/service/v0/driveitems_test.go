package svc_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"time"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	revactx "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/rgrpc/status"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/cs3org/reva/v2/pkg/utils"
	cs3mocks "github.com/cs3org/reva/v2/tests/cs3mocks/mocks"
	"github.com/go-chi/chi/v5"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	libregraph "github.com/owncloud/libre-graph-api-go"
	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
	"github.com/owncloud/ocis/v2/services/graph/mocks"
	"github.com/owncloud/ocis/v2/services/graph/pkg/config"
	"github.com/owncloud/ocis/v2/services/graph/pkg/config/defaults"
	identitymocks "github.com/owncloud/ocis/v2/services/graph/pkg/identity/mocks"
	service "github.com/owncloud/ocis/v2/services/graph/pkg/service/v0"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type itemsList struct {
	Value []*libregraph.DriveItem
}

var _ = Describe("Driveitems", func() {
	var (
		svc             service.Service
		ctx             context.Context
		cfg             *config.Config
		gatewayClient   *cs3mocks.GatewayAPIClient
		gatewaySelector pool.Selectable[gateway.GatewayAPIClient]
		eventsPublisher mocks.Publisher
		identityBackend *identitymocks.Backend

		rr *httptest.ResponseRecorder

		newGroup *libregraph.Group

		currentUser = &userpb.User{
			Id: &userpb.UserId{
				OpaqueId: "user",
			},
		}
	)

	BeforeEach(func() {
		eventsPublisher.On("Publish", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		rr = httptest.NewRecorder()

		pool.RemoveSelector("GatewaySelector" + "com.owncloud.api.gateway")
		gatewayClient = &cs3mocks.GatewayAPIClient{}
		gatewaySelector = pool.GetSelector[gateway.GatewayAPIClient](
			"GatewaySelector",
			"com.owncloud.api.gateway",
			func(cc *grpc.ClientConn) gateway.GatewayAPIClient {
				return gatewayClient
			},
		)

		identityBackend = &identitymocks.Backend{}
		newGroup = libregraph.NewGroup()
		newGroup.SetMembersodataBind([]string{"/users/user1"})
		newGroup.SetId("group1")

		rr = httptest.NewRecorder()
		ctx = context.Background()

		cfg = defaults.FullDefaultConfig()
		cfg.Identity.LDAP.CACert = "" // skip the startup checks, we don't use LDAP at all in this tests
		cfg.TokenManager.JWTSecret = "loremipsum"
		cfg.Commons = &shared.Commons{}
		cfg.GRPCClientTLS = &shared.GRPCClientTLS{}

		svc, _ = service.NewService(
			service.Config(cfg),
			service.WithGatewaySelector(gatewaySelector),
			service.EventsPublisher(&eventsPublisher),
			service.WithIdentityBackend(identityBackend),
		)
	})

	Describe("GetRootDriveChildren", func() {
		It("handles ListStorageSpaces not found ", func() {
			gatewayClient.On("ListStorageSpaces", mock.Anything, mock.Anything).Return(&provider.ListStorageSpacesResponse{
				Status: status.NewNotFound(ctx, "not found"),
			}, nil)

			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/me/drive/root/children", nil)
			r = r.WithContext(revactx.ContextSetUser(ctx, currentUser))
			svc.GetRootDriveChildren(rr, r)
			Expect(rr.Code).To(Equal(http.StatusNotFound))
		})

		It("handles ListStorageSpaces error", func() {
			gatewayClient.On("ListStorageSpaces", mock.Anything, mock.Anything).Return(&provider.ListStorageSpacesResponse{
				Status: status.NewInternal(ctx, "internal error"),
			}, nil)

			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/me/drive/root/children", nil)
			r = r.WithContext(revactx.ContextSetUser(ctx, currentUser))
			svc.GetRootDriveChildren(rr, r)
			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
		})

		It("handles ListContainer not found", func() {
			gatewayClient.On("ListStorageSpaces", mock.Anything, mock.Anything).Return(&provider.ListStorageSpacesResponse{
				Status:        status.NewOK(ctx),
				StorageSpaces: []*provider.StorageSpace{{Owner: currentUser, Root: &provider.ResourceId{}}},
			}, nil)
			gatewayClient.On("ListContainer", mock.Anything, mock.Anything).Return(&provider.ListContainerResponse{
				Status: status.NewNotFound(ctx, "not found"),
			}, nil)

			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/me/drive/root/children", nil)
			r = r.WithContext(revactx.ContextSetUser(ctx, currentUser))
			svc.GetRootDriveChildren(rr, r)
			Expect(rr.Code).To(Equal(http.StatusNotFound))
		})

		It("handles ListContainer permission denied", func() {
			gatewayClient.On("ListStorageSpaces", mock.Anything, mock.Anything).Return(&provider.ListStorageSpacesResponse{
				Status:        status.NewOK(ctx),
				StorageSpaces: []*provider.StorageSpace{{Owner: currentUser, Root: &provider.ResourceId{}}},
			}, nil)
			gatewayClient.On("ListContainer", mock.Anything, mock.Anything).Return(&provider.ListContainerResponse{
				Status: status.NewPermissionDenied(ctx, errors.New("denied"), "denied"),
			}, nil)

			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/me/drive/root/children", nil)
			r = r.WithContext(revactx.ContextSetUser(ctx, currentUser))
			svc.GetRootDriveChildren(rr, r)
			Expect(rr.Code).To(Equal(http.StatusForbidden))
		})

		It("handles ListContainer error", func() {
			gatewayClient.On("ListStorageSpaces", mock.Anything, mock.Anything).Return(&provider.ListStorageSpacesResponse{
				Status:        status.NewOK(ctx),
				StorageSpaces: []*provider.StorageSpace{{Owner: currentUser, Root: &provider.ResourceId{}}},
			}, nil)
			gatewayClient.On("ListContainer", mock.Anything, mock.Anything).Return(&provider.ListContainerResponse{
				Status: status.NewInternal(ctx, "internal"),
			}, nil)

			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/me/drive/root/children", nil)
			r = r.WithContext(revactx.ContextSetUser(ctx, currentUser))
			svc.GetRootDriveChildren(rr, r)
			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
		})

		It("succeeds", func() {
			mtime := time.Now()
			gatewayClient.On("ListStorageSpaces", mock.Anything, mock.Anything).Return(&provider.ListStorageSpacesResponse{
				Status:        status.NewOK(ctx),
				StorageSpaces: []*provider.StorageSpace{{Owner: currentUser, Root: &provider.ResourceId{}}},
			}, nil)
			gatewayClient.On("ListContainer", mock.Anything, mock.Anything).Return(&provider.ListContainerResponse{
				Status: status.NewOK(ctx),
				Infos: []*provider.ResourceInfo{
					{
						Type:  provider.ResourceType_RESOURCE_TYPE_FILE,
						Id:    &provider.ResourceId{StorageId: "storageid", SpaceId: "spaceid", OpaqueId: "opaqueid"},
						Etag:  "etag",
						Mtime: utils.TimeToTS(mtime),
					},
				},
			}, nil)
			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/me/drive/root/children", nil)
			r = r.WithContext(revactx.ContextSetUser(ctx, currentUser))
			svc.GetRootDriveChildren(rr, r)
			Expect(rr.Code).To(Equal(http.StatusOK))
			data, err := io.ReadAll(rr.Body)
			Expect(err).ToNot(HaveOccurred())

			res := itemsList{}

			err = json.Unmarshal(data, &res)
			Expect(err).ToNot(HaveOccurred())

			Expect(len(res.Value)).To(Equal(1))
			Expect(res.Value[0].GetLastModifiedDateTime().Equal(mtime)).To(BeTrue())
			Expect(res.Value[0].GetETag()).To(Equal("etag"))
			Expect(res.Value[0].GetId()).To(Equal("storageid$spaceid!opaqueid"))
		})
	})

	Describe("GetDriveItemChildren", func() {
		It("handles ListContainer not found", func() {
			gatewayClient.On("ListContainer", mock.Anything, mock.Anything).Return(&provider.ListContainerResponse{
				Status: status.NewNotFound(ctx, "not found"),
			}, nil)

			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/drives/storageid$spaceid/items/storageid$spaceid!nodeid/children", nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("driveID", "storageid$spaceid")
			rctx.URLParams.Add("driveItemID", "storageid$spaceid!nodeid")
			r = r.WithContext(context.WithValue(revactx.ContextSetUser(ctx, currentUser), chi.RouteCtxKey, rctx))
			svc.GetDriveItemChildren(rr, r)
			Expect(rr.Code).To(Equal(http.StatusNotFound))
		})

		It("handles ListContainer permission denied as not found", func() {
			gatewayClient.On("ListContainer", mock.Anything, mock.Anything).Return(&provider.ListContainerResponse{
				Status: status.NewPermissionDenied(ctx, errors.New("denied"), "denied"),
			}, nil)

			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/drives/storageid$spaceid/items/storageid$spaceid!nodeid/children", nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("driveID", "storageid$spaceid")
			rctx.URLParams.Add("driveItemID", "storageid$spaceid!nodeid")
			r = r.WithContext(context.WithValue(revactx.ContextSetUser(ctx, currentUser), chi.RouteCtxKey, rctx))
			svc.GetDriveItemChildren(rr, r)
			Expect(rr.Code).To(Equal(http.StatusNotFound))
		})

		It("handles ListContainer error", func() {
			gatewayClient.On("ListContainer", mock.Anything, mock.Anything).Return(&provider.ListContainerResponse{
				Status: status.NewInternal(ctx, "internal"),
			}, nil)

			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/drives/storageid$spaceid/items/storageid$spaceid!nodeid/children", nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("driveID", "storageid$spaceid")
			rctx.URLParams.Add("driveItemID", "storageid$spaceid!nodeid")
			r = r.WithContext(context.WithValue(revactx.ContextSetUser(ctx, currentUser), chi.RouteCtxKey, rctx))
			svc.GetDriveItemChildren(rr, r)
			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
		})

		Context("it succeeds", func() {
			var (
				r     *http.Request
				mtime = time.Now()
			)

			BeforeEach(func() {
				r = httptest.NewRequest(http.MethodGet, "/graph/v1.0/drives/storageid$spaceid/items/storageid$spaceid!nodeid/children", nil)
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("driveID", "storageid$spaceid")
				rctx.URLParams.Add("driveItemID", "storageid$spaceid!nodeid")
				r = r.WithContext(context.WithValue(revactx.ContextSetUser(ctx, currentUser), chi.RouteCtxKey, rctx))
			})

			assertItemsList := func(length int) itemsList {
				svc.GetDriveItemChildren(rr, r)
				Expect(rr.Code).To(Equal(http.StatusOK))
				data, err := io.ReadAll(rr.Body)
				Expect(err).ToNot(HaveOccurred())

				res := itemsList{}

				err = json.Unmarshal(data, &res)
				Expect(err).ToNot(HaveOccurred())

				Expect(len(res.Value)).To(Equal(1))
				Expect(res.Value[0].GetLastModifiedDateTime().Equal(mtime)).To(BeTrue())
				Expect(res.Value[0].GetETag()).To(Equal("etag"))
				Expect(res.Value[0].GetId()).To(Equal("storageid$spaceid!opaqueid"))
				Expect(res.Value[0].GetId()).To(Equal("storageid$spaceid!opaqueid"))

				return res
			}

			It("returns a generic file", func() {
				gatewayClient.On("ListContainer", mock.Anything, mock.Anything).Return(&provider.ListContainerResponse{
					Status: status.NewOK(ctx),
					Infos: []*provider.ResourceInfo{
						{
							Type:              provider.ResourceType_RESOURCE_TYPE_FILE,
							Id:                &provider.ResourceId{StorageId: "storageid", SpaceId: "spaceid", OpaqueId: "opaqueid"},
							Etag:              "etag",
							Mtime:             utils.TimeToTS(mtime),
							ArbitraryMetadata: nil,
						},
					},
				}, nil)

				res := assertItemsList(1)
				Expect(res.Value[0].Audio).To(BeNil())
			})

			It("returns the audio facet if metadata is available", func() {
				gatewayClient.On("ListContainer", mock.Anything, mock.Anything).Return(&provider.ListContainerResponse{
					Status: status.NewOK(ctx),
					Infos: []*provider.ResourceInfo{
						{
							Type:     provider.ResourceType_RESOURCE_TYPE_FILE,
							Id:       &provider.ResourceId{StorageId: "storageid", SpaceId: "spaceid", OpaqueId: "opaqueid"},
							Etag:     "etag",
							Mtime:    utils.TimeToTS(mtime),
							MimeType: "audio/mpeg",
							ArbitraryMetadata: &provider.ArbitraryMetadata{
								Metadata: map[string]string{
									"libre.graph.audio.album":             "Some Album",
									"libre.graph.audio.albumArtist":       "Some AlbumArtist",
									"libre.graph.audio.artist":            "Some Artist",
									"libre.graph.audio.bitrate":           "192",
									"libre.graph.audio.composers":         "Some Composers",
									"libre.graph.audio.copyright":         "Some Copyright",
									"libre.graph.audio.disc":              "2",
									"libre.graph.audio.discCount":         "5",
									"libre.graph.audio.duration":          "225000",
									"libre.graph.audio.genre":             "Some Genre",
									"libre.graph.audio.hasDrm":            "false",
									"libre.graph.audio.isVariableBitrate": "true",
									"libre.graph.audio.title":             "Some Title",
									"libre.graph.audio.track":             "6",
									"libre.graph.audio.trackCount":        "9",
									"libre.graph.audio.year":              "1994",
								},
							},
						},
					},
				}, nil)

				res := assertItemsList(1)
				audio := res.Value[0].Audio

				Expect(audio).ToNot(BeNil())
				Expect(audio.GetAlbum()).To(Equal("Some Album"))
				Expect(audio.GetAlbumArtist()).To(Equal("Some AlbumArtist"))
				Expect(audio.GetArtist()).To(Equal("Some Artist"))
				Expect(audio.GetBitrate()).To(Equal(int64(192)))
				Expect(audio.GetComposers()).To(Equal("Some Composers"))
				Expect(audio.GetCopyright()).To(Equal("Some Copyright"))
				Expect(audio.GetDisc()).To(Equal(int32(2)))
				Expect(audio.GetDiscCount()).To(Equal(int32(5)))
				Expect(audio.GetDuration()).To(Equal(int64(225000)))
				Expect(audio.GetGenre()).To(Equal("Some Genre"))
				Expect(audio.GetHasDrm()).To(Equal(false))
				Expect(audio.GetIsVariableBitrate()).To(Equal(true))
				Expect(audio.GetTitle()).To(Equal("Some Title"))
				Expect(audio.GetTrack()).To(Equal(int32(6)))
				Expect(audio.GetTrackCount()).To(Equal(int32(9)))
				Expect(audio.GetYear()).To(Equal(int32(1994)))
			})
		})
	})
})
