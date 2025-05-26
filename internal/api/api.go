package api

import (
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"slices"
	"strings"
)

const steamMainAPIURL = "http://api.steampowered.com/"
const steamUserURL = "ISteamUser/"
const steamPlayerURL = "IPlayerService/"
const steamAchievementURL = "ISteamUserStats/"

type Player struct {
	SteamID                  string `json:"steamID"`
	CommunityVisibilityState int    `json:"communityVisibilityState"`
	PersonaName              string `json:"personaName"`
	Avatar                   string `json:"avatar"`
	AvatarMedium             string `json:"avatarMedium"`
	AvatarFull               string `json:"avatarFull"`
}

type Summaries struct {
	Players []Player `json:"players"`
}

type SummariesResponse struct {
	Response Summaries `json:"response"`
}

// Make API call to Steam's GetPlayerSummaries endpoint to obtain player data for all steam IDs provided
func (apicfg *ApiConfig) GetPlayerSummaries(steamIDs []string) (Summaries, error) {
	uncachedIDs := []string{}
	cachedPlayers := []Player{}

	for _, steamID := range steamIDs {
		cache, found := apicfg.PlayerCache.ReadCache(steamID)
		if found {
			log.Printf("Cache found for steamID: %s\n", steamID)
			cachedPlayers = append(cachedPlayers, cache)
		} else {
			uncachedIDs = append(uncachedIDs, steamID)
		}
	}

	if len(uncachedIDs) == 0 {
		slices.SortFunc(cachedPlayers, func(i Player, j Player) int {
			return cmp.Compare(i.SteamID, j.SteamID)
		})
		return Summaries{
			Players: cachedPlayers,
		}, nil
	}

	joinedIDs := strings.Join(uncachedIDs, ",")

	baseURL, err := url.Parse(steamMainAPIURL)
	if err != nil {
		return Summaries{}, err
	}

	fullURL := baseURL.JoinPath(steamUserURL, "GetPlayerSummaries", "v0002/")

	query := url.Values{}
	query.Set("key", apicfg.SteamApiKey)
	query.Set("steamids", joinedIDs)

	fullURL.RawQuery = query.Encode()

	resp, err := http.Get(fullURL.String())
	if err != nil {
		return Summaries{}, err
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "application/json") {
		testBody, _ := io.ReadAll(resp.Body)
		fmt.Printf("Unexpected response from Steam API: %s\n", testBody)
		return Summaries{}, errors.New("steam API returned non-JSON response")
	}

	decoder := json.NewDecoder(resp.Body)

	body := SummariesResponse{}
	err = decoder.Decode(&body)
	if err != nil {
		return Summaries{}, err
	}

	allPlayers := body.Response.Players
	for _, player := range allPlayers {
		log.Printf("Adding player to cache with steamID: %s\n", player.SteamID)
		apicfg.PlayerCache.UpdateCache(player.SteamID, player)
	}

	return body.Response, nil
}

// For now, imgIconURL returns img_icon_url for json for since Steam's API uses snake case
type Game struct {
	AppID      int    `json:"appID"`
	Name       string `json:"name"`
	ImgIconURL string `json:"img_icon_url"`
}

type OwnedGames struct {
	SteamID   string
	GameCount int    `json:"game_count"`
	Games     []Game `json:"games"`
}

type OwnedGamesResponse struct {
	Response OwnedGames `json:"response"`
}

// Make API call to Steam's GetOwnedGames endpoint to obtain all owned games for a user
func (apicfg *ApiConfig) GetOwnedGames(steamID string) (OwnedGames, error) {
	_, found := apicfg.OwnedGamesCache.ReadCache(steamID)
	if found {
		log.Printf("OwnedGames cache found for %s\n", steamID)
		return apicfg.OwnedGamesCache.Cache[steamID].Data, nil
	}

	baseURL, err := url.Parse(steamMainAPIURL)
	if err != nil {
		return OwnedGames{}, err
	}

	fullURL := baseURL.JoinPath(steamPlayerURL, "GetOwnedGames", "v0001/")

	query := url.Values{}
	query.Set("key", apicfg.SteamApiKey)
	query.Set("steamid", steamID)
	query.Set("include_appinfo", "true")

	fullURL.RawQuery = query.Encode()

	resp, err := http.Get(fullURL.String())
	if err != nil {
		return OwnedGames{}, err
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "application/json") {
		testBody, _ := io.ReadAll(resp.Body)
		fmt.Printf("Unexpected response from Steam API: %s\n", testBody)
		return OwnedGames{}, errors.New("steam API returned non-JSON response")
	}

	decoder := json.NewDecoder(resp.Body)

	body := OwnedGamesResponse{}
	err = decoder.Decode(&body)
	if err != nil {
		return OwnedGames{}, err
	}

	body.Response.SteamID = steamID

	apicfg.OwnedGamesCache.UpdateCache(steamID, body.Response)

	return body.Response, nil
}

type Friend struct {
	SteamID      string `json:"steamID"`
	Relationship string `json:"relationship"`
	FriendSince  int    `json:"friendSince"`
}

type FriendList struct {
	Friends []Friend `json:"friends"`
}

type FriendListResponse struct {
	Friendlist FriendList `json:"friendslist"`
}

