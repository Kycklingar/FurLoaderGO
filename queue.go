package main

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/kycklingar/FurLoaderGO/dli"
)

type queue struct {
	submissions []dli.Submission
}

var downloadPath = "downloads"

func (q *queue) add(submissions ...dli.Submission) {
	q.submissions = append(q.submissions, submissions...)
}

func (q *queue) start() {
	for _, s := range q.submissions {
		//TODO: Check the datastore if this submission already has been downloaded
		err := q.download(s)
		if err != nil {
			log.Println(err)
			continue
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

	file, err := os.OpenFile(
		filepath.Join(
			downloadPath,
			sub.SiteName(),
			sub.Folder(),
			sub.User().Name(),
			filepath.Base(sub.Filename()),
		),
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
