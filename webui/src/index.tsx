import React, { StrictMode } from "react";
import { BrowserRouter } from "react-router-dom";
import { createRoot } from "react-dom/client";
import { Provider } from "react-redux";
import { PersistGate } from "redux-persist/integration/react";
import { store, persistor } from "@/store";

import App from "@/App";
import "bootstrap/dist/css/bootstrap.min.css";
import "@/assets/scss/style.scss";

const container = document.getElementById("root");
// biome-ignore lint/style/noNonNullAssertion: <explanation>
const root = createRoot(container!);
const app = (
  <BrowserRouter>
    <StrictMode>
      <Provider store={store}>
        <PersistGate loading={null} persistor={persistor}>
          <App />
        </PersistGate>
      </Provider>
    </StrictMode>
  </BrowserRouter>
);
root.render(app);
