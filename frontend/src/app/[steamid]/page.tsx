import styles from "./page.module.css";
import ProfileCard from "../components/ProfileCard/ProfileCard";
import { getPlayerSummaries } from "../api/api";
import FriendsList from "../components/FriendsList/FriendsList";
import { PlayerSummary } from "../definitions/types";
import Link from "next/link";

export default async function GamesPage(props: { params: { steamid: string } }) {
  const { params } = await props;
  const steamid = params.steamid;  

    const resp = await getPlayerSummaries(steamid);

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
            <ProfileCard summary={summary!} games={null} />
            <FriendsList steamid={steamid} />
            </>
            :
            <p className={styles.error}>Invalid id {steamid}</p>
        }
        </main>

    </>
    );
}
