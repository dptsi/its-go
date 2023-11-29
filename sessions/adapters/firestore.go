package adapters

import (
	"context"
	"time"

	"bitbucket.org/dptsi/base-go-libraries/sessions"
	"cloud.google.com/go/firestore"
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

func (f *Firestore) Get(ctx context.Context, id string) (*sessions.Data, error) {
	ref := f.client.Collection(f.collection).Doc(id)
	doc, err := ref.Get(ctx)
	if err != nil && status.Code(err) == codes.NotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var data FirestoreData
	if err := doc.DataTo(&data); err != nil {
		return nil, err
	}
	if data.ExpiredAt.Before(time.Now()) {
		return nil, nil
	}

	sess := sessions.NewData(id, data.CSRFToken, data.Data, data.ExpiredAt)
	return sess, nil
}

func (f *Firestore) Save(ctx context.Context, data *sessions.Data) error {
	ref := f.client.Collection(f.collection).Doc(data.Id())
	_, err := ref.Set(ctx, FirestoreData{data.Data(), data.ExpiredAt(), data.CSRFToken()})
	return err
}

func (f *Firestore) Delete(ctx context.Context, id string) error {
	ref := f.client.Collection(f.collection).Doc(id)
	_, err := ref.Delete(ctx)
	return err
}
