import { useTranslation } from "react-i18next";
import { Link } from "react-router-dom";

import { Box, Button } from "@mui/material";

export default function Footer() {
  const { t } = useTranslation();

  return (
    <Box component="footer" className="flex h-10 items-center justify-center">
    </Box>
  );
}
