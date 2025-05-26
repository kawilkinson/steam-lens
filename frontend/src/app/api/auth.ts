import { backendURL } from "../definitions/urls";

export async function createAccount(email: string, password: string, steamID: string) {
  if (!backendURL) {
    throw new Error("Backend base URL not set up in environment variables");
  }
  const url = new URL("users/create", backendURL);

  try {
    const resp = await fetch(url.toString(), {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ username: email, password, steam_id: steamID }),
    });

    if (resp.status >= 400) {
      const data = await resp.json();
      throw new Error(data.error || `Failed with status ${resp.status}`);
    }
    return await resp.json();
  } catch (err) {
    throw err;
  }
}

export async function login(email: string, password: string) {
  if (!backendURL) {
    throw new Error("Backend base URL not set up in environment variables");
  }
  const url = new URL("users/login", backendURL);

  try {
    const resp = await fetch(url.toString(), {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ username: email, password }),
    //   credentials: "include", (commenting this out for now, but will be used in prod for HttpOnly cookies)
    });

    if (resp.status >= 400) {
      const data = await resp.json();
      throw new Error(data.error || `Failed with status ${resp.status}`);
    }
    return await resp.json();
  } catch (err) {
    throw err;
  }
}
