package service

type SignInput struct {
	Username string
	Password string
}

type NoteInput struct {
	Title   string
	Content string
}

type RefreshInput struct {
	RefreshToken string
}
