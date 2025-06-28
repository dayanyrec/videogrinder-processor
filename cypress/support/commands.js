// cypress/support/commands.js

// Custom command to upload a video file
Cypress.Commands.add('uploadVideo', (fileName) => {
  cy.get('input[type="file"]').selectFile(`cypress/fixtures/${fileName}`, { force: true })
})

// Custom command to wait for processing to complete
Cypress.Commands.add('waitForProcessing', (timeout = 30000) => {
  cy.get('#loading', { timeout }).should('be.visible')
  cy.get('#result', { timeout }).should('be.visible')
  cy.get('#result').should('not.contain', 'Processando')
})

// Custom command to check if processing was successful
Cypress.Commands.add('verifyProcessingSuccess', () => {
  cy.get('#result').should('have.class', 'success')
  cy.get('#result').should('contain', 'Processamento concluído')
})

// Custom command to check if processing failed
Cypress.Commands.add('verifyProcessingError', (errorMessage) => {
  cy.get('#result').should('have.class', 'error')
  if (errorMessage) {
    cy.get('#result').should('contain', errorMessage)
  }
})

// Custom command to check file listing
Cypress.Commands.add('checkFileListing', () => {
  cy.get('#filesList').should('be.visible')
  cy.get('#filesList').should('not.contain', 'Carregando...')
})

// Custom command to download file (check if download starts)
Cypress.Commands.add('downloadFile', (fileName) => {
  cy.get(`a[href*="${fileName}"]`).click()
  // Note: Cypress doesn't handle file downloads naturally
  // We verify the download link exists and is clickable
})

// Custom command to reset application state
Cypress.Commands.add('resetAppState', () => {
  // Clear any existing files and reset form
  cy.visit('/')
  cy.get('h1').should('contain', 'FIAP X - Processador de Vídeos')
  cy.get('input[type="file"]').should('exist')
})

// Custom command to wait for file to appear in listing
Cypress.Commands.add('waitForFileInListing', (fileName) => {
  cy.get('#filesList').should('contain', fileName)
})
