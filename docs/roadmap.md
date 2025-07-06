# :motorway: VideoGrinder - Roadmap de Evolução

Este roadmap descreve os passos planejados para amadurecer o projeto VideoGrinder, evoluindo de um monólito para uma arquitetura modular e escalável.

---

## :white_check_mark: Fase 0 – Fortalecer fundação do projeto ✅ **CONCLUÍDA**

- [x] Criar repositório para POC para manter histórico de fácil acesso
- [x] Documentar propósito e escopo da POC (`README.md`)
- [x] Criar repositório para o projeto onde será feita a evolução
- [x] Documentar roadmap com ideias atuais

---

## :rocket: Fase 1 – Tornar a POC um projeto produtivo ✅ **CONCLUÍDA**

### 1.1 - Desenvolvimento e Qualidade de Código ✅ **CONCLUÍDA**
- [x] Setup de `.editorconfig`, linters e boas práticas
- [x] Melhorar containerização com Docker multistage
- [x] Configurar linter para compliance com Tech Mandates
- [x] Configurar hot reload para desenvolvimento (Air)
- [x] Implementar Makefile simplificado para comandos essenciais
- [x] Configurar ambiente de desenvolvimento Docker-first

### 1.2 - Correções Críticas de Segurança ✅ **CONCLUÍDA**
- [x] **CRÍTICO**: Corrigir vulnerabilidades de path traversal (G304)
- [x] **CRÍTICO**: Corrigir vulnerabilidade de command injection em FFmpeg (G204)
- [x] **CRÍTICO**: Corrigir 9 erros não verificados (errcheck)
- [x] Implementar validação rigorosa de inputs de arquivos
- [x] Adicionar sanitização de nomes de arquivos
- [x] Corrigir permissões inseguras de diretório (G301)
- [x] Refatorar função processVideo para reduzir complexidade ciclomática
- [x] Implementar funções especializadas para validação e limpeza

### 1.5 - Testes e Qualidade ✅ **CONCLUÍDA**
- [x] Corrigir issues de pre-alocação de slices (prealloc)
- [x] Implementar validação automatizada com linters de segurança
- [x] Cobrir código atual com testes unitários básicos
- [x] Cobrir experiência com testes end-to-end (Cypress)

### 1.6 - Pipeline de Qualidade ✅ **CONCLUÍDA**
- [x] Configurar pipeline com GitHub Actions (trunk-based)
- [x] Implementar validação automatizada com Docker multi-stage
- [x] Configurar validação de PR com testes automatizados
- [x] Implementar pipeline completo de qualidade de código
- [x] Configurar Dependabot para atualizações de dependências
- [x] Documentar comandos Make para validação local

---

## :jigsaw: Fase 2 – Modularização (ainda no monolito) ✅ **CONCLUÍDA**

### 2.1 - Estruturação do Código ✅ **CONCLUÍDA**
- [x] Extrair configuração para internal/config (variáveis de ambiente)
- [x] Extrair core de processamento para internal/services/video
- [x] Cobrir o core com testes unitários
- [x] Extrair handlers HTTP para internal/handlers/web
- [x] Cobrir handlers com testes unitários (86 testes passando)
- [x] Extrair frontend para arquivos estáticos
- [x] Cobrir frontend (JS) com testes unitários

### 2.2 - API Design e Qualidade ✅ **CONCLUÍDA**
- [x] Extrair API para módulo interno
- [x] Melhorar implementação da API tornando REST
- [x] Cobrir API com testes unitários
- [x] Cobrir API com testes de integração

---

## :gear: Fase 4 – Arquitetura de microserviços

### 4.0 - Arquitetura Multi-Service ✅ **CONCLUÍDA**

- [x] API e processor não devem se integrar diretamente e sim por HTTP
- [x] Adicionar health check endpoints robustos
- [x] Padronizar formato de health check entre API e Processor (estrutura unificada com checks específicos)
- [x] Remover código referente ao monolito (projeto será modular daqui em diante)
- [x] Criar arquitetura de 3 serviços independentes (Web + API + Processor)
- [x] Separar Web Service (porta 8080) para frontend e arquivos estáticos
- [x] Separar API Service (porta 8081) para REST API e comunicação com Processor
- [x] Manter Processor Service (porta 8082) para processamento de vídeos
- [x] Implementar comunicação HTTP entre API e Processor via client dedicado
- [x] Criar estrutura modular com cada serviço tendo seu próprio internal/
- [x] Migrar models para cada serviço específico (api/internal/models, processor/internal/models)
- [x] Criar configurações independentes por serviço com próprias portas
- [x] Implementar Docker Compose multi-service com hot reload
- [x] Atualizar Makefile com comandos para desenvolvimento multi-service
- [x] Containerizar todos os comandos npm seguindo tech-mandates
- [x] Padronizar endpoints de health (/health) em todos os serviços
- [x] Atualizar documentação (README.md e architecture.md) para refletir nova arquitetura

