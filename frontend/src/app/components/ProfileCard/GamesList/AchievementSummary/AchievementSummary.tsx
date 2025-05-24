import gamesListStyles from "../GamesList.module.css";
import achievementStyles from "./AchievementSummary.module.css";
import { AchievementComparisonData, Game } from "@/app/definitions/types";
import GameIcon from "../GameIcon/GameIcon";

type AchievementSummaryProps = {
  game: Game;
  expandedAppID: number | null;
  handleGameClick: (appID: number) => void;
  loading: number | null;
  achievementData: Record<number, AchievementComparisonData>;
};

export default function AchievementSummary({
  game,
  expandedAppID,
  handleGameClick,
  loading,
  achievementData,
}: AchievementSummaryProps) {
  const isExpanded = expandedAppID === game.appID;
  const data = achievementData[game.appID];

  function renderSummary(data?: AchievementComparisonData) {
    if (!data) return null;
    const playerAch = data.player?.achievements ?? [];
    const friendAch = data.friend?.achievements ?? [];
    const playerCount = playerAch.filter(a => a.achieved).length;
    const friendCount = friendAch.filter(a => a.achieved).length;
    const total = Math.max(playerAch.length, friendAch.length);

    return (
      <div className={achievementStyles.achievementSummary}>
        <p>
          <strong>You:</strong> {playerCount} / {total} achievements
        </p>
        <p>
          <strong>Friend:</strong> {friendCount} / {total} achievements
        </p>
      </div>
    );
  }

  return (
    <>
      <a
        className={`${gamesListStyles.game} ${gamesListStyles.achievementEntry} ${isExpanded ? gamesListStyles.expanded : ""}`}
        href="#"
        onClick={e => {
          e.preventDefault();
          e.stopPropagation();
          handleGameClick(game.appID);
        }}
        tabIndex={0}
        role="button"
        aria-expanded={isExpanded}
      >
        <GameIcon game={game} width={40} height={40} />
        <p>{game.name}</p>
      </a>
      {isExpanded && (
        <div className={achievementStyles.achievementComparison}>
          {loading === game.appID ? (
            <p>Loading achievements...</p>
          ) : (
            renderSummary(data)
          )}
        </div>
      )}
    </>
  );
}
