import React from "react";
import { Navigate, Outlet, Route, Routes } from "react-router-dom";

import { leftMenuBottomItems, leftMenuItems } from "@/menu-items";
import AppLayout from "@/pages/app/layout";
import AuthLayout from "@/pages/auth/layout";
import Loading from "@/pages/loading.tsx";
import NotFound from "@/pages/not-found";
import { MenuItem } from "@/types/types";
import LandingPage from "./components/LandingPage";
import Login from "./components/Login";
import Register from "./components/Register";
import { useAuth } from "./contexts/AuthContext";
import ThemeNoAuth from "./pages/auth/ThemeNoAuth";

// Statically import all possible pages for build
const modules = import.meta.glob("./pages/**/page.tsx");

// Lazy load page components
const lazyLoad = (path: string) => {
  // Handle different paths based on the route
  let key: string;
  if (path === "/") {
    key = "./pages/page.tsx";
  } else if (path.startsWith("/auth")) {
    key = `./pages/auth${path.substring(5)}/page.tsx`; // Remove "/auth"
  } else {
    key = `./pages/app${path}/page.tsx`;
  }

  const importer = modules[key];

  // If file not found fallback to 404
  if (!importer) return <Navigate to="/404" replace />;

  const Component = React.lazy(importer as () => Promise<{ default: React.ComponentType<any> }>);

  return (
    <React.Suspense fallback={<Loading />}>
      <Component />
    </React.Suspense>
  );
};

// Recursively generate routes from menu items
const generateRoutesFromMenuItems = (menuItems: MenuItem[]): React.ReactElement[] => {
  return menuItems.flatMap((item: MenuItem) => {
    const routes: React.ReactElement[] = [];

    // Skip external links
    if (item.isExternalLink || !item.href) {
      return [];
    }

    // Add route for current item
    routes.push(<Route key={item.id} path={item.href} element={lazyLoad(item.href)} />);

    // Add routes for children
    if (item.children && item.children.length > 0) {
      routes.push(...generateRoutesFromMenuItems(item.children));
    }

    return routes;
  });
};

const ProtectedRoute: React.FC = () => {
  const { isAuthenticated, loading } = useAuth();

  if (loading) {
    return <Loading />;
  }

  // If authenticated, render the nested child routes.
  // Otherwise, redirect to the login page.
  return isAuthenticated ? <Outlet /> : <Navigate to="/auth/login" replace />;
};


// Generate auth routes
const generateAuthRoutes = (): React.ReactElement[] => {
  return [
    <Route element={<ThemeNoAuth />}>
     <Route key="sign-in" path="sign-in" element={lazyLoad("/auth/sign-in")} />
     </Route>,
     <Route element={<ThemeNoAuth />}>
    <Route
      path="login"
      element={
        lazyLoad("/auth/sign-in")
      }
    />
    </Route>,
    <Route
      path="register"
      element={
        <Register />
      }
    />,
   
  ];
};

// Generate routes from both menu arrays
const mainRoutes = generateRoutesFromMenuItems(leftMenuItems);
const bottomRoutes = generateRoutesFromMenuItems(leftMenuBottomItems);
const authRoutes = generateAuthRoutes();

// Main Routes component
const AppRoutes = () => {
  return (
    <Routes>
      {/* Landing page route */}

      <Route path="/" element={<LandingPage />} />

      {/* App routes with AppLayout */}
      <Route element={<ProtectedRoute />}>
        <Route element={<AppLayout />}>
          {/* Routes generated from menu items */}
          {mainRoutes}
          {bottomRoutes}
        </Route>
      </Route>

      {/* Auth routes with AuthLayout */}
      <Route path="/auth" element={<AuthLayout />}>
        <Route index element={<Navigate to="/auth/login" replace />} />
        {authRoutes}
      </Route>

      {/* 404 route */}
      <Route path="/404" element={<NotFound />} />
      <Route path="*" element={<Navigate to="/404" replace />} />
    </Routes>
  );
};

export default AppRoutes;
