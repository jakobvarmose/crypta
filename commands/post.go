package commands

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	cid "github.com/ipfs/go-cid"

	"github.com/jakobvarmose/crypta/pathing"
	"github.com/jakobvarmose/crypta/transaction"
	"github.com/jakobvarmose/crypta/userstore"
)

func SetWriters(db transaction.Database, pageAddress string, writers []string) error {
	si, err := transaction.NewSigner(context.TODO(), db, pageAddress)
	if err != nil {
		return err
	}
	m := make(map[string]interface{})
	for _, writer := range writers {
		m[writer] = true
	}
	err = si.Root().Get("writers").Set(m)
	if err != nil {
		return err
	}
	return si.Commit(pageAddress)
}

func CreatePost(db transaction.Database, pageAddress, creatorAddress string, text string, attachments *pathing.Object) (interface{}, error) {
	si, err := transaction.NewSigner(context.TODO(), db, creatorAddress)
	if err != nil {
		return nil, err
	}
	genesis := transaction.New(context.TODO(), db, "")
	now := time.Now().UTC().Format(time.RFC3339)

	attachmentsData := make([]interface{}, 0)
	attachments2 := make([]interface{}, 0)
	attachments.EachSimple(func(_ *pathing.Object, attachment *pathing.Object) error {
		c, err := cid.Decode(attachment.Get("hash").String())
		if err != nil {
			fmt.Println(err)
			return nil
		}
		attachmentsData = append(attachmentsData, map[interface{}]interface{}{
			"t":    attachment.Get("t").String(),
			"link": c,
		})
		attachments2 = append(attachments2, map[string]interface{}{
			"t":    attachment.Get("t").String(),
			"hash": attachment.Get("hash").String(),
		})
		return nil
	})

	err = genesis.Root().Set(map[interface{}]interface{}{
		"t":           "POST",
		"page":        pageAddress,
		"creator":     creatorAddress,
		"text":        text,
		"time":        now,
		"attachments": attachmentsData,
	})
	if err != nil {
		return nil, err
	}
	genesisCid, err := cid.Decode(genesis.Hash())
	if err != nil {
		return nil, err
	}
	err = si.Root().Get("posts").Get(pageAddress).Get(genesis.Hash()).Get(genesis.Hash()).Set(genesisCid)
	if err != nil {
		return nil, err
	}
	go func() {
		err := si.Commit(creatorAddress)
		if err != nil {
			fmt.Println(err)
		}
	}()
	return map[string]interface{}{
		"genesis": map[string]interface{}{
			"hash": genesis.Hash(),
			"creator": map[string]interface{}{
				"address": creatorAddress,
				"name":    si.Root().Get("info").Get("name").String(),
			},
			"text":        text,
			"time":        now,
			"attachments": attachments2,
		},
		"hash": genesis.Hash(),
	}, nil
}

func CreateTextComment(db transaction.Database, pageAddress, creatorAddress, genesisHash, text string) (interface{}, error) {
	si, err := transaction.NewSigner(context.TODO(), db, creatorAddress)
	if err != nil {
		return nil, err
	}
	comment := transaction.New(context.TODO(), db, "")
	now := time.Now().UTC().Format(time.RFC3339)
	comment.Root().Set(map[string]interface{}{
		"t":       "COMMENT",
		"page":    pageAddress,
		"genesis": genesisHash,
		"creator": creatorAddress,
		"text":    text,
		"time":    now,
	})
	commentCid, err := cid.Decode(comment.Hash())
	if err != nil {
		return nil, err
	}
	err = si.Root().Get("posts").Get(pageAddress).Get(genesisHash).Get(comment.Hash()).Set(commentCid)
	if err != nil {
		return nil, err
	}
	go func() {
		err := si.Commit(creatorAddress)
		if err != nil {
			fmt.Println(err)
		}
	}()
	return map[string]interface{}{
		"creator": map[string]interface{}{
			"address": creatorAddress,
			"name":    si.Root().Get("info").Get("name").String(),
		},
		"text": text,
		"time": now,
	}, nil
}

type postSorter []map[string]interface{}

func (p postSorter) Len() int {
	return len(p)
}
func (p postSorter) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
func (p postSorter) Less(i, j int) bool {
	g1, _ := p[i]["genesis"].(map[string]interface{})
	g2, _ := p[j]["genesis"].(map[string]interface{})
	t1, _ := g1["time"].(string)
	t2, _ := g2["time"].(string)
	return t1 > t2
}

type commentSorter []map[string]interface{}

func (p commentSorter) Len() int {
	return len(p)
}
func (p commentSorter) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
func (p commentSorter) Less(i, j int) bool {
	t1, _ := p[i]["time"].(string)
	t2, _ := p[j]["time"].(string)
	return t1 < t2
}

