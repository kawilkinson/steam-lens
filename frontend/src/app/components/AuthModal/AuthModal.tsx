import React, { useState } from "react";
import { useRouter } from "next/navigation";
import styles from "../../page.module.css"
import { createAccount, login } from "../../api/auth";

interface AuthModalProps {
  type: "signup" | "login";
  onClose: () => void;
  setLoading: (val: boolean) => void;
}

const AuthModal: React.FC<AuthModalProps> = ({ type, onClose, setLoading }) => {
  const [fields, setFields] = useState({ username: "", password: "", steam_id: "" });
  const [error, setError] = useState("");
  const router = useRouter();

  const handleChange = (err: React.ChangeEvent<HTMLInputElement>) =>
  setFields({ ...fields, [err.target.name]: err.target.value });

  const handleSubmit = async (err: React.FormEvent) => {
  err.preventDefault();
  setError("");
  setLoading(true);
  try {
      let data;
      if (type === "signup") {
        data = await createAccount(fields.username, fields.password, fields.steam_id);
      } else {
        data = await login(fields.username, fields.password);
      }

      const steamID = data.user?.steam_id;
      if (steamID) {
        onClose();
        router.push(`/${steamID}`);
      } else {
        setError("Missing Steam ID in server response.");
      }
  } catch (err) {
      if (err instanceof Error) {
        setError(err.message);
      } else {
        setError("An unknown error occurred");
      }
    }
  };

  return (
    <div className={styles.modalOverlay}>
      <div className={styles.modalBox}>
        <button className={styles.closeBtn} onClick={onClose}>&times;</button>
        <div className={styles.modalTitle}>{type === "signup" ? "Create Account" : "Log In"}</div>
        <form onSubmit={handleSubmit}>
          <input
            className={styles.modalInput}
            type="email"
            name="username"
            required
            placeholder="Email"
            value={fields.username}
            onChange={handleChange}
            autoComplete="username"
          />
          <input
            className={styles.modalInput}
            type="password"
            name="password"
            required
            placeholder="Password"
            value={fields.password}
            onChange={handleChange}
            autoComplete={type === "signup" ? "new-password" : "current-password"}
          />
          {type === "signup" && (
            <input
              className={styles.modalInput}
              type="text"
              name="steam_id"
              required
              placeholder="Steam ID"
              value={fields.steam_id}
              onChange={handleChange}
            />
          )}
          {error && <div className={styles.modalError}>{error}</div>}
            <button className={styles.modalSubmit} type="submit">
            {type === "signup" ? "Create Account" : "Log In"}
            </button>
        </form>
      </div>
    </div>
  );
};

export default AuthModal;
