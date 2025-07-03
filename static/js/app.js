/* global AppController */
let appController

document.addEventListener('DOMContentLoaded', () => {
  appController = new AppController()
  appController.init()
})

if (typeof module !== 'undefined' && module.exports) {
  module.exports = { appController }
}
