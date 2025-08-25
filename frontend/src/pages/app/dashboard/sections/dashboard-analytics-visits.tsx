import { useMemo, useState } from "react";
import { Box, Button, Card, CardContent, Typography, useTheme } from "@mui/material";
import { LineChart, LineSeriesType } from "@mui/x-charts";
import CustomChartTooltip from "@/components/charts/tooltip/custom-chart-tooltip";
import NiScreen from "@/icons/nexture/ni-screen";
import NiCartEmpty from "@/icons/nexture/ni-cart-empty";
import NiPhone from "@/icons/nexture/ni-phone";
import NiStructure from "@/icons/nexture/ni-structure";
import { colorWithOpacity } from "@/lib/chart-helper";

export default function DashboardAnalyticsVisits() {
  const { palette } = useTheme();
  const [activeIndex, setActiveIndex] = useState<number>(0);

  const series = useMemo(() => {
    const inStore: Omit<LineSeriesType, "type"> = { data: [15, 25, 40, 35, 55, 60, 50], label: "In-Store", stack: "all", area: true, showMark: false, color: palette["accent-1"].main, curve: "bumpX" };
    const online: Omit<LineSeriesType, "type"> = { data: [10, 12, 15, 20, 18, 25, 22], label: "Online", stack: "all", area: true, showMark: false, color: palette.secondary.main, curve: "bumpX" };
    const phone: Omit<LineSeriesType, "type"> = { data: [5, 8, 6, 10, 9, 12, 11], label: "Phone", stack: "all", area: true, showMark: false, color: palette["accent-3"].main, curve: "bumpX" };

    switch (activeIndex) {
      case 1: return [inStore];
      case 2: return [online];
      case 3: return [phone];
      case 0: default: return [inStore, online, phone];
    }
  }, [activeIndex, palette]);

  return (
    <Card className="h-96">
      <CardContent>
        <Typography variant="h5" component="h5" className="card-title">
          Customer Traffic & Peak Hours
        </Typography>
        <Box className="mb-1 flex gap-1">
          <Button variant="outlined" size="medium" startIcon={<NiStructure size="medium" />} color={activeIndex === 0 ? "primary" : "grey"} onClick={() => setActiveIndex(0)}>All</Button>
          <Button variant="outlined" size="medium" startIcon={<NiScreen size="medium" />} color={activeIndex === 1 ? "primary" : "grey"} onClick={() => setActiveIndex(1)}>In-Store</Button>
          <Button variant="outlined" size="medium" startIcon={<NiCartEmpty size="medium" />} color={activeIndex === 2 ? "primary" : "grey"} onClick={() => setActiveIndex(2)}>Online</Button>
          <Button variant="outlined" size="medium" startIcon={<NiPhone size="medium" />} color={activeIndex === 3 ? "primary" : "grey"} onClick={() => setActiveIndex(3)}>Phone</Button>
        </Box>
        <LineChart
          series={series}
          xAxis={[{ disableLine: true, disableTicks: true, data: ["9am", "11am", "1pm", "3pm", "5pm", "7pm", "9pm"], scaleType: "band" }]}
          yAxis={[{ disableLine: true, disableTicks: true, width: 30 }]}
          slots={{ tooltip: CustomChartTooltip }}
          slotProps={{ area: ({ color }) => ({ fill: colorWithOpacity(color) }) }}
          height={250}
          hideLegend
          grid={{ horizontal: true }}
          margin={{ bottom: 0, left: 0, right: 0 }}
        />
      </CardContent>
    </Card>
  );
}