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

O VideoGrinder implementa uma **arquitetura HTTP decoupling** com dois serviços independentes:

### 🎯 API Service (Porta 8080)
- **Responsabilidade**: Interface web, API REST, gerenciamento de arquivos
- **Endpoints**: `/` (web), `/api/v1/videos` (CRUD completo), `/health`
- **Comunicação**: HTTP client para Processor Service
- **Tecnologia**: Go + Gin + HTTP Client
- **Executable**: `cmd/api/main.go`

### ⚙️ Processor Service (Porta 8082)
- **Responsabilidade**: Processamento de vídeos, extração de frames
- **Endpoints**: `/process` (processamento), `/health` (status)
- **Tecnologia**: Go + Gin + FFmpeg
- **Isolamento**: Serviço independente e escalável  
- **Executable**: `cmd/processor/main.go`

### 🔗 Comunicação
- **API → Processor**: HTTP requests via client dedicado
- **Health Checks**: Verificação automática de disponibilidade
- **Timeout**: 5 minutos para processamento de vídeos
- **Error Handling**: Tratamento robusto de falhas de comunicação

### 📊 Benefícios
- ✅ **Escalabilidade**: Processor pode ter múltiplas instâncias
- ✅ **Isolamento**: Falhas em um serviço não afetam o outro
- ✅ **Manutenibilidade**: Desenvolvimento e deploy independentes
- ✅ **Testabilidade**: Testes isolados por serviço
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
make run          # Executar ambos os serviços (API + Processor)
make run-api      # Executar apenas o serviço API
make run-processor # Executar apenas o serviço Processor
```

**Testing:**
```bash
make test         # Executar todos os testes Go (API + Processor)
make test-api     # Executar apenas testes do serviço API
make test-processor # Executar apenas testes do serviço Processor
make test-js      # Executar testes JavaScript
make test-e2e     # Executar testes end-to-end
```

**Operations:**
```bash
make logs         # Ver logs de ambos os serviços
make logs-api     # Ver logs apenas do serviço API
make logs-processor # Ver logs apenas do serviço Processor
make down         # Parar todos os serviços
make help         # Ver todos os comandos disponíveis
```

### 👨‍💻 Fluxo de Desenvolvimento

Para contribuir com o projeto (seguindo nossos [Tech Mandates](./docs/tech-mandates.md)):

```bash
# 1. Executar aplicação com hot reload (auto-build)
make run      # Executar API + Processor services

# 2. Executar testes específicos durante desenvolvimento
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
# Testar desenvolvimento de API
make test-api         # Testes unitários da API
make run-api         # Executar apenas API service
make logs-api        # Ver logs apenas da API

# Testar desenvolvimento de Processor
make test-processor  # Testes unitários do Processor
make run-processor   # Executar apenas Processor service
make logs-processor  # Ver logs apenas do Processor

# Testar integração completa
make test           # Todos os testes (API + Processor + Services)
make run            # Ambos os serviços
make logs           # Logs de ambos os serviços
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
├── cmd/                 # Aplicações executáveis
│   ├── api/             # API Service
│   │   └── main.go      # Aplicação da API
│   ├── processor/       # Processor Service
│   │   └── main.go      # Aplicação do Processor
│   └── web/             # Web Service (opcional)
│       └── main.go      # Aplicação web standalone
├── internal/            # Código interno (não exportado)
│   ├── api/             # API Service handlers
│   │   ├── handlers.go  # Handlers HTTP da API
│   │   └── handlers_test.go # Testes da API
│   ├── processor/       # Processor Service handlers
│   │   ├── handlers.go  # Handlers HTTP do Processor
│   │   └── handlers_test.go # Testes do Processor
│   ├── clients/         # HTTP clients
│   │   └── processor.go # Cliente HTTP para Processor
│   ├── services/        # Lógica de negócio
│   │   ├── video.go     # Serviço de processamento de vídeo
│   │   └── video_test.go # Testes do serviço
│   ├── config/          # Configurações
│   │   └── config.go    # Estruturas de configuração
│   ├── models/          # Modelos de dados
│   │   └── types.go     # Tipos e estruturas
│   ├── utils/           # Utilitários
│   │   ├── validation.go # Validações de segurança
│   │   └── validation_test.go # Testes de validação
│   └── web/             # Web handlers (frontend)
│       └── handlers.go  # Handlers para páginas web
├── static/              # Arquivos estáticos (CSS, JS, HTML)
├── tests/               # Testes JavaScript
├── cypress/             # Testes end-to-end
├── docs/               # Documentação do projeto
│   ├── roadmap.md      # Roadmap de evolução
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
# Executar ambos os serviços
make run      # API (8080) + Processor (8082) em desenvolvimento
make run prod # API (8080) + Processor (8082) em produção

# Executar serviços individualmente
make run-api      # Apenas API service na porta 8080
make run-processor # Apenas Processor service na porta 8082

# Monitoramento
make logs         # Logs de ambos os serviços
make logs-api     # Logs apenas da API
make logs-processor # Logs apenas do Processor
```

### Configurações Atuais
- **API Service**: Porta 8080 (interface externa)
- **Processor Service**: Porta 8082 (processamento interno)
- **Comunicação**: HTTP entre serviços com timeout de 5 minutos
- **Taxa de extração**: 1 frame por segundo (fps=1)  
- **Formatos suportados**: MP4, AVI, MOV, MKV, WMV, FLV, WebM

### Variáveis de Ambiente
```bash
# Configuração do Processor Service
export PROCESSOR_URL=http://localhost:8082  # URL do Processor Service

# Configuração de diretórios (opcional)
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
# Verificar se o Processor está rodando
curl http://localhost:8082/health

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
- Verifique se o Processor Service está acessível: `curl http://localhost:8082/health`

### Portas em uso
```bash
# Porta 8080 (API) ou 8082 (Processor) em uso
make down     # Parar todos os serviços do VideoGrinder

# Verificar processos nas portas
lsof -ti:8080 | xargs kill -9  # API
lsof -ti:8082 | xargs kill -9  # Processor
```

### Problemas com serviços individuais
```bash
# Testar apenas API
make test-api
make run-api

# Testar apenas Processor
make test-processor
make run-processor

# Verificar saúde dos serviços
curl http://localhost:8080/api/v1/videos  # API
curl http://localhost:8082/health         # Processor
```

### Erro de permissão em diretórios
```bash
sudo chmod 755 uploads outputs temp
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
- Interface web básica integrada na API (será separada nas próximas fases)
- Comunicação HTTP entre serviços adiciona latência mínima

## 🎯 Melhorias com Multi-Service Architecture

- ✅ **Escalabilidade**: Múltiplas instâncias do Processor podem processar vídeos simultaneamente
- ✅ **Isolamento**: Falhas no processamento não afetam a API
- ✅ **Manutenção**: Serviços podem ser atualizados independentemente
- ✅ **Monitoramento**: Logs e métricas separados por serviço
- ✅ **Testabilidade**: Testes unitários isolados por responsabilidade

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
3. **🐳 Use Docker** para desenvolvimento (conforme mandates)

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
