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

	visible bool
	focus   tview.Focusable

	// PtrY is used to calculate the total height in the future
	PtrY int

	// attachments contains the string, which is the URL,
	// and its corresponding image
	attachments []image.Image
}

// NewMessage creates a new message
func NewMessage(m *discordgo.Message) (*Message, error) {
	message := &Message{
		visible:     true,
		Message:     m,
		attachments: make([]image.Image, 0, len(m.Attachments)),
	}

	message.focus = message

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
	m.PtrY = m.y + len(lines) + 1

	for _, a := range m.attachments {
		p, err := tviewsixel.NewPicture(a)
		if err != nil {
			// TODO: log
			continue
		}

		// We give the Picture a line limit of 10 to prevent
		// images that are too tall or too large.
		p.SetRect(m.x, m.PtrY, m.width, 15)

		// Draw the image into the screen
		p.Draw(s)

		// Calculate the actual size of the image it was drawn
		_, imgH := p.CalculateSize()

		// Add to ptrY the lines that the image drawn took, by
		// dividing the height in pixels with the height of one
		// cell. 1 is added as a margin/padding of some sort.
		m.PtrY += imgH/p.CharH + 1
	}

	return true
}

// SetVisible sets visibility like CSS.
func (m *Message) SetVisible(v bool) {
	m.visible = v
}

// IsVisible returns visible
func (m *Message) IsVisible() bool {
	return m.visible
}

// GetRect returns the rectangle dimensions
func (m *Message) GetRect() (int, int, int, int) {
	return m.x, m.y, m.width, m.height
}

// SetRect sets the rectangle dimensions
func (m *Message) SetRect(x, y, width, height int) {
	m.x = x
	m.y = y
	m.width = width
	m.height = height
}

// InputHandler sets no input handler, satisfying Primitive
func (m *Message) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return nil
}

// Focus does nothing, really.
func (*Message) Focus(delegate func(tview.Primitive)) {}

// Blur also does nothing.
func (*Message) Blur() {}

// HasFocus always returns false, as you can't focus on this.
func (*Message) HasFocus() bool {
	return false
}

// GetFocusable does whatever the fuck I have no idea
func (m *Message) GetFocusable() tview.Focusable {
	return m.focus
}

func (m *Message) SetOnFocus(func()) {}
func (m *Message) SetOnBlur(func())  {}
