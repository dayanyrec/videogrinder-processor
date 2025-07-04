# 🔍 CI/CD Pipeline - VideoGrinder

## Visão Geral

Pipeline de CI/CD **simplificado e focado** no VideoGrinder com **5 passos essenciais** para garantir qualidade e funcionalidade.

## 🎯 Filosofia: 5 Passos Essenciais

1. **✨ Formatação** - Código está bem formatado?
2. **🔍 Lint** - Qualidade do código está ok?
3. **🧪 Testes** - Funcionalidades estão testadas?
4. **🏗️ Build** - Imagens são geradas corretamente?
5. **🚀 Health** - Serviços funcionam em produção?

## 🚀 Workflow Principal

**Arquivo**: `.github/workflows/ci-cd.yml`

### Jobs do Pipeline

```
🔍 Code Quality Checks (2-3 min)
├── Step 1: make fmt-check (formatação)
├── Step 2: make lint (qualidade)
└── Step 3: make test + make test-js (testes)

🏗️ Build & Health Check (2-3 min)
├── Step 4: make setup prod (build)
└── Step 5: make health (verificação completa)
```

**Tempo total**: ~4-6 minutos

## 🔧 Comandos Principais

### Desenvolvimento Local
```bash
make fmt-check      # Verificar formatação
make lint           # Verificar qualidade
make test           # Testes Go
make test-js        # Testes JavaScript
make setup prod     # Build produção
make health         # Health check completo
```

### Health Check Completo
```bash
make health
# Exemplo de saída com todos os serviços saudáveis:
# 🏥 Checking application health...
# 🌐 Checking Web Service (port 8080)...
# ✅ Web Service: healthy
# 🔌 Checking API Service (port 8081)...
# ✅ API Service: healthy
# ⚙️  Checking Processor Service (port 8082)...
# ✅ Processor Service: healthy
# ✅ All services are healthy!
```

> **Nota:**
> O comando `make health` mostra apenas o status de cada serviço.
> Se algum serviço falhar, a execução para imediatamente e exibe uma dica:
> 
> ```
> ❌ Web Service: failed
> 💡 Dica: rode 'make logs-tail' para ver os logs.
> ```
> Assim, fica fácil identificar e debugar problemas rapidamente.

## ⚡ Benefícios da Simplificação

### Performance
- **Foco no essencial**: Apenas os 5 passos críticos
- **Pipeline linear**: Sem complexidade desnecessária
- **Feedback rápido**: Falha logo no primeiro problema

### Clareza
- **Steps numerados**: 1, 2, 3, 4, 5 - fácil de entender
- **Propósito claro**: Cada step tem objetivo específico
- **Diagnóstico simples**: Sabe exatamente onde falhou

### Manutenibilidade
- **Menos código**: Workflow mais enxuto
- **Menos dependências**: Foco no Docker onde necessário
- **Comandos centralizados**: Tudo via `make`

## 🛡️ Gates de Qualidade

Todo código deve passar pelos **5 Steps**:

| Step | Comando | Verifica |
|------|---------|----------|
| 1 | `make fmt-check` | Formatação Go + JS |
| 2 | `make lint` | Qualidade Go + JS |
| 3 | `make test` + `make test-js` | 70+ testes Go + 59 testes JS |
| 4 | `make setup prod` | Build das 3 imagens |
| 5 | `make health` | 3 serviços funcionando |

## 🔄 Triggers

- **Push para main**: Pipeline completo
- **Pull Requests**: Pipeline completo
- **Falhas**: Logs detalhados para debug

## 📊 Arquitetura Verificada

O pipeline valida toda a arquitetura multi-serviços:

```
🌐 Web Service (8080) ──┐
                        ├── make health
🔌 API Service (8081) ──┤
                        │
⚙️ Processor (8082) ────┘
```

Cada serviço retorna:
- Status de saúde
- Verificação de dependências
- Latência de resposta
- Estado dos diretórios

## 📋 Comandos de Logs

### Logs sem Travar Execução
```bash
make logs-tail           # Últimas 50 linhas de todos os serviços
make logs-web-tail       # Últimas 30 linhas do Web service
make logs-api-tail       # Últimas 30 linhas do API service  
make logs-processor-tail # Últimas 30 linhas do Processor service
```

### Logs com Follow (para desenvolvimento)
```bash
make logs               # Todos os serviços (travado)
make logs-web           # Web service (travado)
make logs-api           # API service (travado)
make logs-processor     # Processor service (travado)
```

**Uso**: `make logs-tail [dev|prod]` - padrão é `dev`

### Exemplo de Uso
```bash
# Para CI/CD - logs rápidos sem travar
make logs-tail prod

# Para debug específico
make logs-api-tail dev

# Para monitoramento contínuo
make logs-web dev
```

### CI/CD Pipeline
O pipeline usa `logs-tail` para não travar a execução:
```bash
# Em caso de falha no health check
make logs-tail prod  # Mostra últimas 50 linhas e continua
```

## 🎯 Próximos Passos

- [ ] **E2E Tests**: Integrar testes Cypress quando estáveis
- [ ] **Deploy**: Adicionar deploy automático para produção
- [ ] **Monitoramento**: Health checks em produção

---

**Estado atual**: Pipeline **otimizado e estável** - 5 steps essenciais funcionando perfeitamente.
