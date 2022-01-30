package handler

import (
	"errors"

	"demodesk/neko/internal/desktop/xorg"
	"demodesk/neko/internal/types"
	"demodesk/neko/internal/types/event"
	"demodesk/neko/internal/types/message"
)

var (
	ErrIsNotAllowedToHost = errors.New("is not allowed to host")
	ErrIsNotTheHost       = errors.New("is not the host")
	ErrIsAlreadyTheHost   = errors.New("is already the host")
	ErrIsAlreadyHosted    = errors.New("is already hosted")
)

func (h *MessageHandlerCtx) controlRelease(session types.Session) error {
	if !session.Profile().CanHost {
		return ErrIsNotAllowedToHost
	}

	if !session.IsHost() {
		return ErrIsNotTheHost
	}

	h.desktop.ResetKeys()
	h.sessions.ClearHost()

	return nil
}

func (h *MessageHandlerCtx) controlRequest(session types.Session) error {
	if !session.Profile().CanHost {
		return ErrIsNotAllowedToHost
	}

	if session.IsHost() {
		return ErrIsAlreadyTheHost
	}

	if !h.sessions.ImplicitHosting() {
		// tell session if there is a host
		if host := h.sessions.GetHost(); host != nil {
			session.Send(
				event.CONTROL_HOST,
				message.ControlHost{
					HasHost: true,
					HostID:  host.ID(),
				})

			return ErrIsAlreadyHosted
		}
	}

	h.sessions.SetHost(session)

	return nil
}

func (h *MessageHandlerCtx) controlKeyPress(session types.Session, payload *message.ControlKey) error {
	if err := h.controlRequest(session); err != nil && !errors.Is(err, ErrIsAlreadyTheHost) {
		return err
	}

	return h.desktop.KeyPress(payload.Keysym)
}

func (h *MessageHandlerCtx) controlKeyDown(session types.Session, payload *message.ControlKey) error {
	if err := h.controlRequest(session); err != nil && !errors.Is(err, ErrIsAlreadyTheHost) {
		return err
	}

	return h.desktop.KeyDown(payload.Keysym)
}

func (h *MessageHandlerCtx) controlKeyUp(session types.Session, payload *message.ControlKey) error {
	if err := h.controlRequest(session); err != nil && !errors.Is(err, ErrIsAlreadyTheHost) {
		return err
	}

	return h.desktop.KeyUp(payload.Keysym)
}

func (h *MessageHandlerCtx) controlCut(session types.Session) error {
	if err := h.controlRequest(session); err != nil && !errors.Is(err, ErrIsAlreadyTheHost) {
		return err
	}

	return h.desktop.KeyPress(xorg.XK_Control_L, xorg.XK_x)
}

func (h *MessageHandlerCtx) controlCopy(session types.Session) error {
	if err := h.controlRequest(session); err != nil && !errors.Is(err, ErrIsAlreadyTheHost) {
		return err
	}

	return h.desktop.KeyPress(xorg.XK_Control_L, xorg.XK_c)
}

func (h *MessageHandlerCtx) controlPaste(session types.Session) error {
	if err := h.controlRequest(session); err != nil && !errors.Is(err, ErrIsAlreadyTheHost) {
		return err
	}

	return h.desktop.KeyPress(xorg.XK_Control_L, xorg.XK_v)
}

func (h *MessageHandlerCtx) controlSelectAll(session types.Session) error {
	if err := h.controlRequest(session); err != nil && !errors.Is(err, ErrIsAlreadyTheHost) {
		return err
	}

	return h.desktop.KeyPress(xorg.XK_Control_L, xorg.XK_a)
}
