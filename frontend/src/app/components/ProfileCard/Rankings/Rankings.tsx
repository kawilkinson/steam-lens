import { MatchingGames } from "@/app/definitions/types"
import styles from "./Rankings.module.css"

export default function Rankings({ rankings }: { rankings: MatchingGames }) {
  const percentageOwned = rankings.friendPercentage;

  return (
    <div className={styles.container}>
        <p className={styles.percentageOwned}>
            You own <span className={styles.percent}>{`${(percentageOwned * 100.0).toFixed(2)}%`}</span> of their games
        </p>
        <p>
        <span className={styles.rankText}>Rank: </span>
        <span className={styles.rankNumber}>{rankings.ranking.toFixed(0)}</span>
        </p>
    </div>
  )
}