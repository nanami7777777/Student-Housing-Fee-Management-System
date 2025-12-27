import http from "./http";

export function getBuildingOccupancy() {
  return http.get("/stats/building-occupancy");
}

export function getBuildingPaymentSummary() {
  return http.get("/stats/building-payments");
}

