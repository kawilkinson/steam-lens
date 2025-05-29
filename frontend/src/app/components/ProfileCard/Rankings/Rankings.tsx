import { MatchingGames } from "@/app/definitions/types"
import styles from "./Rankings.module.css"

function getRankColor(rank: number, maxRank: number): string {
  if (maxRank < 2) return "hsl(120, 100%, 40%)";

  const hue = 120 - ((120 * (rank - 1)) / (maxRank - 1));
  return `hsl(${hue}, 100%, 40%)`;
}

export default function Rankings({ rankings, numOfRanks }: { rankings: MatchingGames | null | undefined, numOfRanks: number }) {
  const isPrivate = !rankings || !rankings.friendGamesCount

  if (isPrivate) {
    return (
      <div className={styles.container}>
        <p className={styles.percentageOwned}></p>
        <p>
          <span className={styles.privateProfile}>Private Profile</span>
        </p>
      </div>
    )
  }

  const percentageOwned = rankings.friendPercentage;
  const rankColor = getRankColor(rankings.ranking, numOfRanks)

  return (
    <div className={styles.container}>
        <p className={styles.percentageOwned}>
            Currently have <span className={styles.percent}>{`${(percentageOwned * 100.0).toFixed(2)}%`}</span> of this player&apos;s games
        </p>
        <p>
        <span className={styles.rankText}>Rank: </span>
        <span className={styles.rankNumber} 
        style={{ color: rankColor }}
        >
          {rankings.ranking.toFixed(0)}
          </span>
        </p>
    </div>
  )
}
