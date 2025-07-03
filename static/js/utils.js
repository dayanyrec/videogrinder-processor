class Utils {
  static formatFileSize(bytes) {
    return Math.round(bytes / 1024)
  }

  static createDownloadLink(zipPath) {
    return '<a href="/download/' + zipPath + '" class="download-btn">📥 Baixar ZIP</a>'
  }

  static validateFile(file) {
    if (!file) {
      return {
        valid: false,
        message: 'Por favor, selecione um arquivo.'
      }
    }
    return {
      valid: true,
      message: ''
    }
  }
}

if (typeof module !== 'undefined' && module.exports) {
  module.exports = Utils
}
