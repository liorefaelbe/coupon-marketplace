import { useEffect, useState } from "react";
import { fetchStoreProducts, purchaseStoreProduct } from "../api/api";
import ProductCard from "../components/ProductCard";
import type { Product, PurchaseResult } from "../types/types";

function CustomerPage() {
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);
  const [errorMessage, setErrorMessage] = useState("");
  const [purchaseResult, setPurchaseResult] = useState<PurchaseResult | null>(null);

  const loadProducts = async () => {
    try {
      setLoading(true);
      setErrorMessage("");
      const data = await fetchStoreProducts();
      setProducts(data);
    } catch (err) {
      setErrorMessage(err instanceof Error ? err.message : "Failed to fetch products");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadProducts();
  }, []);

  const handlePurchase = async (productId: string) => {
    try {
      setErrorMessage("");
      setPurchaseResult(null);
      const result = await purchaseStoreProduct(productId);
      setPurchaseResult(result);
      loadProducts();
    } catch (err) {
      setErrorMessage(err instanceof Error ? err.message : "Purchase failed");
    }
  };

  return (
    <section className="card">
      <div className="sectionHeader">
        <h2>Available Coupons</h2>
        <button onClick={loadProducts}>Refresh</button>
      </div>

      {loading ? (
        <p>Loading products...</p>
      ) : products.length === 0 ? (
        <p>No available coupons.</p>
      ) : (
        <div className="productList">
          {products.map((product) => (
            <ProductCard key={product.id} product={product} onBuy={handlePurchase} />
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
  );
}

export default CustomerPage;
