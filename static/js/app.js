document.addEventListener('DOMContentLoaded', () => {
  loadFilesList()

  document.getElementById('uploadForm').addEventListener('submit', (e) => {
    e.preventDefault()

    const formData = new FormData()
    const fileInput = document.getElementById('videoFile')
    const file = fileInput.files[0]

    if (!file) {
      showResult('Por favor, selecione um arquivo.', 'error')
      return
    }

    formData.append('video', file)

    document.getElementById('loading').style.display = 'block'
    document.getElementById('result').style.display = 'none'

    fetch('/upload', {
      method: 'POST',
      body: formData
    })
      .then(response => response.json())
      .then(data => {
        document.getElementById('loading').style.display = 'none'

        if (data.success) {
          showResult(
            data.message + '<br><br>' +
                    '<a href="/download/' + data.zip_path + '" class="download-btn">📥 Baixar ZIP</a>',
            'success'
          )
          loadFilesList()
        } else {
          showResult('Erro: ' + data.message, 'error')
        }
      })
      .catch(error => {
        document.getElementById('loading').style.display = 'none'
        showResult('Erro de conexão: ' + error.message, 'error')
      })
  })
})

function showResult(message, type) {
  const resultDiv = document.getElementById('result')
  resultDiv.innerHTML = message
  resultDiv.className = 'result ' + type
  resultDiv.style.display = 'block'
}

function loadFilesList() {
  fetch('/api/status')
    .then(response => response.json())
    .then(data => {
      const filesListDiv = document.getElementById('filesList')

      if (data.files && data.files.length > 0) {
        let html = ''
        data.files.forEach(file => {
          html += '<div class="file-item">' +
                           '<span><strong>' + file.filename + '</strong><br>' +
                           '<small>Tamanho: ' + Math.round(file.size / 1024) + ' KB | ' +
                           'Criado: ' + file.created_at + '</small></span>' +
                           '<a href="' + file.download_url + '" class="download-btn">📥 Baixar</a>' +
                           '</div>'
        })
        filesListDiv.innerHTML = html
      } else {
        filesListDiv.innerHTML = '<p style="text-align: center; color: #999;">Nenhum arquivo processado ainda.</p>'
      }
    })
    .catch(_error => {
      document.getElementById('filesList').innerHTML = '<p style="color: red;">Erro ao carregar arquivos.</p>'
    })
}
