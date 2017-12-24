package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/reiver/go-porterstemmer"
)

func Handle(er error) {
	if er != nil {
		log.Fatal(er)
	}
}

type Indexer struct {
	db *sql.DB
}

func (obj *Indexer) OpenCon() {

	USERNAME, PASSWORD, HOST, DBNAME := os.Getenv("MYSQL_USERNAME"), os.Getenv("MYSQL_PASSWORD"), os.Getenv("MYSQL_HOST"), os.Getenv("MYSQL_DBNAME")

	var er error
	obj.db, er = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", USERNAME, PASSWORD, HOST, DBNAME))
	Handle(er)
}

func (obj *Indexer) CloseCon() {
	obj.db.Close()
}

func (obj *Indexer) GetWordsFreq(file *os.File) map[string]int {

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	words := make(map[string]int)

	for scanner.Scan() {
		linetext := scanner.Text()
		line := strings.Split(linetext, " ")

		for _, word := range line {
			stemmed := porterstemmer.StemString(word)
			words[stemmed]++
		}

	}

	return words
}

func (obj *Indexer) InsertDoc(file *os.File) int {

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	cnt := 0
	var url, title, summary string

	for scanner.Scan() {
		linetext := scanner.Text()
		cnt += 1
		first := 0
		if cnt == 1 {
			url = linetext
		} else if cnt == 2 {
			title = linetext
		} else {
			line := strings.Split(linetext, " ")
			for _, word := range line {

				if len(summary) < 300 {
					if first == 1 {
						summary += " "
					}
					summary += word
					first = 1
				}
			}
		}
	}

	er := obj.db.Ping()
	Handle(er)

	ins, er := obj.db.Prepare("insert into documents(url, title, summary) values (?, ?, ?)")
	Handle(er)
	defer ins.Close()

	res, er := ins.Exec(url, title, summary)
	Handle(er)

	doc, er := res.LastInsertId()
	Handle(er)

	return int(doc)
}

func (obj *Indexer) GetWordID(word string) int {

	stmt, er := obj.db.Prepare("select id from words where name = ?")
	Handle(er)
	defer stmt.Close()

	id := -1
	er = stmt.QueryRow(word).Scan(&id)

	if er != nil {
		return -1

	} else {
		return id
	}

}

func (obj *Indexer) ProcFile(path string) {

	file, er := os.Open(path)
	Handle(er)
	defer file.Close()

	docID := obj.InsertDoc(file)
	file.Seek(0, 0)

	freq := obj.GetWordsFreq(file)

	wordsIDs := obj.InsertWords(freq)

	obj.InsertWordsToDoc(freq, wordsIDs, docID)

}

func (obj *Indexer) InsertWordsToDoc(freq, id map[string]int, doc int) {

	er := obj.db.Ping()
	Handle(er)

	ins, er := obj.db.Prepare("insert into words_documents(word_id, document_id, freq) values (?, ?, ?)")
	Handle(er)
	defer ins.Close()

	for word, f := range freq {
		ins.Exec(id[word], doc, f)
	}

}

func (obj *Indexer) InsertWords(mp map[string]int) map[string]int {

	er := obj.db.Ping()
	Handle(er)

	ins, er := obj.db.Prepare("insert into words(name) values (?)")
	Handle(er)
	defer ins.Close()

	ids := make(map[string]int)

	for word, f := range mp {
		_ = f

		id := obj.GetWordID(word)

		if id == -1 {

			res, er := ins.Exec(word)
			Handle(er)

			idd, er := res.LastInsertId()
			Handle(er)
			id = int(idd)
		}

		ids[word] = id
	}

	return ids
}

func (obj *Indexer) ReadFolder(path string) {

	dir, er := os.Open(path)
	Handle(er)
	defer dir.Close()

	filenames, er := dir.Readdirnames(1000)
	Handle(er)

	for _, name := range filenames {
		path := path + "/" + name
		obj.ProcFile(path)
	}

}
