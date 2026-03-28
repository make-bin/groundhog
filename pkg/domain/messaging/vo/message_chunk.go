// @AI_GENERATED
package vo

// MessageChunk represents a single chunk of a message split for channel delivery.
// It is immutable after creation.
type MessageChunk struct {
	index   int
	content string
	isFinal bool
}

// NewMessageChunk creates a new MessageChunk.
func NewMessageChunk(index int, content string, isFinal bool) MessageChunk {
	return MessageChunk{
		index:   index,
		content: content,
		isFinal: isFinal,
	}
}

// Index returns the zero-based position of this chunk in the sequence.
func (m MessageChunk) Index() int { return m.index }

// Content returns the text content of this chunk.
func (m MessageChunk) Content() string { return m.content }

// IsFinal returns true if this is the last chunk in the sequence.
func (m MessageChunk) IsFinal() bool { return m.isFinal }

// @AI_GENERATED: end
