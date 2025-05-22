"use client";

import Image from "next/image"
import { PlayerSummary, MatchingGames, Game } from "@/app/definitions/types";
import styles from "./ProfileCard.module.css"
import { MouseEvent, useState } from "react";
import GamesList from "./GamesList/GamesList";
import Rankings from "./Rankings/Rankings";

export default function ProfileCard({ summary, games, numOfRanks }: {
  summary: PlayerSummary,
  games: MatchingGames | null,
  numOfRanks: number,
}) {
  const [expanded, setExpanded] = useState<boolean>(false);
  const [listType, setListType] = useState<"matching" | "missing">("matching");
  const [displayGames, setDisplayGames] = useState<Game[] | null>(null);

  const handleToggleExpand = () => {
    setExpanded((prev) => !prev);

    if (!expanded && games) {
      setDisplayGames(games.matchingGames);
      setListType("matching");
    }
  };

  const setMatchingGamesDisplay = (ev: MouseEvent<HTMLButtonElement>) => {
    ev.stopPropagation();
    if (games) {
      setDisplayGames(games.matchingGames);
      setListType("matching");
    }
  };

  const setMissingGamesDisplay = (ev: MouseEvent<HTMLButtonElement>) => {
    ev.stopPropagation();
    if (games) {
      setDisplayGames(games.friendOnlyGames);
      setListType("missing");
    }
  };

  return (
    <div className={styles.container} onClick={handleToggleExpand} style={{ cursor: "pointer" }}>
      <div className={styles.header}>
        <div className={styles.profileSection}>
          <Image
            className={styles.avatar}
            src={summary.avatarMedium}
            width={50}
            height={50}
            alt="Profile picture"
          />
          <p className={styles.personaname}>{summary.personaName}</p>
        </div>

          {games == null || games.friendGamesCount == 0 ? <></>
          :
          <Rankings rankings={games} numOfRanks={numOfRanks}/>}

      </div>
      {expanded && games &&
        <GamesList
          games={displayGames}
          listType={listType}
          setMatchingGamesDisplay={setMatchingGamesDisplay}
          setMissingGamesDisplay={setMissingGamesDisplay}
        />
      }
    </div>
  );
}