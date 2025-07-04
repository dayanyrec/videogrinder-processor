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

      cy.get('button[type="submit"]').click()

      cy.get('#loading').should('be.visible')
      cy.get('#loading').should('contain', 'Processando vídeo')

      cy.waitForUploadComplete()

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
      cy.uploadAndProcess('test-video-valid.mp4')

      cy.checkFileListing()
      cy.get('#filesList').should('contain', 'frames_')
      cy.get('#filesList').should('contain', '.zip')

      cy.get('#filesList a[href*="/api/v1/videos/"]').should('exist')
      cy.get('#filesList a[href*="/download"]').should('exist')
    })
  })

  describe('Error Handling', () => {
    it('should reject invalid file types', () => {
      cy.uploadVideo('test-invalid.txt')
      cy.get('button[type="submit"]').click()

      cy.verifyProcessingError('Formato de arquivo não suportado')
      cy.get('#result').should('contain', 'mp4, avi, mov, mkv')
    })

    it('should handle missing file upload', () => {
      cy.get('button[type="submit"]').click()

      cy.get('input[type="file"]').then(($input) => {
        expect($input[0].validity.valid).to.be.false
      })
    })

    it('should handle server errors gracefully', () => {
      cy.intercept('POST', `${API_BASE_URL}/api/v1/videos`, {
        statusCode: 500,
        body: { success: false, message: 'Erro interno do servidor' }
      }).as('uploadError')

      cy.uploadVideo('test-video-valid.mp4')
      cy.get('button[type="submit"]').click()

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
    it('should allow downloading processed files', () => {
      cy.uploadAndProcess('test-video-valid.mp4')

      cy.get('#filesList a[href*="/api/v1/videos/"]').first().then(($link) => {
        const downloadUrl = $link.attr('href').replace('/api/v1/videos/', `${API_BASE_URL}/api/v1/videos/`)
        cy.request(downloadUrl).then((response) => {
          expect(response.status).to.eq(200)
          expect(response.headers['content-type']).to.eq('application/zip')
          expect(response.headers['content-disposition']).to.contain('attachment')
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
      cy.get('button[type="submit"]').click()

      cy.get('#loading').should('be.visible')
      cy.get('#loading').should('contain', 'Processando vídeo')
      cy.get('#loading').should('contain', 'Isso pode levar alguns minutos')

      cy.waitForUploadComplete()
    })

    it('should update file listing dynamically', () => {
      cy.checkFileListing()

      cy.uploadAndProcess('test-video-valid.mp4')

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
