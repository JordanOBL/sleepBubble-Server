# SleepBubble Server

A small Go HTTP server used by the SleepBubble mobile app. It exposes a few endpoints for checking or updating Lennox's sleep status and for subscribing devices to receive push notifications via Expo.

## Requirements

- Go 1.23+
- (optional) Docker

## Running locally

```bash
# run in development
PORT=3000 go run cmd/server/main.go
```

By default the server looks for its CSV database at `/app/cmd/server/sleepbubble.csv`. If running without Docker, ensure this file exists or adjust the path in the code.

## Docker

```bash
# build and run the container
docker build -t sleepbubble-server .
docker run -p 3000:3000 sleepbubble-server
```

## API Endpoints

| Method | Path             | Description                           |
| ------ | ---------------- | ------------------------------------- |
| GET    | `/sleepstatus`   | Current sleep status and a fun quote. |
| POST   | `/join`          | Add an Expo push token.               |
| POST   | `/updatesleep`   | Update sleep status and notify users. |
| GET    | `/v1/healthcheck`| Simple health check.                  |

### `/sleepstatus` response

```json
{
  "sleepStatus": "0", // 0 = awake, 1 = sleeping
  "statement": "Lennox is up and ready to rock…"
}
```

### `/join` body

Plain text Expo token, for example:

```
ExponentPushToken[xxxxxxxxxxxxxxxxxxxxxx]
```

### `/updatesleep` body

`0` for awake or `1` for sleeping. All subscribed tokens will receive a push notification when this endpoint succeeds.

## Project Structure

- `cmd/server` – entry point and CSV database
- `internal/server` – HTTP handlers and helpers

## License

This project currently does not specify a license.
