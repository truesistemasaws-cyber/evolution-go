package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Este arquivo adiciona um handler para listar contatos por etiqueta (label).
// Suposições:
// - O projeto usa Gin como router e GORM para acesso ao banco.
// - Existem tabelas `contacts` e `contact_labels` (contact_id, label).
// Ajuste nomes/queries conforme os models reais do repositório.

// Contact representa os campos retornados para a API.
type Contact struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Email string `json:"email,omitempty"`
}

// LabelHandler agrupa dependências (ex.: DB).
type LabelHandler struct {
	DB *gorm.DB
}

// NewLabelHandler cria um novo LabelHandler.
func NewLabelHandler(db *gorm.DB) *LabelHandler {
	return &LabelHandler{DB: db}
}

// ListContactsByLabel responde GET /labels/:label/contacts
func (h *LabelHandler) ListContactsByLabel(c *gin.Context) {
	label := c.Param("label")
	if label == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "label é obrigatório"})
		return
	}

	var contacts []Contact
	// Query genérica — adapte aos nomes reais de tabelas/colunas do seu schema.
	if err := h.DB.Raw(`
		SELECT c.id, c.name, c.phone, c.email
		FROM contacts c
		JOIN contact_labels cl ON cl.contact_id = c.id
		WHERE cl.label = ?
	`, label).Scan(&contacts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "falha ao buscar contatos", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"label": label, "contacts": contacts})
}

// Exemplo de uso (registre sua rota no main ou onde centraliza rotas):
//
// db := /* seu *gorm.DB */
// r := gin.Default()
// lh := handlers.NewLabelHandler(db)
// r.GET("/labels/:label/contacts", lh.ListContactsByLabel)
//
// Observações:
// - Proteja a rota com autenticação se necessário.
// - Adicione paginação (limit/offset) se houver muitos contatos.
