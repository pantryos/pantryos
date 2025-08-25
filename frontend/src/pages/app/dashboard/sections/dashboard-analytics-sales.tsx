import { useMemo, useState } from "react";
import { Box, Button, Card, CardContent, Typography, useTheme } from "@mui/material";
import { BarChart, BarElement, BarSeriesType } from "@mui/x-charts";
import CustomChartTooltip from "@/components/charts/tooltip/custom-chart-tooltip";
import NiUser from "@/icons/nexture/ni-user";
import NiUsers from "@/icons/nexture/ni-users";
import { withChartElementStyle } from "@/lib/chart-element-hoc";

export default function DashboardAnalyticsSales() {
  const { palette } = useTheme();
  const [activeIndex, setActiveIndex] = useState<number>(0);

  const series = useMemo(() => {
    const john: Omit<BarSeriesType, "type"> = { data: [280, 340, 310, 355, 295, 320, 305], label: "John D.", stack: "all", color: palette.primary.main };
    const jane: Omit<BarSeriesType, "type"> = { data: [180, 230, 170, 220, 185, 210, 190], label: "Jane S.", stack: "all", color: palette.secondary.main };
    const mike: Omit<BarSeriesType, "type"> = { data: [140, 120, 160, 130, 120, 115, 110], label: "Mike R.", stack: "all", color: palette["accent-1"].main };

    switch (activeIndex) {
      case 1: return [john];
      case 2: return [jane];
      case 3: return [mike];
      case 0: default: return [john, jane, mike];
    }
  }, [activeIndex, palette]);

  return (
    <Card className="h-96">
      <CardContent>
        <Typography variant="h5" component="h5" className="card-title">
          Sales by Employee
        </Typography>
        <Box className="mb-1 flex gap-1">
          <Button variant="outlined" size="medium" startIcon={<NiUsers size="medium" />} color={activeIndex === 0 ? "primary" : "grey"} onClick={() => setActiveIndex(0)}>All</Button>
          <Button variant="outlined" size="medium" startIcon={<NiUser size="medium" />} color={activeIndex === 1 ? "primary" : "grey"} onClick={() => setActiveIndex(1)}>John</Button>
          <Button variant="outlined" size="medium" startIcon={<NiUser size="medium" />} color={activeIndex === 2 ? "primary" : "grey"} onClick={() => setActiveIndex(2)}>Jane</Button>
          <Button variant="outlined" size="medium" startIcon={<NiUser size="medium" />} color={activeIndex === 3 ? "primary" : "grey"} onClick={() => setActiveIndex(3)}>Mike</Button>
        </Box>
        <BarChart
          series={series}
          xAxis={[{ disableLine: true, disableTicks: true, data: ["Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"], scaleType: "band", categoryGapRatio: 0.5, }]}
          yAxis={[{ disableLine: true, disableTicks: true, valueFormatter: (value: number) => `$${value}`, width: 40 }]}
          slots={{ tooltip: CustomChartTooltip, bar: withChartElementStyle(BarElement, { rx: 10, ry: 10 }) }}
          height={250}
          hideLegend
          grid={{ horizontal: true }}
          axisHighlight={{ x: "line" }}
          margin={{ bottom: 0, left: 0, right: 0 }}
        />
      </CardContent>
    </Card>
  );
}