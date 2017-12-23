package main
import(
  
)

func main(){

  var obj Indexer = Indexer{}
  path := "/home/mahmod/filestest"
  obj.OpenCon()
  obj.ReadFolder(path)

}
