# 🎬 VideoGrinder

## Sobre o Projeto

O **VideoGrinder** é uma ferramenta desenvolvida especificamente para **jornalistas** que precisam extrair frames de vídeos para criar conteúdo jornalístico, posts em redes sociais ou materiais de apoio para suas matérias.

Esta plataforma permite que os usuários façam upload de vídeos através de uma interface web e recebam um arquivo ZIP contendo todos os frames extraídos em formato PNG, facilitando o processo de seleção de imagens para uso editorial.

> 📋 **Roadmap de Evolução**: Este projeto está em desenvolvimento ativo seguindo nosso [roadmap detalhado](./docs/roadmap.md), que descreve a evolução planejada de monólito para arquitetura de microserviços.
>
> 🏛️ **Tech Mandates**: Todas as decisões técnicas seguem nossos rigorosos [Tech Mandates](./docs/tech-mandates.md), garantindo consistência arquitetural e operacional.

## ✨ Funcionalidades

- 📤 **Upload de vídeos**: Interface web intuitiva para envio de arquivos
- 🖼️ **Extração automática de frames**: Converte vídeos em frames individuais (1 frame por segundo)
- 📦 **Download em ZIP**: Todos os frames são compactados automaticamente
- 🎯 **Formatos suportados**: MP4, AVI, MOV, MKV, WMV, FLV, WebM
- 📊 **Status de processamento**: Acompanhe o andamento da extração
- 📁 **Histórico de arquivos**: Visualize e baixe processamentos anteriores
- 🌐 **Interface em português**: Totalmente localizada para usuários brasileiros

## 🛠️ Tecnologias Utilizadas

- **Backend**: Go (Golang) com framework Gin
- **Processamento de vídeo**: FFmpeg
- **Frontend**: HTML, CSS e JavaScript (integrado)
- **Containerização**: Docker
- **Arquivos**: Manipulação de ZIP nativo

## 🏗️ Arquitetura Multi-Service

O VideoGrinder implementa uma **arquitetura HTTP decoupling** com três serviços independentes:

### 🌐 Web Service (Porta 8080)
- **Responsabilidade**: Interface web, arquivos estáticos, frontend
- **Endpoints**: `/` (página principal), `/static/*` (arquivos estáticos), `/health`
- **Tecnologia**: Go + Gin + Static File Serving
- **Executable**: `web/cmd/main.go`

### 🎯 API Service (Porta 8081)
- **Responsabilidade**: API REST, gerenciamento de arquivos, comunicação com Processor
- **Endpoints**: `/api/v1/videos` (CRUD completo), `/health`
- **Comunicação**: HTTP client para Processor Service
- **Tecnologia**: Go + Gin + HTTP Client
- **Executable**: `api/cmd/main.go`

### ⚙️ Processor Service (Porta 8082)
- **Responsabilidade**: Processamento de vídeos, extração de frames
- **Endpoints**: `/process` (processamento), `/health` (status)
- **Tecnologia**: Go + Gin + FFmpeg
- **Isolamento**: Serviço independente e escalável  
- **Executable**: `processor/cmd/main.go`

### 🔗 Comunicação
- **Web → API**: Frontend JavaScript via AJAX/REST
- **API → Processor**: HTTP requests via client dedicado
- **Health Checks**: Verificação automática de disponibilidade
- **Timeout**: 5 minutos para processamento de vídeos
- **Error Handling**: Tratamento robusto de falhas de comunicação

### 📊 Benefícios
- ✅ **Escalabilidade**: Processor pode ter múltiplas instâncias
- ✅ **Isolamento**: Falhas em um serviço não afetam os outros
- ✅ **Manutenibilidade**: Desenvolvimento e deploy independentes
- ✅ **Testabilidade**: Testes isolados por serviço
- ✅ **Separação Frontend/Backend**: Interface totalmente desacoplada
- ✅ **Microservices Ready**: Preparado para Kubernetes

## 🏛️ Tech Mandates

O VideoGrinder segue um conjunto rigoroso de **[Tech Mandates](./docs/tech-mandates.md)** que definem nossa arquitetura e práticas de desenvolvimento:

- **☁️ AWS como cloud provider exclusivo** - Integração nativa com serviços AWS
- **🏭 Ambiente único de produção** - Desenvolvimento local → Produção direta
- **🐳 Docker-first development** - Zero dependências locais além do Docker
- **🚫 Código autoexplicativo** - Sem comentários desnecessários no código
- **📋 Testes como documentação** - Especificação viva através de testes
- **🔒 Security by design** - AWS Secrets Manager, IAM, KMS integrados
- **📊 Observabilidade completa** - CloudWatch, X-Ray para monitoramento
- **🔧 Infrastructure as Code** - Terraform para toda infraestrutura

