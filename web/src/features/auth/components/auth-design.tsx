"use client";

import { createContext, useContext, type ReactNode } from "react";
import refinedStyles from "./auth.module.css";
import loginStyles from "./login.module.css";

const AuthDesignContext = createContext(refinedStyles);

/** Login retains the original visual theme; shared form behavior stays identical. */
export function AuthDesign({
  registration,
  children,
}: {
  registration: boolean;
  children: ReactNode;
}) {
  return (
    <AuthDesignContext value={registration ? refinedStyles : loginStyles}>
      {children}
    </AuthDesignContext>
  );
}

export function useAuthDesign() {
  return useContext(AuthDesignContext);
}
