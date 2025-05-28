"use client";
import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import styles from "../page.module.css";
import editStyles from "./page.module.css"
import { useAuth } from "../context/AuthContext";
import { editAccount, fetchCurrentUser } from "../api/auth";

export default function EditAccountPage() {
  const { auth, checkAuth } = useAuth();
  const router = useRouter();

  const [fields, setFields] = useState({ username: "", password: "", steam_id: "" });
  const [original, setOriginal] = useState({ username: "", steam_id: "" });
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");

  useEffect(() => {
    if (!auth.loggedIn) {
      setLoading(false);
      setError("You must be logged in to edit your account.");
      return;
    }

    async function fetchData() {
      setLoading(true);
      try {
        const data = await fetchCurrentUser();
        setFields({
          username: data.user?.username ?? "",
          password: "",
          steam_id: data.user?.steam_id ?? ""
        });
        setOriginal({
          username: data.user?.username ?? "",
          steam_id: data.user?.steam_id ?? ""
        });
      } catch (err) {
        throw err
      }
      setLoading(false);
    }
    fetchData();
  }, [auth.loggedIn]);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError("");
    setSuccess("");
    setLoading(true);

    const payload: Record<string, string> = {};
    if (fields.username && fields.username !== original.username) payload.username = fields.username;
    if (fields.password) payload.password = fields.password;
    if (fields.steam_id && fields.steam_id !== original.steam_id) payload.steam_id = fields.steam_id;

    if (Object.keys(payload).length === 0) {
        setError("No changes to save.");
        setLoading(false);
        return;
    }

    try {
        await editAccount(payload);
        setSuccess("Account updated!");
        setFields({ ...fields, password: "" });
        checkAuth();
        setLoading(false);
        router.push("/");
    } catch (err) {
        setError(err instanceof Error ? err.message : "An unknown error occurred");
        setLoading(false);
    }
    }

  if (loading) {
    return (
      <div className={styles.loadingOverlay}>
        <div className={styles.loadingBox}>
          <div className={styles.loadingSpinner}></div>
            Loading account&hellip;
        </div>
      </div>
    );
  }

  if (!auth.loggedIn) {
    return (
      <main className={styles.main}>
        <h2>Edit Account</h2>
        <div className={styles.modalError}>
          You must be logged in to edit your account.
        </div>
      </main>
    );
  }

  return (
    <>
      <header className={styles.header}>
        <h1 className={styles.logoTitle}>Steam Lens</h1>
      </header>
      <main className={styles.main} style={{ marginLeft: 20}}>
        <h2 style={{ marginBottom: 16 }}>Edit Account</h2>
        {error && <div className={styles.modalError}>{error}</div>}
        {success && <div className={styles.modalSuccess}>{success}</div>}
        <form onSubmit={handleSubmit} style={{ maxWidth: 400, width: "100%" }}>
          <label style={{ display: "block", marginBottom: 10 }}>
            Email:
            <input
              className={styles.modalInput}
              type="email"
              name="username"
              value={fields.username}
              onChange={e => setFields({ ...fields, username: e.target.value })}
              required
              style={{ marginTop: 4 }}
            />
          </label>
          <label style={{ display: "block", marginBottom: 10 }}>
            New Password:
            <input
              className={styles.modalInput}
              type="password"
              name="password"
              value={fields.password}
              onChange={e => setFields({ ...fields, password: e.target.value })}
              placeholder="(leave blank to keep current password)"
              style={{ marginTop: 4 }}
            />
          </label>
          <label style={{ display: "block", marginBottom: 16 }}>
            Steam ID:
            <input
              className={styles.modalInput}
              type="text"
              name="steam_id"
              value={fields.steam_id}
              onChange={e => setFields({ ...fields, steam_id: e.target.value })}
              required
              style={{ marginTop: 4 }}
            />
          </label>
            <div style={{ display: "flex", gap: 10, marginTop: 10 }}>
                <button className={styles.formButton} type="submit" disabled={loading}>
                {loading ? "Saving..." : "Save Changes"}
                </button>
                <button
                className={`${styles.formButton} ${editStyles.cancel}`}
                type="button"
                onClick={() => router.push("/")}
                >
                Cancel
                </button>
            </div>
        </form>
      </main>
    </>
  );
}
