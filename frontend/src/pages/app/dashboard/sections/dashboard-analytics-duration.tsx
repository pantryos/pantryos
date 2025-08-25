import { useMemo, useState } from "react";
import { Box, Button, Card, CardContent, Typography, useTheme } from "@mui/material";
import { LineChart, LineSeriesType } from "@mui/x-charts";
import CustomChartTooltip from "@/components/charts/tooltip/custom-chart-tooltip";
import NiCreditCard from "@/icons/nexture/ni-credit-card";
import NiWallet from "@/icons/nexture/ni-wallet";
import NiMoney from "@/icons/nexture/ni-money";
import NiStructure from "@/icons/nexture/ni-structure";
import { colorWithOpacity } from "@/lib/chart-helper";

export default function DashboardAnalyticsDuration() {
  const { palette } = useTheme();
  const [activeIndex, setActiveIndex] = useState<number>(0);

  const series = useMemo(() => {
    const creditCard: Omit<LineSeriesType, "type"> = { data: [1200, 1500, 1400, 1650, 1550, 1800, 1750], label: "Credit Card", stack: "all", area: true, showMark: false, color: palette.primary.main, curve: "bumpX" };
    const cash: Omit<LineSeriesType, "type"> = { data: [450, 500, 480, 520, 500, 550, 530], label: "Cash", stack: "all", area: true, showMark: false, color: palette.secondary.main, curve: "bumpX" };
    const digitalWallet: Omit<LineSeriesType, "type"> = { data: [300, 350, 320, 400, 380, 420, 410], label: "Digital Wallet", stack: "all", area: true, showMark: false, color: palette["accent-2"].main, curve: "bumpX" };

    switch (activeIndex) {
      case 1: return [creditCard];
      case 2: return [cash];
      case 3: return [digitalWallet];
      case 0: default: return [creditCard, cash, digitalWallet];
    }
  }, [activeIndex, palette]);

  return (
    <Card className="h-96">
      <CardContent>
        <Typography variant="h5" component="h5" className="card-title">
          Revenue by Payment Method
        </Typography>
        <Box className="mb-1 flex gap-1">
          <Button variant="outlined" size="medium" startIcon={<NiStructure size="medium" />} color={activeIndex === 0 ? "primary" : "grey"} onClick={() => setActiveIndex(0)}>All</Button>
          <Button variant="outlined" size="medium" startIcon={<NiCreditCard size="medium" />} color={activeIndex === 1 ? "primary" : "grey"} onClick={() => setActiveIndex(1)}>Card</Button>
          <Button variant="outlined" size="medium" startIcon={<NiMoney size="medium" />} color={activeIndex === 2 ? "primary" : "grey"} onClick={() => setActiveIndex(2)}>Cash</Button>
          <Button variant="outlined" size="medium" startIcon={<NiWallet size="medium" />} color={activeIndex === 3 ? "primary" : "grey"} onClick={() => setActiveIndex(3)}>Wallet</Button>
        </Box>
        <LineChart
          series={series}
          xAxis={[{ disableLine: true, disableTicks: true, data: ["Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"], scaleType: "band" }]}
          yAxis={[{ disableLine: true, disableTicks: true, valueFormatter: (value: number) => `$${value/1000}k`, width: 40 }]}
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