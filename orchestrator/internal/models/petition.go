package models

type Petition struct {
	PetitionID  string `json:"id"`
	Description string `json:"description"`
	Location    string `json:"location"`
}
