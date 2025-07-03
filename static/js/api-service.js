class ApiService {
  constructor() {
    this.endpoints = {
      upload: '/upload',
      status: '/api/status'
    }
  }

  async uploadVideo(formData) {
    const response = await fetch(this.endpoints.upload, {
      method: 'POST',
      body: formData
    })
    return await response.json()
  }

  async getFilesList() {
    const response = await fetch(this.endpoints.status)
    return await response.json()
  }

  createFormData(file) {
    const formData = new FormData()
    formData.append('video', file)
    return formData
  }
}

if (typeof module !== 'undefined' && module.exports) {
  module.exports = ApiService
}
