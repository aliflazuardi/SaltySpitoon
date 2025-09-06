package model

type PatchUserModel struct {
	Preference *string `json:"preference,omitempty"`
	Weightunit *string `json:"weightUnit,omitempty"`
	Heightunit *string `json:"heightUnit,omitempty"`
	Weight     *int    `json:"weight,omitempty"`
	Height     *int    `json:"height,omitempty"`
	Name       *string `json:"name,omitempty"`
	Imageuri   *string `json:"imageUri,omitempty"`
}
