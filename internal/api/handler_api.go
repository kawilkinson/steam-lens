package api

import (
	"net/http"
	"strings"
)

type ApiConfig struct {
	SteamApiKey string
}

func (apicfg *ApiConfig) HandlerGetPlayerSummaries(w http.ResponseWriter, req *http.Request) {
	steamIDs := req.URL.Query().Get("steamids")
	if steamIDs == "" {
		RespondWithError(w, http.StatusBadRequest, "'steamids' parameter is required for getting player summaries", nil)
		return
	}

	playerSummaries, err := apicfg.GetPlayerSummaries(strings.Split(steamIDs, ","))
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Unable to perform API call to Steam GetPlayerSummaries endpoint", err)
		return
	}

	RespondWithJSON(w, http.StatusOK, playerSummaries)
}

func (apicfg *ApiConfig) HandlerGetFriendList(w http.ResponseWriter, req *http.Request) {
	steamID := req.URL.Query().Get("steamid")
	if steamID == "" {
		RespondWithError(w, http.StatusBadRequest, "'steamid' parameter is required for getting friend list", nil)
		return
	}

	friends, err := apicfg.GetFriendList(steamID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Unable to perform API call to Steam GetFriendList endpoint", err)
		return
	}

	friendIDs := []string{}
	for _, friend := range friends.Friends {
		friendIDs = append(friendIDs, friend.SteamID)
	}

	summaries, err := apicfg.GetPlayerSummaries(friendIDs)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Unable to get player summaries of friends", err)
		return
	}

	RespondWithJSON(w, http.StatusOK, summaries)
}

func (apicfg *ApiConfig) HandlerGetOwnedGames(w http.ResponseWriter, req *http.Request) {
	steamID := req.URL.Query().Get("steamid")
	if steamID == "" {
		RespondWithError(w, http.StatusBadRequest, "'steamid' parameter is required for getting player owned games", nil)
		return
	}

	ownedGames, err := apicfg.GetOwnedGames(steamID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Unable to perform API call to Steam GetOwnedGames endpoint", err)
		return
	}

	RespondWithJSON(w, http.StatusOK, ownedGames)
}

func (apicfg *ApiConfig) HandlerGetPlayerAchievements(w http.ResponseWriter, req *http.Request) {
	steamID := req.URL.Query().Get("steamid")
	if steamID == "" {
		RespondWithError(w, http.StatusBadRequest, "'steamid' parameter is required for getting player achievements", nil)
		return
	}

	appID := req.URL.Query().Get("appid")
	if appID == "" {
		RespondWithError(w, http.StatusBadRequest, "'appid' parameter is required for getting player achievements", nil)
		return
	}

	achievements, err := apicfg.GetPlayerAchievements(steamID, appID)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Unable to perform API call to Steam GetPlayerAchievements endpoint", err)
		return
	}

	RespondWithJSON(w, http.StatusOK, achievements)
}

func (apicfg *ApiConfig) HandlerCompareOwnedGames(w http.ResponseWriter, req *http.Request) {
	userID := req.URL.Query().Get("userID")
	friendID := req.URL.Query().Get("friendID")
	if userID == "" {
		RespondWithError(w, http.StatusBadRequest, "User Steam ID is required", nil)
		return
	}
	if friendID == "" {
		RespondWithError(w, http.StatusBadRequest, "Friend Steam ID is required", nil)
	}

	listGames := false
	listGamesQuery := req.URL.Query().Get("listGames")
	if listGamesQuery == "true" {
		listGames = true
	}

	userGames, err := apicfg.GetOwnedGames(userID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Unable to perform API calls to Steam GetOwnedGames endpoint for user", err)
		return
	}

	friendGames, err := apicfg.GetOwnedGames(friendID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Unable to perform API calls to Steam GetOwnedGames endpoint for friend", err)
		return
	}

	result := userGames.CompareOwnedGames(friendGames, listGames)
	
	RespondWithJSON(w, http.StatusOK, result)
}
