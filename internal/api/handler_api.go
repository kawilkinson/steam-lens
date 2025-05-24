package api

import (
	"log"
	"net/http"
	"sort"
	"strings"
	"sync"
)

type ApiConfig struct {
	SteamApiKey string

	PlayerCache       Cache[Player]
	FriendListCache   Cache[FriendList]
	OwnedGamesCache   Cache[OwnedGames]
	AchievementsCache Cache[ConvertedPlayerAchievements]
}

func (apicfg *ApiConfig) HandlerGetPlayerSummaries(w http.ResponseWriter, req *http.Request) {
	steamIDs := req.URL.Query().Get("steamIDs")
	if steamIDs == "" {
		RespondWithError(w, http.StatusBadRequest, "'steamIDs' parameter is required for getting player summaries", nil)
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
	steamID := req.URL.Query().Get("steamID")
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
	steamID := req.URL.Query().Get("steamID")
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
	steamID := req.URL.Query().Get("steamID")
	if steamID == "" {
		RespondWithError(w, http.StatusBadRequest, "'steamid' parameter is required for getting player achievements", nil)
		return
	}

	appID := req.URL.Query().Get("appID")
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

func (apicfg *ApiConfig) HandlerCompareAchievements(w http.ResponseWriter, req *http.Request) {
	playerID := req.URL.Query().Get("userID")
	friendID := req.URL.Query().Get("friendID")
	appID := req.URL.Query().Get("appID")

	log.Printf("Comparing achievements: playerID=%s, friendID=%s, appID=%s", playerID, friendID, appID)

	if playerID == "" || friendID == "" || appID == "" {
		RespondWithError(w, http.StatusBadRequest, "Missing one or more required ID params", nil)
		return
	}

	playerAchievements, err := apicfg.GetPlayerAchievements(playerID, appID)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Unable to perform API call to Steam GetPlayerAchievements for user", err)
		return
	}

	friendAchievements, err := apicfg.GetPlayerAchievements(friendID, appID)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Unable to perform API call to Steam GetPlayerAchievements for friend", err)
		return
	}

	result := map[string]interface{}{
		"player": playerAchievements,
		"friend": friendAchievements,
	}

	RespondWithJSON(w, http.StatusOK, result)
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

func (apicfg *ApiConfig) HandlerMatchedGamesRanking(w http.ResponseWriter, req *http.Request) {
	steamid := req.URL.Query().Get("steamID")
	if steamid == "" {
		RespondWithError(w, http.StatusBadRequest, "'steamid' parameter is required for getting player info", nil)
		return
	}

	listGames := false
	listGamesQuery := req.URL.Query().Get("listGames")
	if listGamesQuery == "true" {
		listGames = true
	}

	ownedGames, err := apicfg.GetOwnedGames(steamid)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Unable to perform API calls to Steam GetOwnedGames endpoint", err)
		return
	}

	friendList, err := apicfg.GetFriendList(steamid)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Unable to perform API calls to Steam GetFriendList endpoint", err)
		return
	}

	waitgGroup := sync.WaitGroup{}
	channel := make(chan ComparedMatchedGames, len(friendList.Friends))

	for _, friend := range friendList.Friends {
		waitgGroup.Add(1)

		go func(friend Friend, channel chan ComparedMatchedGames) {
			defer waitgGroup.Done()

			friendGames, err := apicfg.GetOwnedGames(friend.SteamID)
			if err != nil {
				log.Printf("Error getting games for %s: %v", friend.SteamID, err)
				return
			}

			result := ownedGames.CompareOwnedGames(friendGames, listGames)

			channel <- result
		}(friend, channel)
	}

	waitgGroup.Wait()

	results := []ComparedMatchedGames{}
	for range len(friendList.Friends) {
		result := <-channel
		results = append(results, result)
	}

	// Sort results by player's score in descending order
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	for i := range results {
		results[i].Ranking = i + 1
	}

	resp := struct {
		Ranking []ComparedMatchedGames `json:"ranking"`
	}{
		Ranking: results,
	}

	RespondWithJSON(w, http.StatusOK, resp)
}
