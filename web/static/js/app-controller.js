/* global UIManager, ApiService, Utils */
class AppController {
  constructor() {
    this.uiManager = new UIManager()
    this.apiService = new ApiService()
  }

  async init() {
    await this.loadFilesList()
    this.setupEventListeners()
  }

  setupEventListeners() {
    const uploadForm = this.uiManager.getUploadForm()
    if (!uploadForm) {
      return
    }

    // Remove any existing listeners to prevent duplicates (if removeEventListener exists)
    if (uploadForm.removeEventListener) {
      uploadForm.removeEventListener('submit', this.handleUpload)
    }

    // Add the submit handler
    uploadForm.addEventListener('submit', (e) => {
      this.handleUpload(e)
    })
  }

  async handleUpload(e) {
    // Defensive preventDefault
    if (e && e.preventDefault) {
      e.preventDefault()
    }
    if (e && e.stopPropagation) {
      e.stopPropagation()
    }

    const file = this.uiManager.getSelectedFile()
    const validation = Utils.validateFile(file)

    if (!validation.valid) {
      this.uiManager.showResult(validation.message, 'error')
      return
    }

    try {
      this.uiManager.showLoading()

      const formData = this.apiService.createFormData(file)
      const data = await this.apiService.uploadVideo(formData)

      this.uiManager.hideLoading()

      if (data.success) {
        const downloadLink = data.download_url ?
          `<a href="${data.download_url}" class="download-btn">📥 Baixar ZIP</a>` :
          Utils.createDownloadLink(data.zip_path)
        const message = data.message + '<br><br>' + downloadLink
        this.uiManager.showResult(message, 'success')
        await this.loadFilesList()
      } else {
        this.uiManager.showResult('Erro: ' + data.message, 'error')
      }
    } catch (error) {
      this.uiManager.hideLoading()
      this.uiManager.showResult('Erro de conexão: ' + error.message, 'error')
    }
  }

  async loadFilesList() {
    try {
      const data = await this.apiService.getFilesList()
      this.uiManager.displayFilesList(data.videos)
    } catch (_error) {
      this.uiManager.displayFilesError()
    }
  }

  showResult(message, type) {
    this.uiManager.showResult(message, type)
  }
}

if (typeof module !== 'undefined' && module.exports) {
  module.exports = AppController
}

// Expose for testing
window.AppController = AppController
