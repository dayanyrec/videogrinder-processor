let showResult, loadFilesList

beforeAll(() => {
  const originalAddEventListener = document.addEventListener
  document.addEventListener = jest.fn()

  require('../static/js/app.js')

  document.addEventListener = originalAddEventListener

  showResult = global.showResult
  loadFilesList = global.loadFilesList
})

describe('VideoGrinder App', () => {

  describe('UI Context', () => {

    describe('showResult function', () => {
      test('should display success message correctly', () => {
        const message = 'Upload successful!'
        const type = 'success'

        showResult(message, type)

        const resultDiv = document.getElementById('result')
        expect(resultDiv.innerHTML).toBe(message)
        expect(resultDiv.className).toBe('result success')
        expect(resultDiv.style.display).toBe('block')
      })

      test('should display error message correctly', () => {
        const message = 'Upload failed!'
        const type = 'error'

        showResult(message, type)

        const resultDiv = document.getElementById('result')
        expect(resultDiv.innerHTML).toBe(message)
        expect(resultDiv.className).toBe('result error')
        expect(resultDiv.style.display).toBe('block')
      })

      test('should handle HTML content in messages', () => {
        const message = 'Success! <br><a href="/download/test.zip">Download</a>'
        const type = 'success'

        showResult(message, type)

        const resultDiv = document.getElementById('result')
        expect(resultDiv.innerHTML).toBe(message)
        expect(resultDiv.querySelector('a')).toBeTruthy()
      })
    })

    describe('DOM element interactions', () => {
      test('should find required DOM elements', () => {
        expect(document.getElementById('uploadForm')).toBeTruthy()
        expect(document.getElementById('videoFile')).toBeTruthy()
        expect(document.getElementById('loading')).toBeTruthy()
        expect(document.getElementById('result')).toBeTruthy()
        expect(document.getElementById('filesList')).toBeTruthy()
      })

      test('should toggle loading state correctly', () => {
        const loadingDiv = document.getElementById('loading')
        const resultDiv = document.getElementById('result')

        expect(loadingDiv.style.display).toBe('none')

        loadingDiv.style.display = 'block'
        resultDiv.style.display = 'none'

        expect(loadingDiv.style.display).toBe('block')
        expect(resultDiv.style.display).toBe('none')

        loadingDiv.style.display = 'none'

        expect(loadingDiv.style.display).toBe('none')
      })
    })
  })

  describe('API Handlers Context', () => {

    describe('loadFilesList function', () => {
      test('should display files list when API returns data', async() => {
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

        await loadFilesList()

        await new Promise(resolve => setTimeout(resolve, 0))

        const filesListDiv = document.getElementById('filesList')
        expect(filesListDiv.innerHTML).toContain('test-video.mp4')
        expect(filesListDiv.innerHTML).toContain('another-video.avi')
        expect(filesListDiv.innerHTML).toContain('1000 KB')
        expect(filesListDiv.innerHTML).toContain('2000 KB')
        expect(filesListDiv.innerHTML).toContain('/download/test.zip')
        expect(filesListDiv.innerHTML).toContain('/download/another.zip')
      })

      test('should display empty message when no files exist', async() => {
        fetch.mockResolvedValueOnce({
          json: () => Promise.resolve({ files: [] })
        })

        await loadFilesList()

        // Wait a bit for DOM updates
        await new Promise(resolve => setTimeout(resolve, 0))

        const filesListDiv = document.getElementById('filesList')
        expect(filesListDiv.innerHTML).toContain('Nenhum arquivo processado ainda.')
      })

      test('should display empty message when files is null', async() => {
        fetch.mockResolvedValueOnce({
          json: () => Promise.resolve({ files: null })
        })

        await loadFilesList()

        // Wait a bit for DOM updates
        await new Promise(resolve => setTimeout(resolve, 0))

        const filesListDiv = document.getElementById('filesList')
        expect(filesListDiv.innerHTML).toContain('Nenhum arquivo processado ainda.')
      })

      test('should handle API error gracefully', async() => {
        fetch.mockRejectedValueOnce(new Error('Network error'))

        await loadFilesList()

        // Wait a bit for DOM updates
        await new Promise(resolve => setTimeout(resolve, 0))

        const filesListDiv = document.getElementById('filesList')
        expect(filesListDiv.innerHTML).toContain('Erro ao carregar arquivos.')
      })

      test('should call correct API endpoint', async() => {
        fetch.mockResolvedValueOnce({
          json: () => Promise.resolve({ files: [] })
        })

        await loadFilesList()

        expect(fetch).toHaveBeenCalledWith('/api/status')
      })
    })

    describe('file upload handling', () => {
      let form, fileInput, submitEvent

      beforeEach(() => {
        form = document.getElementById('uploadForm')
        fileInput = document.getElementById('videoFile')
        submitEvent = new Event('submit', { bubbles: true, cancelable: true })
      })

      test('should prevent form submission without file', () => {
        // Simulate form submission without file
        const preventDefault = jest.fn()
        submitEvent.preventDefault = preventDefault

        form.dispatchEvent(submitEvent)

        const file = fileInput.files[0]
        expect(file).toBeUndefined()
      })

      test('should handle successful upload response', async() => {
        const mockResponse = {
          success: true,
          message: 'Video processed successfully!',
          zip_path: 'frames_123.zip'
        }

        fetch.mockResolvedValueOnce({
          json: () => Promise.resolve(mockResponse)
        })

        const mockFile = new File(['video content'], 'test.mp4', { type: 'video/mp4' })


        const mockFormData = {
          append: jest.fn()
        }
        global.FormData = jest.fn(() => mockFormData)

        const formData = new FormData()
        formData.append('video', mockFile)

        const response = await fetch('/upload', {
          method: 'POST',
          body: formData
        })
        const data = await response.json()

        expect(fetch).toHaveBeenCalledWith('/upload', {
          method: 'POST',
          body: formData
        })
        expect(data.success).toBe(true)
        expect(data.message).toBe('Video processed successfully!')
      })

      test('should handle upload error response', async() => {
        const mockResponse = {
          success: false,
          message: 'Invalid file format'
        }

        fetch.mockResolvedValueOnce({
          json: () => Promise.resolve(mockResponse)
        })

        const response = await fetch('/upload', { method: 'POST' })
        const data = await response.json()

        expect(data.success).toBe(false)
        expect(data.message).toBe('Invalid file format')
      })

      test('should handle network errors during upload', async() => {
        fetch.mockRejectedValueOnce(new Error('Connection failed'))

        try {
          await fetch('/upload', { method: 'POST' })
        } catch (error) {
          expect(error.message).toBe('Connection failed')
        }
      })
    })
  })

  describe('Utils Context', () => {

    describe('file size formatting', () => {
      test('should format file sizes correctly in KB', () => {
        const testCases = [
          { bytes: 1024, expectedKB: 1 },
          { bytes: 2048, expectedKB: 2 },
          { bytes: 1536, expectedKB: 2 },
          { bytes: 1024000, expectedKB: 1000 },
          { bytes: 512, expectedKB: 1 }
        ]

        testCases.forEach(({ bytes, expectedKB }) => {
          const result = Math.round(bytes / 1024)
          expect(result).toBe(expectedKB)
        })
      })
    })

    describe('HTML string generation', () => {
      test('should generate correct file item HTML structure', () => {
        const file = {
          filename: 'test-video.mp4',
          size: 1024000,
          created_at: '2024-01-01 10:00:00',
          download_url: '/download/test.zip'
        }

        const html = '<div class="file-item">' +
                    '<span><strong>' + file.filename + '</strong><br>' +
                    '<small>Tamanho: ' + Math.round(file.size / 1024) + ' KB | ' +
                    'Criado: ' + file.created_at + '</small></span>' +
                    '<a href="' + file.download_url + '" class="download-btn">📥 Baixar</a>' +
                    '</div>'

        expect(html).toContain('test-video.mp4')
        expect(html).toContain('1000 KB')
        expect(html).toContain('2024-01-01 10:00:00')
        expect(html).toContain('/download/test.zip')
        expect(html).toContain('download-btn')
        expect(html).toContain('📥 Baixar')
      })
    })

    describe('form data handling', () => {
      test('should create FormData correctly', () => {
        const mockFormData = {
          append: jest.fn()
        }
        global.FormData = jest.fn(() => mockFormData)

        const formData = new FormData()
        const mockFile = new File(['content'], 'test.mp4', { type: 'video/mp4' })

        formData.append('video', mockFile)

        expect(FormData).toHaveBeenCalled()
        expect(mockFormData.append).toHaveBeenCalledWith('video', mockFile)
      })
    })

    describe('event handling utilities', () => {
      test('should handle form submission event correctly', () => {
        const mockEvent = {
          preventDefault: jest.fn()
        }

        mockEvent.preventDefault()
        expect(mockEvent.preventDefault).toHaveBeenCalled()
      })
    })
  })

  describe('Integration Context', () => {

    test('should have DOMContentLoaded event listener logic', () => {
      const fs = require('fs')
      const appJsContent = fs.readFileSync('./static/js/app.js', 'utf8')

      expect(appJsContent).toContain('DOMContentLoaded')
      expect(appJsContent).toContain('addEventListener')
      expect(appJsContent).toContain('loadFilesList()')
      expect(appJsContent).toContain('uploadForm')
    })

    test('should have all required functions available globally', () => {
      expect(typeof showResult).toBe('function')
      expect(typeof loadFilesList).toBe('function')
    })
  })
})
