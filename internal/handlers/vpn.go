package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
)

type VPNLink struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Link string `json:"link"`
}

type InsertVPNRequest struct {
	Name string `json:"name"`
	Link string `json:"link"`
}

type UpdateVPNRequest struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Link string `json:"link"`
}

type DeleteVPNRequest struct {
	ID int `json:"id"`
}

func (h *Handler) InsertVPN(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req InsertVPNRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Link == "" {
		http.Error(w, "Name and link are required", http.StatusBadRequest)
		return
	}

	var id int
	err := h.db.QueryRow(context.Background(),
		"INSERT INTO vpn_links (name, link) VALUES ($1, $2) RETURNING id",
		req.Name, req.Link).Scan(&id)
	if err != nil {
		http.Error(w, "Error inserting link", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"id":     id,
	})
}

func (h *Handler) UpdateVPN(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req UpdateVPNRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.ID == 0 || req.Name == "" || req.Link == "" {
		http.Error(w, "ID, name and link are required", http.StatusBadRequest)
		return
	}

	result, err := h.db.Exec(context.Background(),
		"UPDATE vpn_links SET name = $1, link = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3",
		req.Name, req.Link, req.ID)
	if err != nil {
		http.Error(w, "Error updating link", http.StatusInternalServerError)
		return
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Link not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (h *Handler) DeleteVPN(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	result, err := h.db.Exec(context.Background(), "DELETE FROM vpn_links WHERE id = $1", id)
	if err != nil {
		http.Error(w, "Error deleting link", http.StatusInternalServerError)
		return
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Link not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (h *Handler) ListVPN(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rows, err := h.db.Query(context.Background(), "SELECT id, name, link FROM vpn_links ORDER BY id")
	if err != nil {
		http.Error(w, "Error fetching links", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var links []VPNLink
	for rows.Next() {
		var link VPNLink
		if err := rows.Scan(&link.ID, &link.Name, &link.Link); err != nil {
			http.Error(w, "Error scanning links", http.StatusInternalServerError)
			return
		}
		links = append(links, link)
	}

	if links == nil {
		links = []VPNLink{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(links)
}
