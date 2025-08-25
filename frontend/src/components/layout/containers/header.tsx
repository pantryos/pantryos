import Notifications from "../notifications/notifications";
import Search from "../search/search";
import Shortcuts from "../shortcuts/shortcuts";
import User from "../user/user";
import VersionSelect from "../version-select/version-select";
import { useEffect, useState } from "react";
import { Link } from "react-router-dom";

import { Box, Button } from "@mui/material";

import { useLayoutContext } from "@/components/layout/layout-context";
import Logo from "@/components/logo/logo";
import NiMenuSplit from "@/icons/nexture/ni-menu-split";
import { cn } from "@/lib/utils";

export default function Header() {
  const { showLeftInMobile, showLeftMobileButton, leftMenuWidth, leftShowBackdrop } = useLayoutContext();

  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    setMounted(true);
  }, []);

  if (!mounted) {
    return (
      <Box
        component="header"
        className="flex h-14 flex-none flex-row items-center sm:h-16"
        style={{ padding: `0 var(--main-padding)` }}
      >
        <Box className="flex h-full flex-row items-center">
          <Link to="/dashboard">
            <Logo  />
          </Link>
        </Box>
      </Box>
    );
  }

  return (
    <Box className="mui-fixed fixed z-20 h-20 w-full" component="header">
      <Box
        className={cn(
          "bg-background-paper flex h-full w-full flex-none flex-row items-center rounded-br-3xl sm:h-20",
          leftShowBackdrop && "pointer-events-none",
        )}
        style={{ padding: `0 var(--main-padding)` }}
      >
        <Box className="bg-background-paper shadow-darker-xs absolute inset-0 -z-10 rounded-b-3xl"></Box>
        {/* Left menu button */}
        <Button
          variant="text"
          size="large"
          color="text-primary"
          className={cn(
            "icon-only hover-icon-shrink [&.active]:text-primary [&.active]:bg-grey-25 hover:bg-grey-25",
            showLeftMobileButton ? "flex" : "hidden",
            leftMenuWidth.primary > 0 && "active",
          )}
          onClick={() => showLeftInMobile()}
          startIcon={<NiMenuSplit size={24} />}
        />

        <Box className="flex h-full flex-row items-center gap-6">
          {/* Logo */}
          <Link to="/dashboard">
            <Logo  />
          </Link>

        </Box>

        {/* Right buttons */}
        <Box className="ml-auto flex flex-row sm:gap-1">
          <Notifications />
          <Shortcuts />
        </Box>

        {/* User Avatar and Menu */}
        <User />
      </Box>
    </Box>
  );
}
