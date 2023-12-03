# Squash Slack Bot

This bot is meant to serve as a way to manage the ladder and the challenge matches for every season. Players will utilize this bot to see their statistics and current standings on the team.

Currently this is setup to build with docker, will soon be hosted on a dedicated server for the Northeastern Squash Team

# How to run locally:

## Environment Setup:

- Create an .env file with a Slack bot/app token, you have two options. You can create a new slack bot and use your own tokens for your new bot. Or you can contact the owner of this respository (Ryan) if you wish to run the the original bot that was created for Northeastern Squash. (Please see .env.template for initial setup of your .env file)
  - SLACK_BOT_TOKEN
  - SLACK_APP_TOKEN

## Building Image:

- Run the following if you have just cloned this repo:

```
export GO111MODULE="on"
go mod tidy
```

- Then you build the containers and spin up docker, the bot container will await for db to be initialized and ready before spinning up (this may take a minute or two)

```
docker compose build
docker compose up -d
```

- To remove containers and stop running service:

```
docker compose down -v
```

# Want to know how to use this bot?

- See here: ![commands list](./commands/README.md)
