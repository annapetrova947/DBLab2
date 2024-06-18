
package main

import (
    "database/sql"
    "fmt"
    "html/template"
    "log"
    "net/http"

    _ "github.com/go-sql-driver/mysql"
)

type Sector struct {
    ID                int
    Coordinates       string
    LightIntensity    float64
    ForeignObjects    int
    StarObjectsCount  int
    UnknownObjectsCount int
    DefinedObjectsCount int
    Notes             string
}

type Object struct {
    ID       int
    Type     string
    Galaxy   string
    Accuracy float64
}

type Result struct {
    Sectors []Sector
    Objects []Object
}

func dbConn() (db *sql.DB) {
    dbDriver := "mysql"
    dbUser := "user"
    dbPass := "password"
    dbName := "Observatory"
    db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp(127.0.0.1:3306)/"+dbName)
    if err != nil {
        log.Fatal(err)
    }
    return db
}

func main() {
    http.HandleFunc("/", index)
    http.HandleFunc("/delete", deleteRecord)
    http.HandleFunc("/edit", editRecord)
    log.Println("Server started on: http://localhost:8080")
    http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    defer db.Close()

    // Call the stored procedure JoinTables
    rows, err := db.Query("CALL JoinTables(?, ?)", "Sector", "NaturalObjects")
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    var sectors []Sector
    var objects []Object

    // Fetch results and populate the structs
    for rows.Next() {
        var s Sector
        var o Object
        err = rows.Scan(&s.ID, &s.Coordinates, &s.LightIntensity, &s.ForeignObjects,
            &s.StarObjectsCount, &s.UnknownObjectsCount, &s.DefinedObjectsCount,
            &s.Notes, &o.ID, &o.Type, &o.Galaxy, &o.Accuracy)
        if err != nil {
            log.Fatal(err)
        }
        sectors = append(sectors, s)
        objects = append(objects, o)
    }

    res := Result{Sectors: sectors, Objects: objects}

    // Parse and execute the template
    tmpl, err := template.ParseFiles("view/index.html")
    if err != nil {
        log.Fatal(err)
    }
    tmpl.Execute(w, res)
}

func deleteRecord(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Query().Get("id")
    db := dbConn()
    defer db.Close()

    // Execute delete query
    _, err := db.Exec("DELETE FROM Sector WHERE id = ?", id)
    if err != nil {
        log.Fatal(err)
    }

    // Redirect to the index page
    http.Redirect(w, r, "/", 302)
}

func editRecord(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        id := r.FormValue("id")
        coordinates := r.FormValue("coordinates")
        lightIntensity := r.FormValue("light_intensity")
        foreignObjects := r.FormValue("foreign_objects")
        starObjectsCount := r.FormValue("star_objects_count")
        unknownObjectsCount := r.FormValue("unknown_objects_count")
        definedObjectsCount := r.FormValue("defined_objects_count")
        notes := r.FormValue("notes")

        db := dbConn()
        defer db.Close()

        // Execute update query
        _, err := db.Exec(`UPDATE Sector SET coordinates=?, light_intensity=?, foreign_objects=?, 
            star_objects_count=?, unknown_objects_count=?, defined_objects_count=?, notes=? WHERE id=?`,
            coordinates, lightIntensity, foreignObjects, starObjectsCount, unknownObjectsCount, definedObjectsCount, notes, id)
        if err != nil {
            log.Fatal(err)
        }

        // Redirect to the index page
        http.Redirect(w, r, "/", 302)
    } else {
        // Render edit form
        id := r.URL.Query().Get("id")
        db := dbConn()
        defer db.Close()

        var s Sector
        err := db.QueryRow("SELECT * FROM Sector WHERE id=?", id).Scan(&s.ID, &s.Coordinates, &s.LightIntensity, &s.ForeignObjects, 
            &s.StarObjectsCount, &s.UnknownObjectsCount, &s.DefinedObjectsCount, &s.Notes)
        if err != nil {
            log.Fatal(err)
        }

        tmpl, err := template.ParseFiles("view/edit.html")
        if err != nil {
            log.Fatal(err)
        }
        tmpl.Execute(w, s)
    }
}
