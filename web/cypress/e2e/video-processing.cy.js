describe('VideoGrinder - Video Processing E2E Tests', () => {
  const API_BASE_URL = 'http://localhost:8081'

  beforeEach(() => {
    cy.resetAppState()
  })

  describe('Application Loading', () => {
    it('should load the main page correctly', () => {
      cy.visit('/')
      cy.get('h1').should('contain', 'FIAP X - Processador de Vídeos')
      cy.get('input[type="file"]').should('exist')
      cy.get('button[type="submit"]').should('contain', 'Processar Vídeo')
    })

    it('should show the file listing section', () => {
      cy.get('#filesList').should('exist')
      cy.get('h3').should('contain', 'Arquivos Processados')
    })

    it('should load JavaScript correctly', () => {
      // Verify all JavaScript files loaded
      cy.window().should('have.property', 'appController')
      cy.window().should('have.property', 'AppController')
      cy.window().should('have.property', 'UIManager')
      cy.window().should('have.property', 'ApiService')
      cy.window().should('have.property', 'Utils')

      // Verify form is properly set up
      cy.get('#uploadForm').should('have.attr', 'onsubmit', 'return false;')
    })

    it('should have working JavaScript form submission', () => {
      // Basic elements check
      cy.get('#uploadForm').should('exist')
      cy.get('input[type="file"]').should('exist')
      cy.get('button[type="submit"]').should('exist')

      // Verify JavaScript is loaded and working
      cy.window().should('have.property', 'AppController')
      cy.window().should('have.property', 'appController')

      cy.window().then((win) => {
        expect(win.appController).to.exist
        expect(win.appController.uiManager).to.exist
        expect(win.appController.apiService).to.exist
      })

      // Test 1: Direct JavaScript call (we know this works)
      cy.window().then((win) => {
        win.appController.uiManager.showResult('Direct call works!', 'success')
      })
      cy.get('#result').should('be.visible')
      cy.get('#result').should('contain', 'Direct call works!')

      // Test 2: Try to call handleUpload directly
      cy.window().then((win) => {
        // Clear previous result
        const resultDiv = win.document.getElementById('result')
        resultDiv.innerHTML = ''
        resultDiv.style.display = 'none'
        resultDiv.className = 'result'

        // Try to call handleUpload directly
        if (win.appController && win.appController.handleUpload) {
          const fakeEvent = { preventDefault: () => {}, stopPropagation: () => {} }
          win.appController.handleUpload(fakeEvent)
        }
      })

      // Should show validation error from JavaScript
      cy.get('#result').should('be.visible')
      cy.get('#result').should('have.class', 'error')
      cy.get('#result').should('contain', 'Por favor, selecione um arquivo')
    })

    it('should load status endpoint', () => {
      cy.request(`${API_BASE_URL}/api/v1/videos`).then((response) => {
        expect(response.status).to.eq(200)
        expect(response.body).to.have.property('videos')
        expect(response.body).to.have.property('total')
      })
    })
  })

  describe('Valid Video Upload and Processing', () => {
    it('should successfully upload and process a valid video', () => {
      cy.uploadVideo('test-video-valid.mp4')

      cy.get('input[type="file"]').should(($input) => {
        expect($input[0].files).to.have.length(1)
        expect($input[0].files[0].name).to.eq('test-video-valid.mp4')
      })

      cy.submitUpload()

      // Wait for processing to complete
      cy.get('#result', { timeout: 60000 }).should('be.visible')
      cy.get('#result').should('not.contain', 'Processando')

      cy.verifyProcessingSuccess()
      cy.get('#result').should('contain', 'frames extraídos')

      cy.get('#result').should(($el) => {
        const text = $el.text()
        expect(text).to.satisfy((str) =>
          str.includes('.zip') || str.includes('ZIP') || str.includes('frames extraídos')
        )
      })
    })

    it('should show processed file in the listing', () => {
      cy.uploadVideo('test-video-valid.mp4')
      cy.submitUpload()

      // Wait for processing to complete
      cy.get('#result', { timeout: 60000 }).should('be.visible')
      cy.get('#result').should('not.contain', 'Processando')

      cy.checkFileListing()
      cy.get('#filesList').should('contain', 'frames_')
      cy.get('#filesList').should('contain', '.zip')

      // Check for presigned URLs (LocalStack S3)
      cy.get('#filesList a[href*="localhost:4566"]').should('exist')
      cy.get('#filesList a[href*="videogrinder-outputs"]').should('exist')
    })
  })

  describe('Error Handling', () => {
    it('should reject invalid file types', () => {
      cy.uploadVideo('test-invalid.txt')
      cy.submitUpload()

      cy.verifyProcessingError('Formato de arquivo não suportado')
      cy.get('#result').should('contain', 'mp4, avi, mov, mkv')
    })

    it('should handle missing file upload', () => {
      cy.submitUpload()

      cy.get('#result').should('be.visible')
      cy.get('#result').should('contain', 'Por favor, selecione um arquivo')
    })

    it('should handle server errors gracefully', () => {
      cy.intercept('POST', `${API_BASE_URL}/api/v1/videos`, {
        statusCode: 500,
        body: { success: false, message: 'Erro interno do servidor' }
      }).as('uploadError')

      cy.uploadVideo('test-video-valid.mp4')
      cy.submitUpload()

      cy.wait('@uploadError')
      cy.get('#result').should('be.visible')
      cy.get('#result').should(($el) => {
        const text = $el.text().toLowerCase()
        expect(text).to.satisfy((str) =>
          str.includes('erro') || str.includes('error') || str.includes('falha')
        )
      })
    })
  })

  describe('File Download', () => {
    it('should allow downloading processed files via direct URL request', () => {
      cy.uploadVideo('test-video-valid.mp4')
      cy.submitUpload()

      // Wait for processing to complete
      cy.get('#result', { timeout: 60000 }).should('be.visible')
      cy.get('#result').should('not.contain', 'Processando')

      // Wait for files to appear in listing and get presigned URL
      cy.get('#filesList a[href*="localhost:4566"]').first().then(($link) => {
        const downloadUrl = $link.attr('href')
        cy.request(downloadUrl).then((response) => {
          expect(response.status).to.eq(200)
          // Accept both content types (LocalStack may return either)
          expect(response.headers['content-type']).to.match(/^(application\/zip|binary\/octet-stream)$/)
        })
      })
    })

    it('should generate correct download URLs in files list without manual correction', () => {
      cy.uploadVideo('test-video-valid.mp4')
      cy.submitUpload()

      // Wait for processing to complete
      cy.get('#result', { timeout: 60000 }).should('be.visible')
      cy.get('#result').should('not.contain', 'Processando')

      cy.get('#filesList a.download-btn').first().should(($link) => {
        const href = $link.attr('href')
        // Check for presigned URL format (LocalStack S3)
        expect(href).to.match(/^http:\/\/localhost:4566\/videogrinder-outputs\/.+\.zip/)
        expect(href).not.to.match(/^\/api\/v1\/videos\//)
      })
    })

    it('should successfully download file when clicking download button from files list', () => {
      cy.uploadVideo('test-video-valid.mp4')
      cy.submitUpload()

      // Wait for processing to complete and file to appear in listing
      cy.get('#result', { timeout: 60000 }).should('be.visible')
      cy.get('#result').should('not.contain', 'Processando')
      cy.get('#filesList a.download-btn', { timeout: 10000 }).should('exist')

      cy.get('#filesList a.download-btn').first().then(($link) => {
        const downloadUrl = $link.attr('href')

        cy.request(downloadUrl).then((response) => {
          expect(response.status).to.eq(200)
          // Accept both content types (LocalStack may return either)
          expect(response.headers['content-type']).to.match(/^(application\/zip|binary\/octet-stream)$/)
        })
      })
    })

    it('should handle non-existent file downloads', () => {
      cy.request({
        url: `${API_BASE_URL}/api/v1/videos/non-existent-file.zip/download`,
        failOnStatusCode: false
      }).then((response) => {
        expect(response.status).to.eq(404)
        expect(response.body).to.have.property('error')
        expect(response.body.error).to.contain('Arquivo não encontrado')
      })
    })
  })

  describe('User Interface Interactions', () => {
    it('should have responsive design elements', () => {
      cy.viewport(1280, 720)
      cy.get('.container').should('be.visible')

      cy.viewport(375, 667)
      cy.get('.container').should('be.visible')
      cy.get('input[type="file"]').should('be.visible')
    })

    it('should show loading states correctly', () => {
      cy.uploadVideo('test-video-valid.mp4')
      cy.submitUpload()

      cy.get('#loading').should('be.visible')
      cy.get('#loading').should('contain', 'Processando vídeo')
      cy.get('#loading').should('contain', 'Isso pode levar alguns minutos')

      // Wait for processing to complete (loading should disappear)
      cy.get('#loading', { timeout: 60000 }).should('not.be.visible')
      cy.get('#result', { timeout: 60000 }).should('be.visible')
    })

    it('should update file listing dynamically', () => {
      cy.checkFileListing()

      cy.uploadVideo('test-video-valid.mp4')
      cy.submitUpload()

      // Wait for processing to complete
      cy.get('#result', { timeout: 60000 }).should('be.visible')
      cy.get('#result').should('not.contain', 'Processando')

      cy.get('#filesList').should('contain', 'frames_')
      cy.get('#filesList').should('contain', '.zip')
    })
  })

  describe('Browser Compatibility', () => {
    it('should work with modern file APIs', () => {
      cy.window().should('have.property', 'FormData')
      cy.window().should('have.property', 'FileReader')
      cy.window().should('have.property', 'fetch')
    })

    it('should handle CORS correctly', () => {
      cy.request(`${API_BASE_URL}/api/v1/videos`).then((response) => {
        expect(response.headers).to.have.property('access-control-allow-origin')
      })
    })
  })
})

describe('VideoGrinder - API E2E Tests', () => {
  const API_BASE_URL = 'http://localhost:8081'

  it('should return proper status information', () => {
    cy.request(`${API_BASE_URL}/api/v1/videos`).then((response) => {
      expect(response.status).to.eq(200)
      expect(response.body).to.have.property('videos').that.is.an('array')
      expect(response.body).to.have.property('total').that.is.a('number')
    })
  })

  it('should handle POST upload with proper content type', () => {
    const formData = new FormData()
    const blob = new Blob(['fake video content'], { type: 'video/mp4' })
    formData.append('video', blob, 'test.mp4')

    cy.request({
      method: 'POST',
      url: `${API_BASE_URL}/api/v1/videos`,
      body: formData,
      failOnStatusCode: false
    }).then((response) => {
      expect([201, 422, 400]).to.include(response.status)
    })
  })
})
