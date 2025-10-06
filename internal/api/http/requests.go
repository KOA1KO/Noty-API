package handlers

type InputUserRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type InputNoteRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type InputRefresh struct {
	RefreshToken string `json:"refresh_token"`
}
