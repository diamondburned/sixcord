package primitives

import (
	"testing"

	"github.com/Bios-Marcel/tview"
	"github.com/bwmarrin/discordgo"
	"github.com/gdamore/tcell"
)

func TestMessage(t *testing.T) {
	m := &discordgo.Message{
		Content: `This is an essay. Believe me, it is, When I say it is, it usually is.
		
	It's ok, you don't have to trust me, you just need to believe that this is an essay. That's right. It's an essay, It's long, nobody ever reads all this anyway, so just believe that it is an essay.
	
	Is that not convincing enough? How about some statistics? Did you know that 100% of living human beings breathe in oxygen and breathe out carbon dioxide? This is a fact. Very convincing statistics for an essay.
	
	It can be concluded that this is a perfectly fine essay, and you should just go along with it.`,
		Author: &discordgo.User{
			Username:      "diamondburned",
			Discriminator: "1337",
		},
		// Image taken from sixtyten#1559, Turquoise Hexagon
		Attachments: []*discordgo.MessageAttachment{
			&discordgo.MessageAttachment{
				URL:      "https://cdn.discordapp.com/attachments/389480585216786432/568780229078941709/2019-04-19_14-06-44.jpeg",
				ProxyURL: "https://media.discordapp.net/attachments/389480585216786432/568780229078941709/2019-04-19_14-06-44.jpeg?width=400&height=235",
				Width:    400,
				Height:   235,
			},
		},
	}

	mv, err := NewMessage(m)
	if err != nil {
		t.Fatal(err)
	}

	app := tview.NewApplication()
	app.SetBeforeDrawFunc(func(s tcell.Screen) bool {
		s.Clear()
		return false
	})

	if err := app.SetRoot(mv, true).Run(); err != nil {
		panic(err)
	}

	// TODO: complete the Primitive so it can be drawn
	// TODO: add a test view here
}