// Make API call to Steam's GetFriendList endpoint to obtain all friends for a user
func (apicfg *ApiConfig) GetFriendList(steamID string) (FriendList, error) {
	_, found := apicfg.FriendListCache.ReadCache(steamID)
	if found {
		log.Printf("FriendList cache found for %s\n", steamID)
		return apicfg.FriendListCache.Cache[steamID].Data, nil
	}

	baseURL, err := url.Parse(steamMainAPIURL)
	if err != nil {
		return FriendList{}, err
	}

	fullURL := baseURL.JoinPath(steamUserURL, "GetFriendList", "v0001/")

	query := url.Values{}
	query.Set("key", apicfg.SteamApiKey)
	query.Set("steamid", steamID)
	query.Set("relationship", "friend")

	fullURL.RawQuery = query.Encode()

	resp, err := http.Get(fullURL.String())
	if err != nil {
		return FriendList{}, err
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "application/json") {
		testBody, _ := io.ReadAll(resp.Body)
		fmt.Printf("Unexpected response from Steam API: %s\n", testBody)
		return FriendList{}, errors.New("steam API returned non-JSON response")
	}

	decoder := json.NewDecoder(resp.Body)

	body := FriendListResponse{}
	err = decoder.Decode(&body)
	if err != nil {
		return FriendList{}, err
	}

	apicfg.FriendListCache.UpdateCache(steamID, body.Friendlist)

	return body.Friendlist, nil
}

type Achievement struct {
	ApiName  string `json:"apiname"`
	Achieved int    `json:"achieved"`
}

type PlayerAchievements struct {
	Achievements []Achievement `json:"achievements"`
}

type PlayerAchievementsResponse struct {
	PlayerAchievements PlayerAchievements `json:"playerstats"`
}

// Steam returns an int for achieved status, using this struct to convert to bool
type ConvertedAchievement struct {
	ApiName  string `json:"apiName"`
	Achieved bool   `json:"achieved"`
}

type ConvertedPlayerAchievements struct {
	Achievements []ConvertedAchievement `json:"achievements"`
}

// Make API call to Steam's GetPlayerAchievements endpoint to obtain all achievements for a game
func (apicfg *ApiConfig) GetPlayerAchievements(steamID, appID string) (ConvertedPlayerAchievements, error) {
	cacheKey := steamID + "-" + appID
	_, found := apicfg.AchievementsCache.ReadCache(cacheKey)
	if found {
		log.Printf("Achievements cache found for %s\n", cacheKey)
		return apicfg.AchievementsCache.Cache[cacheKey].Data, nil
	}

	baseURL, err := url.Parse(steamMainAPIURL)
	if err != nil {
		return ConvertedPlayerAchievements{}, err
	}

	fullURL := baseURL.JoinPath(steamAchievementURL, "GetPlayerAchievements", "v0001/")

	query := url.Values{}
	query.Set("appid", appID)
	query.Set("key", apicfg.SteamApiKey)
	query.Set("steamid", steamID)

	fullURL.RawQuery = query.Encode()

	resp, err := http.Get(fullURL.String())
	if err != nil {
		return ConvertedPlayerAchievements{}, err
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "application/json") {
		testBody, _ := io.ReadAll(resp.Body)
		fmt.Printf("Unexpected response from Steam API: %s\n", testBody)
		return ConvertedPlayerAchievements{}, errors.New("steam API returned non-JSON response")
	}

	decoder := json.NewDecoder(resp.Body)

	body := PlayerAchievementsResponse{}
	err = decoder.Decode(&body)
	if err != nil {
		return ConvertedPlayerAchievements{}, err
	}

	// Use helper function to convert int values in achieved status to bool value
	convertedAchievements := convertAchievements(body.PlayerAchievements)

	apicfg.AchievementsCache.UpdateCache(cacheKey, convertedAchievements)

	return convertedAchievements, nil
}

// Helper function to converted achievements to bool type
func convertAchievements(input PlayerAchievements) ConvertedPlayerAchievements {
	converted := make([]ConvertedAchievement, len(input.Achievements))
	for i, ach := range input.Achievements {
		converted[i] = ConvertedAchievement{
			ApiName:  ach.ApiName,
			Achieved: ach.Achieved == 1,
		}
	}
	return ConvertedPlayerAchievements{Achievements: converted}
}

type ComparedMatchedGames struct {
	Score            float64 `json:"score"`
	Ranking          int     `json:"ranking"`
	UserID           string  `json:"userID"`
	UserPercentage   float64 `json:"userPercentage"`
	FriendID         string  `json:"friendID"`
	FriendGamesCount int     `json:"friendGamesCount"`
	FriendPercentage float64 `json:"friendPercentage"`
	Matches          int     `json:"matches"`
	MatchingGames    []Game  `json:"matchingGames"`
	FriendOnlyGames  []Game  `json:"friendOnlyGames"`
}

// Run comparisons on user and their friend's games to get overall ranking
func (userGames OwnedGames) CompareOwnedGames(friendGames OwnedGames, listGames bool) ComparedMatchedGames {
	result := ComparedMatchedGames{
		UserID:   userGames.SteamID,
		FriendID: friendGames.SteamID,
	}

	result.MatchingGames = []Game{}
	result.FriendOnlyGames = []Game{}

	for _, game := range friendGames.Games {
		if slices.Contains(userGames.Games, game) {
			result.MatchingGames = append(result.MatchingGames, game)
		} else {
			result.FriendOnlyGames = append(result.FriendOnlyGames, game)
		}
	}
	result.Matches = len(result.MatchingGames)

	if userGames.GameCount > 0 {
		result.UserPercentage = float64(result.Matches) / float64(userGames.GameCount)
	}

	result.FriendGamesCount = friendGames.GameCount
	if friendGames.GameCount > 0 {
		result.FriendPercentage = float64(result.Matches) / float64(friendGames.GameCount)
	}

	matchesWeight := 0.6
	percentWeight := 0.4
	result.Score = float64(result.Matches)*matchesWeight + result.FriendPercentage*100.0*percentWeight

	if !listGames {
		result.MatchingGames = nil
		result.FriendOnlyGames = nil
	}

	return result
}
