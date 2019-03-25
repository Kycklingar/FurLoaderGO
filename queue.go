package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/kycklingar/FurLoaderGO/dli"
)

type queue struct {
	submissions []dli.Submission
}

func Queue(submissions []dli.Submission) *queue {
	return &queue{submissions: submissions}
}

var downloadPath = "downloads"

func (q *queue) add(submissions ...dli.Submission) {
	q.submissions = append(q.submissions, submissions...)
}

func (q *queue) start() {
	for _, s := range q.submissions {
		//TODO: Check the datastore if this submission already has been downloaded
		fmt.Println("Downloading:", s.ID())

		extra, err := s.GetDetails()
		if err != nil {
			log.Println(err)
			continue
		}

		for _, esub := range extra {
			err = q.download(esub)
			if err != nil {
				log.Println(err)
				continue
			}
		}

		err = q.download(s)
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
