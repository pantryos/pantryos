import { Box, Card, CardContent, Chip, Link, Typography } from "@mui/material";

export default function DashboardAnalyticsAppStatus() {
  return (
    <Card className="h-96">
      <CardContent>
        <Typography variant="h5" component="h5" className="card-title">
          Sales Target Progress
        </Typography>
        <Box className="flex flex-col gap-5">
          <Box className="flex flex-col gap-2">
            <Box className="flex flex-row items-center justify-between">
              <Link href="#" variant="subtitle2" color="textPrimary" underline="hover">
                Daily Sales Target
              </Link>
              <Chip label="$1,850 / $2,000" variant="outlined" />
            </Box>
            <Box className="bg-grey-50 h-0.5 w-full">
              <Box className="bg-primary h-0.5" style={{ width: "92.5%" }}></Box>
            </Box>
          </Box>
          <Box className="flex flex-col gap-2">
            <Box className="flex flex-row items-center justify-between">
              <Link href="#" variant="subtitle2" color="textPrimary" underline="hover">
                Weekly Sales Target
              </Link>
              <Chip label="$9,700 / $14,000" variant="outlined" />
            </Box>
            <Box className="bg-grey-50 h-0.5 w-full">
              <Box className="bg-primary h-0.5" style={{ width: "69.2%" }}></Box>
            </Box>
          </Box>
          <Box className="flex flex-col gap-2">
            <Box className="flex flex-row items-center justify-between">
              <Link href="#" variant="subtitle2" color="textPrimary" underline="hover">
                Monthly Sales Target
              </Link>
              <Chip label="$45,100 / $55,000" variant="outlined" />
            </Box>
            <Box className="bg-grey-50 h-0.5 w-full">
              <Box className="bg-primary h-0.5" style={{ width: "82%" }}></Box>
            </Box>
          </Box>
          <Box className="flex flex-col gap-2">
            <Box className="flex flex-row items-center justify-between">
              <Link href="#" variant="subtitle2" color="textPrimary" underline="hover">
                New Customer Goal
              </Link>
              <Chip label="32 / 50" variant="outlined" />
            </Box>
            <Box className="bg-grey-50 h-0.5 w-full">
              <Box className="bg-primary h-0.5" style={{ width: "64%" }}></Box>
            </Box>
          </Box>
          <Box className="flex flex-col gap-2">
            <Box className="flex flex-row items-center justify-between">
              <Link href="#" variant="subtitle2" color="textPrimary" underline="hover">
                Upsell Goal
              </Link>
              <Chip label="85 / 100" variant="outlined" />
            </Box>
            <Box className="bg-grey-50 h-0.5 w-full">
              <Box className="bg-primary h-0.5" style={{ width: "85%" }}></Box>
            </Box>
          </Box>
        </Box>
      </CardContent>
    </Card>
  );
}