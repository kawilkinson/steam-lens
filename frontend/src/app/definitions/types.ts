export interface PlayerSummary {
	steamID: string;
	personaName: string;
	avatar: string;
	avatarMedium: string;
	avatarFull: string;
}

// img_icon_url isn't idomatic with its snake casing since currently the backend is serving json that way
export interface Game {
	appID: number;
	name: string;
	img_icon_url: string;
}

export interface MatchingGames {
	ranking: number;
	score: number;
	UserID: string;
	UserPercentage: number;
    friendID: string;
	friendGamesCount: number;
	friendPercentage: number;
	matches: number;
	matchingGames: Game[];
    friendOnlyGames: Game[];
}

export interface Achievement {
	apiName: string;
	achieved: boolean;
}
