package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/kycklingar/FurLoaderGO/dli"
)

type queue struct {
	next chan dli.Submission
	stop chan bool

	threadCount int
}

func Queue() *queue {
	return &queue{next: make(chan dli.Submission), stop: make(chan bool)}
}

var downloadPath = "downloads"

// TODO: improve this
func (q *queue) addIncDL(call func(int) []dli.Submission) {
	var i = 0
	for {
		subs := call(i)
		i++
		fmt.Println(len(subs))
		if len(subs) <= 0 {
			break
		}

		for _, sub := range subs {
			q.next <- sub
		}

		time.Sleep(time.Second * 2)
	}
}

// stopThread will stop one downloading thread from the queue
func (q *queue) stopThread() {
	q.stop <- true
	q.threadCount--
}

// This will stop all running downloading threads from the queue
func (q *queue) stopAll() {
	for q.threadCount > 0 {
		q.stopThread()
	}
}

// startThread will start one downloading thread on the queue
func (q *queue) startThread() {
	// I'm still learning, no bully
	q.threadCount++
	for {
		select {
		case <-q.stop:
			fmt.Println("Stopping the queue")
			return
		case s := <-q.next:
			{
				fmt.Println("Downloading: ", s.ID())
				dbkey := s.SiteName() + s.ID()
				str := db.Get(dbkey)
				if str != "" {
					fmt.Printf("Found %s in database\n", dbkey)
					continue
				}

				time.Sleep(time.Second * 2)

				extra, err := s.GetDetails()
				if err != nil {
					log.Println(err)
					continue
				}

				for _, esub := range extra {
					fmt.Printf("Downloading extra submission %s\n", esub.ID())
					if db.Get(esub.SiteName()+esub.ID()) != "" {
						continue
					}
					err = q.download(esub)
					if err != nil {
						log.Println(err)
						continue
					}
					db.Store(esub.SiteName()+esub.ID(), esub.FileURL())
					time.Sleep(time.Second * 2)
				}

				err = q.download(s)
				if err != nil {
					log.Println(err)
					continue
				}

				db.Store(dbkey, s.FileURL())
			}
		}
	}
}

func (q *queue) download(sub dli.Submission) error {
	body, err := sub.Download()
	if err != nil {
		log.Println(err)
		return err
	}
	defer body.Close()

	path := filepath.Join(
		downloadPath,
		sub.SiteName(),
		sub.User().Name(),
		sub.Folder(),
		//filepath.Base(sub.Filename()),
	)

	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		log.Println(err)
		return err
	}

	path = filepath.Join(path, filepath.Base(sub.Filename()))

	file, err := os.OpenFile(
		path,
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
		0666,
	)
	if err != nil {
		log.Println(err)
		return err
	}
	defer file.Close()

	if _, err = io.Copy(file, body); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
