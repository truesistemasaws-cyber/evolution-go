package label_repository

import (
	label_model "github.com/evolution-foundation/evolution-go/pkg/label/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type LabelRepository interface {
	InsertLabel(label label_model.Label) error
	UpdateLabel(label label_model.Label) error
	GetLabelByID(id string) (*label_model.Label, error)
	DeleteLabel(id string) error
	GetAllLabelsByInstanceID(instanceID string) ([]label_model.Label, error)
	UpsertLabel(label label_model.Label) error
	// GetContactsByLabel retorna contatos associados a uma label para uma instância (com paginação)
	GetContactsByLabel(instanceID string, label string, limit, offset int) ([]label_model.ContactDTO, error)
}

type labelRepository struct {
	db *gorm.DB
}

func (l *labelRepository) InsertLabel(label label_model.Label) error {
	return l.db.Create(&label).Error
}

func (l *labelRepository) UpdateLabel(label label_model.Label) error {
	return l.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"label_id", "label_name", "label_color"}),
	}).Create(&label).Error
}

func (l *labelRepository) GetLabelByID(id string) (*label_model.Label, error) {
	var label label_model.Label
	err := l.db.Where("id = ?", id).First(&label).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &label, nil
}

func (l *labelRepository) DeleteLabel(id string) error {
	return l.db.Where("id = ?", id).Delete(&label_model.Label{}).Error
}

func (l *labelRepository) GetAllLabelsByInstanceID(instanceID string) ([]label_model.Label, error) {
	var labels []label_model.Label
	err := l.db.Where("instance_id = ?", instanceID).Find(&labels).Error
	if err != nil {
		return nil, err
	}
	return labels, nil
}

func (l *labelRepository) UpsertLabel(label label_model.Label) error {
	return l.db.Where("instance_id = ? AND label_id = ?",
		label.InstanceID,
		label.LabelID,
	).Assign(label_model.Label{
		LabelName:    label.LabelName,
		LabelColor:   label.LabelColor,
		PredefinedId: label.PredefinedId,
	}).FirstOrCreate(&label).Error
}

// GetContactsByLabel faz join entre contact_labels e contacts e retorna DTOs.
func (l *labelRepository) GetContactsByLabel(instanceID string, label string, limit, offset int) ([]label_model.ContactDTO, error) {
	var contacts []label_model.ContactDTO

	if limit <= 0 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	err := l.db.Raw(`
		SELECT c.id, c.name, c.phone, c.email
		FROM contacts c
		JOIN contact_labels cl ON cl.contact_id = c.id
		WHERE cl.instance_id = ? AND cl.label = ?
		ORDER BY c.id
		LIMIT ? OFFSET ?
	`, instanceID, label, limit, offset).Scan(&contacts).Error

	if err != nil {
		return nil, err
	}
	return contacts, nil
}

func NewLabelRepository(db *gorm.DB) LabelRepository {
	return &labelRepository{db: db}
}
