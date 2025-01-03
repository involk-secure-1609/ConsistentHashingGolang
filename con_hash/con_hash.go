package conhash

import (
	
	"errors"
	"fmt"
	"hash/crc32"
	"log"
	"strconv"

	rbt "github.com/emirpasic/gods/trees/redblacktree"
	"github.com/emirpasic/gods/utils"
)
const MOD = 1 << 32 // 2^32

func hash(node_to_insert string) int32 {

	hash := (crc32.ChecksumIEEE([]byte(node_to_insert)))
	// log.Println(hash)
	if int32(hash)<0{
		return (int32(hash)+2147483647)
	}
	return int32(hash)
}


type ConHash struct {
	virtualFactor int
	servers       []string
	tree          *rbt.Tree
}

func NewConhash(virtualFactor int, servers []string) *ConHash {
	return &ConHash{
		servers: servers,
		virtualFactor: virtualFactor,
	}
}
func (conHash *ConHash) Initialize() error {

	conHash.tree = rbt.NewWith(utils.Int32Comparator)
	for _, s := range conHash.servers {
		for i := 0; i < conHash.virtualFactor; i++ {
			node_to_insert := s +"#"+strconv.Itoa(i)
			key:=hash(node_to_insert);
			// log.Println(key)
			conHash.tree.Put(key, s)
		}
	}
	return nil
}

func (conHash *ConHash) LocateNode(node string) (string, error) {

	server, largest := conHash.tree.Ceiling(hash(node))
	if !largest {
		iterator := conHash.tree.Iterator()
		tree_is_empty := iterator.Next()
		if !tree_is_empty {
			return "", errors.New("no servers availible at the current moment")
		}
		server = iterator.Node()
	}

	return fmt.Sprintf("%v", server.Value), nil
}
func (conHash *ConHash) AddServer(server string) error {

	conHash.servers = append(conHash.servers, server)
	for i := 0; i < conHash.virtualFactor; i++ {
		node_to_insert := server + strconv.Itoa(i)
		key:=hash(node_to_insert);
		log.Println(key)
		conHash.tree.Put(hash(node_to_insert), server)
	}
	return nil
}

func (conHash *ConHash) Layout(){
	for _,key:=range(conHash.tree.Keys()){
		value,ok:=conHash.tree.Get(key)
		key1,ok:=key.(int32)
		if !ok{
			log.Fatalln("notOk")
		}
		log.Println(key1,value)
	}
}

func (conHash *ConHash) RemoveServer(server string) error {

	for i := 0; i < conHash.virtualFactor; i++ {
		node_to_insert := server + strconv.Itoa(i)
		conHash.tree.Remove(hash(node_to_insert))
	}
	index := -1
	for i, v := range conHash.servers {
		if v == server {
			index = i
			break
		}
	}

	if index != -1 {
		conHash.servers = append(conHash.servers[:index], conHash.servers[index+1:]...)
	}

	return nil
}
