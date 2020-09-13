package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)

type post struct {
	Id       int
	Title    string
	Page     string
	MarkDown string
	Created  time.Time
}

/*
CREATE TABLE IF NOT EXISTS "posts"
("id" integer not null primary key autoincrement, "uuid" varchar(36) not null,
"title" varchar(150) not null, "slug" varchar(150) not null,
"markdown" text null, "html" text null, "image" text null, "featured" boolean not null default '0', "page" boolean not null default '0',
"status" varchar(150) not null default 'draft', "language" varchar(6) not null default 'en_US',
"meta_title" varchar(150) null,
"meta_description" varchar(200) null, "author_id" integer not null, "created_at" datetime not null, "created_by" integer not null,
"updated_at" datetime null, "updated_by" integer null, "published_at" datetime null, "published_by" integer null,
"visibility" varchar(150) not null default 'public', "mobiledoc" text null, "amp" text null);
*/

func genPostFile(pPost post, pPath string) {

	// Clean up title
	re, _ := regexp.Compile(`[^\w]`)
	t2 := re.ReplaceAllString(pPost.Title, " ")

	re, _ = regexp.Compile(`\ {2,}`)
	t2 = re.ReplaceAllString(t2, "_")

	t2 = strings.Replace(t2, " ", "_", -1)
	t2 = pPost.Created.Format("20060102_") + t2

	t2 = strings.ToLower(t2)
	fic := pPath + "/" + t2 + ".md"

	// Write the file
	log.Println("")
	f, err := os.Create(fic)
	if err != nil {
		fmt.Println("Skipping " + t2 + "error: " + err.Error())
	}
	w := bufio.NewWriter(f)

	/*
	   +++
	   date = "2019-02-27T21:30:00+01:00"
	   draft = false
	   title = "More blogging"
	   tags = ["this", "post", "also", "has", "some", "tags"]

	   +++
	*/
	fmt.Fprintln(w, "---")
	fmt.Fprintln(w, "title: \""+pPost.Title+"\"")
	fmt.Fprintln(w, "date: "+pPost.Created.Format("2006-01-02T15:04:05+02:00"))
	fmt.Fprintln(w, "draft: false")
	fmt.Fprintln(w, "---")
	fmt.Fprintln(w, pPost.MarkDown)

	w.Flush()
	f.Close()

	log.Println("...Generating " + fic)
}

func migratePosts(db *sql.DB, rootPath string) {

	row, err := db.Query("SELECT id, created_at, title, markdown FROM posts ORDER by id  DESC ")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	for row.Next() { // Iterate and fetch the records from result cursor
		var p post
		row.Scan(&p.Id, &p.Created, &p.Title, &p.Page, &p.MarkDown)
		genPostFile(p, rootPath)
	}
}

func main() {

	log.Println("Opening ghost database.")
	sqliteDatabase, _ := sql.Open("sqlite3", "./ghost.db")
	defer sqliteDatabase.Close() // Defer Closing the database

	log.Println("Starting migration")

	r := "./tmp"
	migratePosts(sqliteDatabase, r)

	log.Println("Migration done")
}
