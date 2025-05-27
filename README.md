# PSN & Steam & Xbox APIs in Go 

A backend service aggregating PlayStation Network and Steam gaming data through their official APIs. Provides unified endpoints for tracking game progress, achievements and trophies across platforms.

**Current Status**: 🚧 In active development - Core PSN/Steam integrations working, more features and documentation WIP

> Live Demo using these APIs: [Gaming Profile](https://www.ebbilogue.com/gaming)

## Key Features
- 📦 PSN user authentication & trophy data fetching 
- 🎮 Steam game library & achievements tracking
- ⚔️ Xbox APIs
- ⚡ Redis caching for high-performance queries
- 🌉 Cross-platform gaming progress aggregation
- 🔄 Unified API for accessing gaming data across all platforms

## Development Status
### 2025/04/25 Xbox APIs
Integrated Xbox APIs using [OPENXBL](https://xbl.io/)

### Implemented
✅ PSN Authentication & Trophy Title Listing  
✅ Steam Game Library Retrieval  
✅ Basic Caching Mechanism  
✅ Xbox APIs by [OPENXBL](https://xbl.io/)  
✅ Unified API for all three platforms

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

### Unified Endpoints

#### GetAllRecentlyPlayedGames
```
http://localhost:7071/api/unified/recentlyPlayed
```
Get user's recently played games across all platforms (PSN, Steam, and Xbox). Results are combined and sorted by the most recently played games first.

Query Parameters:
- `psn_account_id`: (string) Optional. PSN account ID. If not provided, defaults to "me" (your own account).
- `steam_id`: (string) Optional. Steam ID. If not provided, uses the default Steam ID from the configuration.
- `time_range`: (string) Optional. Time range to filter games played within this timeframe. Available options are:
  - two_weeks: Games played within the last 2 weeks.
  - one_month: Games played within the last 1 month (default).
  - three_months: Games played within the last 3 months.

This endpoint returns an aggregated list of recently played games from all configured platforms, making it easy to track gaming activity in one place.

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

#### GetRecentlyPlayedGames
```
http://localhost:7071/api/psn/:accountId/recentlyPlayed
```
Get user's recently played PSN games.

Query Parameters:
- `timeRange`: (string) Optional.Time range to filter games played within this timeframe. Available options are
    - two_weeks: Games played within the last 2 weeks.
    - one_month: Games played within the last 1 month (default).
    - three_months: Games played within the last 3 months.

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
    - one_month: Games played within the last 1 month (default).
    - three_months: Games played within the last 3 months.


### Xbox Endpoints 
To use Xbox endpoints, you need to first get an api key from [OPENXBL](https://xbl.io/)  

Noticed that APIs implemented here only can fetch data from **your own account**. 
#### GetPlayerAchievements
```
http://localhost:7071/api/xbox/achievements
```
Get the user's Xbox achievements across games.

#### GetGameStats
```
http://localhost:7071/api/xbox/gameStats/:titleId
```
Get detailed statistics for a specific Xbox game.

#### GetRecentlyPlayedGames
```
http://localhost:7071/api/xbox/recentlyPlayedGames
```
Get user's recently played Xbox games.

Query Parameters:
- `timeRange`: (string) Optional. Time range to filter games played within this timeframe. Available options are:
    - two_weeks: Games played within the last 2 weeks.
    - one_month: Games played within the last 1 month (default).
    - three_months: Games played within the last 3 months.
