const mockUIManager = {
  getUploadForm: jest.fn(),
  getSelectedFile: jest.fn(),
  showResult: jest.fn(),
  showLoading: jest.fn(),
  hideLoading: jest.fn(),
  displayFilesList: jest.fn(),
  displayFilesError: jest.fn()
}

const mockApiService = {
  createFormData: jest.fn(),
  uploadVideo: jest.fn(),
  getFilesList: jest.fn()
}

const mockUtils = {
  validateFile: jest.fn(),
  createDownloadLink: jest.fn()
}

global.UIManager = jest.fn().mockImplementation(() => mockUIManager)
global.ApiService = jest.fn().mockImplementation(() => mockApiService)
global.Utils = mockUtils

const AppController = require('../static/js/app-controller.js')

describe('AppController Class', () => {
  let appController

  beforeEach(() => {
    jest.clearAllMocks()
    appController = new AppController()
  })

  describe('constructor', () => {
    test('should initialize with UIManager and ApiService', () => {
      expect(appController.uiManager).toBeDefined()
      expect(appController.apiService).toBeDefined()
    })
  })

  describe('init', () => {
    test('should load files list and setup event listeners', async() => {
      const mockForm = { addEventListener: jest.fn() }
      mockUIManager.getUploadForm.mockReturnValue(mockForm)
      mockApiService.getFilesList.mockResolvedValue({ files: [] })

      await appController.init()

      expect(mockApiService.getFilesList).toHaveBeenCalled()
      expect(mockUIManager.displayFilesList).toHaveBeenCalledWith([])
      expect(mockUIManager.getUploadForm).toHaveBeenCalled()
      expect(mockForm.addEventListener).toHaveBeenCalledWith('submit', expect.any(Function))
    })

    test('should handle error during initialization', async() => {
      const mockForm = { addEventListener: jest.fn() }
      mockUIManager.getUploadForm.mockReturnValue(mockForm)
      mockApiService.getFilesList.mockRejectedValue(new Error('API error'))

      await appController.init()

      expect(mockUIManager.displayFilesError).toHaveBeenCalled()
    })
  })

  describe('handleUpload', () => {
    test('should handle successful upload', async() => {
      const mockEvent = { preventDefault: jest.fn() }
      const mockFile = new File(['content'], 'test.mp4', { type: 'video/mp4' })
      const mockFormData = new FormData()

      mockUIManager.getSelectedFile.mockReturnValue(mockFile)
      mockUtils.validateFile.mockReturnValue({ valid: true, message: '' })
      mockApiService.createFormData.mockReturnValue(mockFormData)
      mockApiService.uploadVideo.mockResolvedValue({
        success: true,
        message: 'Upload successful!',
        zip_path: 'frames_123.zip'
      })
      mockUtils.createDownloadLink.mockReturnValue('<a href="/download/frames_123.zip">Download</a>')
      mockApiService.getFilesList.mockResolvedValue({ files: [] })

      await appController.handleUpload(mockEvent)

      expect(mockEvent.preventDefault).toHaveBeenCalled()
      expect(mockUIManager.showLoading).toHaveBeenCalled()
      expect(mockApiService.uploadVideo).toHaveBeenCalledWith(mockFormData)
      expect(mockUIManager.hideLoading).toHaveBeenCalled()
      expect(mockUIManager.showResult).toHaveBeenCalledWith(
        'Upload successful!<br><br><a href="/download/frames_123.zip">Download</a>',
        'success'
      )
      expect(mockApiService.getFilesList).toHaveBeenCalled()
    })

    test('should handle validation error', async() => {
      const mockEvent = { preventDefault: jest.fn() }

      mockUIManager.getSelectedFile.mockReturnValue(null)
      mockUtils.validateFile.mockReturnValue({ valid: false, message: 'No file selected' })

      await appController.handleUpload(mockEvent)

      expect(mockEvent.preventDefault).toHaveBeenCalled()
      expect(mockUIManager.showResult).toHaveBeenCalledWith('No file selected', 'error')
      expect(mockUIManager.showLoading).not.toHaveBeenCalled()
    })

    test('should handle upload failure', async() => {
      const mockEvent = { preventDefault: jest.fn() }
      const mockFile = new File(['content'], 'test.mp4', { type: 'video/mp4' })
      const mockFormData = new FormData()

      mockUIManager.getSelectedFile.mockReturnValue(mockFile)
      mockUtils.validateFile.mockReturnValue({ valid: true, message: '' })
      mockApiService.createFormData.mockReturnValue(mockFormData)
      mockApiService.uploadVideo.mockResolvedValue({
        success: false,
        message: 'Upload failed!'
      })

      await appController.handleUpload(mockEvent)

      expect(mockUIManager.showLoading).toHaveBeenCalled()
      expect(mockUIManager.hideLoading).toHaveBeenCalled()
      expect(mockUIManager.showResult).toHaveBeenCalledWith('Erro: Upload failed!', 'error')
    })

    test('should handle network error', async() => {
      const mockEvent = { preventDefault: jest.fn() }
      const mockFile = new File(['content'], 'test.mp4', { type: 'video/mp4' })
      const mockFormData = new FormData()

      mockUIManager.getSelectedFile.mockReturnValue(mockFile)
      mockUtils.validateFile.mockReturnValue({ valid: true, message: '' })
      mockApiService.createFormData.mockReturnValue(mockFormData)
      mockApiService.uploadVideo.mockRejectedValue(new Error('Network error'))

      await appController.handleUpload(mockEvent)

      expect(mockUIManager.showLoading).toHaveBeenCalled()
      expect(mockUIManager.hideLoading).toHaveBeenCalled()
      expect(mockUIManager.showResult).toHaveBeenCalledWith('Erro de conexão: Network error', 'error')
    })
  })

  describe('loadFilesList', () => {
    test('should load and display files list', async() => {
      const mockFiles = [
        { filename: 'test.mp4', size: 1024, created_at: '2024-01-01', download_url: '/download/test.zip' }
      ]
      mockApiService.getFilesList.mockResolvedValue({ files: mockFiles })

      await appController.loadFilesList()

      expect(mockApiService.getFilesList).toHaveBeenCalled()
      expect(mockUIManager.displayFilesList).toHaveBeenCalledWith(mockFiles)
    })

    test('should handle error when loading files list', async() => {
      mockApiService.getFilesList.mockRejectedValue(new Error('API error'))

      await appController.loadFilesList()

      expect(mockUIManager.displayFilesError).toHaveBeenCalled()
    })
  })

  describe('showResult', () => {
    test('should delegate to UIManager showResult', () => {
      const message = 'Test message'
      const type = 'success'

      appController.showResult(message, type)

      expect(mockUIManager.showResult).toHaveBeenCalledWith(message, type)
    })
  })
})