> 📖 **Consulte nossos [Tech Mandates completos](./docs/tech-mandates.md)** para entender as diretrizes técnicas que guiam todas as decisões de arquitetura do projeto.

## 📋 Pré-requisitos

- Docker instalado
- Git (para clonagem do repositório)

## 🚀 Como Executar

1. **Clone o repositório:**
```bash
git clone <url-do-repositorio>
cd videogrinder-processor
```

2. **Execute a aplicação (auto-build):**
```bash
make run      # API + Processor services (desenvolvimento)
make run prod # API + Processor services (produção)
```

3. **Acesse no navegador:**
```
http://localhost:8080    # Interface web + API REST
```

### 🛠️ Comandos Essenciais

**Multi-Service Architecture:**
```bash
make run          # Executar todos os serviços (Web + API + Processor)
make run-web      # Executar apenas o serviço Web
make run-api      # Executar apenas o serviço API
make run-processor # Executar apenas o serviço Processor
```

**Testing:**
```bash
make test         # Executar todos os testes Go (Web + API + Processor)
make test-web     # Executar apenas testes do serviço Web
make test-api     # Executar apenas testes do serviço API
make test-processor # Executar apenas testes do serviço Processor
make test-js      # Executar testes JavaScript
make test-e2e     # Executar testes end-to-end
```

**Operations:**
```bash
make logs         # Ver logs de todos os serviços
make logs-web     # Ver logs apenas do serviço Web
make logs-api     # Ver logs apenas do serviço API
make logs-processor # Ver logs apenas do serviço Processor
make down         # Parar todos os serviços
make help         # Ver todos os comandos disponíveis
```

### 👨‍💻 Fluxo de Desenvolvimento

Para contribuir com o projeto (seguindo nossos [Tech Mandates](./docs/tech-mandates.md)):

```bash
# 1. Executar aplicação com hot reload (auto-build)
make run      # Executar Web + API + Processor services

# 2. Executar testes específicos durante desenvolvimento
make test-web        # Testar apenas Web service
make test-api        # Testar apenas API service
make test-processor  # Testar apenas Processor service
make test           # Testar todos os serviços

# 3. Executar verificações antes de commit
make check    # Executa: format + lint + test (todos os serviços)

# 4. Parar serviços quando terminar
make down
```

### 🧪 Exemplos de Teste por Serviço

```bash
# Testar desenvolvimento de Web
make test-web        # Testes unitários do Web
make run-web         # Executar apenas Web service
make logs-web        # Ver logs apenas do Web

# Testar desenvolvimento de API
make test-api        # Testes unitários da API
make run-api         # Executar apenas API service
make logs-api        # Ver logs apenas da API

# Testar desenvolvimento de Processor
make test-processor  # Testes unitários do Processor
make run-processor   # Executar apenas Processor service
make logs-processor  # Ver logs apenas do Processor

# Testar integração completa
make test           # Todos os testes (Web + API + Processor + Services)
make run            # Todos os serviços
make logs           # Logs de todos os serviços
```

## 📖 Como Usar

1. **Acesse a interface web** em `http://localhost:8080`

2. **Selecione um arquivo de vídeo** clicando em "Selecione um arquivo de vídeo"
   - Formatos aceitos: `.mp4`, `.avi`, `.mov`, `.mkv`, `.wmv`, `.flv`, `.webm`

3. **Clique em "🚀 Processar Vídeo"**
   - O sistema extrairá 1 frame por segundo do vídeo
   - O processamento pode levar alguns minutos dependendo do tamanho do vídeo

4. **Faça o download do ZIP**
   - Após o processamento, um link de download será exibido
   - O arquivo ZIP conterá todos os frames em formato PNG

5. **Visualize o histórico**
   - Na seção "Arquivos Processados" você pode ver e baixar processamentos anteriores

## 📁 Estrutura do Projeto

