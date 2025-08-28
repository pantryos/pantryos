import DashboardAnalyticsAppStatus from "./sections/dashboard-analytics-app-status";
import DashboardAnalyticsCategories from "./sections/dashboard-analytics-categories";
import DashboardAnalyticsCurrencies from "./sections/dashboard-analytics-currencies";
import DashboardAnalyticsDuration from "./sections/dashboard-analytics-duration";
import DashboardAnalyticsOrders from "./sections/dashboard-analytics-orders";
import DashboardAnalyticsOrdersStocks from "./sections/dashboard-analytics-orders-stocks";
import DashboardAnalyticsProgresses from "./sections/dashboard-analytics-progresses";
import DashboardAnalyticsSales from "./sections/dashboard-analytics-sales";
import DashboardAnalyticsStats from "./sections/dashboard-analytics-stats";
import DashboardAnalyticsVisits from "./sections/dashboard-analytics-visits";
import dayjs, { Dayjs } from "dayjs";
import { useState } from "react";
import { Link } from "react-router-dom";

import { Breadcrumbs, Button, FormControl, Tooltip, Typography, Box } from "@mui/material";
import { Grid } from "@mui/material";
import { LocalizationProvider } from "@mui/x-date-pickers";
import { AdapterDayjs } from "@mui/x-date-pickers/AdapterDayjs";
import { DatePicker } from "@mui/x-date-pickers/DatePicker";

import NiCalendar from "@/icons/nexture/ni-calendar";
import NiCellsPlus from "@/icons/nexture/ni-cells-plus";
import NiChevronDownSmall from "@/icons/nexture/ni-chevron-down-small";
import NiChevronLeftSmall from "@/icons/nexture/ni-chevron-left-small";
import NiChevronRightSmall from "@/icons/nexture/ni-chevron-right-small";
import NiCross from "@/icons/nexture/ni-cross";
import NiKnobs from "@/icons/nexture/ni-knobs";
import { cn } from "@/lib/utils";
import { useAuth } from "@/contexts/AuthContext";
import LowStockBanner from "@/components/LowStockBanner";

