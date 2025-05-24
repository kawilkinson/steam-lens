import { useState } from "react";
import { Game } from "@/app/definitions/types";
import styles from "./GamesList.module.css";
import GameIcon from "./GameIcon/GameIcon";
import { MouseEventHandler } from "react";
import { AchievementComparisonData } from "@/app/definitions/types";
import { getAchievementComparison } from "@/app/api/api";
import AchievementSummary from "./AchievementSummary/AchievementSummary";

export default function GamesList({ 
  games, 
  listType, 
  setMatchingGamesDisplay, 
  setMissingGamesDisplay, 
  setAchievementsDisplay,
  numMatchingGames, 
  numMissingGames,
  userID,
  friendID,
}: {
    games: Game[] | null,
    listType: "matching" | "missing" | "achievements",
    setMatchingGamesDisplay: MouseEventHandler<HTMLButtonElement>,
    setMissingGamesDisplay: MouseEventHandler<HTMLButtonElement>,
    setAchievementsDisplay: MouseEventHandler<HTMLButtonElement>,
    numMatchingGames: number,
    numMissingGames: number,
    userID: string,
    friendID: string,
  }) {

  const [expandedAppID, setExpandedAppID] = useState<number | null>(null);
  const [loading, setLoading] = useState<number | null>(null);
  const [achievementData, setAchievementData] = useState<Record<number, AchievementComparisonData>>({});

  const iconWidth = 40;
  const iconHeight = 40;

  const handleGameClick = async (appID: number) => {
    if (expandedAppID === appID) {
      setExpandedAppID(null);
      return;
    }
    setLoading(appID);
    setExpandedAppID(appID);

    if (!achievementData[appID]) {
      const data = await getAchievementComparison(userID, friendID, appID);
      if (data) {
        setAchievementData(prev => ({ ...prev, [appID]: data }));
      }
    }
    setLoading(null);
  };

  if (games != null) {
    games.sort((a: Game, b: Game): number => {
      if (a.name < b.name) {
        return -1;
      } else if (a.name > b.name) {
        return 1;
      }
      return 0;
    });
  }
  return (
    <div className={`${styles.container} ${listType === "missing" ? styles.missingList : ""} ${listType === "achievements" ? styles.achievementsList : ""}`}>
      <div className={styles.buttonsContainer}>
        <button
          onClick={setMatchingGamesDisplay}
          className={`${styles.button} ${listType === "matching" ? styles.active : ""}`}
          disabled={listType === "matching"}
        >
          <h3>Matching Games</h3>
        </button>
        <button
          onClick={setMissingGamesDisplay}
          className={`${styles.button} ${listType === "missing" ? styles.activeMissing : ""}`}
          disabled={listType === "missing"}
        >
          <h3>Missing Games</h3>
        </button>
        <button
          onClick={setAchievementsDisplay}
          className={`${styles.button} ${listType === "achievements" ? styles.activeAchievements : ""}`}
          disabled={listType === "achievements"}
        >
          <h3>Achievements</h3>
        </button>
      </div>
      <div className={styles.counts}>
        {listType === "matching" && (
          <span className={styles.matchingCount}>
            {numMatchingGames} games both players own
          </span>
        )}
        {listType === "missing" && (
          <span className={styles.missingCount}>
            {numMissingGames} games only friend owns
          </span>
        )}
        {listType === "achievements" && (
          <span className={styles.achievementCount}>
            {numMatchingGames} games to compare achievements with
          </span>
        )}
      </div>
        <div className={styles.list}>
          {games && games.length > 0 ? (
            games.map((game) => (
              <div key={game.appID} className={styles.gameWrapper}>
                {(listType === "matching" || listType === "missing") && (
                  <a
                    className={`${styles.game} ${listType === "missing" ? styles.missingEntry : ""}`}
                    href={`https://store.steampowered.com/app/${game.appID}/`}
                    target="_blank"
                    rel="noopener"
                    onClick={ev => ev.stopPropagation()}
                  >
                    <GameIcon game={game} width={iconWidth} height={iconHeight} />
                    <p>{game.name}</p>
                  </a>
                )}
                {listType === "achievements" && (
                  <AchievementSummary
                    game={game}
                    expandedAppID={expandedAppID}
                    handleGameClick={handleGameClick}
                    loading={loading}
                    achievementData={achievementData}
                  />
                )}
            </div>
          ))
        ) : (
          <p>No games found</p>
        )}
      </div>
    </div>
  );
}
