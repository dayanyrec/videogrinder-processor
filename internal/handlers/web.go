package handlers

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"video-processor/internal/config"
	"video-processor/internal/models"
	"video-processor/internal/services"

	"github.com/gin-gonic/gin"
)

type WebHandlers struct {
	videoService *services.VideoService
	config       *config.Config
}

func NewWebHandlers(videoService *services.VideoService, cfg *config.Config) *WebHandlers {
	return &WebHandlers{
		videoService: videoService,
		config:       cfg,
	}
}

func (wh *WebHandlers) HandleVideoUpload(c *gin.Context) {
	file, header, err := c.Request.FormFile("video")
	if err != nil {
		c.JSON(400, models.ProcessingResult{
			Success: false,
			Message: "Erro ao receber arquivo: " + err.Error(),
		})
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("Warning: Failed to close uploaded file: %v", err)
		}
	}()

	if !IsValidVideoFile(header.Filename) {
		c.JSON(400, models.ProcessingResult{
			Success: false,
			Message: "Formato de arquivo não suportado. Use: mp4, avi, mov, mkv",
		})
		return
	}

	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("%s_%s", timestamp, filepath.Base(header.Filename))
	videoPath := filepath.Join(wh.config.UploadsDir, filename)

	cleanVideoPath := filepath.Clean(videoPath)
	uploadsDir, _ := filepath.Abs(wh.config.UploadsDir)
	absVideoPath, _ := filepath.Abs(cleanVideoPath)
	if !strings.HasPrefix(absVideoPath, uploadsDir+string(filepath.Separator)) {
		c.JSON(400, models.ProcessingResult{
			Success: false,
			Message: "Invalid file path",
		})
		return
	}

	out, err := os.Create(filepath.Clean(videoPath))
	if err != nil {
		c.JSON(500, models.ProcessingResult{
			Success: false,
			Message: "Erro ao salvar arquivo: " + err.Error(),
		})
		return
	}
	defer func() {
		if err := out.Close(); err != nil {
			log.Printf("Warning: Failed to close output file: %v", err)
		}
	}()

	_, err = io.Copy(out, file)
	if err != nil {
		c.JSON(500, models.ProcessingResult{
			Success: false,
			Message: "Erro ao salvar arquivo: " + err.Error(),
		})
		return
	}

	result := wh.videoService.ProcessVideo(videoPath, timestamp)

	if result.Success {
		if err := os.Remove(videoPath); err != nil {
			log.Printf("Warning: Failed to remove video file %s: %v", videoPath, err)
		}
	}

	c.JSON(200, result)
}

func (wh *WebHandlers) HandleDownload(c *gin.Context) {
	filename := c.Param("filename")
	filePath := filepath.Join(wh.config.OutputsDir, filename)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(404, gin.H{"error": "Arquivo não encontrado"})
		return
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/zip")

	c.File(filePath)
}

func (wh *WebHandlers) HandleStatus(c *gin.Context) {
	files, err := filepath.Glob(filepath.Join(wh.config.OutputsDir, "*.zip"))
	if err != nil {
		c.JSON(500, gin.H{"error": "Erro ao listar arquivos"})
		return
	}

	results := make([]map[string]interface{}, 0, len(files))
	for _, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			continue
		}

		results = append(results, map[string]interface{}{
			"filename":     filepath.Base(file),
			"size":         info.Size(),
			"created_at":   info.ModTime().Format("2006-01-02 15:04:05"),
			"download_url": "/download/" + filepath.Base(file),
		})
	}

	c.JSON(200, gin.H{
		"files": results,
		"total": len(results),
	})
}

func (wh *WebHandlers) HandleHome(c *gin.Context) {
	c.Header("Content-Type", "text/html")
	c.String(200, GetHTMLForm())
}

func IsValidVideoFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	validExts := []string{".mp4", ".avi", ".mov", ".mkv", ".wmv", ".flv", ".webm"}

	for _, validExt := range validExts {
		if ext == validExt {
			return true
		}
	}
	return false
}

