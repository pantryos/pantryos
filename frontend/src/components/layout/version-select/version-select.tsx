import { SyntheticEvent, useState } from "react";

import { Button, Fade, Menu, MenuItem, PopoverVirtualElement } from "@mui/material";

import NiChevronRightSmall from "@/icons/nexture/ni-chevron-right-small";
import { cn } from "@/lib/utils";
export default function VersionSelect() {
  const [anchorEl, setAnchorEl] = useState<EventTarget | Element | PopoverVirtualElement | null>(null);
  const open = Boolean(anchorEl);
  const handleClick = (event: Event | SyntheticEvent) => {
    setAnchorEl(event.currentTarget);
  };
  const handleClose = () => {
    setAnchorEl(null);
  };
  return (
    <>
      <Button
        className="px-4"
        variant="text"
        color="grey"
        onClick={handleClick}
        endIcon={
          <NiChevronRightSmall size="medium" className={cn("-ml-1 transition-transform", open && "rotate-90")} />
        }
      >
        v6.3.0
      </Button>
      <Menu
        anchorEl={anchorEl as Element}
        open={open}
        onClose={handleClose}
        className="mt-1"
        slots={{
          transition: Fade,
        }}
      >
        <MenuItem onClick={handleClose}>v6.3.0</MenuItem>
        <MenuItem onClick={handleClose}>v6.2.0</MenuItem>
        <MenuItem onClick={handleClose}>v6.1.0</MenuItem>
        <MenuItem onClick={handleClose}>v6.0.0</MenuItem>
      </Menu>
    </>
  );
}
