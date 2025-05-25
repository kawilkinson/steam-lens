'use client'

import { useState } from "react";
import styles from "./page.module.css";
import { useRouter } from "next/navigation";

export default function Home() {
  const [userID, setUserID] = useState<string>("");
  const router = useRouter();

  // Allows for an optional Steam Community profile URL that contains the Steam ID to be submitted instead
  function handleClick() {
    let destination = userID;
    const match = userID.match(/https:\/\/steamcommunity\.com\/profiles\/(\d+)/);
    if (match !== null) {
      destination = match[1];
    }

    router.push(`/${destination}`);
  }

  return (
    <>
      <header className={styles.header}>
        <h1 className={styles.logoTitle}>Steam Lens</h1>
      </header>
      <main className={styles.main}>
        <div className={styles.section}>
          <div className={styles.form}>
            <label htmlFor="userID">Steam ID:</label>
            <input name="userID" type="text" required value={userID} onChange={(val) => setUserID(val.target.value)} />
            <button type="submit" onClick={handleClick}>Submit</button>
          </div>
        </div>
        <div className={styles.section}>
          <h2 className={styles.h2}>What is Steam Lens?</h2>
          <p>
            Steam Lens is a website that allows comparing a user&apos;s owned games and achievements with their friends.
          </p>
          <p>
            Steam Lens gives the ability to see what games a friend may have that the user doesn&apos;t own or to see who has more achievements for a game. 
            You can also click a matched or missing game to go straight to its Steam store page!
          </p>
        </div>
        <div className={styles.section}>
          <h2 className={styles.h2}>How to Use</h2>
          <p>
            Enter a 17 digit numeric Steam ID of your own or a user you want to run comparisons on.
          </p>
          <p>
            Once entered, clicking submit will redirect you over to the main dashboard where you&apos;ll be able to compare games and achievements.
          </p>
          <p>
            Note, if the user or a friend has their profile privated in any way then that can prevent data from the API to return for that account.
            If you notice rankings, games, achievement counts, etc. aren&apos;t loading for a specific person then apart from the daily API limit being reached that would 
            the reason for data not loading.
          </p>
        </div>
        <div className={styles.section}>
          <h2 className={styles.h2}>How rankings are calculated</h2>
          <p>
            The user&apos;s friends are ranked according to the similarity of their game libraries. The more games owned by both the user and friend, the higher the friend will be placed. 
            For example as of right now there is a hidden score that is calculated with the number of matching games making up 60% of the score, and percentage of matching games making up 40% of the score.
          </p>
        </div>
        <div className={styles.section}>
          <h2 className={styles.h2}>Limitations</h2>
          <p>
            This website works only because of Valve&apos;s Steam API that is useable for anyone with a Steam account that requests an API key.
            They have a 100,000 daily API call limit. If you notice errors being thrown or data not loading, then it is likely due to that limit being reached
            since user/friend data, game data, and achievement data are all gotten from Steam&apos;s endpoints.
          </p>
          <p>
            I have added backend caching of the data returned from Steam&apos;s endpoints to reduce API calls as much as possible to make it unlikely this daily limit is reached.
          </p>
        </div>
      </main >
    </>
  );
}
