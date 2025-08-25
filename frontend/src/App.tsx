import { Suspense } from "react";
import { useTranslation } from "react-i18next";
import { BrowserRouter } from "react-router-dom";

import { Box, StyledEngineProvider } from "@mui/material";

import BackgroundWrapper from "@/components/layout/containers/background-wrapper";
import SnackbarWrapper from "@/components/layout/containers/snackbar-wrapper";
import LayoutContextProvider from "@/components/layout/layout-context";
import Loading from "@/pages/loading";
import AppRoutes from "@/routes";
import ThemeProvider from "@/theme/theme-provider";
import { AuthProvider } from "./contexts/AuthContext";
const App = () => {
  const { i18n } = useTranslation();

  return (
    <BrowserRouter>
      <AuthProvider>
        <Suspense fallback={<Loading />}>
          {/* Routes */}
          <AppRoutes />
        </Suspense>
      </AuthProvider>
    </BrowserRouter>
  );
};

export default App;
