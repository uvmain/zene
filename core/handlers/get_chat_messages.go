package handlers

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"
	"zene/core/database"
	"zene/core/logger"
	"zene/core/net"
	"zene/core/subsonic"
	"zene/core/types"
)

func HandleGetChatMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		errorString := fmt.Sprintf("Unsupported method: %s", r.Method)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, errorString, "")
		return
	}

	ctx := r.Context()

	timeSinceParam := r.FormValue("since")
	var timeSince int64
	if timeSinceParam != "" {
		var err error
		timeSince, err = strconv.ParseInt(timeSinceParam, 10, 64)
		if err != nil {
			errorString := fmt.Sprintf("Invalid since parameter: %s", timeSinceParam)
			net.WriteSubsonicError(w, r, types.ErrorGeneric, errorString, "")
			return
		}
	} else {
		timeSince = 0
	}

	chats, err := database.GetChats(ctx, timeSince)
	if err != nil {
		logger.Printf("Error fetching chats from database: %v", err)
		net.WriteSubsonicError(w, r, types.ErrorGeneric, "Failed to fetch chat messages", "")
		return
	}

	logger.Printf("Fetched %d chats since %d", len(chats), timeSince)

	response := types.SubsonicChatMessagesResponse{}
	stdRes := subsonic.GetPopulatedSubsonicResponse(ctx, false)

	response.SubsonicResponse.XMLName = stdRes.SubsonicResponse.XMLName
	response.SubsonicResponse.Xmlns = stdRes.SubsonicResponse.Xmlns
	response.SubsonicResponse.Status = stdRes.SubsonicResponse.Status
	response.SubsonicResponse.Version = stdRes.SubsonicResponse.Version
	response.SubsonicResponse.Type = stdRes.SubsonicResponse.Type
	response.SubsonicResponse.ServerVersion = stdRes.SubsonicResponse.ServerVersion
	response.SubsonicResponse.OpenSubsonic = stdRes.SubsonicResponse.OpenSubsonic

	response.SubsonicResponse.ChatMessages = &types.ChatMessages{
		ChatMessage: chats,
	}

	format := r.FormValue("f")
	if format == "json" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	} else {
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>`))
		xml.NewEncoder(w).Encode(response.SubsonicResponse)
	}
}
