package handlers

import (
	"context"
	"encoding/json"
	"log"
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
	log.Printf("[HANDLER] InsertVPN called")

	if r.Method != http.MethodPost {
		log.Printf("[HANDLER] InsertVPN: method not allowed: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req InsertVPNRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("[HANDLER] InsertVPN: invalid request body: %v", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Link == "" {
		log.Printf("[HANDLER] InsertVPN: name or link is empty")
		http.Error(w, "Name and link are required", http.StatusBadRequest)
		return
	}

	log.Printf("[HANDLER] InsertVPN: inserting link name='%s'", req.Name)
	var id int
	err := h.db.QueryRow(context.Background(),
		"INSERT INTO vpn_links (name, link) VALUES ($1, $2) RETURNING id",
		req.Name, req.Link).Scan(&id)
	if err != nil {
		log.Printf("[HANDLER] InsertVPN: database error: %v", err)
		http.Error(w, "Error inserting link", http.StatusInternalServerError)
		return
	}

	log.Printf("[HANDLER] InsertVPN: success, id=%d", id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"id":     id,
	})
}

func (h *Handler) UpdateVPN(w http.ResponseWriter, r *http.Request) {
	log.Printf("[HANDLER] UpdateVPN called")

	if r.Method != http.MethodPut {
		log.Printf("[HANDLER] UpdateVPN: method not allowed: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req UpdateVPNRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("[HANDLER] UpdateVPN: invalid request body: %v", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.ID == 0 || req.Name == "" || req.Link == "" {
		log.Printf("[HANDLER] UpdateVPN: missing required fields id=%d, name='%s'", req.ID, req.Name)
		http.Error(w, "ID, name and link are required", http.StatusBadRequest)
		return
	}

	log.Printf("[HANDLER] UpdateVPN: updating link id=%d, name='%s'", req.ID, req.Name)
	result, err := h.db.Exec(context.Background(),
		"UPDATE vpn_links SET name = $1, link = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3",
		req.Name, req.Link, req.ID)
	if err != nil {
		log.Printf("[HANDLER] UpdateVPN: database error: %v", err)
		http.Error(w, "Error updating link", http.StatusInternalServerError)
		return
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		log.Printf("[HANDLER] UpdateVPN: link not found id=%d", req.ID)
		http.Error(w, "Link not found", http.StatusNotFound)
		return
	}

	log.Printf("[HANDLER] UpdateVPN: success, rows affected=%d", rowsAffected)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (h *Handler) DeleteVPN(w http.ResponseWriter, r *http.Request) {
	log.Printf("[HANDLER] DeleteVPN called")

	if r.Method != http.MethodDelete {
		log.Printf("[HANDLER] DeleteVPN: method not allowed: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		log.Printf("[HANDLER] DeleteVPN: ID parameter missing")
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("[HANDLER] DeleteVPN: invalid ID: %v", err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	log.Printf("[HANDLER] DeleteVPN: deleting link id=%d", id)
	result, err := h.db.Exec(context.Background(), "DELETE FROM vpn_links WHERE id = $1", id)
	if err != nil {
		log.Printf("[HANDLER] DeleteVPN: database error: %v", err)
		http.Error(w, "Error deleting link", http.StatusInternalServerError)
		return
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		log.Printf("[HANDLER] DeleteVPN: link not found id=%d", id)
		http.Error(w, "Link not found", http.StatusNotFound)
		return
	}

	log.Printf("[HANDLER] DeleteVPN: success, rows affected=%d", rowsAffected)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (h *Handler) ListVPN(w http.ResponseWriter, r *http.Request) {
	log.Printf("[HANDLER] ListVPN called")

	if r.Method != http.MethodGet {
		log.Printf("[HANDLER] ListVPN: method not allowed: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	log.Printf("[HANDLER] ListVPN: querying database")
	rows, err := h.db.Query(context.Background(), "SELECT id, name, link FROM vpn_links ORDER BY id")
	if err != nil {
		log.Printf("[HANDLER] ListVPN: database error: %v", err)
		http.Error(w, "Error fetching links", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var links []VPNLink
	for rows.Next() {
		var link VPNLink
		if err := rows.Scan(&link.ID, &link.Name, &link.Link); err != nil {
			log.Printf("[HANDLER] ListVPN: error scanning row: %v", err)
			http.Error(w, "Error scanning links", http.StatusInternalServerError)
			return
		}
		links = append(links, link)
	}

	if links == nil {
		links = []VPNLink{}
	}

	log.Printf("[HANDLER] ListVPN: success, returning %d links", len(links))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(links)
}
