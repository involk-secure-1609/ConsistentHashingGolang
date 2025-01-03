package main

import conhash "kevin/consistent_hash/con_hash"


func main(){
	// initialize the number of virtual nodes
	virtualFactor:=3;

	// initialize the servers
	physicalNodes:=[]string{"localhost:8080","localhost:8081","localhost:8082"}

	// intialize the consistentHash 
	cHash:=conhash.NewConhash(virtualFactor,physicalNodes)
	cHash.Initialize()

	// see the ordering of the nodes
	cHash.Layout()
}