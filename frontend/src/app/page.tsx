'use client'

import { useState } from "react";
import styles from "./page.module.css";
import { useRouter } from "next/navigation";

export default function Home() {
  const [userID, setUserID] = useState<string>("");
  const router = useRouter();

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
        <h1>Steam Lens</h1>
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
          <h2>What is Steam Lens?</h2>
          <p>Steam Lens is a website that allows comparing a user&apos;s owned games and achievements with their friends.</p>
          <p>The user&apos;s friends are ranked according to the similarity of their game libraries. The more common games the friend has to the user, the higher they will be placed.</p>
          <p>
            Steam Lens also gives the ability to see what games a friend may have that you don&apos;t own, or vice versa.
          </p>
        </div>
        <div className={styles.section}>
          <h2>How to Use</h2>
          <p>
            Enter a 17 digit numeric Steam ID of your own or a user you want to run comparisons on.
          </p>
          <p>
            Once entered, clicking submit will redirect you over to the main dashboard where you&apos;ll be able to compare games and achievements.
          </p>
          <p>
            You can also paste a user&apos;s profile URL from Steam, and the form should extract the numeric part automatically.
          </p>
        </div>
        <div className={styles.section}>
          <h2>Limitations</h2>
          <p>
            Since this website uses the public Steam API to get the users&apos; game libraries and friends lists, it is limited by the 100,000 daily calls they have set up.
          </p>
          <p>
            Some effort was made to cache that information whenever possible to reduce the number of calls made, but it remains a possible issue.
          </p>
        </div>
      </main >
    </>
  );
}
