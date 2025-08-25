import { Card, CardContent, Typography, useTheme } from "@mui/material";
import { PieArc, PieChart } from "@mui/x-charts";

import CustomChartMark from "@/components/charts/mark/custom-chart-mark";
import CustomChartTooltip from "@/components/charts/tooltip/custom-chart-tooltip";
import { withChartElementStyle } from "@/lib/chart-element-hoc";

export default function DashboardAnalyticsCategories() {
  const { palette } = useTheme();

  return (
    <Card className="h-96">
      <CardContent>
        <Typography variant="h5" component="h5" className="card-title">
          Inventory by Category
        </Typography>
        <PieChart
          series={[
            {
              innerRadius: 16,
              paddingAngle: 4,
              cornerRadius: 4,
              data: [
                { label: "Produce", value: 40, color: palette.primary.main, labelMarkType: CustomChartMark },
                { label: "Proteins", value: 30, color: palette.secondary.main, labelMarkType: CustomChartMark },
                { label: "Grains", value: 20, color: palette["accent-1"].main, labelMarkType: CustomChartMark },
                { label: "Dairy", value: 15, color: palette["accent-2"].main, labelMarkType: CustomChartMark },
              ],
            },
          ]}
          height={270}
          slots={{ pieArc: withChartElementStyle(PieArc), tooltip: CustomChartTooltip }}
          slotProps={{ legend: { direction: "horizontal", position: { horizontal: "center", vertical: "bottom" } } }}
        />
      </CardContent>
    </Card>
  );
}