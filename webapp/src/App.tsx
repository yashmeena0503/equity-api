import { BrowserRouter, Navigate, Route, Routes } from "react-router-dom";

import Navbar from "./components/Navbar";
import User from "./pages/User";
import Transactions from "./pages/Transactions";
import Payouts from "./pages/Payouts";

function App() {
  return (
    <BrowserRouter>
      <div className="app">
        <Navbar />

        <main className="main">
          <Routes>
            <Route path="/" element={<Navigate to="/user" replace />} />
            <Route path="/user" element={<User />} />
            <Route path="/transactions" element={<Transactions />} />
            <Route path="/payouts" element={<Payouts />} />
          </Routes>
        </main>
      </div>
    </BrowserRouter>
  );
}

export default App;
