import { SyntheticEvent, useRef, useState } from "react";

import {
  Avatar,
  Box,
  Button,
  Card,
  CardActions,
  ClickAwayListener,
  Fade,
  List,
  ListItem,
  ListItemAvatar,
  ListItemButton,
  ListItemText,
  Popper,
  Tooltip,
  Typography,
} from "@mui/material";

import NiBag from "@/icons/nexture/ni-bag";
import NiCells from "@/icons/nexture/ni-cells";
import NiEllipsisHorizontal from "@/icons/nexture/ni-ellipsis-horizontal";
import NiMessages from "@/icons/nexture/ni-messages";
import NiPercent from "@/icons/nexture/ni-percent";
import NiPlus from "@/icons/nexture/ni-plus";
import NiTelescope from "@/icons/nexture/ni-telescope";
import NiUsers from "@/icons/nexture/ni-users";
import { cn } from "@/lib/utils";

export default function Shortcuts() {
  const [tooltipShow, setTooltipShow] = useState(false);

  const [open, setOpen] = useState(false);
  const anchorRef = useRef<HTMLButtonElement>(null);

  const handleToggle = () => {
    setOpen((prevOpen) => !prevOpen);
  };

  const handleClose = (event: Event | SyntheticEvent) => {
    if (anchorRef.current && anchorRef.current.contains(event.target as HTMLElement)) {
      return;
    }

    setOpen(false);
  };

  return (
    <>
      <Tooltip title="Shortcuts" placement="bottom" arrow open={!open && tooltipShow}>
        <Button
          variant="text"
          size="large"
          color="text-primary"
          className={cn(
            "icon-only hover-icon-shrink [&.active]:text-primary hover:bg-grey-25",
            open && "active bg-grey-25",
          )}
          onClick={handleToggle}
          onMouseEnter={() => setTooltipShow(true)}
          onMouseLeave={() => setTooltipShow(false)}
          ref={anchorRef}
          startIcon={<NiCells variant={open ? "contained" : "outlined"} size={24} />}
        />
      </Tooltip>
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
                <Card className="shadow-darker-sm! w-xs">
                  <Box className="flex flex-1 flex-row items-start justify-between pr-4">
                    <Typography variant="h5" component="h5" className="card-title px-4 pt-4">
                      Shortcuts
                    </Typography>
                    <Button
                      className="icon-only mt-3"
                      size="tiny"
                      color="grey"
                      variant="text"
                      startIcon={<NiPlus size={"small"} />}
                    />
                  </Box>
                  <Box className="mb-4">
                    <List className="max-h-72 overflow-auto">
                      <ListItem className="py-0 pr-4 pl-0">
                        <ListItemButton classes={{ root: "group items-start" }}>
                          <ListItemAvatar>
                            <Avatar className="medium bg-primary-light/10 mr-3">
                              <NiBag size="medium" className="text-primary" />
                            </Avatar>
                          </ListItemAvatar>
                          <ListItemText
                            primary={
                              <Typography component="p" variant="subtitle2" className="leading-4">
                                Add Product
                              </Typography>
                            }
                            secondary="/home/products/add"
                          />
                          <Button
                            className="icon-only hover:text-text-primary hover:bg-grey-100 mt-1 hidden opacity-0 group-hover:flex group-hover:opacity-100"
                            size="tiny"
                            color="grey"
                            variant="text"
                            startIcon={<NiEllipsisHorizontal size={"small"} />}
                          />
                        </ListItemButton>
                      </ListItem>
                    </List>
                  </Box>
                  <CardActions disableSpacing>
                    <Button variant="outlined" size="tiny" color="grey" className="w-full">
                      Add Shortcut
                    </Button>
                  </CardActions>
                </Card>
              </ClickAwayListener>
            </Box>
          </Fade>
        )}
      </Popper>
    </>
  );
}
