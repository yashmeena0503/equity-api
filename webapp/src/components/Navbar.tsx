import { NavLink } from "react-router-dom";

function Navbar() {
  return (
    <nav className="navbar">
      <div className="navbar-brand">
        Magi
      </div>

      <div className="navbar-links">
        <NavLink to="/nodegraph">
          NodeGraph
        </NavLink>

        <NavLink to="/transactions">
          Transactions
        </NavLink>

        <NavLink to="/payouts">
          Payouts
        </NavLink>
      </div>
    </nav>
  );
}

export default Navbar;
