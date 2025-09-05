package server

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

type PatchUserRequest struct {
	Preference *string `json:"preference,omitempty"`
	Weightunit *string `json:"weightUnit,omitempty"`
	Heightunit *string `json:"heightUnit,omitempty"`
	Weight     *int64  `json:"weight,omitempty"`
	Height     *int64  `json:"height,omitempty"`
	Name       *string `json:"name,omitempty"`
	Imageuri   *string `json:"imageUri,omitempty"`
}
