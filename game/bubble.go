package game

import (
	"mazox/runtime"
	"time"
)

type Bubble struct {
	image           runtime.Image
	textLen         int
	displayedAt     time.Time
	displayDuration int64
}

func newBubble(k string) *Bubble {
	t := I18nGet(k)
	b := &Bubble{image: gm.font.DrawString(t, 250, 300, 10, 0)}
	b.textLen = len(t)
	return b
}

func (b *Bubble) display() {
	b.displayedAt = time.Now()
	b.displayDuration = int64(b.textLen) * 150 // 150 ms per char
}

func (b *Bubble) mustHide() bool {
	return time.Since(b.displayedAt).Milliseconds() >= b.displayDuration
}
