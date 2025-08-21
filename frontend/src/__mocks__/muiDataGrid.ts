import React from 'react';

// Complete mock for MUI X Data Grid components
jest.mock('@mui/x-data-grid', () => {
  const React = require('react');
  
  return {
    DataGrid: React.forwardRef((props: any, ref: any) => {
      // THE FIX: Destructure all props that are NOT valid for a <div>
      const { 
        rows = [], 
        columns = [], 
        onRowClick, 
        loading,
        // --- ADD ALL MUI-SPECIFIC PROPS USED IN YOUR APP HERE ---
        initialState,
        pageSizeOptions,
        disableRowSelectionOnClick,
        checkboxSelection,
        slots,
        slotProps,
        // etc.
        // ---------------------------------------------------------
        ...otherProps // Now, otherProps only contains valid DOM attributes
      } = props;
      
      if (loading) {
        return React.createElement('div', { 'data-testid': 'data-grid-loading' }, 'Loading...');
      }
      
      return React.createElement(
        'div',
        { 
          'data-testid': 'data-grid',
          ref,
          ...otherProps // It is now safe to spread this
        },
        React.createElement('div', { 'data-testid': 'data-grid-header' }, 
          columns.map((col: any) => 
            React.createElement('span', { key: col.field }, col.headerName || col.field)
          )
        ),
        React.createElement('div', { 'data-testid': 'data-grid-rows' },
          rows.map((row: any, index: number) => 
            React.createElement(
              'div', 
              { 
                key: row.id || index,
                'data-testid': `data-grid-row-${row.id || index}`,
                onClick: () => onRowClick?.(row)
              },
              columns.map((col: any) => 
                React.createElement('span', { 
                  key: `${row.id}-${col.field}`,
                  'data-testid': `cell-${col.field}`
                }, row[col.field])
              )
            )
          )
        )
      );
    }),
    
    GridToolbar: (props: any) => React.createElement('div', { 'data-testid': 'grid-toolbar', ...props }),
    
    GridActionsCellItem: (props: any) => 
      React.createElement('button', { 
        'data-testid': 'grid-action-item',
        onClick: props.onClick,
        ...props 
      }, props.label || props.children),
      
    gridClasses: {
      root: 'MuiDataGrid-root',
      cell: 'MuiDataGrid-cell',
      row: 'MuiDataGrid-row',
    },
    
    useGridApiRef: () => ({
      current: {
        updateRows: jest.fn(),
        getRowModels: jest.fn(() => new Map()),
        setRows: jest.fn(),
      }
    }),
  };
});

// Mock MUI X internals that cause issues
jest.mock('@mui/x-internals/hash/hash.js', () => ({
  xxh: jest.fn(() => 'mocked-hash-value'),
}));

jest.mock('@mui/x-data-grid/material/variables.js', () => ({
  useMaterialCSSVariables: jest.fn(() => ({})),
}));