# 🚀 Transaction Flow API

## 📌 Sobre o Projeto

A **Transaction Flow API** é um sistema de gerenciamento de contas e transações, permitindo:
- Criar contas associadas a um **document number**.
- Recuperar informações de uma conta existente.
- Registrar transações financeiras de diferentes tipos (**compras à vista, compras parceladas, saques e pagamentos**).

## 🛠️ Tecnologias Utilizadas

Este projeto utiliza as seguintes tecnologias:

- **[Go](https://go.dev/doc/)** - Linguagem de programação.
- **[Chi Router](https://github.com/go-chi/chi)** - Roteamento leve e flexível para APIs HTTP.
- **[PostgreSQL](https://www.postgresql.org/docs/)** - Banco de dados relacional.
- **[Docker](https://docs.docker.com/)** - Containerização do projeto.
- **[golang-migrate](https://github.com/golang-migrate/migrate)** - Gerenciamento de migrações de banco de dados.
- **[Slog](https://pkg.go.dev/log/slog)** - Biblioteca para logging estruturado.

---

## 📂 Estrutura do Projeto

```plaintext
├── cmd/
│   ├── transaction-flow/        # Inicialização
│   │   ├── main.go              
│
├── internal/
│   ├── api/                     # Camada de interface HTTP
│   │   ├── handlers/             # Handlers para requisições HTTP
│   │   ├── middleware/           # Middlewares (logging, recovery)
│   │   ├── dto/                  # Estruturas de requisição/resposta
│   │   ├── response/              # Formatação de erros
│   │
│   ├── application/             # Regras de negócio (use cases)
│   │   ├── usecase/               # Casos de uso (Account, Transaction)
│   │
│   ├── domain/                  # Modelos de entidades e regras de domínio
│   │   ├── account.go
│   │   ├── transaction.go
│   │
│   ├── infra/                   # Camada de infraestrutura
│   │   ├── repository/            # Acesso ao banco de dados
│   │   ├── logger/                # Configuração de logs
│   │
├── migrations/                  # Arquivos SQL para criação e atualização do banco
│
├── config/                      # Configuração de envs do ambiente
│   ├── config.go
│                
```

---

## 🚀 **Como rodar o projeto localmente**

### **📌 Pré-requisitos**

Antes de começar, você precisará ter instalado:

- **[Docker](https://docs.docker.com/get-docker/)**
- **[Docker Compose](https://docs.docker.com/compose/install/)**

### **📌 Passo a passo**

1️⃣ **Clone o repositório**
```bash
git clone https://github.com/seu-usuario/transaction-flow.git
cd transaction-flow
```

2️⃣ **Suba os containers do projeto (API, banco e migrations)**
```bash
docker-compose up -d
```
📌 Isso iniciará todos os serviços necessários, incluindo o banco de dados e a aplicação.

---

## 🔥 **Endpoints da API**

### **📌 Criar uma Conta**
📍 **POST** `/accounts`
```bash
curl -X POST http://localhost:8080/accounts \
     -H "Content-Type: application/json" \
     -d '{"document_number": "12345678900"}'
```
📌 **Resposta (201 Created)**
```json
{
  "id": 1
}
```

### **📌 Recuperar uma Conta**
📍 **GET** `/accounts/{id}`
```bash
curl -X GET http://localhost:8080/accounts/1
```
📌 **Resposta (200 OK)**
```json
{
  "account_id": 1,
  "document_number": "12345678900"
}
```

---

### **📌 Criar uma Transação**
📍 **POST** `/transactions`
```bash
curl --X POST 'http://localhost:8080/accounts' \
     -H "Content-Type: application/json" \
     -d '{
            "account_id": 1,
            "operation_type_id": 4,
            "amount": 123.45
        }'
```
📌 **Resposta (201 Created)**
```json
{
  "id": 10
}
```

---

## 🏗️ **Configuração do Ambiente**

Se precisar configurar variáveis de ambiente, crie um arquivo **`.env`**:
```ini
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=transactions
SERVER_PORT=8080
```
📌 A API já lê automaticamente esse arquivo ao iniciar.

---

## 📜 **Licença**
Este projeto é distribuído sob a licença **MIT**. Veja o arquivo [LICENSE](LICENSE) para mais detalhes.
