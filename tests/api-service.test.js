const ApiService = require('../static/js/api-service.js')

global.fetch = jest.fn()

// Mock window.location
delete global.window.location
global.window = Object.create(window)
global.window.location = {
  port: '8081',
  protocol: 'http:',
  hostname: 'localhost',
  origin: 'http://localhost:8081'
}

describe('ApiService Class', () => {
  let apiService

  beforeEach(() => {
    apiService = new ApiService()
    fetch.mockClear()
  })

  describe('constructor', () => {
    test('should initialize with correct API base URL for production environment', () => {
      // Default mock is production environment (port 8081)
      expect(apiService.baseURL).toBe('http://localhost:8081')
      expect(apiService.endpoints.videos).toBe('http://localhost:8081/api/v1/videos')
    })

    test('should initialize with API base URL for development environment', () => {
      // Mock development environment (port 8080)
      global.window.location.port = '8080'
      global.window.location.origin = 'http://localhost:8080'

      const devApiService = new ApiService()
      expect(devApiService.baseURL).toBe('http://localhost:8081')
      expect(devApiService.endpoints.videos).toBe('http://localhost:8081/api/v1/videos')

      // Reset to production for other tests
      global.window.location.port = '8081'
      global.window.location.origin = 'http://localhost:8081'
    })
  })

  describe('uploadVideo', () => {
    test('should successfully upload video and return processing result', async() => {
      const mockResponse = {
        success: true,
        message: 'Video processed successfully!',
        zip_path: 'frames_123.zip'
      }

      fetch.mockResolvedValueOnce({
        json: () => Promise.resolve(mockResponse)
      })

      const mockFormData = new FormData()
      const result = await apiService.uploadVideo(mockFormData)

      expect(fetch).toHaveBeenCalledWith('http://localhost:8081/api/v1/videos', {
        method: 'POST',
        body: mockFormData
      })
      expect(result).toEqual(mockResponse)
    })

    test('should return failure response when upload fails on server', async() => {
      const mockResponse = {
        success: false,
        message: 'Upload failed!'
      }

      fetch.mockResolvedValueOnce({
        json: () => Promise.resolve(mockResponse)
      })

      const mockFormData = new FormData()
      const result = await apiService.uploadVideo(mockFormData)

      expect(result).toEqual(mockResponse)
      expect(result.success).toBe(false)
    })

    test('should throw error when network request fails', async() => {
      fetch.mockRejectedValueOnce(new Error('Network error'))

      const mockFormData = new FormData()

      await expect(apiService.uploadVideo(mockFormData)).rejects.toThrow('Network error')
    })
  })

  describe('getFilesList', () => {
    test('should fetch and return formatted videos list with metadata', async() => {
      const mockVideos = [
        {
          filename: 'test-video.mp4',
          size: 1024000,
          created_at: '2024-01-01 10:00:00',
          download_url: '/api/v1/videos/test.zip/download'
        },
        {
          filename: 'another-video.avi',
          size: 2048000,
          created_at: '2024-01-02 11:00:00',
          download_url: '/api/v1/videos/another.zip/download'
        }
      ]

      fetch.mockResolvedValueOnce({
        json: () => Promise.resolve({ videos: mockVideos, total: 2 })
      })

      const result = await apiService.getFilesList()

      expect(fetch).toHaveBeenCalledWith('http://localhost:8081/api/v1/videos')
      expect(result.videos).toEqual(mockVideos)
      expect(result.videos).toHaveLength(2)
      expect(result.total).toBe(2)
    })

    test('should return empty list when no processed videos exist', async() => {
      fetch.mockResolvedValueOnce({
        json: () => Promise.resolve({ videos: [], total: 0 })
      })

      const result = await apiService.getFilesList()

      expect(result.videos).toEqual([])
      expect(result.videos).toHaveLength(0)
      expect(result.total).toBe(0)
    })

    test('should handle null videos list from server response', async() => {
      fetch.mockResolvedValueOnce({
        json: () => Promise.resolve({ videos: null, total: 0 })
      })

      const result = await apiService.getFilesList()

      expect(result.videos).toBeNull()
      expect(result.total).toBe(0)
    })

    test('should throw error when API request fails', async() => {
      fetch.mockRejectedValueOnce(new Error('API error'))

      await expect(apiService.getFilesList()).rejects.toThrow('API error')
    })
  })

  describe('deleteVideo', () => {
    test('should successfully delete video and return true when file exists', async() => {
      fetch.mockResolvedValueOnce({
        ok: true
      })

      const result = await apiService.deleteVideo('test.zip')

      expect(fetch).toHaveBeenCalledWith('http://localhost:8081/api/v1/videos/test.zip', {
        method: 'DELETE'
      })
      expect(result).toBe(true)
    })

    test('should return false when deletion fails on server', async() => {
      fetch.mockResolvedValueOnce({
        ok: false
      })

      const result = await apiService.deleteVideo('test.zip')

      expect(result).toBe(false)
    })

    test('should throw error when network request fails during deletion', async() => {
      fetch.mockRejectedValueOnce(new Error('Network error'))

      await expect(apiService.deleteVideo('test.zip')).rejects.toThrow('Network error')
    })
  })

  describe('createFormData', () => {
    test('should create FormData object with video file for mp4 format', () => {
      const mockFile = new File(['video content'], 'test.mp4', { type: 'video/mp4' })

      const formData = apiService.createFormData(mockFile)

      expect(formData).toBeInstanceOf(FormData)
      expect(formData.get('video')).toBe(mockFile)
    })

    test('should create FormData object with video file for avi format', () => {
      const mockFile = new File(['video content'], 'test.avi', { type: 'video/avi' })

      const formData = apiService.createFormData(mockFile)

      expect(formData).toBeInstanceOf(FormData)
      expect(formData.get('video')).toBe(mockFile)
    })

    test('should create FormData object with video file when no MIME type is specified', () => {
      const mockFile = new File(['video content'], 'test.mov')

      const formData = apiService.createFormData(mockFile)

      expect(formData).toBeInstanceOf(FormData)
      expect(formData.get('video')).toBe(mockFile)
    })
  })
})
