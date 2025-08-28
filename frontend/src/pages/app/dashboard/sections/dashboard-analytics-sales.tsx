import React, { useMemo, useState, useEffect } from 'react';

// Define the Delivery type
interface Delivery {
  id: number;
  inventory_item_id: number;
  vendor: string;
  quantity: number;
  cost: number;
  delivery_date: string;
}

// Mock API service
const mockApiService = {
  getDeliveries: (): Promise<Delivery[]> => {
    return new Promise(resolve => {
      setTimeout(() => {
        const mockData: Delivery[] = [
          { id: 1, inventory_item_id: 101, vendor: 'Fresh Produce Co.', quantity: 50, cost: 125.50, delivery_date: '2025-08-25T10:00:00Z' },
          { id: 2, inventory_item_id: 102, vendor: 'Bakery Supplies Inc.', quantity: 20, cost: 250.00, delivery_date: '2025-08-25T11:30:00Z' },
          { id: 3, inventory_item_id: 103, vendor: 'Dairy Farms Ltd.', quantity: 100, cost: 300.75, delivery_date: '2025-08-26T09:00:00Z' },
          { id: 4, inventory_item_id: 101, vendor: 'Fresh Produce Co.', quantity: 70, cost: 175.00, delivery_date: '2025-08-27T10:00:00Z' },
          { id: 5, inventory_item_id: 104, vendor: 'Meat Wholesalers', quantity: 30, cost: 550.25, delivery_date: '2025-08-28T14:00:00Z' },
          { id: 6, inventory_item_id: 102, vendor: 'Bakery Supplies Inc.', quantity: 25, cost: 312.50, delivery_date: '2025-08-29T11:30:00Z' },
          { id: 7, inventory_item_id: 103, vendor: 'Dairy Farms Ltd.', quantity: 120, cost: 360.90, delivery_date: '2025-08-30T09:00:00Z' },
          { id: 8, inventory_item_id: 104, vendor: 'Meat Wholesalers', quantity: 40, cost: 730.00, delivery_date: '2025-08-24T14:00:00Z' },
        ];
        resolve(mockData);
      }, 1000);
    });
  }
};

