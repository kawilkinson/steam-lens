# Steam Lens

[![ci](https://github.com/Khazz0r/steam-lens/actions/workflows/ci.yaml/badge.svg)](https://github.com/Khazz0r/steam-lens/actions/workflows/ci.yaml)

## What is Steam Lens?

Steam Lens is a website that gives you the ability to compare a user's Steam library and achievements against their friends at a moment's notice!

The user's friends are ranked starting from rank 1 at the top, I have the ranking currently weighted slightly more towards the *number* of common games a friend has over their percentage of common games.

Generally, as long as a user or friend doesn't have their profile privated, you should be able to see all of the matching games, the missing games, and the number of achievements both the user and their friend have achieved.

![steam-lens-demo](https://github.com/user-attachments/assets/3f6d970a-cbe0-48fd-a8e8-43c79079664e)

## Motivation

Steam Lens provides a very quick way to run comparisons of Steam libraries, usually this is a more manual process through Steam but this website makes it quick and easy to have all of those comparisons with just a simple input of your own Steam ID or another user's! I wanted to practice CI/CD through a full stack website (you'll find the cd.yaml is commented out due to AWS no longer hosting this website), understand frontend development better through a framework, further master my backend development skills through security and RESTful API building, and also sharpen my database skills working with a PostgreSQL database. So I decided to build Steam-Lens to do this.

## Project Layout

Currently the project can be broken up into 3 parts; it has a frontend, a backend, and a PostgreSQL database. The backend is in the root of the project, the frontend is in the frontend directory, and the database schema is located in the sql/schema directory.

## Installation / Quick Start

1. After cloning the repository, the easiest way to get the entire website running is to download and install Docker Desktop from their official website [here](https://www.docker.com/products/docker-desktop/). This way you don't have to worry about downloading or installing anything extra, Docker containers will handle all the dependencies and requirements with a clean environment to ensure stability.

2. Next set up a .env file in the root for the backend that contains this information:

```env
# Port to run the backend server on, feel free to use 8080 or any other port you want.
PORT=8080 

# Steam API Key is obtainable from Steam for free, you just need a Steam account to request one from the form link here https://steamcommunity.com/dev (usually instant)
STEAM_APIKEY={...} 

# A database URL for the backend to use to talk to the PostgreSQL database, the one below is what I use for testing purposes
DATABASE_URL="postgres://steam_lens:password@db:5432/steam_lens_db?sslmode=disable"

# Type of platform, if not set to "dev" then the backend will assume it is running in production which can change its behavior a bit (such as using SecureOnly cookies)
PLATFORM="dev"

# JWT set for security purposes (implemented for practice), feel free to use anything for this.
JWTSECRET="test"
```

3. Afterwards set up a .env.production file in the frontend directory of the project, this is a simple file that will only contain these variables:

```env
# Copy/paste both of these in a file named .env.production in the frontend directory, this is what the frontend server uses when set to production (the default)
NEXT_PUBLIC_API_URL=http://backend:8080/api/steam/
NEXT_PUBLIC_BACKEND_URL=http://localhost:8080/v1/
```

4. Finally, in the root of the project run this command in your operating system's terminal to build out the entire website:
```bash
docker compose up --build
```

5. **OPTIONAL** If you want to very quickly see how the website runs with just its main functionality without installing Docker (though you will need at least Go installed), then you can simply run these set of commands to get the frontend and backend running, but you will not be able to create or use user accounts since there will be no database:
```bash
# From one terminal window in the root run this
go run .
# From a second terminal window in the frontend directory run this
npm run dev
```

## Contribution
If you have any ideas for this project whether it'd be new features, optimizations, etc. feel free to contribute! I'm always open to new ideas and improvements.

## REST API Endpoints
### Users Endpoints

**UserCreate**

Creates a user profile

Endpoint: POST /v1/users/create

*Path Parameters*
```json
{
    "username": "user1@domain.com",
    "steam_id": "76561197997096401",
    "password": "Password123"
}
```

*Response*
```json
{
    "users": {
        "id": "6e32eed8-c431-4aec-b028-5bcbe1fbe79c",
        "created_at": "2025-03-14 23:15:42.123456789 +0000 UTC",
        "updated_at": "2025-03-14 23:15:42.123456789 +0000 UTC",
        "username": "user1@domain.com",
        "steam_id": "76561197997096401"
    }
}
```

**Login**

Logs user in

Endpoint: POST /v1/users/login

*Path Parameters*
```json
{
    "username": "user1@domain.com",
    "password": "Password123"
}
```
*Response*
```json
{
    "user": {
        "id": "6e32eed8-c431-4aec-b028-5bcbe1fbe79c",
        "created_at": "2025-03-14 23:15:42.123456789 +0000 UTC",
        "updated_at": "2025-03-14 23:15:42.123456789 +0000 UTC",
        "username": "user1@domain.com",
        "steam_id": "76561197997096401"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6Ikp",
    "refresh_token": "eyJzdh1D5kfiO2Dwslt3ODkwIiwibmFt"
}
```

**Logout**

Logs user out if they're already logged in, uses refresh token in cookie to ensure it matches up

Endpoint: POST /v1/users/logout

*Response*
```json
{
    "user": {
        "id": "6e32eed8-c431-4aec-b028-5bcbe1fbe79c",
        "username": "user1@domain.com",
        "steam_id": "76561197997096401"
    }
}
```

**GetMe**

Gets info about user

Endpoint: GET /v1/users/me

*Path Parameters*
```json
{
    "id": "6e32eed8-c431-4aec-b028-5bcbe1fbe79c"
}
```

*Response*
```json
{
    "id": "6e32eed8-c431-4aec-b028-5bcbe1fbe79c",
    "username": "user1@domain.com",
    "steam_id": "76561197997096401"
}
```

**UpdateUser**

Updates user profile with choice of new username, password, and/or Steam ID

Endpoint: PATCH /v1/users/me

*Path Parameters*
```json
{
    "username": "user2@domain.com",
    "password": "NewPassword!",
    "steam_id": "76561197997096419"
}
```
*Response*
```json
{
    "user": {
        "username": "user2@domain.com",
        "steam_id": "76561197997096419"
    }
}
```

### Steam Endpoints
[Here](https://developer.valvesoftware.com/wiki/Steam_Web_API#GetGlobalAchievementPercentagesForApp_.28v0001.29) is where you can view the parameters needed to make api calls to Steam manually.

**GetPlayerSummaries**

Gets basic profile information from a Steam ID

Endpoint: /api/steam/player-summaries

*Response*
```json
{
    "steamID": "76561197997096401",
    "communityVisibilityState": 3,
    "personaName": "user",
    "avatar": "https://avatars.steamstatic.com/example.jpg",
    "avatarMedium": "https://avatars.steamstatic.com/example_medium.jpg",
    "avatarFull": "https://avatars.steamstatic.com/example_full.jpg"
}
```

**GetOwnedGames**

Get all owned games for player from Steam ID

Endpoint: /api/steam/games

*Response*
```json
{
    "game_count": 25,
    "games": [
        {
        "appID": 379720,
        "name": "DOOM",
        "img_icon_url": "https://cdn.fastly.steamstatic.com/steamcommunity/public/images/apps/379720/b6e72ff47d1990cb644700751eeeff14e0aba6dc.jpg"
        }
    ]
}
```

**GetFriendList**

Get all friends of user from their Steam ID

Endpoint: /api/steam/friends

*Response*
```json
{
    "friends": [
        {
            "steamID": "76561197997096401",
            "communityVisibilityState": 3,
            "personaName": "user",
            "avatar": "https://avatars.steamstatic.com/example.jpg",
            "avatarMedium": "https://avatars.steamstatic.com/example_medium.jpg",
            "avatarFull": "https://avatars.steamstatic.com/example_full.jpg"
        }
    ]
}
```

**GetPlayerAchievements**

Get all achievements of a user from their Steam ID

Endpoint: /api/steam/compare-achievements

*Response*
```json
{
    "achievements": [
        {
            "apiName": "jumped_500_times",
            "achieved": true
        }
    ] 
}
```
