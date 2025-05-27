"use client";
import { useRouter } from "next/navigation";
import { logout } from "../../api/auth";

export default function LogoutButton() {
  const router = useRouter();
  return (
    <button
      className="formButton"
      style={{marginLeft: 10}}
      onClick={() => {
        logout().finally(() => {
          router.push("/");
        });
      }}
    >
      Log Out
    </button>
  );
}
