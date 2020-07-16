import React from "react";
import Home from "./pages/Home";
import {
  BrowserRouter as Router,
  Switch,
  Route,
  Link
} from "react-router-dom";

export default function App() {
  return (
    <Router>
      <div>
          <Route path="/">
            <Home />
          </Route>
      </div>
    </Router>
  );
}