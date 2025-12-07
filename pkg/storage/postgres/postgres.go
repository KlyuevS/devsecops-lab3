package postgres

import (
 "context"
 "errors"
 "go-news/pkg/storage"

 "github.com/jackc/pgx/v4/pgxpool"
)

// Хранилище данных.
type Store struct {
 db *pgxpool.Pool
}

// Конструктор объекта хранилища.
func New(connectionString string) (*Store, error) {
 db, err := pgxpool.Connect(context.Background(), connectionString)
 if err != nil {
  return nil, err
 }
 s := Store{
  db: db,
 }
 return &s, nil
}

func (s *Store) Tasks() ([]storage.Task, error) {
 rows, err := s.db.Query(context.Background(), `
  SELECT
   id,
   responsible_id,
   responsible_name,
   context,
   assigned_at,
   due_date
  FROM posts
  ORDER BY id;
 `)
 if err != nil {
  return nil, err
 }
 var posts []storage.Task
 for rows.Next() {
  var p storage.Task
  err = rows.Scan(
   &p.ID,
   &p.ResponsibleID,
   &p.ResponsibleName,
   &p.Context,
   &p.AssignedAt,
   &p.DueDate,
  )
  if err != nil {
   return nil, err
  }
  posts = append(posts, p)
 }
 return posts, rows.Err()
}

func (s *Store) AddTask(p storage.Task) error {
 tx, err := s.db.Begin(context.Background())
 if err != nil {
  return err
 }
 defer tx.Rollback(context.Background())

 _, err = tx.Exec(context.Background(), `
  INSERT INTO posts (id, responsible_id, responsible_name, context, assigned_at, due_date)
  VALUES ($1, $2, $3, $4, $5, $6);
  `,
  p.ID,
  p.ResponsibleID,
  p.ResponsibleName,
  p.Context,
  p.AssignedAt,
  p.DueDate,
 )
 if err != nil {
  return err
 }

 return tx.Commit(context.Background())
}

func (s *Store) UpdateTask(p storage.Task) error {
 tx, err := s.db.Begin(context.Background())
 if err != nil {
  return err
 }
 defer tx.Rollback(context.Background())

 commandTag, err := tx.Exec(context.Background(), `
  UPDATE posts SET
   responsible_id = $1,
   responsible_name = $2,
   context = $3,
   due_date = $4
  WHERE id = $5;
  `,
  p.ResponsibleID,
  p.ResponsibleName,
  p.Context,
  p.DueDate,
  p.ID,
 )

 if err != nil {
  return err
 }

 if commandTag.RowsAffected() != 1 {
  return errors.New("no row found to update")
 }

 return tx.Commit(context.Background())
}

func (s *Store) DeleteTask(p storage.Task) error {
 tx, err := s.db.Begin(context.Background())
 if err != nil {
  return err
 }
 defer tx.Rollback(context.Background())

 commandTag, err := tx.Exec(context.Background(), `
  DELETE FROM posts
  WHERE id = $1;
  `,
  p.ID,
 )

 if err != nil {
  return err
 }

 if commandTag.RowsAffected() != 1 {
  return errors.New("no row found to delete")
 }

 return tx.Commit(context.Background())
}

/*var posts = []storage.Post{
 {
  ID:      1,
  Title:   "Effective Go",
  Content: "Go is a new language. Although it borrows ideas from existing languages, it has unusual properties that make effective Go programs different in character from programs written in its relatives. A straightforward translation of a C++ or Java program into Go is unlikely to produce a satisfactory result—Java programs are written in Java, not Go. On the other hand, thinking about the problem from a Go perspective could produce a successful but quite different program. In other words, to write Go well, it's important to understand its properties and idioms. It's also important to know the established conventions for programming in Go, such as naming, formatting, program construction, and so on, so that programs you write will be easy for other Go programmers to understand.",
 },
 {
  ID:      2,
  Title:   "The Go Memory Model",
  Content: "The Go memory model specifies the conditions under which reads of a variable in one goroutine can be guaranteed to observe values produced by writes to the same variable in a different goroutine.",
 },
}*/

