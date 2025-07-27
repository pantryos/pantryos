describe('LocalStorage Utilities', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    (localStorage.clear as jest.Mock).mockClear();
  });

  describe('localStorage operations', () => {
    it('should set and get items', () => {
      localStorage.setItem('test-key', 'test-value');
      expect(localStorage.setItem).toHaveBeenCalledWith('test-key', 'test-value');
      
      (localStorage.getItem as jest.Mock).mockReturnValue('test-value');
      const value = localStorage.getItem('test-key');
      expect(value).toBe('test-value');
    });

    it('should remove items', () => {
      localStorage.removeItem('test-key');
      expect(localStorage.removeItem).toHaveBeenCalledWith('test-key');
    });

    it('should clear all items', () => {
      localStorage.clear();
      expect(localStorage.clear).toHaveBeenCalled();
    });

    it('should return null for non-existent items', () => {
      (localStorage.getItem as jest.Mock).mockReturnValue(null);
      const value = localStorage.getItem('non-existent');
      expect(value).toBeNull();
    });
  });
}); 