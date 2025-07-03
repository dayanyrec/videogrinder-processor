const UIManager = require('../static/js/ui-manager.js')

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
  })

  describe('constructor', () => {
    test('should initialize with correct DOM elements', () => {
      expect(uiManager.elements.loading).toBeTruthy()
      expect(uiManager.elements.result).toBeTruthy()
      expect(uiManager.elements.filesList).toBeTruthy()
      expect(uiManager.elements.uploadForm).toBeTruthy()
      expect(uiManager.elements.videoFile).toBeTruthy()
    })
  })

  describe('showResult', () => {
    test('should display success message correctly', () => {
      const message = 'Upload successful!'
      const type = 'success'

      uiManager.showResult(message, type)

      expect(uiManager.elements.result.innerHTML).toBe(message)
      expect(uiManager.elements.result.className).toBe('result success')
      expect(uiManager.elements.result.style.display).toBe('block')
    })

    test('should display error message correctly', () => {
      const message = 'Upload failed!'
      const type = 'error'

      uiManager.showResult(message, type)

      expect(uiManager.elements.result.innerHTML).toBe(message)
      expect(uiManager.elements.result.className).toBe('result error')
      expect(uiManager.elements.result.style.display).toBe('block')
    })

    test('should handle HTML content in message', () => {
      const message = '<strong>Success!</strong> File uploaded.'
      const type = 'success'

      uiManager.showResult(message, type)

      expect(uiManager.elements.result.innerHTML).toBe(message)
      expect(uiManager.elements.result.className).toBe('result success')
    })
  })

  describe('showLoading', () => {
    test('should show loading and hide result', () => {
      uiManager.showLoading()

      expect(uiManager.elements.loading.style.display).toBe('block')
      expect(uiManager.elements.result.style.display).toBe('none')
    })
  })

  describe('hideLoading', () => {
    test('should hide loading element', () => {
      uiManager.hideLoading()

      expect(uiManager.elements.loading.style.display).toBe('none')
    })
  })

  describe('displayFilesList', () => {
    test('should display files list when files exist', () => {
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

      uiManager.displayFilesList(mockFiles)

      const filesListDiv = uiManager.elements.filesList
      expect(filesListDiv.innerHTML).toContain('test-video.mp4')
      expect(filesListDiv.innerHTML).toContain('another-video.avi')
      expect(filesListDiv.innerHTML).toContain('1000 KB')
      expect(filesListDiv.innerHTML).toContain('2000 KB')
      expect(filesListDiv.innerHTML).toContain('/download/test.zip')
      expect(filesListDiv.innerHTML).toContain('/download/another.zip')
      expect(filesListDiv.innerHTML).toContain('class="file-item"')
      expect(filesListDiv.innerHTML).toContain('class="download-btn"')
    })

    test('should display empty message when no files exist', () => {
      uiManager.displayFilesList([])

      const filesListDiv = uiManager.elements.filesList
      expect(filesListDiv.innerHTML).toContain('Nenhum arquivo processado ainda.')
    })

    test('should display empty message when files is null', () => {
      uiManager.displayFilesList(null)

      const filesListDiv = uiManager.elements.filesList
      expect(filesListDiv.innerHTML).toContain('Nenhum arquivo processado ainda.')
    })

    test('should display empty message when files is undefined', () => {
      uiManager.displayFilesList(undefined)

      const filesListDiv = uiManager.elements.filesList
      expect(filesListDiv.innerHTML).toContain('Nenhum arquivo processado ainda.')
    })
  })

  describe('displayFilesError', () => {
    test('should display error message for files loading', () => {
      uiManager.displayFilesError()

      const filesListDiv = uiManager.elements.filesList
      expect(filesListDiv.innerHTML).toContain('Erro ao carregar arquivos.')
      expect(filesListDiv.innerHTML).toContain('color: red')
    })
  })

  describe('getSelectedFile', () => {
    test('should return selected file', () => {
      const mockFile = new File(['content'], 'test.mp4', { type: 'video/mp4' })

      Object.defineProperty(uiManager.elements.videoFile, 'files', {
        value: [mockFile],
        writable: false
      })

      const selectedFile = uiManager.getSelectedFile()
      expect(selectedFile).toBe(mockFile)
    })

    test('should return undefined when no file selected', () => {
      Object.defineProperty(uiManager.elements.videoFile, 'files', {
        value: [],
        writable: false
      })

      const selectedFile = uiManager.getSelectedFile()
      expect(selectedFile).toBeUndefined()
    })
  })

  describe('getUploadForm', () => {
    test('should return upload form element', () => {
      const form = uiManager.getUploadForm()
      expect(form).toBe(uiManager.elements.uploadForm)
      expect(form.tagName).toBe('FORM')
    })
  })
})
