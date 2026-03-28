// @AI_GENERATED
package vo

// MediaRef represents a reference to a media attachment in a message.
type MediaRef struct {
	url      string
	mimeType string
}

// NewMediaRef creates a new MediaRef.
func NewMediaRef(url, mimeType string) MediaRef {
	return MediaRef{url: url, mimeType: mimeType}
}

// URL returns the media URL.
func (m MediaRef) URL() string { return m.url }

// MimeType returns the MIME type of the media.
func (m MediaRef) MimeType() string { return m.mimeType }

// ParsedCommand represents a parsed chat command with its name and arguments.
type ParsedCommand struct {
	name string
	args []string
}

// NewParsedCommand creates a new ParsedCommand.
func NewParsedCommand(name string, args []string) ParsedCommand {
	argsCopy := make([]string, len(args))
	copy(argsCopy, args)
	return ParsedCommand{name: name, args: argsCopy}
}

// Name returns the command name (without the leading slash).
func (p ParsedCommand) Name() string { return p.name }

// Args returns the command arguments.
func (p ParsedCommand) Args() []string {
	result := make([]string, len(p.args))
	copy(result, p.args)
	return result
}

// MessageStatus represents the delivery status of a message.
type MessageStatus int

const (
	MessageStatusPending   MessageStatus = iota // Message received, not yet routed
	MessageStatusRouted                         // Message routed to an agent session
	MessageStatusDelivered                      // Response delivered back to the channel
	MessageStatusFailed                         // Delivery failed
)

// String returns the string representation of the message status.
func (s MessageStatus) String() string {
	switch s {
	case MessageStatusPending:
		return "Pending"
	case MessageStatusRouted:
		return "Routed"
	case MessageStatusDelivered:
		return "Delivered"
	case MessageStatusFailed:
		return "Failed"
	default:
		return "unknown"
	}
}

// MessageContent represents the content of an inbound or outbound message.
// It is immutable after creation.
type MessageContent struct {
	text      string
	mediaRefs []MediaRef
	isCommand bool
	command   *ParsedCommand
}

// NewMessageContent creates a new MessageContent.
func NewMessageContent(text string, mediaRefs []MediaRef, isCommand bool, command *ParsedCommand) MessageContent {
	refs := make([]MediaRef, len(mediaRefs))
	copy(refs, mediaRefs)
	return MessageContent{
		text:      text,
		mediaRefs: refs,
		isCommand: isCommand,
		command:   command,
	}
}

// Text returns the plain text of the message.
func (m MessageContent) Text() string { return m.text }

// MediaRefs returns the media attachments of the message.
func (m MessageContent) MediaRefs() []MediaRef {
	result := make([]MediaRef, len(m.mediaRefs))
	copy(result, m.mediaRefs)
	return result
}

// IsCommand returns true if the message is a command.
func (m MessageContent) IsCommand() bool { return m.isCommand }

// Command returns the parsed command, or nil if the message is not a command.
func (m MessageContent) Command() *ParsedCommand { return m.command }

// @AI_GENERATED: end
