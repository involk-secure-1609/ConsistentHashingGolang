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

/*
function which is used for hashing the nodes
however currently the uniformity of the hash is not uniform
*/
func hash(node_to_insert string) int32 {

	hash := (crc32.ChecksumIEEE([]byte(node_to_insert)))
	// log.Println(hash)
	if int32(hash)<0{
		return (int32(hash)+2147483647)
	}
	return int32(hash)
}


type ConHash struct {
	// the number of virtual nodes for each server
	virtualFactor int
	//number of physical nodes
	physicalNodes       []string
	// the underlying RedBlackTree which is used for storing the hash
	tree          *rbt.Tree
}

// returns a NewConHash
func NewConhash(virtualFactor int, physicalNodes []string) *ConHash {
	return &ConHash{
		physicalNodes: physicalNodes,
		virtualFactor: virtualFactor,
	}
}

// Initializes the ConHash struct by hashing the virtualNodes and inseritng them into the RedBlackTree
func (conHash *ConHash) Initialize() error {

	conHash.tree = rbt.NewWith(utils.Int32Comparator)
	for _, s := range conHash.physicalNodes {
		for i := 0; i < conHash.virtualFactor; i++ {
			node_to_insert := s +"#"+strconv.Itoa(i)
			key:=hash(node_to_insert);
			conHash.tree.Put(key, s)
		}
	}
	return nil
}

// Locates the next greatest node for the key in the RedBlackTree
// if there is no next greatest node we return the smalles node in the tree
// else we return an error
func (conHash *ConHash) LocateNode(key string) (string, error) {

	server, largest := conHash.tree.Ceiling(hash(key))
	if !largest {
		iterator := conHash.tree.Iterator()
		tree_is_empty := iterator.Next()

		// if there is no smallest node then the tree is empty so we return an error
		if !tree_is_empty {
			return "", errors.New("no servers availible at the current moment")
		}
		server = iterator.Node()
	}

	return fmt.Sprintf("%v", server.Value), nil
}

// adds a new Physical node to the consistent hash structure
func (conHash *ConHash) AddNode(physicalNode string) error {

	conHash.physicalNodes = append(conHash.physicalNodes, physicalNode)
	for i := 0; i < conHash.virtualFactor; i++ {
		node_to_insert := physicalNode + strconv.Itoa(i)
		key:=hash(node_to_insert);
		log.Println(key)
		conHash.tree.Put(hash(node_to_insert), physicalNode)
	}
	return nil
}

// Returns all the virtual nodes in sorted order along with their physical nodes
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
	for i, v := range conHash.physicalNodes {
		if v == server {
			index = i
			break
		}
	}

	if index != -1 {
		conHash.physicalNodes = append(conHash.physicalNodes[:index], conHash.physicalNodes[index+1:]...)
	}

	return nil
}