**Resultado**: Arquitetura multi-service totalmente funcional com 3 serviços independentes, comunicação HTTP, configurações isoladas, testes completos (39+ testes Go + 59 testes JS + 19 testes E2E) e documentação atualizada. O projeto está preparado para microserviços e segue todos os tech-mandates.

### 4.1 - Integração LocalStack AWS Básica
- [x] **Setup LocalStack Environment**
  - [x] Configurar LocalStack com S3, DynamoDB e SQS
  - [x] Criar docker-compose para LocalStack integrado ao ambiente de desenvolvimento
  - [x] Configurar credenciais AWS locais para desenvolvimento
- [x] **Migração de Armazenamento para S3**
  - [x] Implementar AWS S3 client para upload de vídeos
  - [x] Migrar diretório uploads/ para S3 buckets
  - [x] Migrar diretório outputs/ para S3 buckets
  - [x] Implementar presigned URLs para download de arquivos
- [ ] **Persistência com DynamoDB**
  - [x] Criar schema da tabela `video_jobs` no DynamoDB
  - [ ] Implementar repository pattern para jobs de processamento
  - [ ] Salvar solicitação recebida na API no DynamoDB
  - [ ] Atualizar status do job quando Processor terminar processamento
- [ ] **Comunicação Assíncrona**
  - [ ] Implementar SQS para comunicação API → Processor
  - [ ] Refatorar API para enviar jobs via SQS em vez de HTTP direto
  - [ ] Implementar polling no frontend para atualizar status em tempo real
- [x] **Configuração e Ambiente**
  - [x] Adicionar variáveis de ambiente para AWS services
  - [x] Implementar fallback para desenvolvimento local vs LocalStack
  - [x] Atualizar Makefile com comandos para LocalStack

**Resultado**: Infraestrutura LocalStack totalmente integrada com LocalStack 3.0, S3 client implementado com presigned URLs, schema DynamoDB e SQS configurados, fallback para desenvolvimento local funcional, health checks S3, comandos Makefile para LocalStack e integração completa com CI/CD. Sistema preparado para migração completa para AWS.

### 4.2 - Preparação dos repositórios

- [ ] Criar repositório `videogrinder-frontend`
- [ ] Criar repositório `videogrinder-api` 
- [ ] Criar repositório `videogrinder-handler`
- [ ] Manter `videogrinder-processor` como repositório atual
- [ ] Definir contratos de API entre serviços

### 4.3 - Extração gradual (uma por vez)

- [ ] Extrair e deployar frontend em S3
- [ ] Implementar CI/CD para frontend
- [ ] Testar integração frontend -> API monolítica
- [ ] Extrair API para repositório próprio
- [ ] Implementar CI/CD para API em EKS
- [ ] Testar integração frontend -> API extraída

### 4.4 - Handler e orquestração

- [ ] Criar Lambda handler para orquestração  
- [ ] Implementar CI/CD para handler
- [ ] Ativar stream do DynamoDB -> handler
- [ ] Setup do processor como Job K8S
- [ ] Trigger processor via handler

### 4.5 - Integração e otimização

- [ ] Implementar processamento paralelo de múltiplos vídeos
- [x] Testes de integração end-to-end
- [ ] Implementar rollback strategies
- [ ] Monitoramento e alertas básicos

---

## :seedling: Fase 5 – Futuro

### 5.1 - Otimizações de Arquitetura
- [ ] Adicionar SQS entre Lambda e Job (desacoplamento)
- [ ] Observabilidade com OpenTelemetry (traços e métricas avançados)
- [ ] Implementar processamento paralelo inteligente
- [ ] Otimizar custos com Spot Instances

### 5.2 - Funcionalidades Avançadas
- [ ] Melhorar UX com preview e status visual em tempo real
- [ ] Preview dos frames antes do download
- [ ] Interface mais sofisticada e responsiva
- [ ] Sistema de notificações em tempo real
- [ ] Seleção customizada de taxa de frames (fps)
- [ ] Suporte a mais formatos de saída (JPEG, WebP, GIF)
- [ ] Implementar compressão inteligente de imagens

### 5.3 - Funcionalidades de Negócio
- [ ] Adicionar autenticação robusta (Cognito)
- [ ] Permitir busca pela listagem de vídeos
- [ ] Permitir inserir tags que agrupem vídeos e facilitem a busca
- [ ] Implementar quotas por usuário
- [ ] Dashboard de analytics para jornalistas

---

### Parking lot - Entender se são necessárias
- [x] Adicionar variáveis de ambiente para configuração
- [ ] Implementar timeouts para operações FFmpeg
- [ ] Adicionar rate limiting básico
- [ ] Documentar API com Swagger/OpenAPI automático

