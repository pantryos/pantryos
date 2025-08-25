import UserLanguageSwitch from "./user-language-switch";
import UserModeSwitch from "./user-mode-switch";
import UserThemeSwitch from "./user-theme-switch";
import { SyntheticEvent, useRef, useState } from "react";
import { useTranslation } from "react-i18next";
import { Link, useNavigate } from "react-router-dom";

import {
  Accordion,
  AccordionDetails,
  AccordionSummary,
  Avatar,
  AvatarGroup,
  Box,
  Card,
  CardContent,
  Divider,
  Fade,
  ListItemIcon,
  Typography,
} from "@mui/material";
import Button from "@mui/material/Button";
import ClickAwayListener from "@mui/material/ClickAwayListener";
import MenuItem from "@mui/material/MenuItem";
import MenuList from "@mui/material/MenuList";
import Popper from "@mui/material/Popper";

import NiBuilding from "@/icons/nexture/ni-building";
import NiChevronRightSmall from "@/icons/nexture/ni-chevron-right-small";
import NiDocumentFull from "@/icons/nexture/ni-document-full";
import NiFolder from "@/icons/nexture/ni-folder";
import NiQuestionHexagon from "@/icons/nexture/ni-question-hexagon";
import NiSettings from "@/icons/nexture/ni-settings";
import NiUser from "@/icons/nexture/ni-user";
import NiUsers from "@/icons/nexture/ni-users";
import { cn } from "@/lib/utils";
import { useAuth } from "@/contexts/AuthContext";

export default function User() {
  const [open, setOpen] = useState(false);
  const anchorRef = useRef<HTMLButtonElement>(null);
  const { t } = useTranslation();
  const { user, logout } = useAuth();

  console.log("user from context:", user);
  

  const handleToggle = () => {
    setOpen((prevOpen) => !prevOpen);
  };

  const handleClose = (event: Event | SyntheticEvent) => {
    if (anchorRef.current && anchorRef.current.contains(event.target as HTMLElement)) {
      return;
    }
    setOpen(false);
  };

  const navigate = useNavigate();

  return (
    <>
      <Box ref={anchorRef}>
        {/* Desktop button */}
        <Button
          variant="text"
          color="text-primary"
          className={cn(
            "group hover:bg-grey-25 ml-2 hidden gap-2 rounded-lg py-0! pr-0! hover:py-1! hover:pr-1.5! md:flex",
            open && "active bg-grey-25 py-1! pr-1.5!",
          )}
          onClick={handleToggle}
        >
          <Box>{user?.email}</Box>
          <Avatar
            alt="avatar"
            src="/images/avatars/avatar-1.jpg"
            className={cn(
              "large transition-all group-hover:ml-0.5 group-hover:h-8 group-hover:w-8",
              open && "ml-0.5 h-8! w-8!",
            )}
          />
        </Button>
        {/* Desktop button */}

        {/* Mobile button */}
        <Button
          variant="text"
          size="large"
          color="text-primary"
          className={cn(
            "icon-only hover:bg-grey-25 hover-icon-shrink [&.active]:text-primary group mr-1 ml-1 p-0! hover:p-1.5! md:hidden",
            open && "active bg-grey-25 p-1.5!",
          )}
          onClick={handleToggle}
          startIcon={
            <Avatar
              alt="avatar"
              src="/images/avatars/avatar-1.jpg"
              className={cn("large transition-all group-hover:h-7 group-hover:w-7", open && "h-7! w-7!")}
            />
          }
        />
        {/* Mobile button */}
      </Box>

      <Popper
        open={open}
        anchorEl={anchorRef.current}
        role={undefined}
        placement="bottom-end"
        className="mt-3!"
        transition
      >
        {({ TransitionProps }) => (
          <Fade {...TransitionProps}>
            <Box>
              <ClickAwayListener onClickAway={handleClose}>
                <Card className="shadow-darker-sm!">
                  <CardContent>
                    <Box className="max-w-64 sm:w-72 sm:max-w-none">
                      <Box className="mb-4 flex flex-col items-center">
                        <Avatar alt="avatar" src="/images/avatars/avatar-1.jpg" className="large mb-2" />
                        <Typography variant="subtitle1" component="p">
                          admin
                        </Typography>
                        <Typography variant="body2" component="p" className="text-text-secondary -mt-2">
                          {user?.email}
                        </Typography>
                      </Box>

                      <Divider className="large" />
                      <Box className="my-8"></Box>
                      <Button
                        // 1. Remove the link properties (component and to)
                        onClick={logout} // 2. Add the onClick handler to call the logout function
                        variant="outlined"
                        size="tiny"
                        color="grey"
                        className="w-full"
                      >
                        {t("user-sign-out")}
                      </Button>
                    </Box>
                  </CardContent>
                </Card>
              </ClickAwayListener>
            </Box>
          </Fade>
        )}
      </Popper>
    </>
  );
}
