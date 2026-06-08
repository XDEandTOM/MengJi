package main

import (
	"database/sql"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

//go:embed dist/*
var staticFiles embed.FS

var db *sql.DB

func main() {
	initDB()
	initAdmin()
	http.HandleFunc("/api/", handleAPI)
	http.HandleFunc("/uploads/", handleUploads)
	http.HandleFunc("/", handleStatic)
	port := "3001"
	if len(os.Args) > 1 && os.Args[1] == "-port" && len(os.Args) > 2 {
		port = os.Args[2]
	}
	log.Println("Server on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func initDB() {
	dbPath := "meng.db"
	_, err := os.Stat(dbPath)
	os.MkdirAll("uploads", 0755)
	db, err = sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	db.Exec("PRAGMA journal_mode=WAL")
	db.Exec("PRAGMA foreign_keys=ON")
	tables := []string{
		`CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT UNIQUE NOT NULL, password TEXT NOT NULL, nickname TEXT DEFAULT '', avatar TEXT DEFAULT '', role TEXT DEFAULT 'user', created_at INTEGER DEFAULT 0)`,
		`CREATE TABLE IF NOT EXISTS notes (id TEXT PRIMARY KEY, content TEXT, created_at INTEGER, updated_at INTEGER, pinned INTEGER DEFAULT 0, tags TEXT DEFAULT '[]', username TEXT, avatar TEXT, nickname TEXT)`,
		`CREATE TABLE IF NOT EXISTS settings (key TEXT PRIMARY KEY, value TEXT)`,
	}
	for _, t := range tables {
		if _, err := db.Exec(t); err != nil {
			log.Fatal(err)
		}
	}
}

func initAdmin() {
	var count int
	db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if count == 0 {
		db.Exec("INSERT INTO users (username, password, nickname, role, created_at) VALUES (?, ?, ?, ?, ?)",
			"admin", hashPassword("admin"), "Admin", "admin", time.Now().UnixMilli())
		log.Println("Admin user created: admin / admin")
	}
}

func hashPassword(pwd string) string {
	h := 0
	for _, c := range pwd {
		h = h*31 + int(c)
	}
	return fmt.Sprintf("%x", h)
}

func checkPassword(input, stored string) bool {
	return hashPassword(input) == stored
}

func cors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func jsonResp(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func errResp(w http.ResponseWriter, msg string, code int) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func handleAPI(w http.ResponseWriter, r *http.Request) {
	cors(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(200)
		return
	}
	path := strings.TrimPrefix(r.URL.Path, "/api")
	switch {
	case strings.HasPrefix(path, "/auth/"):
		handleAuth(w, r, path)
	case strings.HasPrefix(path, "/notes"):
		handleNotes(w, r, path)
	case strings.HasPrefix(path, "/settings"):
		handleSettings(w, r)
	case strings.HasPrefix(path, "/admin/"):
		handleAdmin(w, r, path)
	default:
		errResp(w, "not found", 404)
	}
}

func handleAuth(w http.ResponseWriter, r *http.Request, path string) {
	switch {
	case path == "/auth/login" && r.Method == "POST":
		var body struct{ Username, Password string }
		json.NewDecoder(r.Body).Decode(&body)
		var storedPwd, role string
		err := db.QueryRow("SELECT password, role FROM users WHERE username=?", body.Username).Scan(&storedPwd, &role)
		if err != nil || !checkPassword(body.Password, storedPwd) {
			errResp(w, "用户名或密码错误", 401)
			return
		}
		var avatar, nickname string
		db.QueryRow("SELECT avatar, nickname FROM users WHERE username=?", body.Username).Scan(&avatar, &nickname)
		jsonResp(w, map[string]interface{}{"username": body.Username, "avatar": avatar, "nickname": nickname, "role": role})

	case path == "/auth/register" && r.Method == "POST":
		var body struct{ Username, Password string }
		json.NewDecoder(r.Body).Decode(&body)
		if len(body.Username) < 2 || len(body.Password) < 4 {
			errResp(w, "用户名至少2个字符，密码至少4个", 400)
			return
		}
		var allowReg string
		db.QueryRow("SELECT value FROM settings WHERE key='allow_register'").Scan(&allowReg)
		if allowReg == "false" {
			errResp(w, "注册已关闭", 403)
			return
		}
		_, err := db.Exec("INSERT INTO users (username, password, role, created_at) VALUES (?, ?, ?, ?)",
			body.Username, hashPassword(body.Password), "user", time.Now().UnixMilli())
		if err != nil {
			errResp(w, "用户名已存在", 409)
			return
		}
		jsonResp(w, map[string]string{"username": body.Username, "role": "user"})

	case path == "/auth/verify" && r.Method == "GET":
		username := r.URL.Query().Get("username")
		var avatar, nickname, role string
		err := db.QueryRow("SELECT avatar, nickname, role FROM users WHERE username=?", username).Scan(&avatar, &nickname, &role)
		jsonResp(w, map[string]interface{}{"valid": err == nil, "avatar": avatar, "nickname": nickname, "role": role})

	case path == "/auth/avatar" && r.Method == "PATCH":
		var body struct{ Username, Avatar string }
		json.NewDecoder(r.Body).Decode(&body)
		db.Exec("UPDATE users SET avatar=? WHERE username=?", body.Avatar, body.Username)
		jsonResp(w, map[string]string{"success": "ok"})

	case path == "/auth/nickname" && r.Method == "PATCH":
		var body struct{ Username, Nickname string }
		json.NewDecoder(r.Body).Decode(&body)
		var count int
		db.QueryRow("SELECT COUNT(*) FROM users WHERE nickname=? AND username!=?", body.Nickname, body.Username).Scan(&count)
		if count > 0 {
			errResp(w, "用户名已存在", 409)
			return
		}
		db.Exec("UPDATE users SET nickname=? WHERE username=?", body.Nickname, body.Username)
		jsonResp(w, map[string]interface{}{"success": true, "nickname": body.Nickname})

	case path == "/auth/app-icon" && r.Method == "PATCH":
		var body struct{ Username, AppIcon string }
		json.NewDecoder(r.Body).Decode(&body)
		db.Exec("UPDATE users SET avatar=? WHERE username=?", body.AppIcon, body.Username)
		jsonResp(w, map[string]string{"success": "ok"})

	case path == "/auth/password" && r.Method == "PATCH":
		var body struct{ Username, OldPassword, NewPassword string }
		json.NewDecoder(r.Body).Decode(&body)
		var storedPwd string
		db.QueryRow("SELECT password FROM users WHERE username=?", body.Username).Scan(&storedPwd)
		if !checkPassword(body.OldPassword, storedPwd) {
			errResp(w, "用户验证失败", 401)
			return
		}
		db.Exec("UPDATE users SET password=? WHERE username=?", hashPassword(body.NewPassword), body.Username)
		jsonResp(w, map[string]string{"success": "ok"})

	case path == "/auth/avatar/upload" && r.Method == "POST":
		file, header, err := r.FormFile("avatar")
		if err != nil {
			errResp(w, "文件读取失败", 400)
			return
		}
		defer file.Close()
		ext := filepath.Ext(header.Filename)
		name := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
		dst, _ := os.Create(filepath.Join("uploads", name))
		io.Copy(dst, file)
		dst.Close()
		jsonResp(w, map[string]interface{}{"success": true, "url": "/uploads/" + name})
	}
}

func handleNotes(w http.ResponseWriter, r *http.Request, path string) {
	switch {
	case path == "/notes" && r.Method == "GET":
		rows, err := db.Query("SELECT id, content, created_at, updated_at, pinned, tags, username, avatar, nickname FROM notes ORDER BY created_at DESC")
		if err != nil {
			errResp(w, err.Error(), 500)
			return
		}
		defer rows.Close()
		var notes []map[string]interface{}
		for rows.Next() {
			var id, content, username, tags, avatar, nickname string
			var createdAt, updatedAt int64
			var pinned int
			rows.Scan(&id, &content, &createdAt, &updatedAt, &pinned, &tags, &username, &avatar, &nickname)
			var tagList []string
			json.Unmarshal([]byte(tags), &tagList)
			notes = append(notes, map[string]interface{}{
				"id": id, "content": content, "createdAt": createdAt, "updatedAt": updatedAt,
				"pinned": pinned == 1, "tags": tagList, "username": username,
				"avatar": avatar, "nickname": nickname,
			})
		}
		if notes == nil {
			notes = []map[string]interface{}{}
		}
		jsonResp(w, notes)

	case path == "/notes" && r.Method == "POST":
		var n struct {
			Id, Content, Username, Avatar, Nickname string
			CreatedAt, UpdatedAt                    int64
			Tags                                    []string
		}
		json.NewDecoder(r.Body).Decode(&n)
		tagBytes, _ := json.Marshal(n.Tags)
		db.Exec("INSERT INTO notes (id, content, created_at, updated_at, pinned, tags, username, avatar, nickname) VALUES (?,?,?,?,0,?,?,?,?)",
			n.Id, n.Content, n.CreatedAt, n.UpdatedAt, string(tagBytes), n.Username, n.Avatar, n.Nickname)
		jsonResp(w, map[string]string{"success": "ok"})

	case strings.Contains(path, "/upload") && r.Method == "POST":
		file, header, err := r.FormFile("image")
		if err != nil {
			errResp(w, "文件读取失败", 400)
			return
		}
		defer file.Close()
		ext := filepath.Ext(header.Filename)
		name := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
		dst, _ := os.Create(filepath.Join("uploads", name))
		io.Copy(dst, file)
		dst.Close()
		jsonResp(w, map[string]interface{}{"success": true, "url": "/uploads/" + name})

	default:
		parts := strings.Split(strings.TrimPrefix(path, "/notes/"), "/")
		if len(parts) == 1 && r.Method == "PUT" {
			var body struct{ Content string; Tags []string; UpdatedAt int64 }
			json.NewDecoder(r.Body).Decode(&body)
			tagBytes, _ := json.Marshal(body.Tags)
			db.Exec("UPDATE notes SET content=?, tags=?, updated_at=? WHERE id=?", body.Content, string(tagBytes), body.UpdatedAt, parts[0])
			jsonResp(w, map[string]string{"success": "ok"})
		} else if len(parts) == 1 && r.Method == "DELETE" {
			db.Exec("DELETE FROM notes WHERE id=?", parts[0])
			jsonResp(w, map[string]string{"success": "ok"})
		} else if len(parts) == 2 && parts[1] == "pin" && r.Method == "PATCH" {
			db.Exec("UPDATE notes SET pinned = CASE WHEN pinned=0 THEN 1 ELSE 0 END WHERE id=?", parts[0])
			jsonResp(w, map[string]string{"success": "ok"})
		}
	}
}

func handleSettings(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		rows, _ := db.Query("SELECT key, value FROM settings")
		defer rows.Close()
		s := map[string]string{"site_title": "", "allow_register": "true", "site_favicon": ""}
		for rows.Next() {
			var k, v string
			rows.Scan(&k, &v)
			s[k] = v
		}
		jsonResp(w, s)
	} else if r.Method == "POST" {
		var body struct{ Key, Value string }
		json.NewDecoder(r.Body).Decode(&body)
		db.Exec("INSERT OR REPLACE INTO settings (key, value) VALUES (?,?)", body.Key, body.Value)
		jsonResp(w, map[string]string{"success": "ok"})
	}
}

func handleAdmin(w http.ResponseWriter, r *http.Request, path string) {
	switch {
	case path == "/admin/stats":
		var users, notes int
		db.QueryRow("SELECT COUNT(*) FROM users").Scan(&users)
		db.QueryRow("SELECT COUNT(*) FROM notes").Scan(&notes)
		jsonResp(w, map[string]int{"totalUsers": users, "totalNotes": notes})

	case path == "/admin/users":
		rows, _ := db.Query("SELECT id, username, nickname, avatar, role, created_at FROM users ORDER BY id")
		defer rows.Close()
		var users []map[string]interface{}
		for rows.Next() {
			var id int; var username, nickname, avatar, role string; var createdAt int64
			rows.Scan(&id, &username, &nickname, &avatar, &role, &createdAt)
			var memoCount int
			db.QueryRow("SELECT COUNT(*) FROM notes WHERE username=?", username).Scan(&memoCount)
			users = append(users, map[string]interface{}{
				"id": id, "username": username, "nickname": nickname, "avatar": avatar,
				"role": role, "createdAt": createdAt, "memoCount": memoCount,
			})
		}
		if users == nil { users = []map[string]interface{}{} }
		jsonResp(w, users)

	default:
		parts := strings.Split(strings.TrimPrefix(path, "/admin/users/"), "/")
		if len(parts) == 1 && r.Method == "DELETE" {
			var username string
			db.QueryRow("SELECT username FROM users WHERE id=?", parts[0]).Scan(&username)
			db.Exec("DELETE FROM notes WHERE username=?", username)
			db.Exec("DELETE FROM users WHERE id=?", parts[0])
			jsonResp(w, map[string]string{"success": "ok"})
		}
	}
}

func handleUploads(w http.ResponseWriter, r *http.Request) {
	filePath := strings.TrimPrefix(r.URL.Path, "/uploads/")
	fullPath := filepath.Join("uploads", filePath)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, fullPath)
}

func handleStatic(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		data, err := staticFiles.ReadFile("dist/index.html")
		if err != nil {
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(data)
		return
	}
	data, err := staticFiles.ReadFile("dist" + r.URL.Path)
	if err != nil {
		data, _ = staticFiles.ReadFile("dist/index.html")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(data)
		return
	}
	ext := filepath.Ext(r.URL.Path)
	mime := map[string]string{".js": "application/javascript", ".css": "text/css", ".png": "image/png", ".jpg": "image/jpeg", ".svg": "image/svg+xml", ".woff": "font/woff", ".woff2": "font/woff2", ".ico": "image/x-icon"}
	if m, ok := mime[ext]; ok {
		w.Header().Set("Content-Type", m)
	}
	w.Write(data)
}

