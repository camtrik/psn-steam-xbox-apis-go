# EbbiLogue Backend

A backend system for My personal website [EbbiLogue](https://ebbilogue.com)
frontend code [here](https://github.com/camtrik/ebbilogue)

Only PSN features are available for now


## Features
- PSN trophy data fetching (filtering)
- Redis support

## TODO 
[] AI Chat Avatar
[] Comment System

## Prerequisites
- Go 1.22
- Redis

## Quick Start
```bash
git clone git@github.com:camtrik/ebbilogue-backend.git

# Set up env, set your own IDs or tokens
cp .env.example .env

sudo service redis-server start

go run cmd/main.go
```

