import "@/style/global.css";

import { Suspense, useEffect } from "react";
import { Outlet, useLocation } from "react-router-dom";

import ContentWrapper from "@/components/layout/containers/content-wrapper";
import Header from "@/components/layout/containers/header";
import Main from "@/components/layout/containers/main";
import LeftMenu from "@/components/layout/menu/left-menu";
import MenuBackdrop from "@/components/layout/menu/menu-backdrop";
import Loading from "@/pages/loading";

import BackgroundWrapper from "@/components/layout/containers/background-wrapper";
import SnackbarWrapper from "@/components/layout/containers/snackbar-wrapper";
import LayoutContextProvider from "@/components/layout/layout-context";
import { Box, StyledEngineProvider } from "@mui/material";
import { useTranslation } from "react-i18next";
import ThemeProvider from "@/theme/theme-provider";

export default function AppLayout() {
  const { pathname } = useLocation();
  const { i18n } = useTranslation();

  useEffect(() => {
    window.scrollTo(0, 0);
  }, [pathname]);

  return (
    <>
      <StyledEngineProvider enableCssLayer>
        <ThemeProvider>
        <Box lang={i18n.language} className="font-mulish font-urbanist relative overflow-hidden antialiased">
          {/* Initial loader */}
          <div id="initial-loader">
            <div className="spinner"></div>
          </div>
          {/* Initial loader end */}
          <LayoutContextProvider>
            <BackgroundWrapper />
            <SnackbarWrapper>
              <Header />
              <LeftMenu />
              <Main>
                <ContentWrapper>
                  <Suspense fallback={<Loading />}>
                    <Outlet />
                  </Suspense>
                </ContentWrapper>
              </Main>
              <MenuBackdrop />
            </SnackbarWrapper>
          </LayoutContextProvider>
        </Box>
        </ThemeProvider>
      </StyledEngineProvider>
    </>
  );
}
