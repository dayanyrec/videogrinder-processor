import './commands'

Cypress.on('uncaught:exception', (err, _runnable) => {
  if (err.message.includes('ResizeObserver')) {
    return false
  }
  return true
})
