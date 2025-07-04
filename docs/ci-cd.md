# 🔍 CI/CD Pipeline - VideoGrinder

## Visão Geral

Pipeline de CI/CD otimizado para o VideoGrinder com **abordagem híbrida**: comandos nativos para performance e Docker onde necessário.

## 🚀 Workflow Principal

**Arquivo**: `.github/workflows/ci-cd.yml`

### Jobs do Pipeline

```
Quality Checks ⚡ (3-4 min)
├── Setup Go + Node.js (nativo)
└── make check-ci
    ├── fmt-ci (formatação)
    ├── lint-ci (linting)
    ├── test-ci (testes Go)
    └── test-js-ci (testes JS)

Build & Test 🏗️ (2-3 min)
├── Setup Docker
├── make setup prod
├── make run prod
└── make logs prod

E2E Tests 🎭 (1-2 min)
├── make setup dev
├── make run dev
├── make health-ci (aguarda app)
└── make test-e2e
```

**Tempo total**: ~6-9 minutos

## 🔧 Comandos Disponíveis

### Para Desenvolvimento (Docker)
```bash
make check          # Todos os quality checks
make health         # Health check via Docker
make run            # Iniciar aplicação
make logs           # Ver logs
make down           # Parar serviços
```

### Para CI/CD (Nativo - mais rápido)
```bash
make check-ci       # Todos os quality checks
make health-ci      # Health check direto
make fmt-ci         # Formatação
make lint-ci        # Linting
make test-ci        # Testes Go
make test-js-ci     # Testes JS
```

## ⚡ Otimizações Implementadas

### Performance
- **Quality checks nativos**: Go + Node.js instalados diretamente (3x mais rápido)
- **Docker apenas onde necessário**: Build de produção e E2E tests
- **Comandos centralizados**: Tudo via `make` para consistência
- **Health check eficiente**: `make health-ci` sem overhead do Docker

### Estrutura Híbrida
- **CI rápido**: Comandos `-ci` sem Docker para verificações básicas
- **Integração completa**: Docker para testes de produção e E2E
- **Flexibilidade**: Desenvolvedores podem usar Docker localmente

## 🛡️ Gates de Qualidade

Todo código deve passar por:
- ✅ **Formatação** (gofmt + eslint)
- ✅ **Linting** (golangci-lint + eslint)
- ✅ **Testes unitários** (Go + JavaScript)
- ✅ **Build de produção** (Docker)
- ✅ **Testes E2E** (Cypress)

## 🔄 Triggers

- **Push para main**: Pipeline completo
- **Pull Requests**: Pipeline completo
- **Falhas**: Logs e artefatos coletados automaticamente

## 📦 Artefatos

- **Cypress screenshots** (falhas, 30 dias)
- **Cypress videos** (sempre, 30 dias)
- **Logs de aplicação** (disponíveis no workflow)

## 🎯 Próximos Passos

Quando pronto para deploy:
1. Adicionar jobs de deploy aos workflows
2. Configurar infraestrutura AWS
3. Adicionar monitoramento de produção

---

**Estado atual**: Pipeline otimizado e estável, pronto para desenvolvimento contínuo. 
