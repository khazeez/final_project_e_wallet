 package pkg

// import (
// 	"database/sql"
// 	"net/http"

// 	"github.com/KhoirulAziz99/final_project_e_wallet/internal/domain"
// 	"github.com/kataras/go-sessions"
// 	"golang.org/x/crypto/bcrypt"
// )

// func QueryUser(username string) domain.User {
// 	var db *sql.DB
// 	var users = domain.User{}
// 	err := db.QueryRow(`
// 		SELECT id, 
// 		name, 
// 		email, 
// 		password
// 		FROM users WHERE username=?
// 		`, username).
// 		Scan(
// 			&users.ID,
// 			&users.Name,
// 			&users.Email,
// 			&users.Password,
// 		)
// 		if err != nil {
// 			panic(err)
// 		}
// 	return users
// }

// // func login(w http.ResponseWriter, r *http.Request) {
// // 	// session := sessions.Start(w, r)
// // 	// if len(session.GetString("username")) != 0 && checkErr(w, r, err) {
// // 	// 	http.Redirect(w, r, "/", 302)
// // 	// }
// // 	// if r.Method != "POST" {
// // 	// 	http.ServeFile(w, r, "views/login.html")
// // 	// 	return
// // 	}
// // 	// username := r.FormValue("username")
// // 	// password := r.FormValue("password")

// // 	// users := QueryUser(username)

// // 	//deskripsi dan compare password
// // 	var password_tes = bcrypt.CompareHashAndPassword([]byte(users.Password), []byte(password))

// // 	// if password_tes == nil {
// // 	// 	//login success
// // 	// 	session := sessions.Start(w, r)
// // 	// 	session.Set("username", users.Username)
// // 	// 	session.Set("name", users.FirstName)
// // 	// 	http.Redirect(w, r, "/", 302)
// // 	// } else {
// // 	// 	//login failed
// // 	// 	http.Redirect(w, r, "/login", 302)
// // 	// }

// // }
