import { useEffect, useState } from "react";
import "./App.css";

const API_BASE = "http://localhost:8080";

type Product = {
  id: string;
  name: string;
  description: string;
  image_url: string;
  price: number;
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

function getErrorMessage(err: unknown): string {
  if (err instanceof Error) {
    return err.message;
  }
  return "Something went wrong";
}

function App() {
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [purchaseResult, setPurchaseResult] = useState<PurchaseResult | null>(null);
  const [errorMessage, setErrorMessage] = useState<string>("");

  const fetchProducts = async (): Promise<void> => {
    try {
      setLoading(true);
      setErrorMessage("");

      const res = await fetch(`${API_BASE}/store/products`);
      const data: Product[] | ApiError = await res.json();

      if (!res.ok) {
        const errData = data as ApiError;
        throw new Error(errData.message || errData.error || "Failed to fetch products");
      }

      setProducts(Array.isArray(data) ? data : []);
    } catch (err: unknown) {
      setErrorMessage(getErrorMessage(err));
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchProducts();
  }, []);

  const handlePurchase = async (productId: string): Promise<void> => {
    try {
      setErrorMessage("");
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
      fetchProducts();
    } catch (err: unknown) {
      setErrorMessage(getErrorMessage(err));
    }
  };

  return (
    <div className="page">
      <header className="hero">
        <h1>Digital Coupon Marketplace</h1>
        <p>Buy digital coupons directly from our website.</p>
      </header>

      <section className="card">
        <div className="sectionHeader">
          <h2>Available Coupons</h2>
          <button onClick={fetchProducts}>Refresh</button>
        </div>

        {loading ? (
          <p>Loading products...</p>
        ) : products.length === 0 ? (
          <p>No available coupons.</p>
        ) : (
          <div className="productList">
            {products.map((product) => (
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
                    <button onClick={() => handlePurchase(product.id)}>Buy Now</button>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}

        {errorMessage && <div className="message error">{errorMessage}</div>}

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
    </div>
  );
}

export default App;
