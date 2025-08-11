package subsonic

import (
	"context"
	"zene/core/database"
	"zene/core/logger"
	"zene/core/types"
)

func GetPopulatedSubsonicResponse(ctx context.Context, withError bool) types.SubsonicResponse {
	latestVersion, err := database.GetLatestVersion(ctx)
	if err != nil {
		logger.Printf("Failed to get latest version: %v", err)
		return types.SubsonicResponse{
			SubsonicResponse: types.SubsonicStandard{
				Status: "error",
				Error: &types.SubsonicError{
					Code:    types.ErrorGeneric,
					Message: "Failed to get latest version",
				},
			},
		}
	}

	response := types.SubsonicResponse{
		SubsonicResponse: types.SubsonicStandard{
			Status:        "ok",
			Version:       latestVersion.SubsonicApiVersion,
			Type:          "zene",
			ServerVersion: latestVersion.ServerVersion,
			OpenSubsonic:  true,
			Xmlns:         "http://subsonic.org/restapi",
		},
	}

	if withError {
		response.SubsonicResponse.Status = "error"
		response.SubsonicResponse.Error = &types.SubsonicError{
			Code:    types.ErrorGeneric,
			Message: "An error occurred",
		}
	}
	return response
}
