package memdb

import "go-news/pkg/storage"

type Store struct{}

func New() *Store {
	return new(Store)
}

func (s *Store) Tasks() ([]storage.Task, error) {
	return posts, nil
}

func (s *Store) AddTask(p storage.Task) error {
	posts = append(posts, p)
	return nil
}

func (s *Store) UpdateTask(p storage.Task) error {
	for i := range posts {
		if posts[i].ID == p.ID {
			posts[i].ResponsibleID = p.ResponsibleID
			posts[i].ResponsibleName = p.ResponsibleName
			posts[i].Context = p.Context
			posts[i].AssignedAt = p.AssignedAt
			posts[i].DueDate = p.DueDate
			return nil
		}
	}
	return nil
}

func (s *Store) DeleteTask(p storage.Task) error {
	for i := range posts {
		if posts[i].ID == p.ID {
			posts = append(posts[:i], posts[i+1:]...)
			return nil
		}
	}
	return nil
}

var posts = []storage.Task{
	{
		ID:              1,
		ResponsibleID:   10,
		ResponsibleName: "Иван",
		Context:         "Test 1 Content",
		AssignedAt:      0,
		DueDate:         0,
	},
	{
		ID:              2,
		ResponsibleID:   11,
		ResponsibleName: "Пётр",
		Context:         "Test 2 Content",
		AssignedAt:      0,
		DueDate:         0,
	},
}
