import type { Product } from "../types/types";

type ProductCardProps = {
  product: Product;
  onBuy: (id: string) => void;
};

function ProductCard({ product, onBuy }: ProductCardProps) {
  return (
    <div className="productCard">
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
          <button onClick={() => onBuy(product.id)}>Buy Now</button>
        </div>
      </div>
    </div>
  );
}

export default ProductCard;
