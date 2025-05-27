"use client";
import { useAuth } from "@/app/context/AuthContext";
import styles from "../../[steamid]/page.module.css";
import homeStyles from "../../page.module.css"
import Link from "next/link";
import { useRouter } from "next/navigation";
import { logout } from "../../api/auth";

export default function MainHeader() {
  const { auth, setAuth } = useAuth();
  const router = useRouter();

  async function handleLogout() {
    await logout();
    setAuth({ loggedIn: false, steamID: null });
    router.push("/");
  }

  return (
    <header className={styles.header}>
      <h1 className={styles.logoTitle}>Steam Lens</h1>
      <div style={{ display: "flex", gap: 12, alignItems: "center" }}>
        <Link href="/" className={homeStyles.formButton}>Home</Link>
        {auth.loggedIn && (
          <button className={homeStyles.formButton} style={{marginLeft: 0}} onClick={handleLogout}>Log Out</button>
        )}
      </div>
    </header>
  );
}
