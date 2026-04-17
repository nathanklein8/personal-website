package models

type Project struct {
	ID             int      `json:"id,omitempty"`
	Icon           string   `json:"icon"`
	Title          string   `json:"title"`
	Description    string   `json:"description"`
	Technologies   []string `json:"technologies"`
	DeploymentLink *string  `json:"deploymentLink,omitempty"`
	Image          *string  `json:"image,omitempty"`
	AltText        *string  `json:"altText,omitempty"`
}
