import { MenuItem } from "@/types/types";

export const leftMenuItems: MenuItem[] = [
  {
    id: "dashboard",
    icon: "NiHome",
    label: "dashboards",
    color: "text-primary",
    href: "/dashboard", 
  },
    {
    id: "categories",
    icon: "NiCatalog",
    label: "categories",
    color: "text-primary",
    href: "/categories", 
  },
   {
    id: "inventory",
    icon: "NiArchiveCheck",
    label: "inventory",
    color: "text-primary",
    href: "/inventory", 
  },
     {
    id: "menu",
    icon: "NiBook",
    label: "menu",
    color: "text-primary",
    href: "/menu", 
  },
     {
    id: "delivery",
    icon: "NiCar",
    label: "delivery",
    color: "text-primary",
    href: "/delivery", 
  },
  

];

export const leftMenuBottomItems: MenuItem[] = [
  // { id: "settings", label: "menu-settings", href: "/settings", icon: "NiSettings" },
];
