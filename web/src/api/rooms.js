import http from "./http";

export function listRooms(params) {
  return http.get("/rooms", { params });
}

export function createRoom(data) {
  return http.post("/rooms", data);
}

export function updateRoom(id, data) {
  return http.put("/rooms/" + id, data);
}

export function deleteRoom(id) {
  return http.delete("/rooms/" + id);
}
