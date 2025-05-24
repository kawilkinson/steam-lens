import styles from "./page.module.css";
import ProfileCard from "../components/ProfileCard/ProfileCard";
import { getPlayerSummaries } from "../api/api";
import FriendsList from "../components/FriendsList/FriendsList";
import { PlayerSummary } from "../definitions/types";
import Link from "next/link";

export default async function GamesPage({ params }: { params: { steamid: string } }) {
    const steamID = (await params).steamid;

    const resp = await getPlayerSummaries(steamID);

    let summary: PlayerSummary | null = null
    if (resp.length > 0) {
    summary = resp[0];
    }

    return (
    <>
        <header className={styles.header}>
        <h1>Steam Lens</h1>
        <Link href="/" className={styles.homeButton}>Home</Link>
        </header>
        <main className={styles.main}>
        {summary != null ?
            <>
            <ProfileCard summary={summary!} games={null} numOfRanks={0} userID={steamID} />
            <FriendsList steamid={steamID} />
            </>
            :
            <p className={styles.error}>Invalid id {steamID}</p>
        }
        </main>

    </>
    );
}
