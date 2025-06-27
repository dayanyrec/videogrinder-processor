# 🏛️ Tech Mandates

> **Diretrizes técnicas fundamentais do VideoGrinder**  
> *Estas diretrizes são obrigatórias e devem ser seguidas em todas as fases do projeto*

---

## ☁️ **Cloud Provider**

**AWS é nosso cloud provider exclusivo**

- Todos os serviços de infraestrutura devem utilizar AWS
- Configurações e deployments devem ser otimizados para o ecossistema AWS
- Integração nativa com serviços como ECR, S3, DynamoDB, Lambda, etc.

---

## 🏭 **Ambientes**

**Temos somente 1 ambiente produtivo (prod)**

- Não existem ambientes de staging ou homologação
- Desenvolvimento local → Produção
- Foco em CI/CD robusto e testes abrangentes
- Deploy direto para produção com confiança

---

## 🐳 **Desenvolvimento Local**

**Desenvolvimento local deve ser feito somente usando Docker**

- Zero dependências locais além do Docker
- Ambiente de desenvolvimento idêntico à produção
- Todos os comandos executados via containers
- Hot reload e ferramentas de desenvolvimento containerizadas

---

## 🚫 **Comentários no Código**

**Comentários no código são code smells; o código deve ser autoexplicativo**

- Nomes de variáveis, funções e classes descritivos
- Estrutura de código clara e intuitiva
- Refatoração contínua para manter legibilidade
- Documentação separada do código quando necessária

---

## 📋 **Documentação via Testes**

**Testes também são documentação**

- Testes devem servir como especificação viva
- Cenários de teste claros e descritivos
- Cobertura de teste como indicador de qualidade
- Testes como fonte de verdade do comportamento esperado

---

## 🔒 **Security First**

**Segurança é prioridade desde o design**

- Princípio de menor privilégio em todas as configurações
- Secrets gerenciados exclusivamente via AWS Secrets Manager
- AWS IAM para controle de acesso granular
- Validação rigorosa de inputs (especialmente arquivos de mídia)
- HTTPS obrigatório para todas as comunicações
- Containers executados como usuário não-root

---

## ⚡ **Performance & Efficiency**

**Otimização de recursos é obrigatória**

- Processamento de vídeo deve ser otimizado para custo/performance
- Monitoramento via CloudWatch de CPU, memória e I/O
- Implementação de timeouts apropriados
- Cleanup automático de arquivos temporários via S3 Lifecycle
- Dimensionamento baseado em métricas reais do CloudWatch

---

## 📊 **Observabilidade**

**Visibilidade completa do sistema**

- Logs estruturados em JSON via CloudWatch Logs
- Métricas coletadas via CloudWatch Metrics
- Health checks obrigatórios em todos os serviços
- Tracing distribuído utilizando AWS X-Ray
- Alertas proativos via CloudWatch Alarms baseados em SLIs/SLOs

---

## 🎯 **API Design**

**Consistência e simplicidade**

- RESTful APIs com convenções consistentes
- Versionamento de APIs obrigatório
- Rate limiting via AWS API Gateway
- Documentação OpenAPI automática
- Graceful degradation para falhas

---

## 🔧 **Infrastructure as Code**

**Infraestrutura versionada e reproduzível**

- Toda infraestrutura definida como código utilizando Terraform
- Nenhuma configuração manual em produção
- Rollback automatizado em caso de falha
- Ambientes efêmeros para testing
- GitOps para deployment de infraestrutura

---

## 💾 **Data Protection**

**Proteção de conteúdo dos jornalistas**

- Criptografia em trânsito e em repouso via AWS KMS
- Armazenamento de arquivos em S3 com versionamento
- Backup automático via AWS Backup com teste de restore
- Compliance com LGPD para dados brasileiros
- Isolamento de dados entre usuários via DynamoDB