### 1.3 - Observabilidade e Monitoramento
- [x] Implementar logging estruturado em JSON
- [x] Melhorar tratamento de erros com contexto adequado
- [ ] Integrar CloudWatch Logs para desenvolvimento
- [ ] Implementar métricas básicas de performance
- [ ] Adicionar monitoramento de recursos (CPU, memória, I/O)

### 1.4 - Integração AWS Básica
- [ ] Integrar AWS Secrets Manager para configurações sensíveis
- [ ] Configurar AWS IAM roles e políticas básicas
- [x] Setup do LocalStack com S3 e DynamoDB (desenvolvimento local)
- [x] Migrar armazenamento local para S3 (LocalStack)
- [ ] Implementar criptografia básica com AWS KMS (LocalStack)

### 1.7 - Infraestrutura de Deploy na AWS (Futuro)
- [ ] Setup de Terraform para infraestrutura AWS
- [ ] Configurar AWS ECS/Fargate para containers
- [ ] Pipeline de deploy automatizado para produção
- [ ] Implementar health checks e monitoring na AWS
- [ ] Configurar AWS CloudWatch para logs e métricas
- [ ] Setup de domínio e SSL/TLS

### 2.3 - Observabilidade Avançada
- [ ] Implementar tracing distribuído com AWS X-Ray
- [ ] Adicionar métricas customizadas de negócio
- [ ] Implementar alertas proativos baseados em SLIs/SLOs
- [ ] Melhorar logs com contexto de requisição
- [ ] Adicionar dashboards básicos de monitoramento
- [ ] Implementar testes de segurança automatizados (SAST/DAST)

### 2.2 - API Design e Qualidade
- [x] Implementar versionamento de API (v1, v2, etc.)
- [ ] Adicionar middleware de CORS, rate limiting básico
- [ ] Implementar graceful degradation para falhas

## :card_file_box: Fase 3 – Persistência e rastreabilidade

### 3.1 - Persistência de Dados
- [ ] Criar schema da tabela `jobs` no DynamoDB
- [ ] Implementar repository pattern para jobs
- [ ] Persistir pedidos de processamento na base
- [ ] Implementar queries eficientes (evitando scan)
- [ ] Listar processamentos com paginação
- [ ] Associar usuário identificado ou anônimo a cada job

### 3.2 - Proteção de Dados e LGPD
- [ ] Implementar criptografia em trânsito e em repouso via AWS KMS
- [ ] Configurar isolamento de dados entre usuários via DynamoDB
- [ ] Implementar retenção de dados configurável por usuário
- [ ] Adicionar políticas de LGPD para dados brasileiros
- [ ] Implementar direito ao esquecimento (exclusão de dados)
- [ ] Criar auditoria de acesso a dados pessoais
- [ ] Configurar backup automático via AWS Backup

### 3.3 - Operações e Manutenção
- [ ] Implementar cleanup de jobs antigos
- [ ] Adicionar estratégias de backup/recovery
- [ ] Implementar teste automatizado de restore
- [ ] Configurar S3 Lifecycle para cleanup automático
- [ ] Migrar de LocalStack para DynamoDB real na AWS

---

---

## :warning: Tech Debts

### Funcionalidades Depreciadas
- **Filesystem Storage (Legacy)**: O sistema de armazenamento local (filesystem) está marcado para depreciação
  - **Razão**: Migração para S3 com presigned URLs oferece melhor escalabilidade e segurança
  - **Impacto**: Funcionalidade de fallback para desenvolvimento local ainda funcional
  - **Plano de Remoção**: Remover após implementação completa do LocalStack (Fase 4.1)
  - **Arquivos Afetados**: 
    - `api/internal/handlers/handlers.go` (métodos `*FromFilesystem`)
    - `processor/internal/services/video.go` (métodos `*ToFilesystem`)
    - Configurações de diretórios locais

### Melhorias Técnicas Pendentes
- **E2E Tests**: Atualizar testes para trabalhar com presigned URLs em vez de streaming direto
- **Frontend**: Remover lógica de fallback para URLs locais após depreciação do filesystem
- **Configuração**: Simplificar configuração removendo variáveis de ambiente relacionadas ao filesystem

---

## :warning: Notas Importantes

### Compliance com Tech Mandates
Este roadmap foi atualizado após revisão de compliance com nossos [Tech Mandates](./tech-mandates.md). As correções críticas de segurança na **Fase 1.2** são obrigatórias antes de continuar com desenvolvimentos posteriores.

---

:pushpin: Este roadmap será mantido e atualizado conforme o projeto evolui.  
Sugestões e contribuições são bem-vindas via issues ou pull requests.
