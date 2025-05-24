import { useState } from "react";
import { Game } from "@/app/definitions/types";
import styles from "./GamesList.module.css";
import GameIcon from "./GameIcon/GameIcon";
import { MouseEventHandler } from "react";
import { AchievementComparisonData } from "@/app/definitions/types";
import { getAchievementComparison } from "@/app/api/api";


export default function GamesList({ games, 
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
                  onClick={(ev) => ev.stopPropagation()}
                >
                  <GameIcon game={game} width={iconWidth} height={iconHeight} />
                  <p>{game.name}</p>
                </a>
              )}
              {listType === "achievements" && (
                <>
                  <button
                    className={`${styles.game} ${styles.achievementEntry} ${expandedAppID === game.appID ? styles.expanded : ""}`}
                    onClick={e => {
                      e.stopPropagation();
                      handleGameClick(game.appID);
                    }}
                  >
                    <GameIcon game={game} width={iconWidth} height={iconHeight} />
                    <p>{game.name}</p>
                  </button>
                  {expandedAppID === game.appID && (
                    <div className={styles.achievementComparison}>
                      {loading === game.appID ? (
                        <p>Loading achievements...</p>
                      ) : achievementData[game.appID] ? (
                        <div className={styles.achievementSummary}>
                          {(() => {
                            const playerAch = achievementData[game.appID]?.player?.achievements ?? [];
                            const friendAch = achievementData[game.appID]?.friend?.achievements ?? [];
                            const playerCount = playerAch.filter(a => a.achieved).length;
                            const friendCount = friendAch.filter(a => a.achieved).length;
                            const total = Math.max(playerAch.length, friendAch.length);

                            let verdict = "It's a tie";
                            if (playerCount > friendCount) {
                              verdict = "User has more achievements";
                            } else if (friendCount > playerCount) {
                              verdict = "Friend has more achievements";
                            }

                            return (
                              <>
                                <p>
                                  <strong>You:</strong> {playerCount} / {total} achievements
                                </p>
                                <p>
                                  <strong>Friend:</strong> {friendCount} / {total} achievements
                                </p>
                                <p>{verdict}</p>
                              </>
                            );
                          })()}
                        </div>
                      ) : null}
                    </div>
                  )}
                </>
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
