Cypress.Commands.add('uploadVideo', (fileName) => {
  cy.get('input[type="file"]').selectFile(`cypress/fixtures/${fileName}`, { force: true })
})

Cypress.Commands.add('waitForUploadComplete', () => {
  cy.get('#loading', { timeout: 45000 }).should('not.be.visible')
  cy.get('#result', { timeout: 45000 }).should('be.visible')
  cy.get('#result').should('not.contain', 'Processando')
})

Cypress.Commands.add('uploadAndProcess', (fileName) => {
  cy.uploadVideo(fileName)
  cy.submitUpload()
  cy.waitForUploadComplete()
  cy.verifyProcessingSuccess()
})

Cypress.Commands.add('verifyProcessingSuccess', () => {
  cy.get('#result').should('have.class', 'success')
  cy.get('#result').should('contain', 'Processamento concluído')
})

Cypress.Commands.add('verifyProcessingError', (errorMessage) => {
  cy.get('#result').should('have.class', 'error')
  if (errorMessage) {
    cy.get('#result').should('contain', errorMessage)
  }
})

Cypress.Commands.add('checkFileListing', () => {
  cy.get('#filesList').should('be.visible')
  cy.get('#filesList').should('not.contain', 'Carregando...')
})

Cypress.Commands.add('downloadFile', (fileName) => {
  cy.get(`a[href*="${fileName}"]`).click()
})

Cypress.Commands.add('resetAppState', () => {
  cy.visit('/')
  cy.get('h1').should('contain', 'FIAP X - Processador de Vídeos')
  cy.get('input[type="file"]').should('exist')

  // Wait for JavaScript to fully load
  cy.window().should('have.property', 'appController')
  cy.window().should('have.property', 'AppController')
  cy.window().should('have.property', 'UIManager')
  cy.window().should('have.property', 'ApiService')
  cy.window().should('have.property', 'Utils')

  // Wait for app controller to be initialized
  cy.window().then((win) => {
    expect(win.appController).to.exist
    expect(win.appController.uiManager).to.exist
    expect(win.appController.apiService).to.exist
  })
})

Cypress.Commands.add('waitForFileInListing', (fileName) => {
  cy.get('#filesList').should('contain', fileName)
})

Cypress.Commands.add('submitUpload', () => {
  cy.waitForFormReady()

  // Call handleUpload directly (works better in Cypress)
  cy.window().then((win) => {
    const fakeEvent = { preventDefault: () => {}, stopPropagation: () => {} }
    win.appController.handleUpload(fakeEvent)
  })
})

Cypress.Commands.add('waitForFormReady', () => {
  // Wait for form to be ready with event listeners
  cy.get('#uploadForm').should('exist')
  cy.get('#uploadForm').should('have.attr', 'onsubmit', 'return false;')

  // Wait for appController to be initialized
  cy.window().should('have.property', 'appController')
  cy.window().then((win) => {
    expect(win.appController).to.exist
    expect(win.appController.uiManager).to.exist
    expect(win.appController.apiService).to.exist
  })
})
