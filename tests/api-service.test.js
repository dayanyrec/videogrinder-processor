const ApiService = require('../static/js/api-service.js')

global.fetch = jest.fn()

describe('ApiService Class', () => {
  let apiService

  beforeEach(() => {
    apiService = new ApiService()
    fetch.mockClear()
  })

  describe('constructor', () => {
    test('should initialize with correct endpoints', () => {
      expect(apiService.endpoints.upload).toBe('/upload')
      expect(apiService.endpoints.status).toBe('/api/status')
    })
  })

  describe('uploadVideo', () => {
    test('should upload video successfully', async() => {
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

      expect(fetch).toHaveBeenCalledWith('/upload', {
        method: 'POST',
        body: mockFormData
      })
      expect(result).toEqual(mockResponse)
    })

    test('should handle upload failure', async() => {
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

    test('should handle network errors', async() => {
      fetch.mockRejectedValueOnce(new Error('Network error'))

      const mockFormData = new FormData()

      await expect(apiService.uploadVideo(mockFormData)).rejects.toThrow('Network error')
    })
  })

  describe('getFilesList', () => {
    test('should fetch files list successfully', async() => {
      const mockFiles = [
        {
          filename: 'test-video.mp4',
          size: 1024000,
          created_at: '2024-01-01 10:00:00',
          download_url: '/download/test.zip'
        },
        {
          filename: 'another-video.avi',
          size: 2048000,
          created_at: '2024-01-02 11:00:00',
          download_url: '/download/another.zip'
        }
      ]

      fetch.mockResolvedValueOnce({
        json: () => Promise.resolve({ files: mockFiles })
      })

      const result = await apiService.getFilesList()

      expect(fetch).toHaveBeenCalledWith('/api/status')
      expect(result.files).toEqual(mockFiles)
      expect(result.files).toHaveLength(2)
    })

    test('should handle empty files list', async() => {
      fetch.mockResolvedValueOnce({
        json: () => Promise.resolve({ files: [] })
      })

      const result = await apiService.getFilesList()

      expect(result.files).toEqual([])
      expect(result.files).toHaveLength(0)
    })

    test('should handle null files list', async() => {
      fetch.mockResolvedValueOnce({
        json: () => Promise.resolve({ files: null })
      })

      const result = await apiService.getFilesList()

      expect(result.files).toBeNull()
    })

    test('should handle API errors', async() => {
      fetch.mockRejectedValueOnce(new Error('API error'))

      await expect(apiService.getFilesList()).rejects.toThrow('API error')
    })
  })

  describe('createFormData', () => {
    test('should create FormData with video file', () => {
      const mockFile = new File(['video content'], 'test.mp4', { type: 'video/mp4' })

      const formData = apiService.createFormData(mockFile)

      expect(formData).toBeInstanceOf(FormData)
      expect(formData.get('video')).toBe(mockFile)
    })

    test('should handle different file types', () => {
      const mockFile = new File(['video content'], 'test.avi', { type: 'video/avi' })

      const formData = apiService.createFormData(mockFile)

      expect(formData).toBeInstanceOf(FormData)
      expect(formData.get('video')).toBe(mockFile)
    })

    test('should handle file with no type', () => {
      const mockFile = new File(['video content'], 'test.mov')

      const formData = apiService.createFormData(mockFile)

      expect(formData).toBeInstanceOf(FormData)
      expect(formData.get('video')).toBe(mockFile)
    })
  })
})
