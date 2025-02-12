# EbbiLogue Backend

A backend system for My personal website [EbbiLogue](https://ebbilogue.com)
frontend code [here](https://github.com/camtrik/ebbilogue)

## Features
- PSN trophy data fetching (filtering, sorting, pagination)
- Steam game library and achievements tracking
- Redis caching support for better performance
- Cross-platform gaming progress tracking

## TODO 
[] AI Chat Avatar  
[] Comment System

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


## API Documentation

### PSN Endpoints

#### GetUserTitles
```http
http://localhost:6061/api/psn/:accountId/trophyTitles
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
```http
http://localhost:6061/api/steam/:steamId/ownedGames
```
Get user's owned Steam games.

#### GetPlayerAchievements
```http
http://localhost:6061/api/steam/:steamId/playerAchievements/:appId
```
Get achievements for a specific game.

#### GetPlayerGameDetails
```http
http://localhost:6061/api/steam/:steamId/playerGameDetails
```
Get detailed information about player's games including achievements and playtime, game logo & banner art.

Query Parameters:
- `minPlayTime`: Filter games by minimum playtime 
- `sortByTime`: Sort games by playtime (true/false)

