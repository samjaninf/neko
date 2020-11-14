package member

import (
	"github.com/go-chi/chi"

	"demodesk/neko/internal/types"
)

type MemberHandler struct {
	sessions types.SessionManager
}

func New(
	sessions types.SessionManager,
) *MemberHandler {
	// Init

	return &MemberHandler{
		sessions: sessions,
	}
}

func (h *MemberHandler) Route(r chi.Router) {
	
}
