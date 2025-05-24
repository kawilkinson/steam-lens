import { getMatchingGames, getPlayerSummaries } from "@/app/api/api"
import styles from "./FriendsList.module.css"
import ProfileCard from "../ProfileCard/ProfileCard";
import { PlayerSummary } from "@/app/definitions/types";

export default async function FriendsList( { steamid: steamID }: { steamid: string}) {
    const matchedGames = await getMatchingGames(steamID);
	
    const steamIDs = matchedGames.map(entry => entry.friendID);

    interface IPlayerHash {
        [id: string]: PlayerSummary
    }
    const players: IPlayerHash = {};
    const summaries = await getPlayerSummaries(steamIDs.join(","));

    summaries.forEach((summary) => {
        players[summary.steamID] = summary;
    })

	const numOfRanks = summaries.length;

    return (
		<div className={styles.container}>
			<div className={styles.header}>
				<h2 className={styles.title}>Friends</h2>
			</div>
			{matchedGames.length == 0 ?
				<p>No friends found</p>
				: <ul>
					{matchedGames.map((entry) => {
						const summary = players[entry.friendID];
						if (!summary) return null;
						return (
						<li key={entry.friendID}>
							<ProfileCard
							games={entry}
							summary={summary}
							numOfRanks={numOfRanks}
							userID={steamID}
							/>
						</li>
						);
					})}
				</ul>
			}
		</div>
	)
}
