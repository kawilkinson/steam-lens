package api

import (
	"log"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
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

	type job struct {
		friend Friend
		idx    int
	}

	friendCount := len(friendList.Friends)
	if friendCount == 0 {
		RespondWithJSON(w, http.StatusOK, struct {
			Ranking []ComparedMatchedGames `json:"ranking"`
		}{Ranking: []ComparedMatchedGames{}})
		return
	}

	var waitTime int
	if friendCount > 100 {
		waitTime = 15
	} else if friendCount > 50 {
		waitTime = 8
	} else {
		waitTime = 4
	}

	// Use a worker pool (currently 1 worker) to make all API calls over an interval
	// that is determined by how many friends a user has, this is to help prevent 429 errors from Steam's end
	jobs := make(chan job, friendCount)
	results := make([]ComparedMatchedGames, friendCount)
	waitGroup := sync.WaitGroup{}

	interval := time.Duration(float64(waitTime) / float64(friendCount) * float64(time.Second))
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	worker := func() {
		for i := range jobs {
			<-ticker.C
			friendGames, err := apicfg.GetOwnedGames(i.friend.SteamID)
			if err != nil {
				log.Printf("Error getting games for %s: %v", i.friend.SteamID, err)
				waitGroup.Done()
				continue
			}
			result := ownedGames.CompareOwnedGames(friendGames, listGames)
			results[i.idx] = result
			waitGroup.Done()
		}
	}

	go worker()

	// Queue up jobs
	for idx, friend := range friendList.Friends {
		waitGroup.Add(1)
		jobs <- job{friend: friend, idx: idx}
	}
	close(jobs)
	waitGroup.Wait()

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
