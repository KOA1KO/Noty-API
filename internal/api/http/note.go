package handlers

import (
	httpCodes "Noty/internal/api/http_codes"
	"Noty/internal/ctxKey"
	"Noty/internal/service"
	"Noty/pkg/logger/sl"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) PostNewNote(w http.ResponseWriter, r *http.Request) {
	const op = "internal.handlers.note.PostNewNote"

	log := h.Logger.With(
		slog.String("op", op),
	)

	ctx := r.Context()
	userID, err := ctxKey.GetUserIDFromRequest(r)
	if err != nil {
		log.Error("Smth with gettin userID", sl.Err(err))
		httpCodes.CheckOutAndWriteError(w, err)
		return
	}

	var req InputNoteRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("Decoding input error", sl.Err(err))
		httpCodes.CheckOutAndWriteError(w, err)
		return
	}

	input := service.NoteInput{
		Title:   req.Title,
		Content: req.Content,
	}

	response, err := h.NoteService.CreateNewNote(ctx, input, userID)
	if err != nil {
		httpCodes.CheckOutAndWriteError(w, err)
		return
	}

	log.Info("user add new note")
	httpCodes.WriteJSONResponse(w, http.StatusCreated, response)
}

// @Summary     Список заметок пользователя
// @Tags        notes
// @Security    BearerAuth
// @Param       limit   query  int   false  "ограничение размера выдачи" minimum(1) maximum(100)
// @Param       offset  query  int   false  "смещение"
// @Success     200     {array}  models.Note
// @Failure     400     {object} httpCodes.Response  "Какая то ошибка тумтум"
// @Router      /notes [get]
func (h *Handler) GetNotes(w http.ResponseWriter, r *http.Request) {
	const op = "internal.handlers.note.GetNotes"

	log := h.Logger.With(
		slog.String("op", op),
	)

	ctx := r.Context()
	userID, err := ctxKey.GetUserIDFromRequest(r)
	if err != nil {
		log.Error("Smth with gettin userID", sl.Err(err))
		httpCodes.CheckOutAndWriteError(w, err)
		return
	}

	meta, _ := ctxKey.GetNoteMeta(ctx)

	notes, err := h.NoteService.GetNotes(ctx, userID, meta)
	if err != nil {
		httpCodes.CheckOutAndWriteError(w, err)
		return
	}

	log.Info("User get his notes")
	httpCodes.WriteJSONResponse(w, http.StatusOK, *notes)
}

func (h *Handler) GetNoteByID(w http.ResponseWriter, r *http.Request) {
	const op = "internal.handlers.note.GetNoteByID"

	log := h.Logger.With(
		slog.String("op", op),
	)

	userID, err := ctxKey.GetUserIDFromRequest(r)
	if err != nil {
		log.Error("Smth with gettin userID", sl.Err(err))
		httpCodes.CheckOutAndWriteError(w, err)
		return
	}

	noteIDstr := chi.URLParam(r, "note_id")
	noteID, err := strconv.Atoi(noteIDstr)
	if err != nil {
		log.Error("Cannot convert noteIDstr to int", sl.Err(err))
		httpCodes.WriteErrorResponse(w, http.StatusBadRequest, "note_id должен быть целочисленным")
		return
	}

	note, err := h.NoteService.GetNote(r.Context(), userID, noteID)
	if err != nil {
		httpCodes.CheckOutAndWriteError(w, err)
		return
	}

	log.Info("Note is getted")
	httpCodes.WriteJSONResponse(w, http.StatusOK, note)
}

func (h *Handler) PutUpdateNote(w http.ResponseWriter, r *http.Request) {
	const op = "internal.handlers.note.PutUpdateNote"

	log := h.Logger.With(
		slog.String("op", op),
	)

	userID, err := ctxKey.GetUserIDFromRequest(r)
	if err != nil {
		log.Error("Smth with gettin userID", sl.Err(err))
		httpCodes.CheckOutAndWriteError(w, err)
		return
	}

	noteIDstr := chi.URLParam(r, "note_id")
	noteID, err := strconv.Atoi(noteIDstr)
	if err != nil {
		log.Error("Cannot convert noteIDstr to int", sl.Err(err))
		httpCodes.WriteErrorResponse(w, http.StatusBadRequest, "note_id должен быть целочисленным")
		return
	}

	var note InputNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		log.Error("Decoding input error", sl.Err(err))
		httpCodes.CheckOutAndWriteError(w, err)
		return
	}

	response, err := h.NoteService.UpdateNote(r.Context(), userID, noteID, note.Title, note.Content)
	if err != nil {
		httpCodes.CheckOutAndWriteError(w, err)
		return
	}

	log.Info("Пользователь успешно обновил заметку")
	httpCodes.WriteJSONResponse(w, http.StatusOK, response)
}

func (h *Handler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	const op = "internal.handlers.note.DeleteNote"

	log := h.Logger.With(
		slog.String("op", op),
	)

	userID, err := ctxKey.GetUserIDFromRequest(r)
	if err != nil {
		log.Error("Smth with gettin userID", sl.Err(err))
		httpCodes.CheckOutAndWriteError(w, err)
		return
	}

	noteIDstr := chi.URLParam(r, "note_id")
	noteID, err := strconv.Atoi(noteIDstr)
	if err != nil {
		log.Error("Cannot convert noteIDstr to int", sl.Err(err))
		httpCodes.WriteErrorResponse(w, http.StatusBadRequest, "note_id должен быть целочисленным")
		return
	}

	err = h.NoteService.DeleteNote(r.Context(), userID, noteID)
	if err != nil {
		httpCodes.CheckOutAndWriteError(w, err)
		return
	}

	log.Info("Пользователь успешно удалил заметку")
	httpCodes.WriteJSONResponse(w, http.StatusOK, httpCodes.ResponseOK())
}
