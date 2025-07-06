class Utils {
  static formatFileSize(bytes) {
    return Math.round(bytes / 1024)
  }

  static getApiBaseURL() {
    if (window.location.port === '8080') {
      return `${window.location.protocol}//${window.location.hostname}:8081`
    }
    return window.location.origin
  }

  static createDownloadLink(zipPath) {
    const apiBaseURL = this.getApiBaseURL()
    // Use redirect mode for direct download experience
    return `<a href="${apiBaseURL}/api/v1/videos/${zipPath}/download?redirect=true" class="download-btn">📥 Baixar ZIP</a>`
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

// Expose for testing
window.Utils = Utils