```
videogrinder-processor/
├── api/                 # API Service (Porta 8081)
│   ├── cmd/main.go      # Aplicação da API
│   └── internal/        # Código interno da API
│       ├── handlers/    # Handlers HTTP da API
│       ├── clients/     # Cliente HTTP para Processor
│       ├── config/      # Configurações da API
│       └── models/      # Modelos de dados da API
├── processor/           # Processor Service (Porta 8082)
│   ├── cmd/main.go      # Aplicação do Processor
│   └── internal/        # Código interno do Processor
│       ├── handlers/    # Handlers HTTP do Processor
│       ├── services/    # Lógica de processamento de vídeo
│       ├── config/      # Configurações do Processor
│       ├── models/      # Modelos de dados do Processor
│       └── utils/       # Utilitários e validações
├── web/                 # Web Service (Porta 8080)
│   ├── cmd/main.go      # Aplicação do Web
│   ├── internal/        # Código interno do Web
│   │   ├── handlers/    # Handlers HTTP do Web
│   │   └── config/      # Configurações do Web
│   ├── static/          # Arquivos estáticos (CSS, JS, HTML)
│   │   ├── css/styles.css # Estilos CSS
│   │   ├── index.html   # Página principal
│   │   └── js/          # JavaScript modules
│   ├── tests/           # Testes JavaScript
│   ├── cypress/         # Testes end-to-end
│   ├── .eslintrc.js     # Configuração ESLint
│   ├── cypress.config.js # Configuração do Cypress
│   └── package.json     # Dependências Node.js
├── internal/            # Código compartilhado
│   └── config/          # Configurações base compartilhadas
├── docs/               # Documentação do projeto
│   ├── roadmap.md      # Roadmap de evolução
│   ├── architecture.md # Arquitetura detalhada
│   └── tech-mandates.md # Diretrizes técnicas obrigatórias
├── uploads/            # Vídeos enviados (temporário)
├── outputs/            # Arquivos ZIP gerados
├── temp/               # Arquivos temporários durante processamento
├── docker-compose.yml  # Configuração multi-service
├── Dockerfile          # Configuração do Docker
├── Makefile           # Comandos de automação
└── README.md          # Este arquivo
```

## 🔧 Configuração

### Ambiente Multi-Service
```bash
# Executar todos os serviços
make run      # Web (8080) + API (8081) + Processor (8082) em desenvolvimento
make run prod # Web (8080) + API (8081) + Processor (8082) em produção

# Executar serviços individualmente
make run-web      # Apenas Web service na porta 8080
make run-api      # Apenas API service na porta 8081
make run-processor # Apenas Processor service na porta 8082

# Monitoramento
make logs         # Logs de todos os serviços
make logs-web     # Logs apenas do Web
make logs-api     # Logs apenas da API
make logs-processor # Logs apenas do Processor
```

### Configurações Atuais
- **Web Service**: Porta 8080 (interface web)
- **API Service**: Porta 8081 (API REST)
- **Processor Service**: Porta 8082 (processamento interno)
- **Comunicação**: HTTP entre serviços com timeout de 5 minutos
- **Taxa de extração**: 1 frame por segundo (fps=1)  
- **Formatos suportados**: MP4, AVI, MOV, MKV, WMV, FLV, WebM

### Variáveis de Ambiente
```bash
# Web Service (Porta 8080)
export PORT=8080
export API_URL=http://localhost:8081

# API Service (Porta 8081)
export PORT=8081
export PROCESSOR_URL=http://localhost:8082

# Processor Service (Porta 8082)
export PORT=8082

# Configuração de diretórios (compartilhada)
export UPLOADS_DIR=./uploads
export OUTPUTS_DIR=./outputs
export TEMP_DIR=./temp
```

> ⚠️ **Nota**: Configurações adicionais via variáveis de ambiente serão implementadas na Fase 1.4 conforme nosso [roadmap](./docs/roadmap.md).

## 🐛 Solução de Problemas

### Aplicação não inicia
```bash
make down     # Parar serviços existentes
make setup    # Reconfigurar ambiente
make run      # Tentar executar novamente
```

### Verificar logs da aplicação
```bash
make logs           # Ver logs de ambos os serviços
make logs-api       # Ver logs apenas da API
make logs-processor # Ver logs apenas do Processor
```

### Erro de comunicação entre serviços
```bash
# Verificar se o Processor está rodando via Docker
docker-compose run --rm videogrinder-web-dev sh -c "curl http://localhost:8082/health"

# Verificar se a API consegue acessar o Processor
make logs-api | grep "processor"

# Reiniciar ambos os serviços
make down
make run
```

### Vídeo não é processado
- Verifique se o formato é suportado
- Confirme se o arquivo não está corrompido
- Execute `make logs-processor` para ver erros específicos do processamento
- Verifique se o Processor Service está acessível via Docker: `docker-compose run --rm videogrinder-web-dev sh -c "curl http://localhost:8082/health"`

### Portas em uso
```bash
# Portas 8080 (Web), 8081 (API) ou 8082 (Processor) em uso
make down     # Parar todos os serviços do VideoGrinder

# Verificar processos nas portas via Docker (se necessário)
docker-compose run --rm videogrinder-web-dev sh -c "netstat -tulpn | grep :8080"  # Web
docker-compose run --rm videogrinder-web-dev sh -c "netstat -tulpn | grep :8081"  # API
docker-compose run --rm videogrinder-web-dev sh -c "netstat -tulpn | grep :8082"  # Processor
```

