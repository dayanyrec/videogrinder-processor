const UIManager = require('../static/js/ui-manager.js')

delete global.window.location
global.window = Object.create(window)
global.window.location = {
  port: '8081',
  protocol: 'http:',
  hostname: 'localhost',
  origin: 'http://localhost:8081'
}

describe('UIManager Class', () => {
  let uiManager

  beforeEach(() => {
    document.body.innerHTML = `
      <div id="loading" style="display: none;">Loading...</div>
      <div id="result" style="display: none;"></div>
      <div id="filesList"></div>
      <form id="uploadForm">
        <input type="file" id="videoFile" />
      </form>
    `
    uiManager = new UIManager()

    global.window.location.port = '8081'
    global.window.location.origin = 'http://localhost:8081'
  })

  describe('constructor', () => {
    test('should initialize with all required DOM elements properly bound', () => {
      expect(uiManager.elements.loading).toBeTruthy()
      expect(uiManager.elements.result).toBeTruthy()
      expect(uiManager.elements.filesList).toBeTruthy()
      expect(uiManager.elements.uploadForm).toBeTruthy()
      expect(uiManager.elements.videoFile).toBeTruthy()
    })
  })

  describe('getApiBaseURL', () => {
    test('should return API URL for development environment (port 8080)', () => {
      global.window.location.port = '8080'
      global.window.location.origin = 'http://localhost:8080'

      const baseURL = uiManager.getApiBaseURL()
      expect(baseURL).toBe('http://localhost:8081')
    })

    test('should return same origin for production environment (port 8081)', () => {
      global.window.location.port = '8081'
      global.window.location.origin = 'http://localhost:8081'

      const baseURL = uiManager.getApiBaseURL()
      expect(baseURL).toBe('http://localhost:8081')
    })

    test('should return same origin for other ports', () => {
      global.window.location.port = '3000'
      global.window.location.origin = 'http://localhost:3000'

      const baseURL = uiManager.getApiBaseURL()
      expect(baseURL).toBe('http://localhost:3000')
    })
  })

  describe('showResult', () => {
    test('should display success message with correct styling and visibility', () => {
      const message = 'Upload successful!'
      const type = 'success'

      uiManager.showResult(message, type)

      expect(uiManager.elements.result.innerHTML).toBe(message)
      expect(uiManager.elements.result.className).toBe('result success')
      expect(uiManager.elements.result.style.display).toBe('block')
    })

    test('should display error message with correct styling and visibility', () => {
      const message = 'Upload failed!'
      const type = 'error'

      uiManager.showResult(message, type)

      expect(uiManager.elements.result.innerHTML).toBe(message)
      expect(uiManager.elements.result.className).toBe('result error')
      expect(uiManager.elements.result.style.display).toBe('block')
    })

    test('should properly render HTML content within result message', () => {
      const message = '<strong>Success!</strong> File uploaded.'
      const type = 'success'

      uiManager.showResult(message, type)

      expect(uiManager.elements.result.innerHTML).toBe(message)
      expect(uiManager.elements.result.className).toBe('result success')
    })
  })

  describe('showLoading', () => {
    test('should show loading indicator and hide result message', () => {
      uiManager.showLoading()

      expect(uiManager.elements.loading.style.display).toBe('block')
      expect(uiManager.elements.result.style.display).toBe('none')
    })
  })

  describe('hideLoading', () => {
    test('should hide loading indicator from user interface', () => {
      uiManager.hideLoading()

      expect(uiManager.elements.loading.style.display).toBe('none')
    })
  })

  describe('displayFilesList', () => {
    test('should render formatted files list with absolute download URLs when files exist', () => {
      const mockFiles = [
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

      uiManager.displayFilesList(mockFiles)

      const filesListDiv = uiManager.elements.filesList
      expect(filesListDiv.innerHTML).toContain('test-video.mp4')
      expect(filesListDiv.innerHTML).toContain('another-video.avi')
      expect(filesListDiv.innerHTML).toContain('1000 KB')
      expect(filesListDiv.innerHTML).toContain('2000 KB')
      expect(filesListDiv.innerHTML).toContain('http://localhost:8081/api/v1/videos/test.zip/download')
      expect(filesListDiv.innerHTML).toContain('http://localhost:8081/api/v1/videos/another.zip/download')
      expect(filesListDiv.innerHTML).toContain('class="file-item"')
      expect(filesListDiv.innerHTML).toContain('class="download-btn"')
    })

    test('should convert relative URLs to absolute URLs for development environment', () => {
      global.window.location.port = '8080'
      global.window.location.origin = 'http://localhost:8080'

      const mockFiles = [
        {
          filename: 'dev-video.mp4',
          size: 1024000,
          created_at: '2024-01-01 10:00:00',
          download_url: '/api/v1/videos/dev.zip/download'
        }
      ]

      uiManager.displayFilesList(mockFiles)

      const filesListDiv = uiManager.elements.filesList
      expect(filesListDiv.innerHTML).toContain('http://localhost:8081/api/v1/videos/dev.zip/download')
      expect(filesListDiv.innerHTML).not.toContain('http://localhost:8080/api/v1/videos/dev.zip/download')
    })

    test('should handle absolute URLs without modification', () => {
      const mockFiles = [
        {
          filename: 'external-video.mp4',
          size: 1024000,
          created_at: '2024-01-01 10:00:00',
          download_url: 'https://external-api.com/videos/external.zip/download'
        }
      ]

      uiManager.displayFilesList(mockFiles)

      const filesListDiv = uiManager.elements.filesList
      expect(filesListDiv.innerHTML).toContain('https://external-api.com/videos/external.zip/download')
    })

    test('should display empty state message when no processed files exist', () => {
      uiManager.displayFilesList([])

      const filesListDiv = uiManager.elements.filesList
      expect(filesListDiv.innerHTML).toContain('Nenhum arquivo processado ainda.')
    })

    test('should display empty state message when files list is null', () => {
      uiManager.displayFilesList(null)

      const filesListDiv = uiManager.elements.filesList
      expect(filesListDiv.innerHTML).toContain('Nenhum arquivo processado ainda.')
    })

    test('should display empty state message when files list is undefined', () => {
      uiManager.displayFilesList(undefined)

      const filesListDiv = uiManager.elements.filesList
      expect(filesListDiv.innerHTML).toContain('Nenhum arquivo processado ainda.')
    })
  })

  describe('displayFilesError', () => {
    test('should display error message when files loading fails', () => {
      uiManager.displayFilesError()

      const filesListDiv = uiManager.elements.filesList
      expect(filesListDiv.innerHTML).toContain('Erro ao carregar arquivos.')
      expect(filesListDiv.innerHTML).toContain('color: red')
    })
  })

  describe('getSelectedFile', () => {
    test('should return the file selected by user from file input', () => {
      const mockFile = new File(['content'], 'test.mp4', { type: 'video/mp4' })

      Object.defineProperty(uiManager.elements.videoFile, 'files', {
        value: [mockFile],
        writable: false
      })

      const selectedFile = uiManager.getSelectedFile()
      expect(selectedFile).toBe(mockFile)
    })

    test('should return undefined when no file is selected', () => {
      Object.defineProperty(uiManager.elements.videoFile, 'files', {
        value: [],
        writable: false
      })

      const selectedFile = uiManager.getSelectedFile()
      expect(selectedFile).toBeUndefined()
    })
  })

  describe('getUploadForm', () => {
    test('should return the upload form DOM element for event binding', () => {
      const form = uiManager.getUploadForm()
      expect(form).toBe(uiManager.elements.uploadForm)
      expect(form.tagName).toBe('FORM')
    })
  })
})
