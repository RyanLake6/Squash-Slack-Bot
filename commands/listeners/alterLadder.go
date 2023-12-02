package listeners

import (
	"github.com/shomali11/slacker"
)

func (c *Client) AlterLadder() {
	c.Bot.Command("alter ladder <position> <newPlayer>", &slacker.CommandDefinition{
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			user := botCtx.Event().UserName
			if user != "Resistance Leader" {
				response.Reply("Hello there, looks like you're not good enough to do this " + user + "\n Get good and try next year!")
			} else {
				// TODO implement a alter ladder ability here
			}
		},
	})
}