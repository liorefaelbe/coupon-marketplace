type Mode = "customer" | "admin";

type TabsProps = {
  mode: Mode;
  onCustomerClick: () => void;
  onAdminClick: () => void;
};

function Tabs({ mode, onCustomerClick, onAdminClick }: TabsProps) {
  return (
    <div className="tabs">
      <button className={mode === "customer" ? "tab active" : "tab"} onClick={onCustomerClick}>
        Customer
      </button>
      <button className={mode === "admin" ? "tab active" : "tab"} onClick={onAdminClick}>
        Admin Access
      </button>
    </div>
  );
}

export default Tabs;
