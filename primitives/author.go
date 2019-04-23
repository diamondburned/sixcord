package primitives

import (
	"errors"
	"image"

	"github.com/bwmarrin/discordgo"
)

// avatarStore contains a map of images, which
// corresponds to an user's ID
var avatarStore map[string]image.Image

var (
	// ErrNoAuthor is returned when Author is nil
	ErrNoAuthor = errors.New("Author not found")
)

// Author draws the author primitive
type Author struct {
	*discordgo.User

	firstChild *Message
	children   []*Message
}

// NewAuthor makes a new author primitive
func NewAuthor(m *discordgo.Message) (*Author, error) {
	if m.Author == nil {
		return nil, ErrNoAuthor
	}

	mPrimitive, err := NewMessage(m)
	if err != nil {
		return nil, err
	}

	a := &Author{
		User:       m.Author,
		firstChild: mPrimitive,
		children:   []*Message{mPrimitive},
	}

	if _, ok := avatarStore[a.ID]; !ok {
		i, err := DownloadImage(m.Author.AvatarURL("32"))
		if err != nil {
			return a, nil
		}

		avatarStore[a.ID] = i
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

// from https://github.com/golang/go/wiki/SliceTricks, this function
// deletes without memory leaking
func deleteMessage(a []*Message, i int) []*Message {
	copy(a[i:], a[i+1:])
	a[len(a)-1] = nil
	a = a[:len(a)-1]

	return a
}
