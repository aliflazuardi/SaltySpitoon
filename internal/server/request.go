package server

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

type CreateActivityRequest struct {
	ActivityType      string `json:"activityType" validate:"required,activity_type_enum"`
	DoneAt            string `json:"doneAt" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	DurationInMinutes int    `json:"durationInMinutes" validate:"required,min=1"`
}
