import { Box, Card, CardContent, CircularProgress, Grid, Typography } from "@mui/material";

export default function DashboardAnalyticsProgresses() {
  return (
    <Grid container size={12} spacing={2.5}>
      <Grid size={{ lg: 6, xs: 12 }}>
                <Card className="h-24">
          <CardContent className="flex flex-row items-start justify-between">
            <Box>
              <Typography variant="subtitle2" className="text-text-secondary leading-5 transition-colors">Inventory Accuracy</Typography>
              <Typography variant="h5" className="text-leading-5">Good</Typography>
            </Box>
            <Box className="relative inline-flex w-10">
              <CircularProgress variant="determinate" value={96} className="relative z-1 h-10! w-10! text-info" />
              <Box className="absolute top-0 right-0 bottom-0 left-0 flex items-center justify-center"><Typography variant="caption" component="div" className="text-text-secondary">100%</Typography></Box>
              <Box className="outline-grey-100 absolute top-0 right-0 bottom-0 left-0 z-0 rounded-full outline -outline-offset-2"></Box>
            </Box>
          </CardContent>
        </Card>
        {/* <Card className="h-24">
          <CardContent className="flex flex-row items-start justify-between">
            <Box>
              <Typography variant="subtitle2" className="text-text-secondary leading-5 transition-colors">Customer Satisfaction</Typography>
              <Typography variant="h5" className="text-leading-5">Excellent</Typography>
            </Box>
            <Box className="relative inline-flex w-10">
              <CircularProgress variant="determinate" value={92} className="relative z-1 h-10! w-10! text-success" />
              <Box className="absolute top-0 right-0 bottom-0 left-0 flex items-center justify-center"><Typography variant="caption" component="div" className="text-text-secondary">92%</Typography></Box>
              <Box className="outline-grey-100 absolute top-0 right-0 bottom-0 left-0 z-0 rounded-full outline -outline-offset-2"></Box>
            </Box>
          </CardContent>
        </Card> */}
      </Grid>
      <Grid size={{ lg: 6, xs: 12 }}>
        <Card className="h-24">
          {/* <CardContent className="flex flex-row items-start justify-between">
            <Box>
              <Typography variant="subtitle2" className="text-text-secondary leading-5 transition-colors">POS Uptime</Typography>
              <Typography variant="h5" className="text-leading-5">Stable</Typography>
            </Box>
            <Box className="relative inline-flex w-10">
              <CircularProgress variant="determinate" value={99} className="relative z-1 h-10! w-10! text-success" />
              <Box className="absolute top-0 right-0 bottom-0 left-0 flex items-center justify-center"><Typography variant="caption" component="div" className="text-text-secondary">99%</Typography></Box>
              <Box className="outline-grey-100 absolute top-0 right-0 bottom-0 left-0 z-0 rounded-full outline -outline-offset-2"></Box>
            </Box>
          </CardContent> */}
        </Card>
      </Grid>
      <Grid size={{ lg: 6, xs: 12 }}>

      </Grid>
      <Grid size={{ lg: 6, xs: 12 }}>
        <Card className="h-24">
          {/* <CardContent className="flex flex-row items-start justify-between">
            <Box>
              <Typography variant="subtitle2" className="text-text-secondary leading-5 transition-colors">Staff on Duty</Typography>
              <Typography variant="h5" className="text-leading-5">4 / 5</Typography>
            </Box>
            <Box className="relative inline-flex w-10">
              <CircularProgress variant="determinate" value={80} className="relative z-1 h-10! w-10! text-warning" />
              <Box className="absolute top-0 right-0 bottom-0 left-0 flex items-center justify-center"><Typography variant="caption" component="div" className="text-text-secondary">80%</Typography></Box>
              <Box className="outline-grey-100 absolute top-0 right-0 bottom-0 left-0 z-0 rounded-full outline -outline-offset-2"></Box>
            </Box>
          </CardContent> */}
        </Card>
      </Grid>
    </Grid>
  );
}