# Production Deployment Guide

## AWS Configuration

Para executar o VideoGrinder em produção, você precisa configurar as credenciais AWS e criar os buckets S3 necessários.

### 1. Configuração das Variáveis de Ambiente

Copie o arquivo `env.example` para `.env` e configure suas credenciais AWS:

```bash
cp env.example .env
```

Edite o arquivo `.env` com suas credenciais:

```env
# AWS Configuration for Production
AWS_REGION=us-east-1
AWS_ACCESS_KEY_ID=your-actual-access-key-id
AWS_SECRET_ACCESS_KEY=your-actual-secret-access-key
AWS_ENDPOINT_URL=

# S3 Bucket Names
S3_UPLOADS_BUCKET=videogrinder-uploads
S3_OUTPUTS_BUCKET=videogrinder-outputs
```

### 2. Criação dos Buckets S3

Crie os buckets S3 necessários na sua conta AWS:

```bash
# Usando AWS CLI
aws s3 mb s3://videogrinder-uploads --region us-east-1
aws s3 mb s3://videogrinder-outputs --region us-east-1

# Configure as permissões adequadas
aws s3api put-bucket-policy --bucket videogrinder-uploads --policy file://bucket-policy.json
aws s3api put-bucket-policy --bucket videogrinder-outputs --policy file://bucket-policy.json
```

### 3. Execução em Produção

Para executar todos os serviços em modo produção:

```bash
# Carregue as variáveis de ambiente
source .env

# Execute os serviços
make run prod
```

Ou usando docker-compose diretamente:

```bash
docker-compose --profile prod up -d
```

### 4. Verificação do Status

Verifique se os serviços estão rodando:

```bash
make ps
```

Teste os endpoints de saúde:

```bash
curl http://localhost:8080/health  # Web service
curl http://localhost:8081/health  # API service
curl http://localhost:8082/health  # Processor service
```

### 5. Configuração de IAM

Certifique-se de que sua conta AWS tem as seguintes permissões:

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "s3:GetObject",
                "s3:PutObject",
                "s3:DeleteObject",
                "s3:ListBucket"
            ],
            "Resource": [
                "arn:aws:s3:::videogrinder-uploads",
                "arn:aws:s3:::videogrinder-uploads/*",
                "arn:aws:s3:::videogrinder-outputs",
                "arn:aws:s3:::videogrinder-outputs/*"
            ]
        }
    ]
}
```

### 6. Troubleshooting

**Erro "NoCredentialProviders"**: Verifique se as variáveis de ambiente AWS estão definidas corretamente.

**Erro "BucketNotFound"**: Certifique-se de que os buckets S3 existem na região especificada.

**Erro de permissão**: Verifique se as credenciais AWS têm as permissões necessárias para acessar os buckets S3.

## Monitoramento

Os serviços expõem endpoints de health check:

- Web: `http://localhost:8080/health`
- API: `http://localhost:8081/health` 
- Processor: `http://localhost:8082/health`

Estes endpoints podem ser usados para configurar health checks em sistemas de orquestração como Kubernetes ou load balancers. 
