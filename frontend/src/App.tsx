import type { ChangeEvent, SubmitEvent } from "react";
import { useEffect, useMemo, useState } from "react";
import "./App.css";

const API_BASE = "http://localhost:8080";

type Mode = "customer" | "admin";

type Product = {
  id: string;
  name: string;
  description: string;
  image_url: string;
  price?: number;
  type?: string;
  cost_price?: number;
  margin_percentage?: number;
  minimum_sell_price?: number;
  is_sold?: boolean;
  value_type?: string;
  value?: string;
};

type PurchaseResult = {
  product_id: string;
  final_price: number;
  value_type: string;
  value: string;
};

type ApiError = {
  error_code?: string;
  message?: string;
  error?: string;
};

type CouponForm = {
  name: string;
  description: string;
  image_url: string;
  cost_price: string;
  margin_percentage: string;
  value_type: string;
  value: string;
};

const initialForm: CouponForm = {
  name: "",
  description: "",
  image_url: "",
  cost_price: "",
  margin_percentage: "",
  value_type: "STRING",
  value: "",
};

function getErrorMessage(err: unknown): string {
  if (err instanceof Error) {
    return err.message;
  }
  return "Something went wrong";
}

function buildBasicAuth(username: string, password: string): string {
  return `Basic ${btoa(`${username}:${password}`)}`;
}

