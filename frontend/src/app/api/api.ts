"use server";

import { MatchingGames, PlayerSummary, AchievementComparisonData } from "../definitions/types";
import { apiURL } from "../definitions/urls";

export async function getPlayerSummaries(steamIDs: string): Promise<PlayerSummary[]> {
	if (!apiURL) {
		console.log("Backend base URL not set up in environment variables")
		return [];
	}

	const url = new URL("player-summaries", apiURL);
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
	} catch (err) {
		console.log(`Failed to get player summaries: ${err}`);
	}

	return summaries;
}

export async function getFriendsList(steamID: string): Promise<PlayerSummary[]> {
	if (!apiURL) {
		console.log("Backend base URL not set up in environment variables")
		return [];
	}

	const url = new URL("friends", apiURL);
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

	} catch (err) {
		console.log(`Failed to get friends list: ${err}`);
	}

	return friendsList
}

export async function getMatchingGames(steamID: string): Promise<MatchingGames[]> {
    if (!apiURL) {
        console.log("Backend base URL not set up in environment variable")
        return [];
    }

    const url = new URL("friends/matchGames", apiURL);
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
    } catch (err) {
        console.log(`Failed to get matched games ranking: ${err}`);
    }

    return matchedGamesRanking
}

export async function getAchievementComparison(
  userID: string,
  friendID: string,
  appID: number
): Promise<AchievementComparisonData | null> {
  if (!apiURL) {
    console.log("Backend base URL not set up in environment variables");
    return null;
  }

  const url = new URL("compare-achievements", apiURL);
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
  } catch (err) {
    console.log(`Failed to get achievement comparison: ${err}`);
    return null;
  }
}
