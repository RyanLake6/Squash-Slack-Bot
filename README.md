## Squash Slack Bot

Currently this is only set to run locally but will hopefully be hosted and funded by the Northeastern Squash Team

# How to run locally:

- On your terminal enter into the correct folder of this repo
- Run the following (dependencies):

```
export GO111MODULE="on"
go mod tidy
```

- Then you can run the bot (assuming sql server is set up):

```
go run main.go
```

# Want to know how to use this bot?

- See here: ![commands list](./commands/README.md)
