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
	Preference *string `json:"preference,omitempty" validate:"required,oneof=CARDIO WEIGHT"`
	Weightunit *string `json:"weightUnit,omitempty" validate:"required,oneof=KG LBS"`
	Heightunit *string `json:"heightUnit,omitempty" validate:"required,oneof=CM INCH"`
	Weight     *int64  `json:"weight,omitempty" validate:"required,gte=10,lte=1000"`
	Height     *int64  `json:"height,omitempty" validate:"required,gte=3,lte=250"`
	Name       *string `json:"name,omitempty" validate:"omitempty,min=2,max=60"`
	Imageuri   *string `json:"imageUri,omitempty" validate:"omitempty,uri"`
}
