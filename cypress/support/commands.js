Cypress.Commands.add('uploadVideo', (fileName) => {
  cy.get('input[type="file"]').selectFile(`cypress/fixtures/${fileName}`, { force: true })
})

Cypress.Commands.add('waitForProcessing', (timeout = 30000) => {
  cy.get('#loading', { timeout }).should('be.visible')
  cy.get('#result', { timeout }).should('be.visible')
  cy.get('#result').should('not.contain', 'Processando')
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
})

Cypress.Commands.add('waitForFileInListing', (fileName) => {
  cy.get('#filesList').should('contain', fileName)
})
