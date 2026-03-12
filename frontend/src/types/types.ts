export type Product = {
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

export type PurchaseResult = {
  product_id: string;
  final_price: number;
  value_type: string;
  value: string;
};

export type ApiError = {
  error_code?: string;
  message?: string;
  error?: string;
};

export type CouponForm = {
  name: string;
  description: string;
  image_url: string;
  cost_price: string;
  margin_percentage: string;
  value_type: string;
  value: string;
};
