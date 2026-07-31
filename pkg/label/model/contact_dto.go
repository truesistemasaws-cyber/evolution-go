package label_model

// ContactDTO é a estrutura retornada pela API para cada contato.
type ContactDTO struct {
	ID    int64  `json:"id"`
	Name  string `json:"name,omitempty"`
	Phone string `json:"phone,omitempty"`
	Email string `json:"email,omitempty"`
}
