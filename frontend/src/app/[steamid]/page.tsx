import styles from "./page.module.css";
import ProfileCard from "../components/ProfileCard/ProfileCard";
import { getPlayerSummaries } from "../api/api";
import FriendsList from "../components/FriendsList/FriendsList";
import { PlayerSummary } from "../definitions/types";
import MainHeader from "../components/Header/MainHeader";

type Params = Promise<{steamid: string}>

export default async function GamesPage({ params }: { params: Params }) {
  const steamID = (await params).steamid;

    const resp = await getPlayerSummaries(steamID);

    let summary: PlayerSummary | null = null
    if (resp.length > 0) {
    summary = resp[0];
    }

    return (
    <>
        <MainHeader />
        <main className={styles.main}>
        {summary != null ?
            <>
            <p className={styles.title}>User</p>
            <ProfileCard summary={summary!} games={null} numOfRanks={0} userID={steamID} isUserProfile={true} />
            <FriendsList steamid={steamID} />
            </>
            :
            <p className={styles.error}>Invalid id {steamID}</p>
        }
        </main>
    </>
    );
}
