package main
import(
  "fmt"
  "os"
  "log"
  "bufio"
  "strings"
  "github.com/reiver/go-porterstemmer"
  "gopkg.in/fatih/set.v0"
)

func main(){
  path := "/home/yazan/itsaperfectdayforsomemayhem/filestest"
  //path:="../../../../filestest"

  dir,er:=os.Open(path)
  defer dir.Close()

  if er!=nil{
    log.Fatal(er)
  }

  filenames,er:=dir.Readdirnames(2)

  if er!=nil{
    log.Fatal(er)
  }

  mp:=make(map[string][]string)

  for i:=0; i<len(filenames); i++{
    filepath:= path + "/" + filenames[i]

    file,er:=os.Open(filepath)
    defer file.Close()

    if er!=nil{
      log.Fatal(er)
    }

    scanner := bufio.NewScanner(file)
    scanner.Split(bufio.ScanLines)

    //var words[]string
    words:=set.New()

    for scanner.Scan(){
      line := scanner.Text()
      //words = append(words,strings.Split(line," ")...)
      str:=strings.Split(line, " ")

      for j:=0;j<len(str);j++{
        words.Add(str[j])
      }

    }


  }

  wordsset:=set.StringSlice(words)

  for i:=0;i<len(words);i++{
    wordsset[i] = porterstemmer.StemString(wordsset[i])
    fmt.Println(wordsset[i])
  }
}
