package svc

import (
	"context"
	"net/http"

	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	collaboration "github.com/cs3org/go-cs3apis/cs3/sharing/collaboration/v1beta1"
	link "github.com/cs3org/go-cs3apis/cs3/sharing/link/v1beta1"
	"github.com/cs3org/reva/v2/pkg/share"
	"github.com/cs3org/reva/v2/pkg/storagespace"
	"github.com/go-chi/render"
	libregraph "github.com/owncloud/libre-graph-api-go"
	"github.com/owncloud/ocis/v2/services/graph/pkg/service/v0/errorcode"
)

type driveItemsByResourceID map[string]libregraph.DriveItem

// GetSharedByMe implements the Service interface (/me/drives/sharedByMe endpoint)
func (g Graph) GetSharedByMe(w http.ResponseWriter, r *http.Request) {
	g.logger.Debug().Msg("Calling GetRootDriveChildren")
	ctx := r.Context()

	driveItems := make(driveItemsByResourceID)
	var err error
	driveItems, err = g.listUserShares(ctx, driveItems)
	if err != nil {
		errorcode.RenderError(w, r, err)
		return
	}

	_, err = g.listPublicShares(ctx, driveItems)
	if err != nil {
		errorcode.RenderError(w, r, err)
		return
	}

	res := make([]libregraph.DriveItem, 0, len(driveItems))
	for _, v := range driveItems {
		res = append(res, v)
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, &ListResponse{Value: res})
}

func (g Graph) listUserShares(ctx context.Context, driveItems driveItemsByResourceID) (driveItemsByResourceID, error) {
	gatewayClient, err := g.gatewaySelector.Next()
	if err != nil {
		g.logger.Error().Err(err).Msg("could not select next gateway client")
		return driveItems, errorcode.New(errorcode.GeneralException, err.Error())
	}

	filters := []*collaboration.Filter{
		share.UserGranteeFilter(),
		share.GroupGranteeFilter(),
	}
	lsUserSharesRequest := collaboration.ListSharesRequest{
		Filters: filters,
	}

	lsUserSharesResponse, err := gatewayClient.ListShares(ctx, &lsUserSharesRequest)
	if err != nil {
		return driveItems, errorcode.New(errorcode.GeneralException, err.Error())
	}
	if statusCode := lsUserSharesResponse.GetStatus().GetCode(); statusCode != rpc.Code_CODE_OK {
		return driveItems, errorcode.New(cs3StatusToErrCode(statusCode), lsUserSharesResponse.Status.Message)
	}
	driveItems, err = g.cs3UserSharesToDriveItems(ctx, lsUserSharesResponse.Shares, driveItems)
	if err != nil {
		return driveItems, errorcode.New(errorcode.GeneralException, err.Error())
	}
	return driveItems, nil
}

func (g Graph) listPublicShares(ctx context.Context, driveItems driveItemsByResourceID) (driveItemsByResourceID, error) {

	gatewayClient, err := g.gatewaySelector.Next()
	if err != nil {
		g.logger.Error().Err(err).Msg("could not select next gateway client")
		return driveItems, errorcode.New(errorcode.GeneralException, err.Error())
	}

	filters := []*link.ListPublicSharesRequest_Filter{}

	req := link.ListPublicSharesRequest{
		Filters: filters,
	}

	lsPublicSharesResponse, err := gatewayClient.ListPublicShares(ctx, &req)
	if err != nil {
		return driveItems, errorcode.New(errorcode.GeneralException, err.Error())
	}
	if statusCode := lsPublicSharesResponse.GetStatus().GetCode(); statusCode != rpc.Code_CODE_OK {
		return driveItems, errorcode.New(cs3StatusToErrCode(statusCode), lsPublicSharesResponse.Status.Message)
	}
	driveItems, err = g.cs3PublicSharesToDriveItems(ctx, lsPublicSharesResponse.Share, driveItems)
	if err != nil {
		return driveItems, errorcode.New(errorcode.GeneralException, err.Error())
	}
	return driveItems, nil

}

func (g Graph) cs3UserSharesToDriveItems(ctx context.Context, shares []*collaboration.Share, driveItems driveItemsByResourceID) (driveItemsByResourceID, error) {
	for _, s := range shares {
		g.logger.Debug().Interface("CS3 UserShare", s).Msg("Got Share")
		resIDStr := storagespace.FormatResourceID(*s.ResourceId)
		item, ok := driveItems[resIDStr]
		if !ok {
			item = libregraph.DriveItem{
				Id: libregraph.PtrString(resIDStr),
			}
		}
		driveItems[resIDStr] = item
	}

	return driveItems, nil
}

func (g Graph) cs3PublicSharesToDriveItems(ctx context.Context, shares []*link.PublicShare, driveItems driveItemsByResourceID) (driveItemsByResourceID, error) {
	for _, s := range shares {
		g.logger.Debug().Interface("CS3 PublicShare", s).Msg("Got Share")
		resIDStr := storagespace.FormatResourceID(*s.ResourceId)
		item, ok := driveItems[resIDStr]
		if !ok {
			item = libregraph.DriveItem{
				Id: libregraph.PtrString(resIDStr),
			}
		}
		driveItems[resIDStr] = item
	}

	return driveItems, nil
}

func cs3StatusToErrCode(code rpc.Code) (errcode errorcode.ErrorCode) {
	switch code {
	case rpc.Code_CODE_UNAUTHENTICATED:
		errcode = errorcode.Unauthenticated
	case rpc.Code_CODE_PERMISSION_DENIED:
		errcode = errorcode.AccessDenied
	case rpc.Code_CODE_NOT_FOUND:
		errcode = errorcode.ItemNotFound
	default:
		errcode = errorcode.GeneralException
	}
	return errcode
}
