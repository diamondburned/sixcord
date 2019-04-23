package primitives

import (
	"errors"
	"image"

	"github.com/Bios-Marcel/tview"
	"github.com/bwmarrin/discordgo"
	"github.com/gdamore/tcell"
	tviewsixel "gitlab.com/diamondburned/tview-sixel"
)

const avatarSize = 48 // size of avatar in px

var (
	// ErrNoAuthor is returned when Author is nil
	ErrNoAuthor = errors.New("Author not found")
)

// Author draws the author primitive
type Author struct {
	*discordgo.User

	x, y, width, height int

	visible bool
	focus   tview.Focusable

	// PtrY is used to calculate the total height in the future
	PtrY int

	avatar image.Image

	children []*Message
}

// NewAuthor makes a new author primitive
func NewAuthor(u *discordgo.User) (*Author, error) {
	if u == nil {
		return nil, ErrNoAuthor
	}

	a := &Author{
		User:     u,
		children: []*Message{},
	}

	i, err := DownloadImage(u.AvatarURL("64"))
	if err == nil {
		a.avatar = i
	}

	return a, nil
}

// AddMessage adds a message into the buffer. User should call Draw
// after this.
func (a *Author) AddMessage(m *discordgo.Message) error {
	mPrimitive, err := NewMessage(m)
	if err != nil {
		return err
	}

	a.children = append(a.children, mPrimitive)
	return nil
}

// DeleteMessage deletes a message from the buffer. True is returned
// if the message is removed. Again, call Draw.
func (a *Author) DeleteMessage(m *discordgo.Message) bool {
	for i, mp := range a.children {
		if mp.ID == m.ID {
			a.children = deleteMessage(a.children, i)
			return true
		}
	}

	return false
}

// EditMessage is similar to DeleteMessage, wherein true is returned
// if there's an actual edit. An extra error is introduced.
// Again, call Draw.
func (a *Author) EditMessage(m *discordgo.Message) (bool, error) {
	for _, mp := range a.children {
		if mp.ID == m.ID {
			p, err := NewMessage(m)
			if err != nil {
				return false, err
			}

			// mp is a pointer in the array
			mp = p
			return true, nil
		}
	}

	return false, nil
}

// ShouldRemove returns true if the user should remove this Author
// primitive from the container. If this is true, nothing will be
// drawn when Draw() is called.
func (a *Author) ShouldRemove() bool {
	return len(a.children) == 0
}

// Draw draws the author down
func (a *Author) Draw(s tcell.Screen) bool {
	if a.width <= 0 || a.height <= 0 {
		return false
	}

	if a.ShouldRemove() {
		return false
	}

	icellW := avatarSize / tviewsixel.CharW
	icellH := avatarSize / tviewsixel.CharH

	i, err := tviewsixel.NewPicture(a.avatar)
	if err == nil {
		i.SetRect(a.x, a.y, icellW, icellH)
		i.Draw(s)
	}

	for i, r := range []rune(a.Username) {
		s.SetContent(a.x+icellW+i+1, a.y, r, nil, 0)
	}

	a.PtrY = a.y + 1
	x := a.x + icellW + 1

	for _, m := range a.children {
		m.SetRect(x, a.PtrY, a.width-icellW, a.height-a.PtrY)
		m.Draw(s)

		a.PtrY = m.PtrY

		if a.PtrY > a.height {
			break
		}
	}

	return true
}

// SetVisible sets visibility like CSS.
func (a *Author) SetVisible(v bool) {
	a.visible = v
}

// IsVisible returns visible
func (a *Author) IsVisible() bool {
	return a.visible
}

// GetRect returns the rectangle dia.nsions
func (a *Author) GetRect() (int, int, int, int) {
	return a.x, a.y, a.width, a.height
}

// SetRect sets the rectangle dia.nsions
func (a *Author) SetRect(x, y, width, height int) {
	a.x = x
	a.y = y
	a.width = width
	a.height = height
}

// InputHandler sets no input handler, satisfying Pria.tive
func (a *Author) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return nil
}

// Focus does nothing, really.
func (*Author) Focus(delegate func(tview.Primitive)) {}

// Blur also does nothing.
func (*Author) Blur() {}

// HasFocus always returns false, as you can't focus on this.
func (*Author) HasFocus() bool {
	return false
}

// GetFocusable does whatever the fuck I have no idea
func (a *Author) GetFocusable() tview.Focusable {
	return a.focus
}

func (a *Author) SetOnFocus(func()) {}
func (a *Author) SetOnBlur(func())  {}

func intDivCeil(i, j int) int {
	return int((float64(i) + 0.5) / float64(j))
}

// from https://github.com/golang/go/wiki/SliceTricks, this function
// deletes without memory leaking
func deleteMessage(a []*Message, i int) []*Message {
	copy(a[i:], a[i+1:])
	a[len(a)-1] = nil
	a = a[:len(a)-1]

	return a
}
