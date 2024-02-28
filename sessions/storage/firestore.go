package storage

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/dptsi/its-go/contracts"
	"github.com/dptsi/its-go/sessions"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FirestoreData struct {
	Data      map[string]interface{} `firestore:"data"`
	ExpiredAt time.Time              `firestore:"expired_at"`
	CSRFToken string                 `firestore:"csrf_token"`
}

type Firestore struct {
	client     *firestore.Client
	collection string
}

func NewFirestore(client *firestore.Client, collection string) *Firestore {
	return &Firestore{client, collection}
}

func (s *Firestore) Get(ctx context.Context, id string) (contracts.SessionData, error) {
	var data FirestoreData
	if uuid.Validate(id) != nil {
		return nil, nil
	}

	dsnap, err := s.client.Collection(s.collection).Doc(id).Get(ctx)
	if status.Code(err) == codes.NotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if err := dsnap.DataTo(&data); err != nil {
		return nil, err
	}

	sess := sessions.NewData(id, data.CSRFToken, data.Data, data.ExpiredAt)
	return sess, nil
}

func (s *Firestore) Save(ctx context.Context, data contracts.SessionData) error {
	fData := FirestoreData{data.Data(), data.ExpiredAt(), data.CSRFToken()}
	_, err := s.client.Collection(s.collection).Doc(data.Id()).
		Set(ctx, fData)

	return err
}

func (s *Firestore) Delete(ctx context.Context, id string) error {
	_, err := s.client.Collection(s.collection).Doc(id).Delete(ctx)
	return err
}
