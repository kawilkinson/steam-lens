import { getMatchingGames, getPlayerSummaries } from "@/app/api/api"
import styles from "./FriendsList.module.css"
import ProfileCard from "../ProfileCard/ProfileCard";
import { PlayerSummary } from "@/app/definitions/types";

export default async function FriendsList( { steamid }: { steamid: string}) {
    const matchedGames = await getMatchingGames(steamid);
    const steamids: string[] = [];

    interface IPlayerHash {
        [id: string]: PlayerSummary
    }
    const players: IPlayerHash = {};
    const summaries = await getPlayerSummaries(steamids.join(","));

    summaries.forEach((summary) => {
        players[summary.steamid] = summary;
    })

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

						return (summary ?
							<li key={entry.friendID}>
								<ProfileCard
									games={entry}
									summary={players[entry.friendID]} />
							</li>
							: <></>)
					})}
				</ul>
			}
		</div>
	)
}