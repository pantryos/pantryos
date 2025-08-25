import dayjs from "dayjs";
import duration from "dayjs/plugin/duration";
import relativeTime from "dayjs/plugin/relativeTime";
import { Box, Button, Card, CardContent, Link, Toolbar, Typography } from "@mui/material";
import { DataGrid, GridColDef, GridRenderCellParams } from "@mui/x-data-grid";
import NiArrowDown from "@/icons/nexture/ni-arrow-down";
import NiArrowUp from "@/icons/nexture/ni-arrow-up";
import NiCheckSquare from "@/icons/nexture/ni-check-square";
import NiEllipsisVertical from "@/icons/nexture/ni-ellipsis-vertical";
import NiExclamationSquare from "@/icons/nexture/ni-exclamation-square";

dayjs.extend(duration);
dayjs.extend(relativeTime);

export default function DashboardAnalyticsOrders() {
  function GridCustomToolbar() {
    return (
      <Toolbar className="flex flex-row items-center">
        <Typography variant="h5" component="h5" className="card-title flex-1">
          Recent Transactions
        </Typography>
      </Toolbar>
    );
  }

  return (
    <Card>
      <CardContent>
        <Box className="h-[668px]">
          <DataGrid
            rows={rows}
            columns={columns}
            hideFooter
            disableColumnMenu
            columnHeaderHeight={40}
            disableRowSelectionOnClick
            className="border-none"
            showToolbar
            slots={{
              toolbar: GridCustomToolbar,
              columnSortedDescendingIcon: () => <NiArrowDown size={"small"} />,
              columnSortedAscendingIcon: () => <NiArrowUp size={"small"} />,
            }}
          />
        </Box>
      </CardContent>
    </Card>
  );
}

const currencyFormatter = new Intl.NumberFormat("en-US", { style: "currency", currency: "USD" });

const columns: GridColDef<(typeof rows)[number]>[] = [
  { field: "id", headerName: "Transaction ID", width: 130 },
  {
    field: "customer", headerName: "Customer", width: 150,
    renderCell: (params) => <Link href="#" variant="body1" underline="hover" className="text-text-primary">{params.value}</Link>,
  },
  { field: "items", headerName: "Items", type: "number", width: 80, align: "center", headerAlign: "center" },
  {
    field: "total", headerName: "Total", type: "number", width: 100,
    valueFormatter: (value) => currencyFormatter.format(value),
  },
  {
    field: "status", headerName: "Status", flex: 1, minWidth: 120,
    renderCell: (params) => {
      const status = params.value;
      if (status === "Completed") {
        return <Button className="pointer-events-none self-center" size="tiny" color="success" variant="pastel" startIcon={<NiCheckSquare size={"tiny"} />}>{status}</Button>;
      }
      return <Button className="pointer-events-none self-center" size="tiny" color="error" variant="pastel" startIcon={<NiExclamationSquare size={"tiny"} />}>{status}</Button>;
    },
  },
];

const rows = [
  { id: "TXN-00875", customer: "John Doe", items: 3, total: 85.50, status: "Completed" },
  { id: "TXN-00874", customer: "Jane Smith", items: 1, total: 12.99, status: "Completed" },
  { id: "TXN-00873", customer: "Mike Johnson", items: 5, total: 152.75, status: "Completed" },
  { id: "TXN-00872", customer: "Sarah Brown", items: 2, total: 45.00, status: "Refunded" },
  { id: "TXN-00871", customer: "David Wilson", items: 1, total: 9.99, status: "Completed" },
  { id: "TXN-00870", customer: "Emily Davis", items: 4, total: 99.80, status: "Completed" },
  { id: "TXN-00869", customer: "Chris Miller", items: 2, total: 32.50, status: "Completed" },
  { id: "TXN-00868", customer: "Jessica Garcia", items: 1, total: 78.00, status: "Completed" },
  { id: "TXN-00867", customer: "Daniel Martinez", items: 3, total: 61.25, status: "Completed" },
  { id: "TXN-00866", customer: "Laura Rodriguez", items: 1, total: 15.00, status: "Completed" },
  { id: "TXN-00865", customer: "Kevin White", items: 2, total: 24.99, status: "Refunded" },
];