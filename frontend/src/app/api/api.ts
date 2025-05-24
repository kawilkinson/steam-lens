"use server";

import { MatchingGames, PlayerSummary, AchievementComparisonData } from "../definitions/types";
import { baseURL } from "../definitions/urls";

export async function getPlayerSummaries(steamIDs: string): Promise<PlayerSummary[]> {
	if (!baseURL) {
		console.log("Backend base URL not set up in environment variables")
		return [];
	}

	const url = new URL("player-summaries", baseURL);
	url.searchParams.set("steamIDs", steamIDs);

	let summaries: PlayerSummary[] = [];

	try {
		console.log(`Sending request for Player Summaries: GET ${url.toString()}`)
		const resp = await fetch(url);

		if (resp.status >= 400) {
			throw new Error(`Failed http request with status ${resp.status}`)
		}

		const json = await resp.json();

		summaries = json.players
	} catch (e) {
		console.log(`Failed to get player summaries: ${e}`);
	}

	return summaries;
}

export async function getFriendsList(steamID: string): Promise<PlayerSummary[]> {
	if (!baseURL) {
		console.log("Backend base URL not set up in environment variables")
		return [];
	}

	const url = new URL("friends", baseURL);
	url.searchParams.set("steamID", steamID)

	let friendsList: PlayerSummary[] = [];

	try {
		console.log(`Sending request for Friends List: GET ${url.toString()}`)
		const resp = await fetch(url);

		if (resp.status >= 400) {
			throw new Error(`Failed http request with status ${resp.status}`);
		}

		const json = await resp.json();

		friendsList = json.players;

	} catch (e) {
		console.log(`Failed to get friends list: ${e}`);
	}

	return friendsList
}

export async function getMatchingGames(steamID: string): Promise<MatchingGames[]> {
    if (!baseURL) {
        console.log("Backend base URL not set up in environment variable")
        return [];
    }

    const url = new URL("friends/matchGames", baseURL);
    url.searchParams.set("steamID", steamID);
    url.searchParams.set("listGames", "true");

    let matchedGamesRanking: MatchingGames[] = [];

    try {
        console.log(`Sending request for matched games ranking: GET ${url.toString()}`)
        const resp = await fetch(url);

        if (resp.status >= 400) {
            throw new Error(`Failed http request with status ${resp.status}`);
        }

        const json = await resp.json();

        if (json.ranking != undefined) {
            matchedGamesRanking = json.ranking;
        }
    } catch (e) {
        console.log(`Failed to get matched games ranking: ${e}`);
    }

    return matchedGamesRanking
}

export async function getAchievementComparison(
  userID: string,
  friendID: string,
  appID: number
): Promise<AchievementComparisonData | null> {
  if (!baseURL) {
    console.log("Backend base URL not set up in environment variables");
    return null;
  }

  const url = new URL("compare-achievements", baseURL);
  url.searchParams.set("userID", userID);
  url.searchParams.set("friendID", friendID);
  url.searchParams.set("appID", appID.toString());

  try {
    console.log(`Sending request for Achievement Comparison: GET ${url.toString()}`);
	console.log("Fetching achievement comparison", userID, friendID, appID);
    const resp = await fetch(url);

    if (resp.status >= 400) {
      throw new Error(`Failed http request with status ${resp.status}`);
    }

    // This should match the AchievementComparisonData interface
    const data: AchievementComparisonData = await resp.json();
    return data;
  } catch (e) {
    console.log(`Failed to get achievement comparison: ${e}`);
    return null;
  }
}
