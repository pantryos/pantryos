// Test utility functions
describe('Utility Functions', () => {
  describe('String utilities', () => {
    it('should capitalize first letter', () => {
      const capitalize = (str: string) => str.charAt(0).toUpperCase() + str.slice(1);
      
      expect(capitalize('hello')).toBe('Hello');
      expect(capitalize('world')).toBe('World');
      expect(capitalize('')).toBe('');
    });

    it('should format currency', () => {
      const formatCurrency = (amount: number) => `$${amount.toFixed(2)}`;
      
      expect(formatCurrency(10)).toBe('$10.00');
      expect(formatCurrency(10.5)).toBe('$10.50');
      expect(formatCurrency(10.99)).toBe('$10.99');
    });

    it('should truncate long strings', () => {
      const truncate = (str: string, maxLength: number) => 
        str.length > maxLength ? str.slice(0, maxLength) + '...' : str;
      
      expect(truncate('Hello World', 5)).toBe('Hello...');
      expect(truncate('Short', 10)).toBe('Short');
      expect(truncate('', 5)).toBe('');
    });
  });

  describe('Date utilities', () => {
    it('should format date', () => {
      const formatDate = (date: Date) => date.toLocaleDateString();
      
      const testDate = new Date('2024-01-15');
      const formatted = formatDate(testDate);
      expect(formatted).toMatch(/^\d{1,2}\/\d{1,2}\/\d{4}$/);
    });

    it('should check if date is today', () => {
      const isToday = (date: Date) => {
        const today = new Date();
        return date.toDateString() === today.toDateString();
      };
      
      expect(isToday(new Date())).toBe(true);
      expect(isToday(new Date('2024-01-15'))).toBe(false);
    });

    it('should get days between dates', () => {
      const getDaysBetween = (date1: Date, date2: Date) => {
        const oneDay = 24 * 60 * 60 * 1000;
        return Math.round(Math.abs((date1.getTime() - date2.getTime()) / oneDay));
      };
      
      const date1 = new Date('2024-01-15');
      const date2 = new Date('2024-01-20');
      expect(getDaysBetween(date1, date2)).toBe(5);
    });
  });

  describe('Array utilities', () => {
    it('should remove duplicates', () => {
      const removeDuplicates = <T>(arr: T[]) => Array.from(new Set(arr));
      
      expect(removeDuplicates([1, 2, 2, 3, 3, 4])).toEqual([1, 2, 3, 4]);
      expect(removeDuplicates(['a', 'b', 'a', 'c'])).toEqual(['a', 'b', 'c']);
    });

    it('should group by property', () => {
      const groupBy = <T>(arr: T[], key: keyof T) => {
        return arr.reduce((groups, item) => {
          const group = String(item[key]);
          groups[group] = groups[group] || [];
          groups[group].push(item);
          return groups;
        }, {} as Record<string, T[]>);
      };
      
      const items = [
        { id: 1, category: 'A' },
        { id: 2, category: 'B' },
        { id: 3, category: 'A' },
      ];
      
      const grouped = groupBy(items, 'category');
      expect(grouped.A).toHaveLength(2);
      expect(grouped.B).toHaveLength(1);
    });

    it('should sort by property', () => {
      const sortBy = <T>(arr: T[], key: keyof T, ascending = true) => {
        return [...arr].sort((a, b) => {
          const aVal = a[key];
          const bVal = b[key];
          if (aVal < bVal) return ascending ? -1 : 1;
          if (aVal > bVal) return ascending ? 1 : -1;
          return 0;
        });
      };
      
      const items = [
        { id: 3, name: 'Charlie' },
        { id: 1, name: 'Alice' },
        { id: 2, name: 'Bob' },
      ];
      
      const sortedById = sortBy(items, 'id');
      expect(sortedById[0].id).toBe(1);
      expect(sortedById[2].id).toBe(3);
      
      const sortedByName = sortBy(items, 'name');
      expect(sortedByName[0].name).toBe('Alice');
      expect(sortedByName[2].name).toBe('Charlie');
    });
  });

  describe('Validation utilities', () => {
    it('should validate email', () => {
      const isValidEmail = (email: string) => {
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        return emailRegex.test(email);
      };
      
      expect(isValidEmail('test@example.com')).toBe(true);
      expect(isValidEmail('invalid-email')).toBe(false);
      expect(isValidEmail('test@')).toBe(false);
      expect(isValidEmail('@example.com')).toBe(false);
    });

    it('should validate password strength', () => {
      const isStrongPassword = (password: string) => {
        return password.length >= 8 && 
               /[A-Z]/.test(password) && 
               /[a-z]/.test(password) && 
               /\d/.test(password);
      };
      
      expect(isStrongPassword('StrongPass123')).toBe(true);
      expect(isStrongPassword('weak')).toBe(false);
      expect(isStrongPassword('nouppercase123')).toBe(false);
      expect(isStrongPassword('NOLOWERCASE123')).toBe(false);
      expect(isStrongPassword('NoNumbers')).toBe(false);
    });

    it('should validate required fields', () => {
      const validateRequired = (obj: Record<string, any>, fields: string[]) => {
        const errors: string[] = [];
        fields.forEach(field => {
          if (!obj[field] || obj[field].toString().trim() === '') {
            errors.push(`${field} is required`);
          }
        });
        return errors;
      };
      
      const data = { name: 'John', email: '', age: 0 };
      const errors = validateRequired(data, ['name', 'email', 'age']);
      
      expect(errors).toContain('email is required');
      expect(errors).not.toContain('name is required');
      expect(errors).toContain('age is required'); // 0 is falsy in JavaScript
    });
  });
}); 