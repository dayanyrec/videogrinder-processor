// Jest setup file for DOM testing
// This file is executed before each test file

// Mock fetch globally
global.fetch = jest.fn()

// Mock console methods to reduce noise in tests
global.console = {
  ...console
  // Uncomment to ignore specific console methods during tests
  // log: jest.fn(),
  // warn: jest.fn(),
  // error: jest.fn(),
}

// Setup DOM
beforeEach(() => {
  // Reset DOM
  document.body.innerHTML = ''
  document.head.innerHTML = ''

  // Reset fetch mock
  fetch.mockClear()

  // Create basic HTML structure for tests
  document.body.innerHTML = `
    <form id="uploadForm">
      <input type="file" id="videoFile" />
      <button type="submit">Upload</button>
    </form>
    <div id="loading" style="display: none;">Loading...</div>
    <div id="result" style="display: none;"></div>
    <div id="filesList"></div>
  `
})

afterEach(() => {
  // Clean up after each test
  jest.clearAllMocks()
})
