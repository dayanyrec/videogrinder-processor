# 🎬 VideoGrinder POC

## Sobre o Projeto

O **VideoGrinder POC** é uma prova de conceito desenvolvida especificamente para **jornalistas** que precisam extrair frames de vídeos para criar conteúdo jornalístico, posts em redes sociais ou materiais de apoio para suas matérias.

Esta ferramenta permite que os usuários façam upload de vídeos através de uma interface web simples e recebam um arquivo ZIP contendo todos os frames extraídos em formato PNG, facilitando o processo de seleção de imagens para uso editorial.

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

## 📋 Pré-requisitos

- Docker instalado
- Git (para clonagem do repositório)

## 🚀 Como Executar

1. **Clone o repositório:**
```bash
git clone <url-do-repositorio>
cd videogrinder-poc
```

2. **Build e execute o container:**
```bash
docker build -t videogrinder .
docker run -p 8080:8080 videogrinder
```

3. **Acesse no navegador:**
```
http://localhost:8080
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
videogrinder-poc/
├── main.go           # Aplicação principal
├── go.mod           # Dependências do Go
├── go.sum           # Checksums das dependências
├── Dockerfile       # Configuração do Docker
├── uploads/         # Vídeos enviados (temporário)
├── outputs/         # Arquivos ZIP gerados
├── temp/            # Arquivos temporários durante processamento
└── README.md        # Este arquivo
```

## 🔧 Configuração

### Portas
- **Porta padrão**: 8080
- Para alterar a porta, modifique a linha `r.Run(":8080")` no arquivo `main.go`

### Extração de Frames
- **Taxa padrão**: 1 frame por segundo (`fps=1`)
- Para alterar, modifique o parâmetro `-vf "fps=1"` na função `processVideo()`

### Formatos Suportados
Os formatos de vídeo são validados na função `isValidVideoFile()`. Para adicionar novos formatos, edite o array `validExts`.

## 🐛 Solução de Problemas

### Erro "FFmpeg não encontrado"
- **Linux/Mac**: `brew install ffmpeg` ou `apt-get install ffmpeg`
- **Windows**: Baixe o FFmpeg e adicione ao PATH do sistema

### Erro de permissão em diretórios
```bash
sudo chmod 755 uploads outputs temp
```

### Vídeo não é processado
- Verifique se o formato é suportado
- Confirme se o arquivo não está corrompido
- Verifique os logs do terminal para erros específicos

### Porta 8080 em uso
- Altere a porta no código ou termine o processo que está usando a porta:
```bash
lsof -ti:8080 | xargs kill -9
```

## 🎯 Casos de Uso para Jornalistas

- **Matérias esportivas**: Extrair momentos-chave de jogos
- **Eventos políticos**: Capturar gestos e expressões importantes
- **Coberturas ao vivo**: Gerar imagens para posts em tempo real
- **Análise de conteúdo**: Estudar sequências de vídeo frame por frame
- **Redes sociais**: Criar carrosséis de imagens para Instagram/Twitter
- **Documentação**: Arquivo visual de eventos importantes

## ⚠️ Limitações Conhecidas

- O processamento é sequencial (um vídeo por vez)
- Arquivos muito grandes podem consumir bastante espaço em disco
- O tempo de processamento é proporcional ao tamanho e duração do vídeo
- Interface web básica (POC)

## 🔮 Melhorias Futuras

- [ ] Processamento paralelo de múltiplos vídeos
- [ ] Seleção customizada de taxa de frames (fps)
- [ ] Preview dos frames antes do download
- [ ] Autenticação de usuários
- [ ] API REST completa
- [ ] Interface mais sofisticada
- [ ] Suporte a mais formatos de saída
- [ ] Compressão inteligente de imagens
- [ ] Integração com serviços de nuvem

## 📄 Licença

Este é um projeto de prova de conceito (POC) desenvolvido para fins educacionais e demonstrativos.

## 🤝 Contribuição

Como este é um POC, contribuições são bem-vindas! Sinta-se à vontade para:

- Reportar bugs
- Sugerir melhorias
- Enviar pull requests
- Compartilhar casos de uso

## 📞 Suporte

Para dúvidas ou problemas:
1. Verifique a seção "Solução de Problemas" 
2. Consulte os logs da aplicação
3. Abra uma issue no repositório
