package models

type LandingCard struct {
	Bio      string     `json:"bio"`
	Email    string     `json:"email"`
	Linkedin string     `json:"linkedin"`
	Github   string     `json:"github"`
	Skills   [][]string `json:"skills"`
}
