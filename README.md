# tiktok-fav-publisher
Reposts to Telegram channel liked videos from TikTok without any watermarks.

# How to use
1. Fill with ENVs `docker-compose.yml`
   1. For `TG_TOKEN` env: create bot via `@BotFather` and paste the token
   2. For `CHANNEL_ID` env: create telegram channel, post something to channel and forward post to `@JsonDumpBot` bot. `id` from `forward_from_chat` is your `CHANNEL_ID`
   3. For `TIKTOK_USERNAME` env: paste your TikTok username (don't forget to open your liked list)
2. Run `docker compose build server`
3. Run `docker-compose pull`
4. Run `docker-compose up -d`

# Features
1. Uses TikTok api (~1 req per 10 sec limitation)
2. Pastes video without watermark