package job

import (
	"context"
	"sync"
	"time"

	"jamdrop/logger"
	"jamdrop/model"
	"jamdrop/requestid"

	"cloud.google.com/go/firestore"
)

func (j *Job) ScanUserPlayers(ctx context.Context) error {
	docs, err := j.DB.Collection(model.CollectionUsers).Documents(ctx).GetAll()
	if err != nil {
		return err
	}

	now := time.Now()
	batch := j.DB.Batch()
	var wg sync.WaitGroup
	docChan := make(chan *firestore.DocumentSnapshot)
	for i := 0; i < j.App.Config.ScanWorkers; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for doc := range docChan {
				user := new(model.User)
				if err := doc.DataTo(user); err != nil {
					j.Logger.Err(err, "failed to read user data", logger.Fields{"user_id": doc.Ref.ID}, requestid.LogContext(ctx))
					continue
				}

				isPlaying, err := j.Spotify.GetCurrentlyPlaying(user)
				if err != nil {
					j.Logger.Err(err, "failed to get currently playing", requestid.LogContext(ctx), logger.Fields{
						"user_id":      user.ID,
						"access_token": user.AccessToken,
					})

					time.Sleep(10 * time.Millisecond)
					continue
				}

				if isPlaying {
					batch.Update(doc.Ref, []firestore.Update{{Path: "last_playing", Value: now}})
				}
			}
		}()
	}

	for _, doc := range docs {
		docChan <- doc
	}

	close(docChan)
	wg.Wait()

	if _, err := batch.Commit(ctx); err != nil {
		if err.Error() == "firestore: cannot commit empty WriteBatch" {
			j.Logger.Error("empty batch")
			return nil
		}

		j.Logger.Err(err, "error committing batch")
		return err
	}

	return nil
}
