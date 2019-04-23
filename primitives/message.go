package primitives

import (
	"image"

	"github.com/Bios-Marcel/tview"
	"github.com/bwmarrin/discordgo"
	"github.com/gdamore/tcell"
	tviewsixel "gitlab.com/diamondburned/tview-sixel"
)

// Message is the primitive for each Discord message.
type Message struct {
	*discordgo.Message
	x, y, width, height int

	// attachments contains the string, which is the URL,
	// and its corresponding image
	attachments []image.Image
}

// NewMessage creates a new message
func NewMessage(m *discordgo.Message) (*Message, error) {
	message := &Message{
		Message:     m,
		attachments: make([]image.Image, 0, len(m.Attachments)),
	}

	for _, a := range m.Attachments {
		// an image doesn't have 0 for either of these
		if a.Width != 0 && a.Height != 0 {
			i, err := DownloadImage(a.ProxyURL)
			if err != nil {
				// Silently ignore error
				// TODO: log this
				continue
			}

			// add it to cache
			message.attachments = append(message.attachments, i)
		}
	}

	return message, nil
}

// Draw draws the whole message, manually
func (m *Message) Draw(s tcell.Screen) bool {
	if m.width <= 0 || m.height <= 0 {
		return false
	}

	var lines = tview.WordWrap(m.Content, m.width)

	// i is incremented vertically, meaning this loop traverses
	// over each line in the screen
	for i, l := range lines {
		// j is incremented horizontally, meaning this loop
		// traverses over each cell in a line
		for j, r := range []rune(l) {
			// x is horizontal, thus +j
			// y is vertical, thus +i

			// TODO: 0 to actual style
			s.SetContent(m.x+j, m.y+i, r, nil, 0)
		}
	}

	// ptr(X|Y) is used to coordinate where the attachments
	// are drawn. +1 is to space out a line
	ptrX, ptrY := m.x, m.y+len(lines)+1

	for _, a := range m.attachments {
		p, err := tviewsixel.NewPicture(a)
		if err != nil {
			// TODO: log
			continue
		}

		// We give the Picture a line limit of 10 to prevent
		// images that are too tall or too large.
		p.SetRect(ptrX, ptrY, m.width, 10)

		// Draw the image into the screen
		p.Draw(s)

		// Calculate the actual size of the image it was drawn
		_, imgH := p.CalculateSize()

		// Add to ptrY the lines that the image drawn took, by
		// dividing the height in pixels with the height of one
		// cell. 1 is added as a margin/padding of some sort.
		ptrY += imgH/p.CharH + 1
	}

	return true
}
