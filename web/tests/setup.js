global.fetch = jest.fn()

global.console = {
  ...console
}

beforeEach(() => {
  document.body.innerHTML = ''
  document.head.innerHTML = ''

  fetch.mockClear()

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
  jest.clearAllMocks()
})