### Problemas com serviços individuais
```bash
# Testar apenas Web
make test-web
make run-web

# Testar apenas API
make test-api
make run-api

# Testar apenas Processor
make test-processor
make run-processor

# Verificar saúde dos serviços via Docker
docker-compose run --rm videogrinder-web-dev sh -c "curl http://localhost:8080/health"  # Web
docker-compose run --rm videogrinder-web-dev sh -c "curl http://localhost:8081/health"  # API
docker-compose run --rm videogrinder-web-dev sh -c "curl http://localhost:8082/health"  # Processor
```

### Erro de permissão em diretórios
```bash
# Ajustar permissões via Docker (se necessário)
docker-compose run --rm videogrinder-web-dev sh -c "chmod 755 uploads outputs temp"
```

### Problemas com Docker
```bash
make docker-clean    # Limpar recursos Docker
make setup          # Recriar ambiente
```

## 🎯 Casos de Uso para Jornalistas

- **Matérias esportivas**: Extrair momentos-chave de jogos
- **Eventos políticos**: Capturar gestos e expressões importantes
- **Coberturas ao vivo**: Gerar imagens para posts em tempo real
- **Análise de conteúdo**: Estudar sequências de vídeo frame por frame
- **Redes sociais**: Criar carrosséis de imagens para Instagram/Twitter
- **Documentação**: Arquivo visual de eventos importantes

## ⚠️ Limitações Atuais

- O processamento é sequencial (um vídeo por vez por instância de Processor)
- Arquivos muito grandes podem consumir bastante espaço em disco
- O tempo de processamento é proporcional ao tamanho e duração do vídeo
- Comunicação HTTP entre serviços adiciona latência mínima
- Armazenamento local (será migrado para S3 na Fase 3)

## 🎯 Melhorias com Multi-Service Architecture

- ✅ **Escalabilidade**: Múltiplas instâncias do Processor podem processar vídeos simultaneamente
- ✅ **Isolamento**: Falhas em um serviço não afetam os outros
- ✅ **Manutenção**: Serviços podem ser atualizados independentemente
- ✅ **Monitoramento**: Logs e métricas separados por serviço
- ✅ **Testabilidade**: Testes unitários isolados por responsabilidade
- ✅ **Separação Frontend/Backend**: Interface totalmente desacoplada

## 🗺️ Roadmap de Evolução

Este projeto está em constante evolução seguindo um roadmap estruturado que visa transformar o VideoGrinder de um monólito em uma arquitetura de microserviços escalável:

- **Fase 1**: Tornar o projeto produtivo com testes, CI/CD e infraestrutura
- **Fase 2**: Modularização interna (ainda no monólito)
- **Fase 3**: Persistência e rastreabilidade com DynamoDB
- **Fase 4**: Arquitetura de microserviços completa

Para detalhes completos sobre as fases, cronograma e entregas, consulte nosso **[Roadmap Detalhado](./docs/roadmap.md)**.

### Próximas Entregas (Fase 1)
- [x] Setup de linters e boas práticas
- [x] Melhorar containerização com Docker multistage
- [x] **HTTP Decoupling**: Arquitetura multi-service implementada (API + Processor)
- [x] **Testes Unitários**: Cobertura completa para ambos os serviços
- [x] **Makefile Atualizado**: Comandos para desenvolvimento multi-service
- [x] **Limpeza Monolítica**: Remoção completa do código monolítico legado
- [ ] **CRÍTICO**: Corrigir vulnerabilidades de segurança (G304, G204, errcheck)
- [ ] Adicionar variáveis de ambiente para configuração
- [ ] Implementar logging estruturado em JSON
- [ ] Implementar testes end-to-end
- [ ] Configurar CI/CD com GitHub Actions

## 🤝 Contribuição

Contribuições são bem-vindas! Antes de contribuir:

1. **📋 Consulte nosso [roadmap](./docs/roadmap.md)** para entender a direção do projeto
2. **🏛️ Leia nossos [Tech Mandates](./docs/tech-mandates.md)** para seguir nossas diretrizes técnicas
3. **🐳 Use Docker exclusivamente** - Todos os comandos devem ser executados via containers (npm, go, curl, etc.)
4. **🚫 Zero dependências locais** - Apenas Docker e Git são necessários na máquina do desenvolvedor

Sinta-se à vontade para:
- Reportar bugs
- Sugerir melhorias
- Enviar pull requests
- Compartilhar casos de uso

## 📞 Suporte

Para dúvidas ou problemas:
1. Verifique a seção "Solução de Problemas" 
2. Consulte nosso [roadmap](./docs/roadmap.md) para entender o status do projeto
3. Revise nossos [Tech Mandates](./docs/tech-mandates.md) para questões arquiteturais
4. Consulte os logs da aplicação
5. Abra uma issue no repositório

---

**Desenvolvido com ❤️**
