import { BrowserRouter, Navigate, Route, Routes } from "react-router-dom";

import Navbar from "./components/Navbar";
import NodeGraph from "./pages/NodeGraph";
import Transactions from "./pages/Transactions";
import Payouts from "./pages/Payouts";

function App() {
  return (
    <BrowserRouter>
      <div className="app">
        <Navbar />

        <main className="main">
          <Routes>
            <Route path="/" element={<Navigate to="/nodegraph" replace />} />
            <Route path="/nodegraph" element={<NodeGraph />} />
            <Route path="/transactions" element={<Transactions />} />
            <Route path="/payouts" element={<Payouts />} />
          </Routes>
        </main>
      </div>
    </BrowserRouter>
  );
}

export default App;
