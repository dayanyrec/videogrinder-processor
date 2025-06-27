
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

- [ ] Setup de `.editorconfig`, linters e boas práticas
- [ ] Melhorar containerização com Docker multistage
- [ ] Adicionar variáveis de ambiente para configuração
- [ ] Melhorar tratamento de erros e logging básico
- [ ] Adicionar health check endpoints
- [ ] Cobrir código atual com testes unitários básicos
- [ ] Setup do LocalStack com S3 e DynamoDB (desenvolvimento local)
- [ ] Persistir arquivos em buckets locais (LocalStack)
- [ ] Cobrir experiência com testes end-to-end (Cypress)
- [ ] Configurar pipeline com GitHub Actions (trunk-based)
- [ ] Setup de Kubernetes local (Minikube ou K3s)
- [ ] Setup de Terraform para infraestrutura base
- [ ] Pipeline de deploy automatizado na AWS

---

## :jigsaw: Fase 2 – Modularização (ainda no monolito)

- [ ] Criar estrutura de pacotes Go (cmd/, internal/, pkg/)
- [ ] Extrair frontend para arquivos estáticos
- [ ] Cobrir frontend (JS) com testes unitários
- [ ] Extrair API para módulo interno
- [ ] Melhorar implementação da API tornando REST
- [ ] Adicionar middleware de CORS, rate limiting básico
- [ ] Documentar API com Swagger/OpenAPI
- [ ] Cobrir API com testes unitários
- [ ] Cobrir API com testes de integração
- [ ] Extrair core de processamento para módulo interno
- [ ] Cobrir o core com testes unitários
- [ ] Cobrir o core com testes de integração

---

## :card_file_box: Fase 3 – Persistência e rastreabilidade

- [ ] Criar schema da tabela `jobs` no DynamoDB
- [ ] Implementar repository pattern para jobs
- [ ] Persistir pedidos de processamento na base
- [ ] Implementar queries eficientes (evitando scan)
- [ ] Listar processamentos com paginação
- [ ] Associar usuário identificado ou anônimo a cada job
- [ ] Implementar cleanup de jobs antigos
- [ ] Adicionar backup/recovery strategy
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

- [ ] Adicionar SQS entre Lambda e Job (desacoplamento)
- [ ] Logs estruturados e centralizados (ex: CloudWatch ou Grafana Loki)
- [ ] Observabilidade com OpenTelemetry (traços e métricas)
- [ ] Adicionar autenticação e rate limiting (se necessário)
- [ ] Melhorar UX com preview e status visual em tempo real
- [ ] Preview dos frames antes do download
- [ ] Interface mais sofisticada e responsiva
- [ ] Sistema de notificações em tempo real
- [ ] Seleção customizada de taxa de frames (fps)
- [ ] Suporte a mais formatos de saída (JPEG, WebP, GIF)
- [ ] Implementar compressão inteligente de imagens
- [ ] Permitir busca pela listagem de vídeos
- [ ] Permitir inserir tags que agrupem vídeos e facilitem a busca

---

:pushpin: Este roadmap será mantido e atualizado conforme o projeto evolui.  
Sugestões e contribuições são bem-vindas via issues ou pull requests.
