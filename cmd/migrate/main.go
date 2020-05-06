package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/denisenkom/go-mssqldb"
	_ "github.com/denisenkom/go-mssqldb"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
	"time"
)

func main() {
	// format: server=%s;user id=%s;password=%s;port=%d;database=%s;
	sqlServerConn := os.Getenv("BLOGGY_SQL_SERVER_CONN_STR")
	if sqlServerConn == "" {
		log.Fatal("BLOGGY_SQL_SERVER_CONN_STR needs to be provided")
	}
	// Create connection pool
	db, err := sql.Open("sqlserver", sqlServerConn)
	if err != nil {
		log.Fatal("Error creating connection pool: ", err.Error())
	}
	defer func() {
		closeErr := db.Close()
		if closeErr != nil {
			log.Fatal(closeErr.Error())
		}
	}()
	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("Connected!\n")
	posts, readErr := ReadTables(db)
	if readErr != nil {
		log.Fatal(readErr.Error())
	}

	sort.SliceStable(posts, func(i, j int) bool {
		return posts[i].createdOnUtc.Unix() > posts[j].createdOnUtc.Unix()
	})

	for _, post := range posts {
		fmt.Printf("Processing '%s' (%s)\n", post.title, post.id)
		tags, err := ReadTags(db, post.id)
		if err != nil {
			log.Fatal(err.Error())
		}

		slugs, err := ReadSlugs(db, post.id)
		if err != nil {
			log.Fatal(err.Error())
		}

		var builder strings.Builder
		builder.WriteString(fmt.Sprintln("---"))
		d, err := yaml.Marshal(&struct {
			Title string `yaml:"title"`
			CreatedOn string `yaml:"created_at"`
			Tags []string `yaml:"tags"`
			Slugs []string `yaml:"slugs"`
		}{
			Title: post.title,
			CreatedOn: post.createdOnUtc.UTC().String(),
			Tags: func() []string {
				r := make([]string, len(tags))
				for i, t := range tags {
					r[i] = t.name
				}
				return r
			}(),
			Slugs: slugs,
		})
		if err != nil {
			log.Fatal(err.Error())
		}
		builder.Write(d)
		builder.WriteString(fmt.Sprintln("---"))

		err = ioutil.WriteFile(fmt.Sprintf("../../web/posts/%s_%s.md", post.createdOnUtc.Format("2006-01-02_15-04-05"), slugs[0]), []byte(builder.String()), 0755)
		if err != nil {
			fmt.Printf("Unable to write file: %v", err)
		}
	}
}

type post struct {
	id string
	title string
	abstract string
	content string
	createdOnUtc time.Time
}

type tag struct {
	name string
	slug string
}

func ReadSlugs(db *sql.DB, postID string) ([]string, error) {
	ctx := context.Background()

	// Check if database is alive.
	err := db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	// don't worry about sql injection, input is trusted
	tsql := fmt.Sprintf("SELECT Path FROM PostSlugEntity WHERE OwnedById = '%s' ORDER BY IsDefault DESC;", postID)

	// Execute query
	rows, err := db.QueryContext(ctx, tsql)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var result []string

	// Iterate through the result set.
	for rows.Next() {
		var slug string

		// Get values from row.
		err := rows.Scan(&slug)
		if err != nil {
			return nil, err
		}

		result = append(result, slug)
	}

	return result, nil
}

func ReadTags(db *sql.DB, postID string) ([]tag, error) {
	ctx := context.Background()

	// Check if database is alive.
	err := db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	// don't worry about sql injection, input is trusted
	tsql := fmt.Sprintf("SELECT t.Name, t.Slug FROM PostTagEntity pte INNER JOIN Tags t ON t.Name = pte.TagName WHERE pte.PostId = '%s';", postID)

	// Execute query
	rows, err := db.QueryContext(ctx, tsql)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var result []tag

	// Iterate through the result set.
	for rows.Next() {
		var name, slug string

		// Get values from row.
		err := rows.Scan(&name, &slug)
		if err != nil {
			return nil, err
		}

		result = append(result, tag{
			name: name,
			slug: slug,
		})
	}

	return result, nil
}

func ReadTables(db *sql.DB) ([]post, error) {
	ctx := context.Background()

	// Check if database is alive.
	err := db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	tsql := fmt.Sprintf("SELECT Id, Abstract, Content, CreatedOnUtc, Title FROM Posts;")

	// Execute query
	rows, err := db.QueryContext(ctx, tsql)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var result []post

	// Iterate through the result set.
	for rows.Next() {
		var id mssql.UniqueIdentifier
		var abstract, content, title string
		var createdOnUtc time.Time

		// Get values from row.
		err := rows.Scan(&id, &abstract, &content, &createdOnUtc, &title)
		if err != nil {
			return nil, err
		}

		result = append(result, post{
			id: id.String(),
			title: title,
			abstract: abstract,
			content: content,
			createdOnUtc: createdOnUtc,
		})
	}

	return result, nil
}