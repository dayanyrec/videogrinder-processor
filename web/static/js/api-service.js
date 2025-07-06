class ApiService {
  constructor() {
    this.baseURL = `${window.location.protocol}//${window.location.hostname}:8081`
    this.endpoints = {
      videos: `${this.baseURL}/api/v1/videos`
    }
  }

  async uploadVideo(formData) {
    const response = await fetch(this.endpoints.videos, {
      method: 'POST',
      body: formData
    })
    return await response.json()
  }

  async getFilesList() {
    const response = await fetch(this.endpoints.videos)
    return await response.json()
  }

  async deleteVideo(filename) {
    const response = await fetch(`${this.endpoints.videos}/${filename}`, {
      method: 'DELETE'
    })
    return response.ok
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

// Expose for testing
window.ApiService = ApiService
