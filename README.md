# Squash Slack Bot

Currently this is only set to run locally but will hopefully be hosted and funded by the Northeastern Squash Team

# How to run locally:

## Environment Setup:

- Create an env file with the following sql variables associated with you own hosted sql server:
  - DBUSER
  - DBPASS
- You will also need the Slack bot/app token, you have two options. YOu can create a new slack bot and use your own tokens for your new bot. Or you can contact the owner of this respository (Ryan) if you wish to run the the original bot that was created for Northeastern Squash.
  - SLACK_BOT_TOKEN
  - SLACK_APP_TOKEN

## Building Image:

- Run the following if you have just cloned this repo:

```
export GO111MODULE="on"
go mod tidy
```

- Then you build the containers and spin up docker

```
docker compose build
docker compose up -d
```

- If you wish to spin eveything down run:

```
docker compose down -v
```

# Want to know how to use this bot?

- See here: ![commands list](./commands/README.md)
