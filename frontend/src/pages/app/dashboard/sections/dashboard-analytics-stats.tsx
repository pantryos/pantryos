import { Box, Card, CardContent, Grid, Typography } from "@mui/material";
import NiBasket from "@/icons/nexture/ni-basket";
import NiCatalog from "@/icons/nexture/ni-catalog";
import NiExclamationSquare from "@/icons/nexture/ni-exclamation-square";
import NiCheckSquare from "@/icons/nexture/ni-check-square";
import NiHeart from "@/icons/nexture/ni-heart";
import NiArchiveCheck from "@/icons/nexture/ni-archive-check";
import NiCar from "@/icons/nexture/ni-car";
import { useEffect, useState } from "react";
import apiService from "@/services/api";
import NiBook from "@/icons/nexture/ni-book";

export default function DashboardAnalyticsStats() {
    const [stats, setStats] = useState({
      inventoryItems: 0,
      menuItems: 0,
      recentDeliveries: 0,
    });

      useEffect(() => {
        const loadStats = async () => {
          try {
            const [inventory, menu, deliveries] = await Promise.all([
              apiService.getInventoryItems(),
              apiService.getMenuItems(),
              apiService.getDeliveries(),
            ]);
    
            // Add a check to ensure deliveries is a valid array before filtering
            const filteredDeliveries = Array.isArray(deliveries)
              ? deliveries.filter(
                  (d) =>
                    new Date(d.delivery_date) >
                    new Date(Date.now() - 7 * 24 * 60 * 60 * 1000)
                )
              : [];
    
            setStats({
              inventoryItems: inventory.length,
              menuItems: menu.length,
              recentDeliveries: filteredDeliveries.length,
            });
          } catch (error) {
            console.error("Failed to load dashboard stats:", error);
          } finally {
          }
        };
    
        loadStats();
      }, []);
    
  return (
    <>
      <Grid size={{ xs: 6, sm: 3 }}>
        <Card component="a" href="#" className="flex flex-col p-1 transition-transform hover:scale-[1.02]">
          <Box className="bg-primary-light/10 flex h-[4.5rem] w-full flex-none items-center justify-center rounded-2xl">
            <NiArchiveCheck className="text-primary" size={"large"} />
          </Box>
          <CardContent className="text-center">
            <Typography variant="body1" className="text-text-secondary leading-5 transition-colors">Items in inventory</Typography>
            <Typography variant="h5" className="text-leading-5">{stats.inventoryItems}</Typography>
          </CardContent>
        </Card>
      </Grid>
      <Grid size={{ xs: 6, sm: 3 }}>
        <Card component="a" href="#" className="flex flex-col p-1 transition-transform hover:scale-[1.02]">
          <Box className="bg-secondary-light/10 flex h-[4.5rem] w-full flex-none items-center justify-center rounded-2xl">
            <NiExclamationSquare className="text-secondary" size={"large"} />
          </Box>
          <CardContent className="text-center">
            <Typography variant="body1" className="text-text-secondary leading-5 transition-colors">Expiring Soon</Typography>
            <Typography variant="h5" className="text-leading-5">8</Typography>
          </CardContent>
        </Card>
      </Grid>
      <Grid size={{ xs: 6, sm: 3 }}>
        <Card component="a" href="#" className="flex flex-col p-1 transition-transform hover:scale-[1.02]">
          <Box className="bg-primary-light/10 flex h-[4.5rem] w-full flex-none items-center justify-center rounded-2xl">
            <NiCar className="text-primary" size={"large"} />
          </Box>
          <CardContent className="text-center">
            <Typography variant="body1" className="text-text-secondary leading-5 transition-colors">Log Delivery</Typography>
            <Typography variant="h5" className="text-leading-5">{stats.recentDeliveries}</Typography>
          </CardContent>
        </Card>
      </Grid>
      <Grid size={{ xs: 6, sm: 3 }}>
        <Card component="a" href="#" className="flex flex-col p-1 transition-transform hover:scale-[1.02]">
          <Box className="bg-secondary-light/10 flex h-[4.5rem] w-full flex-none items-center justify-center rounded-2xl">
            <NiBook className="text-secondary" size={"large"} />
          </Box>
          <CardContent className="text-center">
            <Typography variant="body1" className="text-text-secondary leading-5 transition-colors">Menu</Typography>
            <Typography variant="h5" className="text-leading-5">{stats.menuItems}</Typography>
          </CardContent>
        </Card>
      </Grid>
    </>
  );
}