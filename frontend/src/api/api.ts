import type { ApiError, CouponForm, Product, PurchaseResult } from "../types/types";

const API_BASE = "http://localhost:8080";

function extractErrorMessage(data: ApiError): string {
  return data.message || data.error || "Request failed";
}

function buildBasicAuth(username: string, password: string): string {
  return `Basic ${btoa(`${username}:${password}`)}`;
}

export async function fetchStoreProducts(): Promise<Product[]> {
  const res = await fetch(`${API_BASE}/store/products`);
  const data: Product[] | ApiError = await res.json();

  if (!res.ok) {
    throw new Error(extractErrorMessage(data as ApiError));
  }

  return Array.isArray(data) ? data : [];
}

export async function purchaseStoreProduct(productId: string): Promise<PurchaseResult> {
  const res = await fetch(`${API_BASE}/store/products/${productId}/purchase`, {
    method: "POST",
  });

  const data: PurchaseResult | ApiError = await res.json();

  if (!res.ok) {
    throw new Error(extractErrorMessage(data as ApiError));
  }

  return data as PurchaseResult;
}

export async function fetchAdminProducts(username: string, password: string): Promise<Product[]> {
  const res = await fetch(`${API_BASE}/admin/products`, {
    headers: {
      Authorization: buildBasicAuth(username, password),
    },
  });

  const data: Product[] | ApiError = await res.json();

  if (!res.ok) {
    throw new Error(extractErrorMessage(data as ApiError));
  }

  return Array.isArray(data) ? data : [];
}

export async function createCoupon(username: string, password: string, form: CouponForm): Promise<void> {
  const payload = {
    ...form,
    cost_price: Number(form.cost_price),
    margin_percentage: Number(form.margin_percentage),
  };

  const res = await fetch(`${API_BASE}/admin/coupons`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: buildBasicAuth(username, password),
    },
    body: JSON.stringify(payload),
  });

  const data = await res.json();

  if (!res.ok) {
    throw new Error(extractErrorMessage(data as ApiError));
  }
}

export async function deleteAdminProduct(username: string, password: string, id: string): Promise<void> {
  const res = await fetch(`${API_BASE}/admin/products/${id}`, {
    method: "DELETE",
    headers: {
      Authorization: buildBasicAuth(username, password),
    },
  });

  if (res.status === 204) {
    return;
  }

  const data = await res.json();
  if (!res.ok) {
    throw new Error(extractErrorMessage(data as ApiError));
  }
}
