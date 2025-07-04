const Utils = require('../static/js/utils.js')

describe('Utils Class', () => {
  describe('formatFileSize', () => {
    test('should convert bytes to kilobytes with proper rounding', () => {
      const testCases = [
        { bytes: 1024, expectedKB: 1 },
        { bytes: 1536, expectedKB: 2 },
        { bytes: 1024000, expectedKB: 1000 },
        { bytes: 512, expectedKB: 1 }
      ]

      testCases.forEach(({ bytes, expectedKB }) => {
        const formattedSize = Utils.formatFileSize(bytes)
        expect(formattedSize).toBe(expectedKB)
      })
    })

    test('should return zero when input is zero bytes', () => {
      expect(Utils.formatFileSize(0)).toBe(0)
    })

    test('should handle negative byte values correctly', () => {
      expect(Utils.formatFileSize(-1024)).toBe(-1)
    })
  })

  describe('createDownloadLink', () => {
    test('should generate HTML anchor element with correct download URL and styling', () => {
      const zipPath = 'frames_123456.zip'
      const link = Utils.createDownloadLink(zipPath)

      expect(link).toContain('/download/frames_123456.zip')
      expect(link).toContain('class="download-btn"')
      expect(link).toContain('📥 Baixar ZIP')
    })

    test('should create valid HTML structure even with empty zip path', () => {
      const link = Utils.createDownloadLink('')
      expect(link).toContain('/download/')
      expect(link).toContain('class="download-btn"')
    })
  })

  describe('validateFile', () => {
    test('should return valid result when file object is provided', () => {
      const mockFile = new File(['content'], 'test.mp4', { type: 'video/mp4' })
      const result = Utils.validateFile(mockFile)

      expect(result.valid).toBe(true)
      expect(result.message).toBe('')
    })

    test('should return invalid result with error message when file is null', () => {
      const result = Utils.validateFile(null)

      expect(result.valid).toBe(false)
      expect(result.message).toBe('Por favor, selecione um arquivo.')
    })

    test('should return invalid result with error message when file is undefined', () => {
      const result = Utils.validateFile(undefined)

      expect(result.valid).toBe(false)
      expect(result.message).toBe('Por favor, selecione um arquivo.')
    })

    test('should return invalid result with error message when file is empty string', () => {
      const result = Utils.validateFile('')

      expect(result.valid).toBe(false)
      expect(result.message).toBe('Por favor, selecione um arquivo.')
    })
  })
})
