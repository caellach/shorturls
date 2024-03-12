import "@/extensions/string";
import React, { StrictMode } from "react";
import { BrowserRouter } from "react-router-dom";
import { createRoot } from "react-dom/client";
import { Provider } from "react-redux";
import { PersistGate } from "redux-persist/integration/react";
import { store, persistor } from "@/store";

import App from "@/App";
import "bootstrap/dist/css/bootstrap.min.css";
import "react-virtualized/styles.css";
import "react-toastify/dist/ReactToastify.css";
import "@/assets/scss/style.scss";
import { ToastContainer } from "react-toastify";

const container = document.getElementById("root");
// biome-ignore lint/style/noNonNullAssertion: <explanation>
const root = createRoot(container!);
const app = (
  <BrowserRouter>
    <ToastContainer />
    <Provider store={store}>
      <PersistGate loading={null} persistor={persistor}>
        <App />
      </PersistGate>
    </Provider>
  </BrowserRouter>
);
root.render(app);