func Home(us *userstore.Userstore, db transaction.Database, userAddr string) (map[string]interface{}, error) {
	var mu sync.Mutex
	posts := make([]map[string]interface{}, 0)
	userData, err := us.GetUser(userAddr)
	if err != nil {
		return nil, err
	}
	hasErrors := false
	var wg sync.WaitGroup
	err = userData.Get("subscriptions").EachSimple(func(key *pathing.Object, _ *pathing.Object) error {
		pageAddr := key.String()
		wg.Add(1)
		go func(pageAddr string) {
			defer wg.Done()
			res, err := PostList(us, db, pageAddr, userAddr)
			if err != nil {
				fmt.Println(err)
				hasErrors = true
				return
			}
			mu.Lock()
			defer mu.Unlock()
			for _, post := range res["posts"].([]map[string]interface{}) {
				posts = append(posts, post)
			}
		}(pageAddr)
		return nil
	})
	if err != nil {
		return nil, err
	}
	wg.Wait()
	sort.Sort(postSorter(posts))
	res := map[string]interface{}{
		"posts":     posts,
		"hasErrors": hasErrors,
	}
	return res, nil
}

func PostList(us *userstore.Userstore, db transaction.Database, addr, user string) (map[string]interface{}, error) {
	si, err := transaction.NewSigner(context.TODO(), db, addr)
	if err != nil {
		return nil, err
	}
	root := si.Root()
	ch := make(chan error)
	pageObj := map[string]interface{}{
		"address": addr,
		"name":    si.Root().Get("info").Get("name").String(),
	}
	posts2 := make(map[string]map[string]interface{}, 0)
	hasErrors := false
	count := 0
	writers := make(map[string]bool)
	var mu sync.Mutex
	err = root.Get("writers").EachSimple(func(writerAddress *transaction.Object, _ *transaction.Object) error {
		writers[writerAddress.String()] = true
		count++
		go func(writerAddress string) {
			si, err := transaction.NewSigner(context.TODO(), db, writerAddress)
			if err != nil {
				ch <- err
				return
			}
			creatorObj := map[string]interface{}{
				"address": writerAddress,
				"name":    si.Root().Get("info").Get("name").String(),
			}
			err = si.Root().Get("posts").Get(addr).EachSimple(func(genesisHash *transaction.Object, val *transaction.Object) error {
				return val.EachSimple(func(commentHash *transaction.Object, val2 *transaction.Object) error {
					if val2.Cid().String() != commentHash.String() {
						return nil
					}
					switch val2.Get("t").String() {
					case "POST":
						if commentHash.String() != genesisHash.String() {
							return nil
						}
						attachments := make([]interface{}, 0)
						val2.Get("attachments").EachSimple(func(_ *transaction.Object, attachment *transaction.Object) error {
							attachments = append(attachments, map[string]interface{}{
								"t":    attachment.Get("t").String(),
								"hash": attachment.Get("link").Cid().String(),
							})
							return nil
						})
						genesis := map[string]interface{}{
							"hash":        genesisHash.String(),
							"text":        val2.Get("text").String(),
							"time":        val2.Get("time").String(),
							"creator":     creatorObj,
							"attachments": attachments,
						}
						mu.Lock()
						defer mu.Unlock()
						if posts2[genesisHash.String()] == nil {
							posts2[genesisHash.String()] = map[string]interface{}{
								"comments": []map[string]interface{}{},
							}
						}
						posts2[genesisHash.String()]["hash"] = genesisHash.String()
						posts2[genesisHash.String()]["page"] = pageObj
						posts2[genesisHash.String()]["genesis"] = genesis
					case "COMMENT":
						comment := map[string]interface{}{
							"text":    val2.Get("text").String(),
							"time":    val2.Get("time").String(),
							"creator": creatorObj,
						}
						mu.Lock()
						defer mu.Unlock()
						if posts2[genesisHash.String()] == nil {
							posts2[genesisHash.String()] = map[string]interface{}{
								"comments": []map[string]interface{}{},
							}
						}
						posts2[genesisHash.String()]["comments"] = append(
							posts2[genesisHash.String()]["comments"].([]map[string]interface{}),
							comment,
						)
					}
					return nil
				})
			})
			if err != nil {
				ch <- err
			}
			ch <- nil
		}(writerAddress.String())
		return nil
	})
	if err != nil {
		fmt.Println(err)
		hasErrors = true
	}
	for i := 0; i < count; i++ {
		err := <-ch
		if err != nil {
			fmt.Println(err)
			hasErrors = true
		}
	}
	posts := make([]map[string]interface{}, 0)
	for _, post := range posts2 {
		sort.Sort(commentSorter(post["comments"].([]map[string]interface{})))
		posts = append(posts, post)
	}
	sort.Sort(postSorter(posts))
	infos, err := GetInfos(db, addr)
	if err != nil {
		fmt.Println(err)
		hasErrors = true
	}
	subscribed := false
	theyCanPost := false
	if user != "" {
		obj, err := us.GetUser(user)
		if err != nil {
			return nil, err
		}
		subscribed = obj.Get("subscriptions").Get(addr).Bool()
		userSi, err := transaction.NewSigner(context.TODO(), db, user)
		if err != nil {
			return nil, err
		}
		theyCanPost = userSi.Root().Get("writers").Get(addr).Bool()
	}
	res := map[string]interface{}{
		"address":     addr,
		"hash":        si.Hash(),
		"info":        infos,
		"posts":       posts,
		"posts2":      posts2,
		"writers":     writers,
		"subscribed":  subscribed,
		"theyCanPost": theyCanPost,
		"hasErrors":   hasErrors,
	}
	return res, nil
}
