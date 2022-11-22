package grpc_play

// Following guide
// https://thedevelopercafe.com/articles/rpc-in-go-using-twitchs-twirp-3dcb78ece775

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/twitchtv/twirp"
	"ttin.com/play2022/grpc_play/notes"
)

//---------------------------------------------
//	Service imp
//---------------------------------------------

type notesService struct {
	Notes     []notes.Note
	CurrentId int32
}

func (s *notesService) GetAllNotes(ctx context.Context, params *notes.GetAllNotesParams) (*notes.GetAllNotesResult, error) {
	allNotes := make([]*notes.Note, 0)

	for _, note := range s.Notes {
		n := note
		allNotes = append(allNotes, &n)
	}

	return &notes.GetAllNotesResult{
		Notes: allNotes,
	}, nil
}

func (s *notesService) CreateNote(ctx context.Context, params *notes.CreateNoteParams) (*notes.Note, error) {
	if len(params.Text) < 4 {
		return nil, twirp.InvalidArgument.Error("Text should be min 4 characters.")
	}

	note := notes.Note{
		Id:        s.CurrentId,
		Text:      params.Text,
		CreatedAt: time.Now().UnixMilli(),
	}

	s.Notes = append(s.Notes, note)

	s.CurrentId++

	return &note, nil
}

//---------------------------------------------
//	Play area
//---------------------------------------------

func Play_Grpc_Twirp() {
	wg := new(sync.WaitGroup)
	go Grpc_Twirp_HTTP_Server(wg)

	time.Sleep(2 * time.Second)

	go Grpc_Twirp_HTTP_Client(wg)

	time.Sleep(2 * time.Second)

	wg.Wait()
}

func Grpc_Twirp_HTTP_Server(wg *sync.WaitGroup) {
	notesServer := notes.NewNotesServiceServer(&notesService{})
	mux := http.NewServeMux()
	mux.Handle(notesServer.PathPrefix(), notesServer)

	log.Println("Http server listen at :8000")
	err := http.ListenAndServe(":8000", notesServer)
	if err != nil {
		panic(err)
	}

}

func Grpc_Twirp_HTTP_Client(wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	httpClient := http.Client{}
	ctx := context.Background()

	client := notes.NewNotesServiceProtobufClient("http://localhost:8000", &httpClient)

	tmp, err := client.CreateNote(ctx, &notes.CreateNoteParams{Text: "Hello World 1"})
	log.Println("CreatedNote:", err, tmp)
	if err != nil {
		log.Fatal(err)
	}

	tmp, err = client.CreateNote(ctx, &notes.CreateNoteParams{Text: "Hello World 2"})
	log.Println("CreatedNote:", err, tmp)
	if err != nil {
		log.Fatal(err)
	}

	tmpNotes, err := client.GetAllNotes(ctx, &notes.GetAllNotesParams{})
	log.Println("ListNotes:", err, tmpNotes)
	if err != nil {
		log.Fatal(err)
	}
}
