import {
  BusinessType,
  isStandaloneAccount,
  isMultiLocationAccount,
  isEnterpriseAccount,
  getAccountDisplayName,
  Account,
} from '../api';

describe('API Types and Helpers', () => {
  describe('BusinessType constants', () => {
    it('should have correct business type values', () => {
      expect(BusinessType.SINGLE_LOCATION).toBe('single_location');
      expect(BusinessType.MULTI_LOCATION).toBe('multi_location');
      expect(BusinessType.ENTERPRISE).toBe('enterprise');
    });
  });

  describe('Account helper functions', () => {
    const standaloneAccount: Account = {
      id: 1,
      name: 'Coffee Shop',
      location: '123 Main St',
      phone: '555-0123',
      email: 'coffee@example.com',
      business_type: BusinessType.SINGLE_LOCATION,
      status: 'active',
      created_at: '2024-01-01T00:00:00Z',
      updated_at: '2024-01-01T00:00:00Z',
    };

    const multiLocationAccount: Account = {
      id: 2,
      organization_id: 1,
      name: 'Coffee Chain',
      location: '456 Oak Ave',
      phone: '555-0456',
      email: 'chain@example.com',
      business_type: BusinessType.MULTI_LOCATION,
      status: 'active',
      created_at: '2024-01-01T00:00:00Z',
      updated_at: '2024-01-01T00:00:00Z',
    };

    const enterpriseAccount: Account = {
      id: 3,
      organization_id: 2,
      name: 'Enterprise Coffee',
      location: '789 Pine St',
      phone: '555-0789',
      email: 'enterprise@example.com',
      business_type: BusinessType.ENTERPRISE,
      status: 'active',
      created_at: '2024-01-01T00:00:00Z',
      updated_at: '2024-01-01T00:00:00Z',
    };

    describe('isStandaloneAccount', () => {
      it('should return true for standalone account', () => {
        expect(isStandaloneAccount(standaloneAccount)).toBe(true);
      });

      it('should return false for organization-based account', () => {
        expect(isStandaloneAccount(multiLocationAccount)).toBe(false);
        expect(isStandaloneAccount(enterpriseAccount)).toBe(false);
      });

      it('should handle account with undefined organization_id', () => {
        const accountWithUndefined = { ...standaloneAccount, organization_id: undefined };
        expect(isStandaloneAccount(accountWithUndefined)).toBe(true);
      });
    });

    describe('isMultiLocationAccount', () => {
      it('should return true for multi-location account', () => {
        expect(isMultiLocationAccount(multiLocationAccount)).toBe(true);
      });

      it('should return false for other business types', () => {
        expect(isMultiLocationAccount(standaloneAccount)).toBe(false);
        expect(isMultiLocationAccount(enterpriseAccount)).toBe(false);
      });
    });

    describe('isEnterpriseAccount', () => {
      it('should return true for enterprise account', () => {
        expect(isEnterpriseAccount(enterpriseAccount)).toBe(true);
      });

      it('should return false for other business types', () => {
        expect(isEnterpriseAccount(standaloneAccount)).toBe(false);
        expect(isEnterpriseAccount(multiLocationAccount)).toBe(false);
      });
    });

    describe('getAccountDisplayName', () => {
      it('should return just name for standalone account', () => {
        expect(getAccountDisplayName(standaloneAccount)).toBe('Coffee Shop');
      });

      it('should return name and location for organization-based account', () => {
        expect(getAccountDisplayName(multiLocationAccount)).toBe('Coffee Chain (456 Oak Ave)');
        expect(getAccountDisplayName(enterpriseAccount)).toBe('Enterprise Coffee (789 Pine St)');
      });

      it('should handle account with empty location', () => {
        const accountWithEmptyLocation = { ...multiLocationAccount, location: '' };
        expect(getAccountDisplayName(accountWithEmptyLocation)).toBe('Coffee Chain ()');
      });
    });
  });
}); 