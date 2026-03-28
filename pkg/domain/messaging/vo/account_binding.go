// @AI_GENERATED
package vo

// AccountBinding represents the binding between a channel account and an agent session.
// It maps a (channelID, accountID) pair to a session for routing purposes.
// It is immutable after creation.
type AccountBinding struct {
	channelID ChannelID
	accountID AccountID
	sessionID string
}

// NewAccountBinding creates a new AccountBinding.
func NewAccountBinding(channelID ChannelID, accountID AccountID, sessionID string) AccountBinding {
	return AccountBinding{
		channelID: channelID,
		accountID: accountID,
		sessionID: sessionID,
	}
}

// ChannelID returns the channel ID of the binding.
func (a AccountBinding) ChannelID() ChannelID { return a.channelID }

// AccountID returns the account ID of the binding.
func (a AccountBinding) AccountID() AccountID { return a.accountID }

// SessionID returns the session ID this binding routes to.
func (a AccountBinding) SessionID() string { return a.sessionID }

// Equals returns true if a and other represent the same binding.
func (a AccountBinding) Equals(other AccountBinding) bool {
	return a.channelID.Equals(other.channelID) &&
		a.accountID.Equals(other.accountID) &&
		a.sessionID == other.sessionID
}

// @AI_GENERATED: end
