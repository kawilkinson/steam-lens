export interface PlayerSummary {
	steamid: string;
	personaName: string;
	avatar: string;
	avatarMedium: string;
	avatarFull: string;
}

export interface Game {
	appid: number;
	name: string;
	imgIconURL: string;
}

export interface MatchingGames {
    matchingGames: Game[];
    matches: number;
    friendID: string;
    friendOnlyGames: Game[];
    friendGamesCount: number;
}