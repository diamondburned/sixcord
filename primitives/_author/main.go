package main

import (
	"github.com/Bios-Marcel/tview"
	"github.com/bwmarrin/discordgo"
	"github.com/gdamore/tcell"
	"gitlab.com/diamondburned/sixcord/primitives"
)

var author = map[string]*discordgo.User{
	"diamond": &discordgo.User{
		ID:            "170132746042081280",
		Username:      "diamondburned",
		Discriminator: "1337",
		Avatar:        "5a8a9f925f78304ac3e9d313f836cf26",
	},
	"ym": &discordgo.User{
		ID:            "274573296199270401",
		Username:      "ym555",
		Discriminator: "1339",
		Avatar:        "39337a9d32a1bb766a1a331fdc168c41",
	},
}

// contains dumb test samples
var messages = []*discordgo.Message{
	&discordgo.Message{
		Content: `This is an essay. Believe me, it is, When I say it is, it usually is.
		
	It's ok, you don't have to trust me, you just need to believe that this is an essay. That's right. It's an essay, It's long, nobody ever reads all this anyway, so just believe that it is an essay.
	
	Is that not convincing enough? How about some statistics? Did you know that 100% of living human beings breathe in oxygen and breathe out carbon dioxide? This is a fact. Very convincing statistics for an essay.
	
	It can be concluded that this is a perfectly fine essay, and you should just go along with it.`,
		Author: author["diamond"],
		// Image taken from sixtyten#1559, Turquoise Hexagon
		Attachments: []*discordgo.MessageAttachment{
			&discordgo.MessageAttachment{
				URL:      "https://cdn.discordapp.com/attachments/389480585216786432/568780229078941709/2019-04-19_14-06-44.jpeg",
				ProxyURL: "https://media.discordapp.net/attachments/389480585216786432/568780229078941709/2019-04-19_14-06-44.jpeg?width=400&height=235",
				Width:    400,
				Height:   235,
			},
		},
	},
	&discordgo.Message{
		Content: "yeah that \"essay\" was probably the dumbest thing i've written for a test lamo",
		Author:  author["diamond"],
	},
	&discordgo.Message{
		Content: "have anyone seen ym? @ym555#1339",
		Author:  author["diamond"],
	},
	&discordgo.Message{
		Content: "wat",
		Author:  author["ym"],
	},
	&discordgo.Message{
		Content: "help pls",
		Author:  author["diamond"],
	},
	&discordgo.Message{
		Content: "my penis is broken, how do I send a GET request to nhentai for a hello teapot response?",
		Author:  author["diamond"],
	},
	&discordgo.Message{
		Content: "what? :thonkcerned:",
		Author:  author["ym"],
	},
}

func main() {
	/*
		var msgs = make([]*discordgo.Message, 0, len(messages))
		for _, m := range messages {
			if m.Author.Username == "diamondburned" {
				msgs = append(msgs, m)
			}
		}
	*/

	var msgs = messages

	a, err := primitives.NewAuthor(msgs[0].Author)
	if err != nil {
		panic(err)
	}

	for _, m := range msgs {
		if err := a.AddMessage(m); err != nil {
			panic(err)
		}
	}

	app := tview.NewApplication()
	app.SetBeforeDrawFunc(func(s tcell.Screen) bool {
		s.Clear()
		return false
	})

	if err := app.SetRoot(a, true).Run(); err != nil {
		panic(err)
	}
}
