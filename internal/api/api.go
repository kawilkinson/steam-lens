package api

import (
	"encoding/json"
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
	baseURL, err := url.Parse(steamMainAPIURL)
	if err != nil {
		return Summaries{}, err
	}

	fullURL := baseURL.JoinPath(steamUserURL, "GetPlayerSummaries", "v0002/")

	query := url.Values{}
	query.Set("key", apicfg.SteamApiKey)
	query.Set("steamids", strings.Join(steamIDs, ","))

	fullURL.RawQuery = query.Encode()

	resp, err := http.Get(fullURL.String())
	if err != nil {
		return Summaries{}, err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	body := SummariesResponse{}
	err = decoder.Decode(&body)
	if err != nil {
		return Summaries{}, err
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

	decoder := json.NewDecoder(resp.Body)

	body := OwnedGamesResponse{}
	err = decoder.Decode(&body)
	if err != nil {
		return OwnedGames{}, err
	}

	body.Response.SteamID = steamID

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

	decoder := json.NewDecoder(resp.Body)

	body := FriendListResponse{}
	err = decoder.Decode(&body)
	if err != nil {
		return FriendList{}, err
	}

	return body.Friendlist, nil
}

type Achievement struct {
	ApiName    string `json:"apiName"`
	Achieved   bool   `json:"achieved"`
	UnlockTime string `json:"unlockTime"`
}

type PlayerAchievements struct {
	Achievements []Achievement `json:"achievements"`
}

type PlayerAchievementsResponse struct {
	PlayerAchievements PlayerAchievements
}

// Make API call to Steam's GetPlayerAchievements endpoint to obtain all achievements for a game
func (apicfg *ApiConfig) GetPlayerAchievements(steamID, appid string) (PlayerAchievements, error) {
	baseURL, err := url.Parse(steamMainAPIURL)
	if err != nil {
		return PlayerAchievements{}, err
	}

	fullURL := baseURL.JoinPath(steamAchievementURL, "GetPlayerAchievements", "v0001/")

	query := url.Values{}
	query.Set("key", apicfg.SteamApiKey)
	query.Set("steamid", steamID)
	query.Set("appid", appid)

	fullURL.RawQuery = query.Encode()

	resp, err := http.Get(fullURL.String())
	if err != nil {
		return PlayerAchievements{}, err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	body := PlayerAchievementsResponse{}
	err = decoder.Decode(&body)
	if err != nil {
		return PlayerAchievements{}, err
	}

	return body.PlayerAchievements, nil
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

	matchesWeight := 0.7
	percentWeight := 0.3
	result.Score = float64(result.Matches)*matchesWeight + result.FriendPercentage*100.0*percentWeight

	if !listGames {
		result.MatchingGames = nil
		result.FriendOnlyGames = nil
	}

	return result
}
