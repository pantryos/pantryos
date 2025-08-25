import dayjs from "dayjs";
import { Box, Card, CardContent, Grid, Typography, useTheme } from "@mui/material";
import { SparkLineChart } from "@mui/x-charts";
import useHighlightedSparkline from "@/hooks/use-highlighted-sparkline";
import NiTriangleDown from "@/icons/nexture/ni-triangle-down";
import NiTriangleUp from "@/icons/nexture/ni-triangle-up";

const revenueData = [1800, 1950, 2100, 2050, 2200, 2300, 2250, 2400, 2350.50];
const profitData = [450, 480, 510, 500, 550, 580, 570, 610, 605.25];
const avgTransactionData = [25.50, 26.20, 25.80, 27.10, 26.90, 28.00, 27.50, 28.20, 27.95];

export default function DashboardAnalyticsCurrencies() {
  const { palette } = useTheme();

  const revenueSparkline = useHighlightedSparkline({ data: revenueData, plotType: "line", color: palette.primary.main });
  const profitSparkline = useHighlightedSparkline({ data: profitData, plotType: "line", color: palette.primary.main });
  const avgTransactionSparkline = useHighlightedSparkline({ data: avgTransactionData, plotType: "line", color: palette.primary.main });

  return (
    <Grid container size={12} spacing={2.5}>
      <Grid size={{ xs: 12 }}>
        <Card className="h-24">
          <CardContent className="flex items-center gap-5">
            <Box className="flex-shrink-0">
              <Typography variant="body1" className="w-54">
                Total Revenue
              </Typography>
              <Box className="flex items-center gap-1">
                <Typography variant="h5" className="text-text-primary">
                  ${revenueSparkline.item.value.toLocaleString()}
                </Typography>
                <ChangeStatus change={revenueSparkline.item.change} />
              </Box>
            </Box>
            <Box className="-my-1 grow">
              <SparkLineChart {...revenueSparkline.props} />
            </Box>
          </CardContent>
        </Card>
      </Grid>
      <Grid size={{ xs: 12 }}>
        <Card className="h-24">
          <CardContent className="flex items-center gap-5">
            <Box className="flex-shrink-0">
              <Typography variant="body1" className="w-54">
                Net Profit
              </Typography>
              <Box className="flex items-center gap-1">
                <Typography variant="h5" className="text-text-primary">
                  ${profitSparkline.item.value.toLocaleString()}
                </Typography>
                <ChangeStatus change={profitSparkline.item.change} />
              </Box>
            </Box>
            <Box className="-my-1 grow">
              <SparkLineChart {...profitSparkline.props} />
            </Box>
          </CardContent>
        </Card>
      </Grid>
      <Grid size={{ xs: 12 }}>
        <Card className="h-24">
          <CardContent className="flex items-center gap-5">
            <Box className="flex-shrink-0">
              <Typography variant="body1" className="w-54">
                Avg. Transaction
              </Typography>
              <Box className="flex items-center gap-1">
                <Typography variant="h5" className="text-text-primary">
                  ${avgTransactionSparkline.item.value.toFixed(2)}
                </Typography>
                <ChangeStatus change={avgTransactionSparkline.item.change} />
              </Box>
            </Box>
            <Box className="-my-1 grow">
              <SparkLineChart {...avgTransactionSparkline.props} />
            </Box>
          </CardContent>
        </Card>
      </Grid>
    </Grid>
  );
}

const ChangeStatus = ({ change }: { change: number | string }) => {
  return (
    <Box className="flex">
      {Number(change) < 0 ? (
        <NiTriangleDown size="tiny" className="text-error" />
      ) : (
        <NiTriangleUp size="tiny" className="text-success" />
      )}
      <Typography variant="body2" className={Number(change) < 0 ? "text-error" : "text-success"}>
        {Math.abs(Number(change))}%
      </Typography>
    </Box>
  );
};