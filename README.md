# Rate Limiter com Go e Redis

## 📌 Objetivo

Este projeto implementa um **Rate Limiter** em Go com suporte a limitação de requisições por:

- **Endereço IP**
- **Token de acesso** (via header `API_KEY: <TOKEN>`)

O rate limiter é configurável por meio de variáveis de ambiente e utiliza o **Redis** como mecanismo de persistência. Ele atua como **middleware** para ser facilmente integrado em servidores HTTP.

> ⚠️ A limitação por token sobrepõe a limitação por IP, caso ambos estejam presentes.

---

## ⚙️ Tecnologias

- [Go](https://golang.org/)
- [Redis](https://redis.io/)
- [Docker Compose](https://docs.docker.com/compose/)
- [miniredis](https://github.com/alicebob/miniredis) (para testes)
- [ab (Apache Bench)](https://httpd.apache.org/docs/2.4/programs/ab.html) (para testes de carga)

---

## 🚀 Como executar com Docker Compose

### 1. Clone o repositório
```bash
git clone https://github.com/seu-usuario/rate-limiter-go.git
cd rate-limiter-go
```

### 2. Configure o .env
```bash
DEFAULT_RATE_LIMIT=5
DEFAULT_BLOCK_DURATION=300
REDIS_HOST=localhost:6379
REDIS_PASSWORD=sysdba
REDIS_DB=0
RATE_LIMIT_TOKEN=100
BLOCK_DURATION_TOKEN=300
PORT=8080
```

### 3. Suba aplicação
```bash
docker-compose up --build
```
A aplicação estará disponível em: http://localhost:8080/hello

## 🧪 Executando Testes de Carga

### 1. Instale o Apache Bench (caso não tenha)
Ubuntu:
```bash
sudo apt install apache2-utils
```
macOS:
```bash
brew install httpd
```

### 2. Exemplo: 5000 requisições com até 10 simultâneas
```bash
ab -n 5000 -c 10 http://localhost:8080/hello
```

## 🧪 Testes unitários

### 1. Instale o Apache Bench (caso não tenha)
Para rodar os testes unitários do projeto:
```bash
go test ./... -v
```
