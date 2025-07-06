/* global AppController */
let appController

function initializeApp() {
  if (appController) return // Already initialized

  appController = new AppController()
  appController.init()

  // Expose for testing
  window.appController = appController
}

// Initialize when DOM is ready
document.addEventListener('DOMContentLoaded', initializeApp)

// For Cypress: Also initialize immediately if DOM is already ready
if (document.readyState === 'loading') {
  // DOM is still loading, wait for DOMContentLoaded
} else {
  // DOM is already ready, initialize immediately
  initializeApp()
}

if (typeof module !== 'undefined' && module.exports) {
  module.exports = { appController }
}
