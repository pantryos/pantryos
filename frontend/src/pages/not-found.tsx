import { Link } from "react-router-dom";

import { Box, Button, Paper, Typography } from "@mui/material";

import Logo from "@/components/logo/logo";
import NiHome from "@/icons/nexture/ni-home";
import { cn } from "@/lib/utils";

import BackgroundWrapper from "@/components/layout/containers/background-wrapper";
import SnackbarWrapper from "@/components/layout/containers/snackbar-wrapper";
import LayoutContextProvider from "@/components/layout/layout-context";
import ThemeProvider from "@/theme/theme-provider";
import { StyledEngineProvider } from "@mui/material";
import { useTranslation } from "react-i18next";

export default function Page() {
  const { i18n } = useTranslation();

  return (
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
            <Box className="flex min-h-screen w-full items-center justify-center p-4">
              <Paper
                elevation={3}
                className={cn(
                  "bg-background-paper shadow-darker-xs min-h-[400px] max-w-full min-w-full items-center justify-center rounded-4xl bg-center py-14 md:min-w-[800px]",
                )}
              >
                <Box className="flex flex-col gap-4 px-8 sm:px-14">
                  <Box className="flex flex-col">
                    <Box className="mb-14 flex justify-center">
                      <Logo classNameMobile="hidden" />
                    </Box>

                    <Box className="flex flex-col items-center gap-4">
                      <Typography variant="h1" component="h1">
                        Page not found!️
                      </Typography>
                      <Typography variant="body1" color="text.secondary">
                        Error Code: 404
                      </Typography>
                      <Button variant="outlined" startIcon={<NiHome />} to="/dashboard" component={Link}>
                        dashboard
                      </Button>
                    </Box>
                  </Box>
                </Box>
              </Paper>
            </Box>
          </SnackbarWrapper>
        </LayoutContextProvider>

      </Box>
      </ThemeProvider>
    </StyledEngineProvider>
  );
}
