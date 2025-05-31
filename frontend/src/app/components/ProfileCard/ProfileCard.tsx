"use client";

import Image from "next/image"
import { PlayerSummary, MatchingGames, Game } from "@/app/definitions/types";
import styles from "./ProfileCard.module.css"
import { MouseEvent, useState } from "react";
import GamesList from "./GamesList/GamesList";
import Rankings from "./Rankings/Rankings";

export default function ProfileCard({ summary, games, numOfRanks, userID, isUserProfile = false }: {
  summary: PlayerSummary,
  games: MatchingGames | null,
  numOfRanks: number,
  userID: string
  isUserProfile: boolean,
}) {
  const [expanded, setExpanded] = useState<boolean>(false);
  const [listType, setListType] = useState<"matching" | "missing" | "achievements">("matching");
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

  const setAchievementsDisplay = (ev: MouseEvent<HTMLButtonElement>) => {
    ev.stopPropagation();
    if (games) {
      setDisplayGames(games.matchingGames);
      setListType("achievements")
    }
  }

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
            unoptimized={true}
          />
          <p className={styles.personaname}>{summary.personaName}</p>
        </div>
        {!isUserProfile && (
          <Rankings rankings={games} numOfRanks={numOfRanks}/>
        )}
        </div>
      {expanded && games &&
        <GamesList
          games={displayGames}
          listType={listType}
          setMatchingGamesDisplay={setMatchingGamesDisplay}
          setMissingGamesDisplay={setMissingGamesDisplay}
          setAchievementsDisplay={setAchievementsDisplay}
          numMatchingGames={games.matchingGames.length}
          numMissingGames={games.friendOnlyGames.length}
          userID={userID}
          friendID={games.friendID}
        />
      }
    </div>
  );
}
