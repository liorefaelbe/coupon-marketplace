import { useState } from "react";
import AdminLogin from "./components/AdminLogin";
import Tabs from "./components/Tabs";
import AdminPage from "./pages/AdminPage";
import CustomerPage from "./pages/CustomerPage";
import { fetchAdminProducts } from "./api/api";
import "./App.css";

type Mode = "customer" | "admin";

function App() {
  const [mode, setMode] = useState<Mode>("customer");
  const [adminAuthenticated, setAdminAuthenticated] = useState(false);
  const [showAdminLogin, setShowAdminLogin] = useState(false);
  const [adminUsername, setAdminUsername] = useState("");
  const [adminPassword, setAdminPassword] = useState("");

  const handleOpenAdmin = () => {
    if (adminAuthenticated) {
      setMode("admin");
      return;
    }

    setShowAdminLogin(true);
  };

  const handleAdminLogin = async (username: string, password: string) => {
    await fetchAdminProducts(username, password);

    setAdminUsername(username);
    setAdminPassword(password);
    setAdminAuthenticated(true);
    setShowAdminLogin(false);
    setMode("admin");
  };

  const handleAdminCancel = () => {
    setShowAdminLogin(false);
  };

  return (
    <div className="page">
      <header className="hero">
        <h1>Digital Coupon Marketplace</h1>
        <p>Find The Gift Card You've Been Looking For.</p>
      </header>

      <Tabs mode={mode} onCustomerClick={() => setMode("customer")} onAdminClick={handleOpenAdmin} />

      {mode === "customer" ? <CustomerPage /> : <AdminPage username={adminUsername} password={adminPassword} />}

      {showAdminLogin && <AdminLogin onLogin={handleAdminLogin} onCancel={handleAdminCancel} />}
    </div>
  );
}

export default App;
