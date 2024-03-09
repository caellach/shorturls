import React from "react";
import { RouteObject } from "react-router-dom";
import Layouts from "@/pages/layout";
import Main from "@/pages/main";
import Apps from "@/pages/apps";
import Shorturls from "@/pages/shorturls";
import NotFound from "@/pages/NotFound";
import requireAuth from "@/components/requireAuth";

export default [
  {
    path: "/",
    element: <Layouts.AnonLayout />,
    children: [{ index: true, element: <Main /> }],
    requireAuth: false,
  },
  {
    path: "/apps",
    element: React.createElement(requireAuth(Layouts.LoggedInLayout)),
    children: [{ index: true, element: <Apps /> }],
    requireAuth: true,
  },
  {
    path: "/shorturls",
    element: React.createElement(requireAuth(Layouts.LoggedInLayout)),
    children: [{ index: true, element: <Shorturls /> }],
  },
  {
    path: "/*",
    element: <NotFound />,
  },
] as RouteObject[];
