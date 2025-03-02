# Wheatley

[![codecov](https://codecov.io/gh/tecnologer/wheatley/graph/badge.svg?token=4FBM21AOBC)](https://codecov.io/gh/tecnologer/wheatley)

Wheatley is a telegram bot that notifies when a twitch streamer is online.

![Wheatley, from Portal 2][1]


You can use it on Telegram by adding the bot [@TwitchNotificationsMxBot](https://t.me/@TwitchNotificationsMxBot).

## Usage

1. Download the latest release from the [releases page][2].
2. Run the binary with the required flags.

```bash
./wheatley --telegram-token your_telegram_bot_token \
              --db-name your_db_name \
              --db-password your_db_password \
              --twitch-client-id your_twitch_client_id \
              --twitch-client-secret your_twitch_client_secret
```

## Build

### Docker

Create a .env file with the required environment variables.

```
WHEATLEY_TWITCH_CLIENT_ID=your_twitch_client_id
WHEATLEY_TWITCH_CLIENT_SECRET=your_twitch_client_secret
WHEATLEY_DB_PASSWORD=your_db_password
WHEATLEY_DB_NAME=your_db_name
WHEATLEY_TELEGRAM_BOT_TOKEN=your_telegram_bot_token
WHEATLEY_INTERVAL=1 # in minutes
WHEATLEY_RESEND_INTERVAL=6 # in hours
```

Build the docker image.
```bash
docker build -t wheatley .
```

Run the docker container.
```bash
docker run --env-file .env wheatley
```

> Note: add flags `-d` to run the container in detached mode and `--restart unless-stopped` to restart the container if it stops.
> ```bash
> docker run -d --restart unless-stopped --env-file .env wheatley
> ```

### Local

Build the binary.

```bash
go build -o wheatley cmd/main.go
```

Run the binary.

```bash
./wheatley --telegram-token your_telegram_bot_token \
              --db-name your_db_name \
              --db-password your_db_password \
              --twitch-client-id your_twitch_client_id \
              --twitch-client-secret your_twitch_client_secret
```

> Note: add flags `--interval` to set the interval in minutes to check if a streamer is live and `--resend-interval` to set the interval in hours to resend a notification.

The `--help` flag is available to show the available commands and flags.

```bash
./wheatley --help
```

## ToDo

- [ ] configure client id
- [ ] configure client secret
- [x] get twitch API token
- [x] get twitch user info
- [x] get twitch stream info
- [x] resent notification every `--resend-interval`
- [ ] command to add new admins
- [ ] update delay between notifications
- [ ] multi-language (ESP/ENG)

[1]: https://i1.theportalwiki.net/img/thumb/9/94/Wheatley.png/300px-Wheatley.png
[2]: https://github.com/tecnologer/wheatley/releases