function App() {
  const [mode, setMode] = useState<Mode>("customer");

  const [customerProducts, setCustomerProducts] = useState<Product[]>([]);
  const [customerLoading, setCustomerLoading] = useState<boolean>(true);
  const [customerError, setCustomerError] = useState<string>("");
  const [purchaseResult, setPurchaseResult] = useState<PurchaseResult | null>(null);

  const [adminProducts, setAdminProducts] = useState<Product[]>([]);
  const [adminLoading, setAdminLoading] = useState<boolean>(false);
  const [adminMessage, setAdminMessage] = useState<string>("");

  const [adminUsername, setAdminUsername] = useState<string>("admin");
  const [adminPassword, setAdminPassword] = useState<string>("admin123");

  const [form, setForm] = useState<CouponForm>(initialForm);

  const adminAuthHeader = useMemo(() => buildBasicAuth(adminUsername, adminPassword), [adminUsername, adminPassword]);

  const fetchCustomerProducts = async (): Promise<void> => {
    try {
      setCustomerLoading(true);
      setCustomerError("");

      const res = await fetch(`${API_BASE}/store/products`);
      const data: Product[] | ApiError = await res.json();

      if (!res.ok) {
        const errData = data as ApiError;
        throw new Error(errData.message || errData.error || "Failed to fetch products");
      }

      setCustomerProducts(Array.isArray(data) ? data : []);
    } catch (err: unknown) {
      setCustomerError(getErrorMessage(err));
    } finally {
      setCustomerLoading(false);
    }
  };

  const fetchAdminProducts = async (): Promise<void> => {
    try {
      setAdminLoading(true);
      setAdminMessage("");

      const res = await fetch(`${API_BASE}/admin/products`, {
        headers: {
          Authorization: adminAuthHeader,
        },
      });

      const data: Product[] | ApiError = await res.json();

      if (!res.ok) {
        const errData = data as ApiError;
        throw new Error(errData.message || errData.error || "Failed to fetch admin products");
      }

      setAdminProducts(Array.isArray(data) ? data : []);
    } catch (err: unknown) {
      setAdminMessage(getErrorMessage(err));
      setAdminProducts([]);
    } finally {
      setAdminLoading(false);
    }
  };

  useEffect(() => {
    fetchCustomerProducts();
  }, []);

  useEffect(() => {
    if (mode === "admin") {
      fetchAdminProducts();
    }
  }, [mode]);

  const handleCustomerPurchase = async (productId: string): Promise<void> => {
    try {
      setCustomerError("");
      setPurchaseResult(null);

      const res = await fetch(`${API_BASE}/store/products/${productId}/purchase`, {
        method: "POST",
      });

      const data: PurchaseResult | ApiError = await res.json();

      if (!res.ok) {
        const errData = data as ApiError;
        throw new Error(errData.message || errData.error || "Purchase failed");
      }

      setPurchaseResult(data as PurchaseResult);
      fetchCustomerProducts();
    } catch (err: unknown) {
      setCustomerError(getErrorMessage(err));
    }
  };

  const handleFormChange = (e: ChangeEvent<HTMLInputElement | HTMLSelectElement>): void => {
    setAdminMessage("");
    setForm((prev) => ({
      ...prev,
      [e.target.name]: e.target.value,
    }));
  };

  const handleCreateCoupon = async (e: SubmitEvent<HTMLFormElement>): Promise<void> => {
    e.preventDefault();

    try {
      setAdminMessage("");

      const payload = {
        ...form,
        cost_price: Number(form.cost_price),
        margin_percentage: Number(form.margin_percentage),
      };

      const res = await fetch(`${API_BASE}/admin/coupons`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: adminAuthHeader,
        },
        body: JSON.stringify(payload),
      });

      const data = await res.json();

      if (!res.ok) {
        const errData = data as ApiError;
        throw new Error(errData.message || errData.error || "Failed to create coupon");
      }

      setAdminMessage("Coupon created successfully.");
      setForm(initialForm);
      fetchAdminProducts();
      fetchCustomerProducts();
    } catch (err: unknown) {
      setAdminMessage(getErrorMessage(err));
    }
  };

  const handleDeleteProduct = async (id: string): Promise<void> => {
    try {
      setAdminMessage("");

      const res = await fetch(`${API_BASE}/admin/products/${id}`, {
        method: "DELETE",
        headers: {
          Authorization: adminAuthHeader,
        },
      });

      if (!res.ok) {
        const data: ApiError = await res.json();
        throw new Error(data.message || data.error || "Failed to delete product");
      }

      setAdminMessage("Product deleted successfully.");
      fetchAdminProducts();
      fetchCustomerProducts();
    } catch (err: unknown) {
      setAdminMessage(getErrorMessage(err));
    }
  };

  return (
    <div className="page">
      <header className="hero">
        <h1>Digital Coupon Marketplace</h1>
        <p>Customer storefront + protected admin management panel.</p>
      </header>

      <div className="tabs">
        <button className={mode === "customer" ? "tab active" : "tab"} onClick={() => setMode("customer")}>
          Customer
        </button>
        <button className={mode === "admin" ? "tab active" : "tab"} onClick={() => setMode("admin")}>
          Admin
        </button>
      </div>

      {mode === "customer" ? (
        <section className="card">
          <div className="sectionHeader">
            <h2>Available Coupons</h2>
            <button onClick={fetchCustomerProducts}>Refresh</button>
          </div>

          {customerLoading ? (
            <p>Loading products...</p>
          ) : customerProducts.length === 0 ? (
            <p>No available coupons.</p>
          ) : (
            <div className="productList">
              {customerProducts.map((product) => (
                <div className="productCard" key={product.id}>
                  <img
                    src={product.image_url}
                    alt={product.name}
                    onError={(e) => {
                      e.currentTarget.src = "https://via.placeholder.com/300x200?text=Coupon";
                    }}
                  />
                  <div className="productBody">
                    <h3>{product.name}</h3>
                    <p>{product.description}</p>
                    <div className="productFooter">
                      <span>Price: ${product.price}</span>
                      <button onClick={() => handleCustomerPurchase(product.id)}>Buy Now</button>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}

          {customerError && <div className="message error">{customerError}</div>}

          {purchaseResult && (
            <div className="message success">
              <strong>Purchase successful</strong>
              <div>Product ID: {purchaseResult.product_id}</div>
              <div>Final Price: {purchaseResult.final_price}</div>
              <div>Value Type: {purchaseResult.value_type}</div>
              <div>Coupon Value: {purchaseResult.value}</div>
            </div>
          )}
        </section>
      ) : (
        <section className="card">
          <div className="sectionHeader">
            <h2>Admin Panel</h2>
            <button onClick={fetchAdminProducts}>Refresh</button>
          </div>

          <div className="adminAuthBox">
            <input
              placeholder="Admin username"
              value={adminUsername}
              onChange={(e) => setAdminUsername(e.target.value)}
            />
            <input
              type="password"
              placeholder="Admin password"
              value={adminPassword}
              onChange={(e) => setAdminPassword(e.target.value)}
            />
          </div>

          <form className="form" onSubmit={handleCreateCoupon}>
            <input name="name" placeholder="Name" value={form.name} onChange={handleFormChange} required />
            <input name="description" placeholder="Description" value={form.description} onChange={handleFormChange} />
            <input
              name="image_url"
              placeholder="Image URL"
              value={form.image_url}
              onChange={handleFormChange}
              required
            />
            <input
              name="cost_price"
              type="number"
              step="0.01"
              placeholder="Cost Price"
              value={form.cost_price}
              onChange={handleFormChange}
              required
            />
            <input
              name="margin_percentage"
              type="number"
              step="0.01"
              placeholder="Margin Percentage"
              value={form.margin_percentage}
              onChange={handleFormChange}
              required
            />
            <select name="value_type" value={form.value_type} onChange={handleFormChange}>
              <option value="STRING">STRING</option>
            </select>
            <input name="value" placeholder="Coupon Value" value={form.value} onChange={handleFormChange} required />

            <button type="submit">Create Coupon</button>
          </form>

          {adminMessage && <div className="message info">{adminMessage}</div>}

          <div className="adminListHeader">
            <h3>All Products</h3>
          </div>

          {adminLoading ? (
            <p>Loading admin products...</p>
          ) : adminProducts.length === 0 ? (
            <p>No products found.</p>
          ) : (
            <div className="adminTableWrapper">
              <table className="adminTable">
                <thead>
                  <tr>
                    <th>Name</th>
                    <th>Cost</th>
                    <th>Margin %</th>
                    <th>Min Price</th>
                    <th>Sold</th>
                    <th>Action</th>
                  </tr>
                </thead>
                <tbody>
                  {adminProducts.map((product) => (
                    <tr key={product.id}>
                      <td>{product.name}</td>
                      <td>{product.cost_price}</td>
                      <td>{product.margin_percentage}</td>
                      <td>{product.minimum_sell_price}</td>
                      <td>{product.is_sold ? "Yes" : "No"}</td>
                      <td>
                        <button className="dangerButton" onClick={() => handleDeleteProduct(product.id)}>
                          Delete
                        </button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </section>
      )}
    </div>
  );
}

export default App;
