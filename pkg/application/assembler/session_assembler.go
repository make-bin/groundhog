package assembler

import (
	"github.com/make-bin/groundhog/pkg/application/dto"
	"github.com/make-bin/groundhog/pkg/domain/conversation/aggregate/agent_session"
)

// ToSessionResponse converts an AgentSession domain aggregate to a SessionResponse DTO.
func ToSessionResponse(session *agent_session.AgentSession) *dto.SessionResponse {
	turns := make([]dto.TurnResponse, 0, len(session.Turns()))
	for _, t := range session.Turns() {
		turns = append(turns, dto.TurnResponse{
			ID:          t.ID(),
			UserInput:   t.UserInput(),
			Response:    t.Response(),
			ModelUsed:   t.ModelUsed(),
			StartedAt:   t.StartedAt(),
			CompletedAt: t.CompletedAt(),
		})
	}
	return &dto.SessionResponse{
		ID:           session.ID().Value(),
		AgentID:      session.AgentID().Value(),
		UserID:       session.UserID(),
		State:        session.State().String(),
		ActiveModel:  session.ActiveModel().ModelName(),
		Turns:        turns,
		CreatedAt:    session.CreatedAt(),
		LastActiveAt: session.LastActiveAt(),
	}
}

// ToSessionResponseList converts a slice of AgentSession aggregates to SessionResponse DTOs.
func ToSessionResponseList(sessions []*agent_session.AgentSession) []*dto.SessionResponse {
	result := make([]*dto.SessionResponse, 0, len(sessions))
	for _, s := range sessions {
		result = append(result, ToSessionResponse(s))
	}
	return result
}
