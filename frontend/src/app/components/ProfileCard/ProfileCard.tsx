"use client";

import Image from "next/image"
import { PlayerSummary, MatchingGames, Game } from "@/app/definitions/types";
import styles from "./ProfileCard.module.css"
import { MouseEvent, useEffect, useState } from "react";
import GamesList from "./GamesList/GamesList";

export default function ProfileCard({ summary, games }
    : {
    summary: PlayerSummary,
    games: MatchingGames | null,
    }) {

    const [expanded, setExpanded] = useState<boolean>(false);
    const [displayGames, setDisplayGames] = useState<Game[] | null>(null);
    const [listType, setListType] = useState<"matching" | "missing">("matching");

    useEffect(() => {
    if (games != null) {
        setDisplayGames(games.matchingGames);
        }
    }, [games]);

    const toggleExpand = (ev: MouseEvent<HTMLDivElement>) => {
    if (!expanded) {
      const element = (ev.target as Element)
      element.scrollIntoView({
        behavior: "smooth",
        block: "start",
        inline: "nearest",
        });
    }
    }

    if (games && games?.friendGamesCount > 0) {
      setExpanded(!expanded)
    }

    const setMatchingGamesDisplay = (ev: MouseEvent<HTMLButtonElement>) => {
    ev.stopPropagation();
    if (games) {
      setDisplayGames(games.matchingGames);
      setListType("matching");
    }
    }

    const setMissingGamesDisplay = (ev: MouseEvent<HTMLButtonElement>) => {
    ev.stopPropagation();
    if (games) {
      setDisplayGames(games.friendOnlyGames);
      setListType("missing");
    }
  }

  return (
    <div className={styles.container} onClick={toggleExpand}>
      <div className={styles.header}>
        <div className={styles.profileSection}>
          <Image className={styles.avatar} src={summary.avatarMedium} width={50} height={50} alt="Profile picture" />
          <p className={styles.personaname}>{summary.personaName}</p>
        </div>
      </div>
      {expanded && games &&
        <GamesList
          games={displayGames}
          listType={listType}
          setMatchingGamesDisplay={setMatchingGamesDisplay}
          setMissingGamesDisplay={setMissingGamesDisplay}
        />}
    </div>
  )
}
