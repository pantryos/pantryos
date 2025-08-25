import { Card, CardContent, Typography, useTheme } from "@mui/material";
import { RadarChart } from "@mui/x-charts";
import CustomChartMark from "@/components/charts/mark/custom-chart-mark";
import CustomChartTooltip from "@/components/charts/tooltip/custom-chart-tooltip";

export default function DashboardAnalyticsOrdersStocks() {
  const { palette } = useTheme();

  return (
    <Card className="h-96">
      <CardContent>
        <Typography variant="h5" component="h5" className="card-title">
          Top Product Performance
        </Typography>
        <RadarChart
          series={[
            { label: "Espresso Machine", data: [85, 75, 90], color: palette.primary.main, fillArea: true, labelMarkType: CustomChartMark },
            { label: "Coffee Beans", data: [95, 60, 70], color: palette.secondary.main, fillArea: true, labelMarkType: CustomChartMark },
          ]}
          className="radar-chart order-first min-w-50"
          shape="sharp"
          radar={{ labelGap: 6, max: 100, metrics: ["Sales Volume", "Profit Margin", "Stock Level"] }}
          divisions={3}
          height={270}
          margin={{ left: 30, right: 34 }}
          slots={{ tooltip: CustomChartTooltip }}
          stripeColor={null}
        />
      </CardContent>
    </Card>
  );
}