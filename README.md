# Chargeback: Construindo uma Arquitetura Event-Driven Resiliente com Go, RabbitMQ, MinIO e FTP

[![Go](https://img.shields.io/badge/Go-1.20+-blue)](https://golang.org)
[![RabbitMQ](https://img.shields.io/badge/RabbitMQ-Messaging-orange)](https://www.rabbitmq.com/)
[![Postgres](https://img.shields.io/badge/Postgres-15-blue)](https://www.postgresql.org/)
[![Cassandra](https://img.shields.io/badge/Cassandra-4.1-blue)](https://cassandra.apache.org/)
[![MinIO](https://img.shields.io/badge/MinIO-Storage-red)](https://min.io/)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

Este projeto é uma **arquitetura completa e resiliente para processar chargebacks**, simulando o fluxo real de estornos em bancos e instituições financeiras, incluindo o envio de arquivos NDJSON via FTP para redes de pagamento (ex.: Mastercard).

## Tecnologias Utilizadas

- **Golang**: alto desempenho e concorrência nativa.
- **RabbitMQ**: mensageria para desacoplamento entre microserviços.
- **PostgreSQL**: persistência de batches e controle de envio.
- **Cassandra**: armazenamento de chargebacks com alta disponibilidade.
- **MinIO**: simulação local de S3 para armazenar arquivos NDJSON.
- **FTP (Pure-FTPd)**: simulação de envio de arquivos para a Mastercard.
- **K6**: ferramenta para estresse e carga em APIs.

## Arquitetura

- **Chargeback API**: Serviço HTTP que recebe pedidos de chargeback e publica eventos `chargeback-opened` no RabbitMQ.
- **Chargeback Processor**: Worker que consome os eventos e processa os chargebacks, persistindo no Cassandra e escrevendo em arquivos NDJSON.
- **Chargeback Batch**: Serviço que gerencia o envio dos arquivos NDJSON para a Mastercard via FTP, com lógica de retries e controle diário.
- **FTP Server**: Simula o servidor da Mastercard para receber os arquivos.
- **MinIO**: Armazena os arquivos NDJSON simulando um bucket S3.

## Exemplo da Estrutura de Diretórios
```bash
chargeback-api/
├── cmd/
│   └── main.go
├── internal/
│   ├── application/
│   │   └── application.go
│   ├── domain/
│   │   ├── models/
│   │   │   └── chargeback.go
│   │   └── repositories/
│   │       └── chargeback_repository.go
│   ├── infrastructure/
│   │   ├── respositories/
│   │   │   └── cassandra/
│   │   │       └── chargeback_repository.go
│   │   ├── rabbitmq/
│   │   │   ├── producers/
│   │   │   │   └── chargeback-opened-producers.go
│   │   │   ├── connection.go
│   │   │   └── event.go
│   │   ├── logging/
│   │   │   └── logging.go
│   │   └── middlewares/
│   ├── interfaces/
│   │   ├── http/
│   │   │   ├── handlers/
│   │   │   │   └── chargeback_handler.go
│   │   │   └── routes/
│   │   │       └── routes.go
│   │   └── messaging/
│   │       └── publisher.go
│   ├── presentation/
│   ├── usecases/
│   └── config/
│       └── config.go
├── docs/ # Swagger docs
├── go.mod
├── chargeback-api.dockerfile

```

## Execução do Projeto

### Build e Subida dos Contêineres:

```bash
cd project
make up_build
```
Esse comando irá:

Fazer go build para cada microserviço.

Gerar a documentação Swagger (para o Chargeback API).

Subir todos os serviços via docker-compose.

## Para subir sem rebuild:
```bash
make up
```

## Parar os containers:
```bash
make down
```

## URLs Importantes
Serviço	URL
- Chargeback API	http://localhost:8080/swagger/index.html
- RabbitMQ Management	http://localhost:15672 (guest/guest)
- MinIO Console	http://localhost:9001 (admin/password)
- FTP Server (Mastercard)	ftp://localhost (admin/admin)
- Postgres (Batch)	postgres://admin:admin@localhost:5432/batch
- Cassandra	cqlsh cassandra:9042

## Configurações
Todos os microserviços usam variáveis de ambiente para configuração. Exemplos personalizáveis:

- NEW_RELIC_ENABLED
- CHARGEBACK_MAX_RECORDS, CHARGEBACK_MAX_DURATION_VALUE, CHARGEBACK_MAX_DURATION_UNIT
- FTP_HOST, FTP_PORT
- SCHEDULER_INTERVAL_VALUE, SCHEDULER_INTERVAL_UNIT, BATCH_MAX_FILES_PER_DAY

Confira o docker-compose.yml para ver como aplicar facilmente.

## Testes de Carga
Utilizamos o K6 para simular carga alta de chargebacks.

Exemplo para rodar o teste:
```bash
k6 run chargeback-load-test.js
```
k6 run chargeback-load-test.js
