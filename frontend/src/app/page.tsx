'use client'

import { useState } from "react";
import styles from "./page.module.css";
import { useRouter } from "next/navigation";
import AuthModal from "./components/AuthModal/AuthModal";
import { logout } from "./api/auth";
import { useAuth } from "./context/AuthContext"

export default function Home() {
  const [userID, setUserID] = useState("");
  const [showLogin, setShowLogin] = useState(false);
  const [showSignup, setShowSignup] = useState(false);
  const [showSignupSuccess, setShowSignupSuccess] = useState(false);
  const [showLoginSuccess, setShowLoginSuccess] = useState(false);
  const [loading, setLoading] = useState(false);

  const { auth, setAuth, checkAuth } = useAuth();
  const router = useRouter();

  function handleLoginSuccess() {
    checkAuth();
    setShowLogin(false);
    setShowSignup(false);
    setShowLoginSuccess(true);
    setTimeout(() => setShowLoginSuccess(false), 5000)
  }

  function handleSignupSuccess() {
    setShowSignupSuccess(true);
    setTimeout(() => setShowSignupSuccess(false), 5000)
  }

  function handleLogout() {
    setLoading(true);
    logout().finally(() => {
      setAuth({ loggedIn: false, steamID: null });
      setLoading(false);
      router.push("/");
    });
  }

  const handleViewStats = () => {
    if (auth.steamID) {
      setLoading(true);
      router.push(`/${auth.steamID}`);
    }
  }

  // Allows for an optional Steam Community profile URL that contains the Steam ID to be submitted instead
  function handleSteamIDSubmit() {
    let destination = userID;
    const match = userID.match(/https:\/\/steamcommunity\.com\/profiles\/(\d+)/);
    if (match !== null) {
      destination = match[1];
    }
    setLoading(true);
    router.push(`/${destination}`);
  }

  return (
    <>
      <header className={styles.header}>
        <h1 className={styles.logoTitle}>Steam Lens</h1>
        <div>
          {auth.loggedIn ? (
          <>
            <button className={styles.formButton} onClick={handleViewStats}>View Stats</button>
            <button className={styles.formButton} onClick={() => router.push("/edit-account")} style={{ marginLeft: 12 }}>Edit Account</button>
            <button className={styles.formButton} onClick={handleLogout} style={{ marginLeft: 12 }}>Log Out</button>
          </>
          ) : (
            <>
              <button className={styles.formButton} onClick={() => setShowLogin(true)}>Log In</button>
              <button className={styles.formButton} onClick={() => setShowSignup(true)} style={{marginLeft: 12}}>Create Account</button>
            </>
          )}
        </div>
      </header>

      {showLogin && (
        <AuthModal type="login" onClose={() => setShowLogin(false)} onSuccess={handleLoginSuccess} />
      )}
      {showSignup && (
        <AuthModal type="signup" onClose={() => setShowSignup(false)} onSuccess={handleSignupSuccess} />
      )}

      {loading && (
        <div className={styles.loadingOverlay}>
          <div className={styles.loadingBox}>
            <div className={styles.loadingSpinner}></div>
            Loading stats&hellip;
          </div>
        </div>
      )}

      {showSignupSuccess && (
        <div className={styles.signupSuccessPopup}>
          Account created successfully
        </div>
      )}

      {showLoginSuccess && (
        <div className={styles.loginSuccessPopup}>
          Successfully logged in
        </div>
      )}

      <main className={styles.main}>
        <div className={styles.section}>
          <div className={styles.form}>
            <label htmlFor="userID">Steam ID:</label>
            <input name="userID" type="text" required value={userID} onChange={(val) => setUserID(val.target.value)} />
            <button type="submit" onClick={handleSteamIDSubmit}>Submit</button>
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
            Once entered, clicking submit will redirect you over to the main dashboard where you&apos;ll be able to compare games and achievements 
            (there is a load time depending on how many friends the user has, see under &quot;limitations&quot; at the bottom of this page.).
          </p>
          <p>
            Note, if the user or a friend has their profile privated in any way then that can prevent data from the API to return for that account.
            If you notice data isn&apos;t loading for a specific person then apart from Steam&apos;s API limits being reached that would 
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
            They have a 100,000 daily API call limit and I have also observed a short-term limit if you make too many API calls at once in a short period of time. 
            If you notice errors being thrown or data not loading, then it is likely due to that limit being reached since user/friend data, game data, and achievement data are all gotten from Steam&apos;s endpoints.
          </p>
          <p>
            I have added backend caching of the data returned from Steam&apos;s endpoints to reduce API calls as much as possible to make it unlikely this daily limit is reached.
          </p>
          <p>
            On top of this, I have also added a wait time for stats to load depending on how many friends a user has. For example if a user has over 100 friends
            then it will take 20 seconds to load bare minimum, if more than 50 friends but less than 100 it will take at least 10 seconds, and if less than 50 friends it will take 5 seconds. 
            This is to help prevent &quot;too many request&quot; errors from Steam&apos;s API.
          </p>
        </div>
      </main >
    </>
  );
}
