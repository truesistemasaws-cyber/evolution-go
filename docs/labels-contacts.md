# Listar contatos por etiqueta (Label)

Descrição
---------
Este endpoint permite listar todos os contatos associados a uma determinada etiqueta (label) para a instância autenticada.

Rota
----
GET /label/contacts

Parâmetros
----------
- label (query, required) — nome da etiqueta a ser consultada.
- limit (query, optional) — número máximo de contatos retornados. Default: 100.
- offset (query, optional) — deslocamento para paginação. Default: 0.

Autorização
-----------
A rota está protegida pelo middleware de autenticação do projeto (mesmo usado nas demais rotas do grupo `/label`). Inclua o header/forma de autenticação habitual (p.ex. `Authorization: Bearer <TOKEN>`).

Exemplo de requisição (curl)
----------------------------
curl -H "Authorization: Bearer <TOKEN>" \
  "http://<HOST>:<PORT>/label/contacts?label=clientes&limit=50&offset=0"

Resposta (200)
--------------
JSON com o nome da label e a lista de contatos:

```json
{
  "label": "clientes",
  "contacts": [
    {"id": 1, "name": "João", "phone": "+551199999999", "email": "joao@ex.com"},
    {"id": 2, "name": "Maria", "phone": "+551198888888", "email": "maria@ex.com"}
  ]
}
```

Erros comuns
------------
- 400 Bad Request: quando o parâmetro `label` não for informado.
- 401/403: quando o token/autorização for inválido (conforme política do projeto).
- 500 Internal Server Error: erro no DB ou na lógica do serviço (ver logs).

Migração / Banco de Dados
-------------------------
Foi adicionada a tabela `contact_labels` para relacionar contatos e labels. Arquivo de migração:
- migrations/20260731_create_contact_labels.sql

Conteúdo SQL presente no arquivo de migração (veja também migrations/20260731_create_contact_labels.sql):

```sql
CREATE TABLE IF NOT EXISTS contact_labels (
  id BIGSERIAL PRIMARY KEY,
  instance_id UUID NOT NULL,
  contact_id BIGINT NOT NULL,
  label VARCHAR(255) NOT NULL,
  created_at TIMESTAMPTZ DEFAULT now(),
  UNIQUE (instance_id, contact_id, label)
);

CREATE INDEX IF NOT EXISTS idx_contact_labels_label ON contact_labels(label);
CREATE INDEX IF NOT EXISTS idx_contact_labels_instance_label ON contact_labels(instance_id, label);
```

Observações sobre schema
------------------------
- A query atual assume que existe uma tabela `contacts` com colunas: `id` (BIGINT), `name`, `phone`, `email`.
- Se o seu schema usar `id` como UUID/STRING ou armazena o identificador do contato em outra coluna (ex.: `jid` ou `wa_id`), ajuste a migration e a query em `pkg/label/repository/label_repository.go`.
- A coluna `instance_id` foi usada para manter isolamento por instância (multitenant) — isto segue o padrão do model `Label`. Se o seu projeto não usar esse campo, remova o filtro `instance_id` na query.

Integração com o código
-----------------------
- Rota: GET /label/contacts?label=<nome>&limit=<n>&offset=<n>
- Handler: pkg/label/handler/label_handler.go (GetContactsByLabel)
- Service: pkg/label/service/label_service.go (GetContactsByLabel)
- Repository: pkg/label/repository/label_repository.go (GetContactsByLabel)
- Model: pkg/label/model/contact_label.go
- DTO: pkg/label/model/contact_dto.go
- migrate() em cmd/evolution-go/main.go deve incluir &label_model.ContactLabel{}

Testes recomendados
------------------
- Inserir dados de exemplo em `contacts` e `contact_labels` e chamar o endpoint autenticado.
- Testes unitários para repository / handler.

Rollback
--------
Para reverter a migração:

```sql
DROP TABLE IF EXISTS contact_labels;
```

(Execute apenas após avaliação/backup.)

Contato
------
Se precisar que eu ajuste a migration para usar `instance_id` como VARCHAR em vez de UUID, ou que suporte `contact_id` como UUID/string, diga o esquema real e eu atualizo os arquivos e queries.
