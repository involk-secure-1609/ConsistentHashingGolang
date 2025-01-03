package main

import conhash "kevin/consistent_hash/con_hash"


func main(){
	virtualFactor:=3;
	servers:=[]string{"localhost:8080","localhost:8081","localhost:8082"}
	cHash:=conhash.NewConhash(virtualFactor,servers)
	cHash.Initialize()
	cHash.Layout()
}