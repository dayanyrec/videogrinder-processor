# :motorway: VideoGrinder - Roadmap de Evolução

Este roadmap descreve os passos planejados para amadurecer o projeto VideoGrinder, evoluindo de um monólito para uma arquitetura modular e escalável.

---

## :white_check_mark: Fase 0 – Fortalecer fundação do projeto

- [x] Criar repositório para POC para manter histórico de fácil acesso
- [x] Documentar propósito e escopo da POC (`README.md`)
- [x] Criar repositório para o projeto onde será feita a evolução
- [x] Documentar roadmap com ideias atuais

---

## :rocket: Fase 1 – Tornar a POC um projeto produtivo

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

### 1.3 - Observabilidade e Monitoramento
- [ ] Implementar logging estruturado em JSON
- [ ] Melhorar tratamento de erros com contexto adequado
- [ ] Adicionar health check endpoints robustos
- [ ] Integrar CloudWatch Logs para desenvolvimento
- [ ] Implementar métricas básicas de performance
- [ ] Adicionar monitoramento de recursos (CPU, memória, I/O)

### 1.4 - Integração AWS Básica
- [ ] Integrar AWS Secrets Manager para configurações sensíveis
- [ ] Configurar AWS IAM roles e políticas básicas
- [ ] Setup do LocalStack com S3 e DynamoDB (desenvolvimento local)
- [ ] Migrar armazenamento local para S3 (LocalStack)
- [ ] Implementar criptografia básica com AWS KMS (LocalStack)

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

### 1.7 - Infraestrutura de Deploy na AWS (Futuro)
- [ ] Setup de Terraform para infraestrutura AWS
- [ ] Configurar AWS ECS/Fargate para containers
- [ ] Pipeline de deploy automatizado para produção
- [ ] Implementar health checks e monitoring na AWS
- [ ] Configurar AWS CloudWatch para logs e métricas
- [ ] Setup de domínio e SSL/TLS

---

## :jigsaw: Fase 2 – Modularização (ainda no monolito)

### 2.1 - Estruturação do Código
- [x] Extrair configuração para internal/config (variáveis de ambiente)
- [ ] Criar estrutura de pacotes Go (cmd/, internal/, pkg/)
- [ ] Extrair frontend para arquivos estáticos
- [ ] Cobrir frontend (JS) com testes unitários
- [ ] Extrair core de processamento para módulo interno
- [ ] Cobrir o core com testes unitários
- [ ] Cobrir o core com testes de integração

### 2.2 - API Design e Qualidade
- [ ] Extrair API para módulo interno
- [ ] Melhorar implementação da API tornando REST
- [ ] Implementar versionamento de API (v1, v2, etc.)
- [ ] Adicionar middleware de CORS, rate limiting básico
- [ ] Implementar graceful degradation para falhas
- [ ] Documentar API com Swagger/OpenAPI automático
- [ ] Cobrir API com testes unitários
- [ ] Cobrir API com testes de integração

### 2.3 - Observabilidade Avançada
- [ ] Implementar tracing distribuído com AWS X-Ray
- [ ] Adicionar métricas customizadas de negócio
- [ ] Implementar alertas proativos baseados em SLIs/SLOs
- [ ] Melhorar logs com contexto de requisição
- [ ] Adicionar dashboards básicos de monitoramento
- [ ] Implementar testes de segurança automatizados (SAST/DAST)

---

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

## :gear: Fase 4 – Arquitetura de microserviços

### 4.1 - Preparação dos repositórios

- [ ] Criar repositório `videogrinder-frontend`
- [ ] Criar repositório `videogrinder-api` 
- [ ] Criar repositório `videogrinder-handler`
- [ ] Manter `videogrinder-processor` como repositório atual
- [ ] Definir contratos de API entre serviços

### 4.2 - Extração gradual (uma por vez)

- [ ] Extrair e deployar frontend em S3
- [ ] Implementar CI/CD para frontend
- [ ] Testar integração frontend -> API monolítica
- [ ] Extrair API para repositório próprio
- [ ] Implementar CI/CD para API em EKS
- [ ] Testar integração frontend -> API extraída

### 4.3 - Handler e orquestração

- [ ] Criar Lambda handler para orquestração  
- [ ] Implementar CI/CD para handler
- [ ] Ativar stream do DynamoDB -> handler
- [ ] Setup do processor como Job K8S
- [ ] Trigger processor via handler

### 4.4 - Integração e otimização

- [ ] Implementar processamento paralelo de múltiplos vídeos
- [ ] Testes de integração end-to-end
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
- [ ] Adicionar variáveis de ambiente para configuração
- [ ] Implementar timeouts para operações FFmpeg
- [ ] Adicionar rate limiting básico

---

## :warning: Notas Importantes

### Compliance com Tech Mandates
Este roadmap foi atualizado após revisão de compliance com nossos [Tech Mandates](./tech-mandates.md). As correções críticas de segurança na **Fase 1.2** são obrigatórias antes de continuar com desenvolvimentos posteriores.

---

:pushpin: Este roadmap será mantido e atualizado conforme o projeto evolui.  
Sugestões e contribuições são bem-vindas via issues ou pull requests.