export default function Page() {
  const [startDate, setStartDate] = useState<Dayjs | null>(dayjs().startOf('month'));
  const [endDate, setEndDate] = useState<Dayjs | null>(dayjs().endOf('month'));
  const { user } = useAuth();

  const handleStartDateChange = (newValue: Dayjs | null) => {
    setStartDate(newValue);
    // If start date is after end date, adjust end date
    if (newValue && endDate && newValue.isAfter(endDate)) {
      setEndDate(newValue);
    }
  };

  const handleEndDateChange = (newValue: Dayjs | null) => {
    setEndDate(newValue);
    // If end date is before start date, adjust start date
    if (newValue && startDate && newValue.isBefore(startDate)) {
      setStartDate(newValue);
    }
  };

  return (
    <Grid container spacing={5}>
      <Grid container spacing={2.5} className="w-full" size={12}>
        <Grid size={{ xs: 12, md: "grow" }}>
          <Typography variant="h1" component="h1" className="mb-0">
            Welcome back, to pantryOS 
          </Typography>
          <Breadcrumbs>
            <Link to="/dashboard">Business</Link>
            <Typography>Sales Dashboard</Typography>
          </Breadcrumbs>
        </Grid>



        <Grid size={{ xs: 12, md: "auto" }} className="flex flex-row items-start gap-2">
          <Box className="flex flex-row gap-2 items-center">
            <FormControl variant="standard" className="surface-standard mb-0 w-full md:w-auto">
              <LocalizationProvider dateAdapter={AdapterDayjs}>
                <DatePicker
                  label="Start Date"
                  value={startDate}
                  onChange={handleStartDateChange}
                  slots={{
                    openPickerIcon: (props) => <NiCalendar {...props} className={cn(props.className, "text-text-secondary")} />,
                    switchViewIcon: (props) => <NiChevronDownSmall {...props} className={cn(props.className, "text-text-secondary")} />,
                    leftArrowIcon: (props) => <NiChevronLeftSmall {...props} className={cn(props.className, "text-text-secondary")} />,
                    rightArrowIcon: (props) => <NiChevronRightSmall {...props} className={cn(props.className, "text-text-secondary")} />,
                    clearIcon: (props) => <NiCross {...props} className={cn(props.className, "text-text-secondary")} />,
                  }}
                  slotProps={{
                    textField: { size: "small", variant: "standard" },
                    desktopPaper: { className: "outlined" },
                  }}
                />
              </LocalizationProvider>
            </FormControl>
            
            <FormControl variant="standard" className="surface-standard mb-0 w-full md:w-auto">
              <LocalizationProvider dateAdapter={AdapterDayjs}>
                <DatePicker
                  label="End Date"
                  value={endDate}
                  onChange={handleEndDateChange}
                  minDate={startDate ?? undefined}
                  slots={{
                    openPickerIcon: (props) => <NiCalendar {...props} className={cn(props.className, "text-text-secondary")} />,
                    switchViewIcon: (props) => <NiChevronDownSmall {...props} className={cn(props.className, "text-text-secondary")} />,
                    leftArrowIcon: (props) => <NiChevronLeftSmall {...props} className={cn(props.className, "text-text-secondary")} />,
                    rightArrowIcon: (props) => <NiChevronRightSmall {...props} className={cn(props.className, "text-text-secondary")} />,
                    clearIcon: (props) => <NiCross {...props} className={cn(props.className, "text-text-secondary")} />,
                  }}
                  slotProps={{
                    textField: { size: "small", variant: "standard" },
                    desktopPaper: { className: "outlined" },
                  }}
                />
              </LocalizationProvider>
            </FormControl>
          </Box>
          
          <Tooltip title="Dashboard Settings">
            <Button
              className="icon-only surface-standard flex-none"
              size="medium"
              color="grey"
              variant="surface"
              startIcon={<NiKnobs size={"medium"} />}
            />
          </Tooltip>
          <Tooltip title="New Sale">
            <Button
              className="icon-only surface-standard flex-none"
              size="medium"
              color="primary"
              variant="surface"
              startIcon={<NiCellsPlus size={"medium"} />}
            />
          </Tooltip>
        </Grid>
      </Grid>
      
      {/* ... The rest of the dashboard components ... */}
        <Grid size={{ xs: 12,  }}>
           <LowStockBanner  maxItems={5} />
        </Grid>
       <Grid container size={12}>
        <Grid container size={{ lg: 6, xs: 12 }} className="items-start">
          <Grid container size={12} spacing={2.5} className="flex-none">
            <DashboardAnalyticsStats />
          </Grid>
          <Grid size={12}>
            
            <DashboardAnalyticsProgresses />
          </Grid>
        </Grid>
        <Grid size={{ lg: 6, xs: 12 }}>
          {/* <DashboardAnalyticsOrders /> */}
          <DashboardAnalyticsCategories />
        </Grid>
      </Grid>

      

      <Grid container size={12}>
        <Grid size={{ lg: 6, xs: 12 }}>
          {/* <DashboardAnalyticsCurrencies /> */}
          <DashboardAnalyticsSales />
        </Grid>
        <Grid size={{ lg: 6, xs: 12 }}>
          {/* <DashboardAnalyticsProgresses /> */}
           {/* <DashboardAnalyticsDuration /> */}
        </Grid>
        <Grid size={{ lg: 3, xs: 12 }}>
          
        </Grid>
        <Grid size={{ lg: 3, xs: 12 }}>
          {/* <DashboardAnalyticsOrdersStocks /> */}
        </Grid>
        <Grid size={{ lg: 6, xs: 12 }}>
          {/* <DashboardAnalyticsAppStatus /> */}
         
          
        </Grid>
        <Grid size={{ lg: 6, xs: 12 }}>
          
        </Grid>
        <Grid size={{ lg: 6, xs: 12 }}>
          {/* <DashboardAnalyticsVisits /> */}
        </Grid>
      </Grid>
    </Grid>
  );
}