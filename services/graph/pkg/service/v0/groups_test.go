package svc_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	userv1beta1 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	revactx "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/go-chi/chi/v5"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/test-go/testify/mock"

	libregraph "github.com/owncloud/libre-graph-api-go"
	ogrpc "github.com/owncloud/ocis/v2/ocis-pkg/service/grpc"
	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
	"github.com/owncloud/ocis/v2/services/graph/mocks"
	"github.com/owncloud/ocis/v2/services/graph/pkg/config"
	"github.com/owncloud/ocis/v2/services/graph/pkg/config/defaults"
	identitymocks "github.com/owncloud/ocis/v2/services/graph/pkg/identity/mocks"
	service "github.com/owncloud/ocis/v2/services/graph/pkg/service/v0"
	"github.com/owncloud/ocis/v2/services/graph/pkg/service/v0/errorcode"
)

type groupList struct {
	Value []*libregraph.Group
}

var _ = Describe("Groups", func() {
	var (
		svc             service.Service
		ctx             context.Context
		cfg             *config.Config
		gatewayClient   *mocks.GatewayClient
		eventsPublisher mocks.Publisher
		identityBackend *identitymocks.Backend

		rr *httptest.ResponseRecorder

		newGroup    *libregraph.Group
		currentUser = &userv1beta1.User{
			Id: &userv1beta1.UserId{
				OpaqueId: "user",
			},
		}
	)

	BeforeEach(func() {
		eventsPublisher.On("Publish", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		identityBackend = &identitymocks.Backend{}
		gatewayClient = &mocks.GatewayClient{}
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

		_ = ogrpc.Configure(ogrpc.GetClientOptions(cfg.GRPCClientTLS)...)
		svc, _ = service.NewService(
			service.Config(cfg),
			service.WithGatewayClient(gatewayClient),
			service.EventsPublisher(&eventsPublisher),
			service.WithIdentityBackend(identityBackend),
		)
	})

	Describe("GetGroups", func() {
		It("handles invalid ODATA parameters", func() {
			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/groups?§foo=bar", nil)
			svc.GetGroups(rr, r)

			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})

		It("handles invalid sorting queries", func() {
			identityBackend.On("GetGroups", ctx, mock.Anything).Return([]*libregraph.Group{newGroup}, nil)

			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/groups?$orderby=invalid", nil)
			svc.GetGroups(rr, r)

			Expect(rr.Code).To(Equal(http.StatusBadRequest))
			data, err := io.ReadAll(rr.Body)
			Expect(err).ToNot(HaveOccurred())

			odataerr := libregraph.OdataError{}
			err = json.Unmarshal(data, &odataerr)
			Expect(err).ToNot(HaveOccurred())
			Expect(odataerr.Error.Code).To(Equal("invalidRequest"))
		})

		It("handles unknown backend errors", func() {
			identityBackend.On("GetGroups", ctx, mock.Anything).Return(nil, errors.New("failed"))

			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/groups", nil)
			svc.GetGroups(rr, r)
			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
			data, err := io.ReadAll(rr.Body)
			Expect(err).ToNot(HaveOccurred())

			odataerr := libregraph.OdataError{}
			err = json.Unmarshal(data, &odataerr)
			Expect(err).ToNot(HaveOccurred())
			Expect(odataerr.Error.Code).To(Equal("generalException"))
		})

		It("handles backend errors", func() {
			identityBackend.On("GetGroups", ctx, mock.Anything).Return(nil, errorcode.New(errorcode.AccessDenied, "access denied"))

			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/groups", nil)
			svc.GetGroups(rr, r)

			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
			data, err := io.ReadAll(rr.Body)
			Expect(err).ToNot(HaveOccurred())

			odataerr := libregraph.OdataError{}
			err = json.Unmarshal(data, &odataerr)
			Expect(err).ToNot(HaveOccurred())
			Expect(odataerr.Error.Code).To(Equal("accessDenied"))
		})

		It("renders an empty list of groups", func() {
			identityBackend.On("GetGroups", ctx, mock.Anything).Return([]*libregraph.Group{}, nil)

			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/groups", nil)
			svc.GetGroups(rr, r)

			Expect(rr.Code).To(Equal(http.StatusOK))
			data, err := io.ReadAll(rr.Body)
			Expect(err).ToNot(HaveOccurred())

			res := service.ListResponse{}
			err = json.Unmarshal(data, &res)
			Expect(err).ToNot(HaveOccurred())
			Expect(res.Value).To(Equal([]interface{}{}))
		})

		It("renders a list of groups", func() {
			identityBackend.On("GetGroups", ctx, mock.Anything).Return([]*libregraph.Group{newGroup}, nil)

			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/groups", nil)
			svc.GetGroups(rr, r)

			Expect(rr.Code).To(Equal(http.StatusOK))
			data, err := io.ReadAll(rr.Body)
			Expect(err).ToNot(HaveOccurred())

			res := groupList{}
			err = json.Unmarshal(data, &res)
			Expect(err).ToNot(HaveOccurred())

			Expect(len(res.Value)).To(Equal(1))
			Expect(res.Value[0].GetId()).To(Equal("group1"))
		})
	})

	Describe("GetGroup", func() {
		It("handles missing or empty group id", func() {
			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/groups", nil)
			svc.GetGroup(rr, r)

			Expect(rr.Code).To(Equal(http.StatusBadRequest))

			r = httptest.NewRequest(http.MethodGet, "/graph/v1.0/groups", nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("groupID", "")
			r = r.WithContext(context.WithValue(revactx.ContextSetUser(ctx, nil), chi.RouteCtxKey, rctx))
			svc.GetGroup(rr, r)

			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})

		Context("with an existing group", func() {
			BeforeEach(func() {
				identityBackend.On("GetGroup", mock.Anything, mock.Anything, mock.Anything).Return(newGroup, nil)
			})

			It("gets the group", func() {
				r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/groups/"+*newGroup.Id, nil)
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("groupID", *newGroup.Id)
				r = r.WithContext(context.WithValue(revactx.ContextSetUser(ctx, nil), chi.RouteCtxKey, rctx))

				svc.GetGroup(rr, r)

				Expect(rr.Code).To(Equal(http.StatusOK))
			})
		})
	})

	Describe("PostGroup", func() {
		It("handles invalid body", func() {
			r := httptest.NewRequest(http.MethodPost, "/graph/v1.0/groups/", bytes.NewBufferString("{invalid"))

			svc.PostGroup(rr, r)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})

		It("handles missing display name", func() {
			newGroup = libregraph.NewGroup()
			newGroup.SetId("disallowed")
			newGroup.SetMembersodataBind([]string{"/non-users/user"})
			newGroupJson, err := json.Marshal(newGroup)
			Expect(err).ToNot(HaveOccurred())

			r := httptest.NewRequest(http.MethodPost, "/graph/v1.0/groups/", bytes.NewBuffer(newGroupJson))

			svc.PostGroup(rr, r)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})

		It("disallows group create ids", func() {
			newGroup = libregraph.NewGroup()
			newGroup.SetId("disallowed")
			newGroup.SetDisplayName("New Group")
			newGroup.SetMembersodataBind([]string{"/non-users/user"})
			newGroupJson, err := json.Marshal(newGroup)
			Expect(err).ToNot(HaveOccurred())

			r := httptest.NewRequest(http.MethodPost, "/graph/v1.0/groups/", bytes.NewBuffer(newGroupJson))

			svc.PostGroup(rr, r)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})

		It("creates the group", func() {
			newGroup = libregraph.NewGroup()
			newGroup.SetDisplayName("New Group")
			newGroupJson, err := json.Marshal(newGroup)
			Expect(err).ToNot(HaveOccurred())

			identityBackend.On("CreateGroup", mock.Anything, mock.Anything).Return(newGroup, nil)

			r := httptest.NewRequest(http.MethodPost, "/graph/v1.0/groups/", bytes.NewBuffer(newGroupJson))

			svc.PostGroup(rr, r)

			Expect(rr.Code).To(Equal(http.StatusOK))
		})
	})
	Describe("PatchGroup", func() {
		It("handles invalid body", func() {
			r := httptest.NewRequest(http.MethodPatch, "/graph/v1.0/groups/", bytes.NewBufferString("{invalid"))
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("groupID", *newGroup.Id)
			r = r.WithContext(context.WithValue(revactx.ContextSetUser(ctx, currentUser), chi.RouteCtxKey, rctx))
			svc.PatchGroup(rr, r)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})

		It("handles missing or empty group id", func() {
			r := httptest.NewRequest(http.MethodPatch, "/graph/v1.0/groups", nil)
			svc.PatchGroup(rr, r)

			Expect(rr.Code).To(Equal(http.StatusBadRequest))

			r = httptest.NewRequest(http.MethodPatch, "/graph/v1.0/groups", nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("groupID", "")
			r = r.WithContext(context.WithValue(revactx.ContextSetUser(ctx, currentUser), chi.RouteCtxKey, rctx))
			svc.PatchGroup(rr, r)

			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})

		Context("with an existing group", func() {
			BeforeEach(func() {
				identityBackend.On("GetGroup", mock.Anything, mock.Anything, mock.Anything).Return(newGroup, nil)
			})

			It("fails when the number of users is exceeded - spec says 20 max", func() {
				updatedGroup := libregraph.NewGroup()
				updatedGroup.SetDisplayName("group1 updated")
				updatedGroup.SetMembersodataBind([]string{
					"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18",
					"19", "20", "21",
				})
				updatedGroupJson, err := json.Marshal(updatedGroup)
				Expect(err).ToNot(HaveOccurred())

				r := httptest.NewRequest(http.MethodPatch, "/graph/v1.0/groups", bytes.NewBuffer(updatedGroupJson))
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("groupID", *newGroup.Id)
				r = r.WithContext(context.WithValue(revactx.ContextSetUser(ctx, currentUser), chi.RouteCtxKey, rctx))
				svc.PatchGroup(rr, r)

				resp, err := ioutil.ReadAll(rr.Body)
				Expect(err).ToNot(HaveOccurred())

				Expect(string(resp)).To(ContainSubstring("Request is limited to 20"))
				Expect(rr.Code).To(Equal(http.StatusBadRequest))
			})

			It("succeeds when the number of users is over 20 but the limit is raised to 21", func() {
				updatedGroup := libregraph.NewGroup()
				updatedGroup.SetDisplayName("group1 updated")
				updatedGroup.SetMembersodataBind([]string{
					"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18",
					"19", "20", "21",
				})
				updatedGroupJson, err := json.Marshal(updatedGroup)
				Expect(err).ToNot(HaveOccurred())

				cfg.API.GroupMembersPatchLimit = 21
				svc, _ = service.NewService(
					service.Config(cfg),
					service.WithGatewayClient(gatewayClient),
					service.EventsPublisher(&eventsPublisher),
					service.WithIdentityBackend(identityBackend),
				)

				r := httptest.NewRequest(http.MethodPatch, "/graph/v1.0/groups", bytes.NewBuffer(updatedGroupJson))
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("groupID", *newGroup.Id)
				r = r.WithContext(context.WithValue(revactx.ContextSetUser(ctx, currentUser), chi.RouteCtxKey, rctx))
				svc.PatchGroup(rr, r)

				resp, err := ioutil.ReadAll(rr.Body)
				Expect(err).ToNot(HaveOccurred())

				Expect(string(resp)).To(ContainSubstring("Error parsing member@odata.bind values"))
				Expect(rr.Code).To(Equal(http.StatusBadRequest))
			})

			It("fails on invalid user refs", func() {
				updatedGroup := libregraph.NewGroup()
				updatedGroup.SetDisplayName("group1 updated")
				updatedGroup.SetMembersodataBind([]string{"invalid"})
				updatedGroupJson, err := json.Marshal(updatedGroup)
				Expect(err).ToNot(HaveOccurred())

				r := httptest.NewRequest(http.MethodPatch, "/graph/v1.0/groups", bytes.NewBuffer(updatedGroupJson))
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("groupID", *newGroup.Id)
				r = r.WithContext(context.WithValue(revactx.ContextSetUser(ctx, currentUser), chi.RouteCtxKey, rctx))
				svc.PatchGroup(rr, r)

				Expect(rr.Code).To(Equal(http.StatusBadRequest))
			})

			It("fails when the adding non-users users", func() {
				updatedGroup := libregraph.NewGroup()
				updatedGroup.SetDisplayName("group1 updated")
				updatedGroup.SetMembersodataBind([]string{"/non-users/user1"})
				updatedGroupJson, err := json.Marshal(updatedGroup)
				Expect(err).ToNot(HaveOccurred())

				r := httptest.NewRequest(http.MethodPatch, "/graph/v1.0/groups", bytes.NewBuffer(updatedGroupJson))
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("groupID", *newGroup.Id)
				r = r.WithContext(context.WithValue(revactx.ContextSetUser(ctx, currentUser), chi.RouteCtxKey, rctx))
				svc.PatchGroup(rr, r)

				Expect(rr.Code).To(Equal(http.StatusBadRequest))
			})

			It("adds members to the group", func() {
				identityBackend.On("AddMembersToGroup", mock.Anything, mock.Anything, mock.Anything).Return(nil)

				updatedGroup := libregraph.NewGroup()
				updatedGroup.SetDisplayName("group1 updated")
				updatedGroup.SetMembersodataBind([]string{"/users/user1"})
				updatedGroupJson, err := json.Marshal(updatedGroup)
				Expect(err).ToNot(HaveOccurred())

				r := httptest.NewRequest(http.MethodPatch, "/graph/v1.0/groups", bytes.NewBuffer(updatedGroupJson))
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("groupID", *newGroup.Id)
				r = r.WithContext(context.WithValue(revactx.ContextSetUser(ctx, currentUser), chi.RouteCtxKey, rctx))
				svc.PatchGroup(rr, r)

				Expect(rr.Code).To(Equal(http.StatusNoContent))
				identityBackend.AssertNumberOfCalls(GinkgoT(), "AddMembersToGroup", 1)
			})
		})
	})

	Describe("DeleteGroup", func() {
		Context("with an existing group", func() {
			BeforeEach(func() {
				identityBackend.On("GetGroup", mock.Anything, mock.Anything, mock.Anything).Return(newGroup, nil)
			})
		})

		It("deletes the group", func() {
			identityBackend.On("DeleteGroup", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			r := httptest.NewRequest(http.MethodPatch, "/graph/v1.0/groups", nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("groupID", *newGroup.Id)
			r = r.WithContext(context.WithValue(revactx.ContextSetUser(ctx, currentUser), chi.RouteCtxKey, rctx))
			svc.DeleteGroup(rr, r)

			Expect(rr.Code).To(Equal(http.StatusNoContent))
			identityBackend.AssertNumberOfCalls(GinkgoT(), "DeleteGroup", 1)
			eventsPublisher.AssertNumberOfCalls(GinkgoT(), "Publish", 1)
		})
	})

	Describe("GetGroupMembers", func() {
		It("gets the list of members", func() {
			user := libregraph.NewUser()
			user.SetId("user")
			identityBackend.On("GetGroupMembers", mock.Anything, mock.Anything, mock.Anything).Return([]*libregraph.User{user}, nil)

			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/groups/{groupID}/members", nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("groupID", *newGroup.Id)
			r = r.WithContext(context.WithValue(revactx.ContextSetUser(ctx, currentUser), chi.RouteCtxKey, rctx))
			svc.GetGroupMembers(rr, r)
			Expect(rr.Code).To(Equal(http.StatusOK))

			data, err := io.ReadAll(rr.Body)
			Expect(err).ToNot(HaveOccurred())

			var members []*libregraph.User
			err = json.Unmarshal(data, &members)
			Expect(err).ToNot(HaveOccurred())

			Expect(len(members)).To(Equal(1))
			Expect(members[0].GetId()).To(Equal("user"))
		})
	})

	Describe("PostGroupMembers", func() {
		It("fails on invalid body", func() {
			r := httptest.NewRequest(http.MethodPost, "/graph/v1.0/groups/{groupID}/members", bytes.NewBufferString("{invalid"))
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("groupID", *newGroup.Id)
			r = r.WithContext(context.WithValue(revactx.ContextSetUser(ctx, currentUser), chi.RouteCtxKey, rctx))
			svc.PostGroupMember(rr, r)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})

		It("fails on missing member refs", func() {
			member := libregraph.NewMemberReference()
			data, err := json.Marshal(member)
			Expect(err).ToNot(HaveOccurred())

			r := httptest.NewRequest(http.MethodPost, "/graph/v1.0/groups/{groupID}/members", bytes.NewBuffer(data))
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("groupID", *newGroup.Id)
			r = r.WithContext(context.WithValue(revactx.ContextSetUser(ctx, currentUser), chi.RouteCtxKey, rctx))
			svc.PostGroupMember(rr, r)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})

		It("fails on invalid member refs", func() {
			member := libregraph.NewMemberReference()
			member.SetOdataId("/invalidtype/user")
			data, err := json.Marshal(member)
			Expect(err).ToNot(HaveOccurred())

			r := httptest.NewRequest(http.MethodPost, "/graph/v1.0/groups/{groupID}/members", bytes.NewBuffer(data))
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("groupID", *newGroup.Id)
			r = r.WithContext(context.WithValue(revactx.ContextSetUser(ctx, currentUser), chi.RouteCtxKey, rctx))
			svc.PostGroupMember(rr, r)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})

		It("adds a new member", func() {
			member := libregraph.NewMemberReference()
			member.SetOdataId("/users/user")
			data, err := json.Marshal(member)
			Expect(err).ToNot(HaveOccurred())
			identityBackend.On("AddMembersToGroup", mock.Anything, mock.Anything, mock.Anything).Return(nil)

			r := httptest.NewRequest(http.MethodPost, "/graph/v1.0/groups/{groupID}/members", bytes.NewBuffer(data))
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("groupID", *newGroup.Id)
			r = r.WithContext(context.WithValue(revactx.ContextSetUser(ctx, currentUser), chi.RouteCtxKey, rctx))
			svc.PostGroupMember(rr, r)
			Expect(rr.Code).To(Equal(http.StatusNoContent))

			identityBackend.AssertNumberOfCalls(GinkgoT(), "AddMembersToGroup", 1)
		})
	})

	Describe("DeleteGroupMembers", func() {
		It("handles missing or empty member id", func() {
			r := httptest.NewRequest(http.MethodDelete, "/graph/v1.0/groups/{groupID}/members/{memberID}/$ref", nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("groupID", *newGroup.Id)
			r = r.WithContext(context.WithValue(revactx.ContextSetUser(ctx, currentUser), chi.RouteCtxKey, rctx))
			svc.DeleteGroupMember(rr, r)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})
		It("handles missing or empty member id", func() {
			r := httptest.NewRequest(http.MethodDelete, "/graph/v1.0/groups/{groupID}/members/{memberID}/$ref", nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("memberID", "/users/user")
			r = r.WithContext(context.WithValue(revactx.ContextSetUser(ctx, currentUser), chi.RouteCtxKey, rctx))
			svc.DeleteGroupMember(rr, r)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})

		It("deletes members", func() {
			identityBackend.On("RemoveMemberFromGroup", mock.Anything, mock.Anything, mock.Anything).Return(nil)

			r := httptest.NewRequest(http.MethodDelete, "/graph/v1.0/groups/{groupID}/members/{memberID}/$ref", nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("groupID", *newGroup.Id)
			rctx.URLParams.Add("memberID", "/users/user1")
			r = r.WithContext(context.WithValue(revactx.ContextSetUser(ctx, currentUser), chi.RouteCtxKey, rctx))
			svc.DeleteGroupMember(rr, r)
			Expect(rr.Code).To(Equal(http.StatusNoContent))

			identityBackend.AssertNumberOfCalls(GinkgoT(), "RemoveMemberFromGroup", 1)
		})
	})
})
