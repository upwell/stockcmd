package store

import (
	"encoding/json"

	"go.etcd.io/bbolt"
)

type Group struct {
	Name  string            `json:"name"`
	Codes map[string]string `json:"codes"`
}

func ListGroup() []string {
	groups := make([]string, 0, 32)

	DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(GroupBucketName))

		b.ForEach(func(k, v []byte) error {
			groups = append(groups, string(k))
			return nil
		})
		return nil
	})
	return groups
}

func CheckGroupExist(name string) bool {
	ret := false
	DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(GroupBucketName))
		if b.Get([]byte(name)) == nil {
			ret = false
		} else {
			ret = true
		}
		return nil
	})
	return ret
}

// TODO add lock
func AddGroup(name string) {
	if CheckGroupExist(name) {
		return
	}

	group := &Group{
		Name:  name,
		Codes: make(map[string]string),
	}
	group.Save()
}

func DeleteGroup(name string) {
	DB.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(GroupBucketName))
		b.Delete([]byte(name))
		return nil
	})
}

func GetGroup(name string) *Group {
	if !CheckGroupExist(name) {
		return nil
	}

	var group Group
	DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(GroupBucketName))
		groupBytes := b.Get([]byte(name))
		json.Unmarshal(groupBytes, &group)
		return nil
	})
	return &group
}

func (g *Group) AddStock(code string, name string) {
	g.Codes[code] = name
	g.Save()
}

func (g *Group) RemoveStock(code string) {
	delete(g.Codes, code)
	g.Save()
}

func (g *Group) Save() {
	DB.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(GroupBucketName))
		groupBytes, _ := json.Marshal(g)
		b.Put([]byte(g.Name), groupBytes)
		return nil
	})
}