func GetHTMLForm() string {
	return `
<!DOCTYPE html>
<html lang="pt-BR">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>FIAP X - Processador de Vídeos</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            max-width: 800px;
            margin: 50px auto;
            padding: 20px;
            background-color: #f5f5f5;
        }
        .container {
            background: white;
            padding: 30px;
            border-radius: 10px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        h1 {
            color: #333;
            text-align: center;
            margin-bottom: 30px;
        }
        .upload-form {
            border: 2px dashed #ddd;
            padding: 30px;
            text-align: center;
            border-radius: 10px;
            margin: 20px 0;
        }
        input[type="file"] {
            margin: 20px 0;
            padding: 10px;
        }
        button {
            background: #007bff;
            color: white;
            padding: 12px 30px;
            border: none;
            border-radius: 5px;
            cursor: pointer;
            font-size: 16px;
        }
        button:hover { background: #0056b3; }
        .result {
            margin-top: 20px;
            padding: 15px;
            border-radius: 5px;
            display: none;
        }
        .success { background: #d4edda; color: #155724; border: 1px solid #c3e6cb; }
        .error { background: #f8d7da; color: #721c24; border: 1px solid #f5c6cb; }
        .loading {
            text-align: center;
            display: none;
            margin: 20px 0;
        }
        .files-list {
            margin-top: 30px;
        }
        .file-item {
            background: #f8f9fa;
            padding: 10px;
            margin: 5px 0;
            border-radius: 5px;
            display: flex;
            justify-content: space-between;
            align-items: center;
        }
        .download-btn {
            background: #28a745;
            color: white;
            padding: 5px 15px;
            text-decoration: none;
            border-radius: 3px;
            font-size: 14px;
        }
        .download-btn:hover { background: #218838; }
    </style>
</head>
<body>
    <div class="container">
        <h1>🎬 FIAP X - Processador de Vídeos</h1>
        <p style="text-align: center; color: #666;">
            Faça upload de um vídeo e receba um ZIP com todos os frames extraídos!
        </p>

        <div class="upload-form">
            <h3>📤 Upload do Vídeo</h3>
            <form id="uploadForm" enctype="multipart/form-data">
                <input type="file" id="videoFile" name="video" accept=".mp4,.avi,.mov,.mkv,.wmv,.flv,.webm" required>
                <br>
                <button type="submit">Processar Vídeo</button>
            </form>
        </div>

        <div class="loading" id="loading">
            <p>🎬 Processando vídeo... isso pode levar alguns minutos.</p>
        </div>

        <div class="result" id="result"></div>

        <div class="files-list">
            <h3>📁 Arquivos Disponíveis</h3>
            <div id="filesList">Carregando...</div>
        </div>
    </div>

    <script>
        document.addEventListener('DOMContentLoaded', function() {
            loadFilesList();

            document.getElementById('uploadForm').addEventListener('submit', function(e) {
                e.preventDefault();

                const formData = new FormData();
                const fileInput = document.getElementById('videoFile');
                const file = fileInput.files[0];

                if (!file) {
                    showResult('Por favor, selecione um arquivo.', 'error');
                    return;
                }

                formData.append('video', file);

                document.getElementById('loading').style.display = 'block';
                document.getElementById('result').style.display = 'none';

                fetch('/upload', {
                    method: 'POST',
                    body: formData
                })
                .then(response => response.json())
                .then(data => {
                    document.getElementById('loading').style.display = 'none';

                    if (data.success) {
                        showResult(
                            data.message + '<br><br>' +
                            '<a href="/download/' + data.zip_path + '" class="download-btn">📥 Baixar ZIP</a>',
                            'success'
                        );
                        loadFilesList();
                    } else {
                        showResult('Erro: ' + data.message, 'error');
                    }
                })
                .catch(error => {
                    document.getElementById('loading').style.display = 'none';
                    showResult('Erro de conexão: ' + error.message, 'error');
                });
            });
        });

        function showResult(message, type) {
            const resultDiv = document.getElementById('result');
            resultDiv.innerHTML = message;
            resultDiv.className = 'result ' + type;
            resultDiv.style.display = 'block';
        }

        function loadFilesList() {
            fetch('/api/status')
                .then(response => response.json())
                .then(data => {
                    const filesListDiv = document.getElementById('filesList');

                    if (data.files && data.files.length > 0) {
                        let html = '';
                        data.files.forEach(file => {
                            html += '<div class="file-item">' +
                                   '<span><strong>' + file.filename + '</strong><br>' +
                                   '<small>Tamanho: ' + Math.round(file.size / 1024) + ' KB | ' +
                                   'Criado: ' + file.created_at + '</small></span>' +
                                   '<a href="' + file.download_url + '" class="download-btn">📥 Baixar</a>' +
                                   '</div>';
                        });
                        filesListDiv.innerHTML = html;
                    } else {
                        filesListDiv.innerHTML = '<p style="text-align: center; color: #999;">Nenhum arquivo processado ainda.</p>';
                    }
                })
                .catch(error => {
                    document.getElementById('filesList').innerHTML = '<p style="color: red;">Erro ao carregar arquivos.</p>';
                });
        }
    </script>
</body>
</html>`
}
