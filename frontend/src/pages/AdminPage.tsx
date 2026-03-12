import { useEffect, useState } from "react";
import type { FormEvent } from "react";
import { createCoupon, deleteAdminProduct, fetchAdminProducts } from "../api/api";
import type { CouponForm, Product } from "../types/types";

type AdminPageProps = {
  username: string;
  password: string;
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

function AdminPage({ username, password }: AdminPageProps) {
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState("");
  const [form, setForm] = useState<CouponForm>(initialForm);

  const loadProducts = async () => {
    try {
      setLoading(true);
      setMessage("");
      const data = await fetchAdminProducts(username, password);
      setProducts(data);
    } catch (err) {
      setMessage(err instanceof Error ? err.message : "Failed to fetch admin products");
      setProducts([]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadProducts();
  }, [username, password]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    setMessage("");
    setForm((prev) => ({
      ...prev,
      [e.target.name]: e.target.value,
    }));
  };

  const handleCreate = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    try {
      setMessage("");
      await createCoupon(username, password, form);
      setMessage("Coupon created successfully.");
      setForm(initialForm);
      loadProducts();
    } catch (err) {
      setMessage(err instanceof Error ? err.message : "Failed to create coupon");
    }
  };

  const handleDelete = async (id: string) => {
    try {
      setMessage("");
      await deleteAdminProduct(username, password, id);
      setMessage("Product deleted successfully.");
      loadProducts();
    } catch (err) {
      setMessage(err instanceof Error ? err.message : "Failed to delete product");
    }
  };

  return (
    <section className="card">
      <div className="sectionHeader">
        <h2>Admin Panel</h2>
        <button onClick={loadProducts}>Refresh</button>
      </div>

      <form className="form" onSubmit={handleCreate}>
        <input name="name" placeholder="Name" value={form.name} onChange={handleChange} required />
        <input name="description" placeholder="Description" value={form.description} onChange={handleChange} />
        <input name="image_url" placeholder="Image URL" value={form.image_url} onChange={handleChange} required />
        <input
          name="cost_price"
          type="number"
          step="0.01"
          placeholder="Cost Price"
          value={form.cost_price}
          onChange={handleChange}
          required
        />
        <input
          name="margin_percentage"
          type="number"
          step="0.01"
          placeholder="Margin Percentage"
          value={form.margin_percentage}
          onChange={handleChange}
          required
        />
        <select name="value_type" value={form.value_type} onChange={handleChange}>
          <option value="STRING">STRING</option>
        </select>
        <input name="value" placeholder="Coupon Value" value={form.value} onChange={handleChange} required />
        <button type="submit">Create Coupon</button>
      </form>

      {message && <div className="message info">{message}</div>}

      <div className="adminListHeader">
        <h3>All Products</h3>
      </div>

      {loading ? (
        <p>Loading admin products...</p>
      ) : products.length === 0 ? (
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
              {products.map((product) => (
                <tr key={product.id}>
                  <td>{product.name}</td>
                  <td>{product.cost_price}</td>
                  <td>{product.margin_percentage}</td>
                  <td>{product.minimum_sell_price}</td>
                  <td>{product.is_sold ? "Yes" : "No"}</td>
                  <td>
                    <button className="dangerButton" onClick={() => handleDelete(product.id)}>
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
  );
}

export default AdminPage;