export default function DashboardDeliveriesChart() {
  const [activeVendor, setActiveVendor] = useState<string>('All');
  const [deliveries, setDeliveries] = useState<Delivery[]>([]);
  const [loading, setLoading] = useState(true);
  const [hoveredBar, setHoveredBar] = useState<{day: number, vendor: string} | null>(null);

  useEffect(() => {
    const fetchData = async () => {
      try {
        setLoading(true);
        const deliveriesData = await mockApiService.getDeliveries();
        setDeliveries(deliveriesData || []);
      } catch (error) {
        console.error("Failed to fetch deliveries for chart:", error);
      } finally {
        setLoading(false);
      }
    };
    fetchData();
  }, []);

  const vendors = useMemo(() => {
    const uniqueVendors = [...new Set(deliveries.map(d => d.vendor))];
    return ['All', ...uniqueVendors];
  }, [deliveries]);

  const chartData = useMemo(() => {
    if (!deliveries || deliveries.length === 0) {
      return { labels: [], datasets: [] };
    }
    
    const uniqueVendors = vendors.slice(1);
    const labels = ["Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"];
    const vendorData: { [key: string]: number[] } = uniqueVendors.reduce((acc, vendor) => {
      acc[vendor] = Array(7).fill(0);
      return acc;
    }, {} as { [key: string]: number[] });

    deliveries.forEach(delivery => {
      const deliveryDate = new Date(delivery.delivery_date);
      const dayIndex = deliveryDate.getDay();
      const mappedIndex = dayIndex === 0 ? 6 : dayIndex - 1;
      if (vendorData[delivery.vendor]) {
        vendorData[delivery.vendor][mappedIndex] += delivery.cost;
      }
    });

    const colors = ['#3b82f6', '#10b981', '#f59e0b', '#ef4444', '#8b5cf6', '#06b6d4'];
    const datasets = uniqueVendors.map((vendor, index) => ({
      label: vendor,
      data: vendorData[vendor],
      color: colors[index % colors.length]
    }));

    const visibleDatasets = activeVendor === 'All' ? datasets : datasets.filter(d => d.label === activeVendor);
    return { labels, datasets: visibleDatasets };
  }, [deliveries, activeVendor, vendors]);

  const yMax = useMemo(() => {
    if (activeVendor !== 'All' || chartData.datasets.length === 0) {
      const maxVal = Math.max(...chartData.datasets.flatMap(d => d.data), 0);
      return maxVal > 0 ? maxVal * 1.2 : 100;
    }
    const dailyTotals = chartData.labels.map((_, dayIndex) => 
      chartData.datasets.reduce((sum, dataset) => sum + dataset.data[dayIndex], 0)
    );
    const maxTotal = Math.max(...dailyTotals, 0);
    return maxTotal > 0 ? maxTotal * 1.2 : 100;
  }, [chartData, activeVendor]);

  if (loading) {
    return (
      <div className="w-full h-96 bg-white rounded-lg shadow-sm border border-gray-200 flex items-center justify-center">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  return (
    <div className="w-full h-96 bg-white rounded-lg shadow-sm border border-gray-200 p-6">
      <h2 className="text-xl font-semibold text-gray-800 mb-4">Deliveries by Vendor</h2>
      
      {/* Vendor Filter Buttons */}
      <div className="flex flex-wrap gap-2 mb-6">
        {vendors.map((vendor) => (
          <button
            key={vendor}
            onClick={() => setActiveVendor(vendor)}
            className={`px-3 py-1.5 text-sm rounded-md border transition-colors duration-200 flex items-center gap-1 ${
              activeVendor === vendor
                ? 'bg-blue-50 border-blue-300 text-blue-700'
                : 'bg-gray-50 border-gray-300 text-gray-600 hover:bg-gray-100'
            }`}
          >
            <span className="text-xs">
              {vendor === 'All' ? '🚚' : '🏪'}
            </span>
            {vendor}
          </button>
        ))}
      </div>

      {/* Chart */}
      <div className="h-64 relative">
        <svg width="100%" height="100%" viewBox="0 0 500 250" className="overflow-visible">
          {/* Grid Lines */}
          {[...Array(6)].map((_, i) => (
            <line
              key={i}
              x1="40"
              y1={30 + i * (180 / 5)}
              x2="460"
              y2={30 + i * (180 / 5)}
              stroke="#f3f4f6"
              strokeWidth="1"
            />
          ))}

          {/* Y-Axis Labels */}
          {[...Array(6)].map((_, i) => (
            <text
              key={i}
              x="35"
              y={210 - i * (180 / 5)}
              textAnchor="end"
              fontSize="12"
              fill="#6b7280"
              className="font-medium"
            >
              ${Math.round(yMax / 5 * i)}
            </text>
          ))}

          {/* Bars */}
          {chartData.labels.map((label, dayIndex) => {
            let yOffset = 210;
            const barWidth = 40;
            const columnWidth = 420 / chartData.labels.length;
            const x = 50 + dayIndex * columnWidth + (columnWidth - barWidth) / 2;

            return (
              <g key={label}>
                {chartData.datasets.map((dataset) => {
                  const value = dataset.data[dayIndex];
                  const barHeight = (value / yMax) * 180;
                  
                  if (barHeight <= 0) return null;
                  
                  const currentY = yOffset - barHeight;
                  yOffset -= barHeight;
                  
                  const isHovered = hoveredBar?.day === dayIndex && hoveredBar?.vendor === dataset.label;

                  return (
                    <g key={dataset.label}>
                      <rect
                        x={x}
                        y={currentY}
                        width={barWidth}
                        height={barHeight}
                        fill={`${dataset.color}15`} // Very transparent fill
                        stroke={dataset.color}
                        strokeWidth="2"
                        rx="12"
                        ry="12"
                        className="cursor-pointer transition-all duration-200"
                        style={{
                          filter: isHovered ? 'brightness(1.1)' : 'none',
                          strokeOpacity: 0.7
                        }}
                        onMouseEnter={() => setHoveredBar({day: dayIndex, vendor: dataset.label})}
                        onMouseLeave={() => setHoveredBar(null)}
                      />
                      {isHovered && (
                        <g>
                          <rect
                            x={x - 30}
                            y={currentY - 25}
                            width={barWidth + 60}
                            height="20"
                            fill="#1f2937"
                            rx="4"
                            ry="4"
                            fillOpacity="0.9"
                          />
                          <text
                            x={x + barWidth / 2}
                            y={currentY - 10}
                            textAnchor="middle"
                            fontSize="12"
                            fill="white"
                            className="font-medium"
                          >
                            {dataset.label}: ${value.toFixed(2)}
                          </text>
                        </g>
                      )}
                    </g>
                  );
                })}

                {/* X-Axis Labels */}
                <text
                  x={x + barWidth / 2}
                  y="235"
                  textAnchor="middle"
                  fontSize="12"
                  fill="#6b7280"
                  className="font-medium"
                >
                  {label}
                </text>
              </g>
            );
          })}
        </svg>
      </div>
    </div>
  );
}