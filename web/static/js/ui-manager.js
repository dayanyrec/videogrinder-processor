class UIManager {
  constructor() {
    this.elements = {
      loading: document.getElementById('loading'),
      result: document.getElementById('result'),
      filesList: document.getElementById('filesList'),
      uploadForm: document.getElementById('uploadForm'),
      videoFile: document.getElementById('videoFile')
    }
  }

  getApiBaseURL() {
    if (window.location.port === '8080') {
      return `${window.location.protocol}//${window.location.hostname}:8081`
    }
    return window.location.origin
  }

  showResult(message, type) {
    const resultDiv = this.elements.result
    resultDiv.innerHTML = message
    resultDiv.className = 'result ' + type
    resultDiv.style.display = 'block'
  }

  showLoading() {
    this.elements.loading.style.display = 'block'
    this.elements.result.style.display = 'none'
  }

  hideLoading() {
    this.elements.loading.style.display = 'none'
  }

  displayFilesList(files) {
    const filesListDiv = this.elements.filesList

    if (files && files.length > 0) {
      const apiBaseURL = this.getApiBaseURL()
      let html = ''
      files.forEach(file => {
        const downloadUrl = file.download_url.startsWith('/')
          ? `${apiBaseURL}${file.download_url}`
          : file.download_url

        html += '<div class="file-item">' +
                '<span><strong>' + file.filename + '</strong><br>' +
                '<small>Tamanho: ' + Math.round(file.size / 1024) + ' KB | ' +
                'Criado: ' + file.created_at + '</small></span>' +
                '<a href="' + downloadUrl + '" class="download-btn">📥 Baixar</a>' +
                '</div>'
      })
      filesListDiv.innerHTML = html
    } else {
      filesListDiv.innerHTML = '<p style="text-align: center; color: #999;">Nenhum arquivo processado ainda.</p>'
    }
  }

  displayFilesError() {
    this.elements.filesList.innerHTML = '<p style="color: red;">Erro ao carregar arquivos.</p>'
  }

  getSelectedFile() {
    return this.elements.videoFile.files[0]
  }

  getUploadForm() {
    return this.elements.uploadForm
  }
}

if (typeof module !== 'undefined' && module.exports) {
  module.exports = UIManager
}

// Expose for testing
window.UIManager = UIManager
