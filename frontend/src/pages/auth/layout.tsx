// src/layouts/AuthLayout.tsx (or similar path)

import React, { Suspense } from "react";
import { Navigate, Outlet } from "react-router-dom";

import Loading from "@/pages/loading";
import { useAuth } from "@/contexts/AuthContext";

const AuthLayout: React.FC = () => { // No need for a `children` prop
  const { isAuthenticated, loading } = useAuth(); // Assuming you also have a loading state

  // Show a loading indicator while checking auth status
  if (loading) {
    return <Loading />;
  }

  // If the user is already authenticated, redirect them away from auth pages
  if (isAuthenticated) {
    return <Navigate to="/dashboard" replace />;
  }

  // If the user is not authenticated, render the child route (e.g., Login, Register)
  // The <Outlet /> is a placeholder for the child component.
  return (
    <div>
      {/* You can add layout elements here, like a header or footer for the auth pages */}
      <Suspense fallback={<Loading />}>
        <Outlet />
      </Suspense>
    </div>
  );
};

export default AuthLayout;