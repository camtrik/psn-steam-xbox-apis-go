# PSN Steam API in Go 

A backend service aggregating PlayStation Network and Steam gaming data through their official APIs. Provides unified endpoints for tracking game progress, achievements and trophies across platforms.

**Current Status**: 🚧 In active development - Core PSN/Steam integrations working, more features and documentation WIP

> Live Demo using these APIs (frontend [here](https://github.com/camtrik/psn-steam-api-go)): [Gaming Profile](https://www.ebbilogue.com/gaming)

## Key Features
- 📦 PSN user authentication & trophy data fetching 
- 🎮 Steam game library & achievements tracking
- ⚡ Redis caching for high-performance queries
- 🌉 Cross-platform gaming progress aggregation

## Development Status

### Implemented
✅ PSN Authentication & Trophy Title Listing  
✅ Steam Game Library Retrieval  
✅ Basic Caching Mechanism

### In Progress
🛠️ More APIs implementation  
🛠️ Rate Limiting for API Endpoints  
🛠️ Comprehensive API Documentation  
🛠️ Better Error Handling

## Prerequisites
- Go 1.22
- Redis

## Quick Start
Install redis first:
```bash
sudo apt update
sudo apt install redis-server
sudo service redis-server start
```

Clone the repository and run the application:
```bash
git clone git@github.com:camtrik/ebbilogue-backend.git

# Set up env, set your own IDs or tokens
cp .env.example .env

make run 
```

How to get PSN refresh token: [here](https://www.ebbilogue.com/blog/notes/psn-api-use)

How to get steam api key: [here](https://steamcommunity.com/dev)  

How to get Xbox api key: [here](https://xbl.io/console) (Using OPENXBL for Xbox APIs)


## API Documentation

### PSN Endpoints

#### GetUserTitles
```
http://localhost:7071/api/psn/:accountId/trophyTitles
```
Get user's PSN trophy titles with filtering and pagination options.

Query Parameters:
- `limit`: Number of items per page
- `offset`: Starting position
- `platform`: Filter by platform (e.g., PS4, PS5)
- `minProgress`: Filter by minimum trophy completion percentage
- `sortBy`: Sort results by ("lastUpdated" or "progress")

### Steam Endpoints

#### GetOwnedGames
```
http://localhost:7071/api/steam/:steamId/ownedGames
```
Get user's owned Steam games.

#### GetPlayerAchievements
```
http://localhost:7071/api/steam/:steamId/playerAchievements/:appId
```
Get achievements for a specific game.

#### GetPlayerGameDetails
```
http://localhost:7071/api/steam/:steamId/playerGameDetails
```
Get detailed information about player's games including achievements and playtime, game logo & banner art.

Query Parameters:
- `minPlayTime`: Filter games by minimum playtime 
- `sortByTime`: Sort games by playtime (true/false)

#### GetRecentlyPlayedGames
```
http://localhost:7071/api/steam/:steamId/recentlyPlayedGames
```
Get user's recently played Steam games.

Query Parameters:
- `timeRange`: (string) Optional.Time range to filter games played within this timeframe. Available options are
    - two_weeks: Games played within the last 2 weeks.
    - one_month: Games played within the last 1 month.
    - three_months: Games played within the last 3 months.