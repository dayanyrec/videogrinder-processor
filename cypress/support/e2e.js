import './commands'

// Global before hook - only for the main test suite
// The individual tests will handle their own visits

// Global configuration
Cypress.on('uncaught:exception', (err, runnable) => {
  // Prevent Cypress from failing on uncaught exceptions
  // that might come from the application
  if (err.message.includes('ResizeObserver')) {
    return false
  }

  // Let other exceptions fail the test
  return true
})
