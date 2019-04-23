package primitives

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// Container contains multiple authors
type Container struct {
	*tview.Box

	authors []*Author
}

func (c *Container) Draw(s tcell.Screen) bool {
	if len(c.authors) == 0 {
		return false
	}

	for i := len(c.authors) - 1; i >= 0; i-- {
		a := c.authors[i]

		// TODO: set sim screen size
		// TODO: add mutex for sim screen, or just make
		//       a custom one

		// TODO: optimize this somehow. Drawing twice
		//       is VERY costly.
		a.Draw(sim)
	}

	return true
}
