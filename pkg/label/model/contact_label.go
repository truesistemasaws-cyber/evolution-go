package label_model

import "time"

// ContactLabel representa a relação entre um contato e uma label.
// Usa InstanceID para isolar por instância (multitenant) conforme o padrão do projeto.
type ContactLabel struct {
	ID         uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	InstanceID string    `json:"instance_id" gorm:"type:uuid;index;not null"`
	ContactID  int64     `json:"contact_id" gorm:"index;not null"`
	Label      string    `json:"label" gorm:"type:varchar(255);index;not null"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
}
