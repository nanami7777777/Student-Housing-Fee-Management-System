import http from "./http";

export function listPayments(params) {
  return http.get("/payments", { params });
}

export function createPayment(data) {
  return http.post("/payments", data);
}

export function updatePayment(id, data) {
  return http.put("/payments/" + id, data);
}

export function deletePayment(id) {
  return http.delete("/payments/" + id);
}
