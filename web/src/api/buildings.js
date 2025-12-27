import http from "./http";

export function listBuildings(params) {
  return http.get("/buildings", { params });
}

export function createBuilding(data) {
  return http.post("/buildings", data);
}

export function updateBuilding(id, data) {
  return http.put("/buildings/" + id, data);
}

export function deleteBuilding(id) {
  return http.delete("/buildings/" + id);
}
