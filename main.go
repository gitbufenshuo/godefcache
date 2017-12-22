package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"

	"gopkg.in/mgo.v2/bson"

	mgo "gopkg.in/mgo.v2"
)

/*
 */
var readStdin = flag.Bool("i", false, "read file from stdin")
var offset = flag.Int("o", -1, "file offset of identifier in stdin")
var debug = flag.Bool("debug", false, "debug mode")
var tflag = flag.Bool("t", false, "print type information")
var aflag = flag.Bool("a", false, "print public type and member information")
var Aflag = flag.Bool("A", false, "print all type and members information")
var fflag = flag.String("f", "", "Go source filename")
var acmeFlag = flag.Bool("acme", false, "use current acme window")
var jsonFlag = flag.Bool("json", false, "output location in JSON format (-t flag is ignored)")
var renamegodef = flag.String("s", "godef", "in case you want to rename you godef,,,en or use the other tool instead of godef")
var prefire = flag.Bool("p", false, "in case you want to store all symbols of a specified package")
var predefinedPackage = flag.String("h", "", "important to boost, when you confirm that one identifier is the name of one package")

//mongo
var (
	URL        = "127.0.0.1:27017"
	mgoSession *mgo.Session
	dataBase   = "godefcache"
	collection = "godefcache"
)

func getSession() *mgo.Session {
	if mgoSession == nil {
		var err error
		mgoSession, err = mgo.Dial(URL)
		if err != nil {
			panic("dial wrong")
		}
	}
	return mgoSession.Clone()
}

// md5 as key
var gFlagMD5 string

func modifyOffset(src []byte) {
	for idx := *offset - 1; idx != -1; idx-- {
		if (src[idx] >= 'a' && src[idx] <= 'z') || (src[idx] >= 'A' && src[idx] <= 'z') || (src[idx] >= '0' && src[idx] <= '9') || src[idx] == '_' {
		} else {
			*offset = idx + 1
			break
		}
	}
}
func Find() {

}
func getpackagename(src []byte) string {

	pIdx := bytes.Index(src, []byte{'p', 'a', 'c', 'k', 'a', 'g', 'e'})
	packagenamebytes := make([]byte, 0, 0)
	for idx := pIdx + 8; src[idx] != '\r' && src[idx] != '\n'; idx++ {
		packagenamebytes = append(packagenamebytes, src[idx])
	}
	return string(packagenamebytes)
}
func modifyMD5(src []byte) string {
	packagename := getpackagename(src)
	// go ahead
	headIdx := 0
	for idx := *offset; idx != len(src); idx++ {
		if notLetter(src[idx]) {
			headIdx = idx
			break
		}
	}
	// go back
	backIdx := 0
	for idx := *offset - 1; idx != -1; idx-- {
		if notLetter(src[idx]) && src[idx] != '.' {
			backIdx = idx + 1
			break
		}
	}
	if backIdx == headIdx {
		return ""
	}
	qualifiedidName := string(src[backIdx:headIdx])

	dotNum := strings.Count(qualifiedidName, ".")
	if dotNum == 0 {
		if isCapital(qualifiedidName[0]) {
			return packagename + "." + qualifiedidName
		}
		if src[headIdx] == '(' {
			// a local-package function call
			return packagename + "." + qualifiedidName
		}
		// a local-variable
		return "......"
	}
	if dotNum == 1 {
		segs := strings.Split(qualifiedidName, ".")
		if n, _ := getSession().DB(dataBase).C(collection).Find(bson.M{"_id": segs[0]}).Count(); n == 0 {
			// not a predefined package
			return "......"
		}
		return qualifiedidName
	}
	return "......"
}

func isCapital(c byte) bool {
	if c >= 'a' && c <= 'z' {
		return true
	}
	return false
}

// https://golang.org/ref/spec#letter
func notLetter(c byte) bool {
	if c == '_' {
		return false
	}
	if c >= 'a' && c <= 'z' {
		return false
	}
	if c >= 'A' && c <= 'Z' {
		return false
	}
	if c >= '0' && c <= '9' {
		return false
	}
	return true
}
func md55(input string, inputbytes []byte) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(input))
	md5Ctx.Write(inputbytes)
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}
func genFlagMD5() {
	flagStr := ""
	flagStr += fmt.Sprintf("%v\n", *readStdin)
	flagStr += fmt.Sprintf("%v\n", *offset)
	flagStr += fmt.Sprintf("%v\n", *debug)
	flagStr += fmt.Sprintf("%v\n", *tflag)

	flagStr += fmt.Sprintf("%v\n", *aflag)
	flagStr += fmt.Sprintf("%v\n", *Aflag)
	flagStr += fmt.Sprintf("%v\n", *fflag)
	flagStr += fmt.Sprintf("%v\n", *acmeFlag)
	flagStr += fmt.Sprintf("%v\n", *jsonFlag)
	gFlagMD5 = md55(flagStr, []byte{})
}

type godefcache struct {
	Raw string
}
type godefname struct {
	Name string
}

func fprefire() {
	if *prefire {
		srcbuf := bytes.NewBuffer([]byte{})
		abspath := ""
		if !strings.HasSuffix(*fflag, "/") { // relative path
			abspath = path.Join(os.Getenv("PWD"), *fflag)
		}
		fmt.Printf(">>>>>> firing [%v]\n", abspath)
		if srcfile, err := os.Open(abspath); err != nil {
			fail("%v", err.Error())
		} else {
			srccode, _ := ioutil.ReadAll(srcfile)
			srcbuf.Write(srccode)
			var wg sync.WaitGroup
			for jdx := 0; jdx != srcbuf.Len()/5+1; jdx++ {
				for idx := jdx; idx < srcbuf.Len(); idx += (srcbuf.Len()/5 + 1) {
					wg.Add(1)
					go func(jdx, idx int) {
						defer wg.Done()
						defer fmt.Printf(">>>>>> [%v]--[%v]:[%v]|[%v]\n", jdx, abspath, srcbuf.Len(), idx)
						godefcmd := exec.Command("godef", "-i", "-t", "-o", fmt.Sprintf("%v", idx), "-f", fmt.Sprintf("%v", abspath))
						godefcmd.Run()
					}(jdx, idx)
				}
				wg.Wait()
				fmt.Println("...---...", jdx, srcbuf.Len()/5+1)
			}
		}
		success(abspath + "[ok]")
	}

}
func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: godefcache [flags] [expr]\n")
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() > 1 {
		flag.Usage()
		os.Exit(2)
	}
	fprefire()
	if *renamegodef != "godef" {
		getSession().DB(dataBase).C(collection).Upsert(bson.M{"_id": "toolname"}, bson.M{"$set": bson.M{"name": *renamegodef}})
		success(*renamegodef)
	}
	if *predefinedPackage != "" {
		if strings.HasPrefix(*predefinedPackage, "add:") {
			packagename := strings.TrimPrefix(*predefinedPackage, "add:")
			getSession().DB(dataBase).C(collection).Insert(struct {
				ID string `bson:"_id"`
			}{packagename})
			success(*predefinedPackage)
		}
		if strings.HasPrefix(*predefinedPackage, "del:") {
			packagename := strings.TrimPrefix(*predefinedPackage, "del:")
			getSession().DB(dataBase).C(collection).RemoveId(packagename)
			success(*predefinedPackage)
		}
		fail("add:[packagename]  or  del:[packagename]")
	}
	var src []byte
	if !(*readStdin) {
		fail("%v", "Only support stdin now....")
	} else {
		src, _ = ioutil.ReadFile(*fflag)
	}
	src = append(src, ' ') // for robust
	modifyOffset(src)
	unique := modifyMD5(src)
	if unique == "" {
		success("xxx")
	}
	if !(unique == "......") {
		gFlagMD5 = md55(unique, nil)
	} else {
		genFlagMD5()
		gFlagMD5 = md55(gFlagMD5, src)
	}
	// success(unique)
	// dealWithPredefinedpackage(src)
	var result godefcache
	if err := getSession().DB(dataBase).C(collection).Find(bson.M{"_id": gFlagMD5}).One(&result); err != nil {
		var cmdstr string
		cmdstdin := bytes.NewBuffer(src)
		cmdstdout := bytes.NewBuffer([]byte{})
		cmdstderr := bytes.NewBuffer([]byte{})
		var godefnameresult godefname
		if err := getSession().DB(dataBase).C(collection).Find(bson.M{"_id": "toolname"}).One(&godefnameresult); err != nil {
			cmdstr = "godef"
		} else {
			cmdstr = godefnameresult.Name
		}
		godefcmd := exec.Command(cmdstr, "-i", "-t", "-o", fmt.Sprintf("%v", *offset), "-f", fmt.Sprintf("%v", *fflag))
		godefcmd.Stdin = cmdstdin
		godefcmd.Stdout = cmdstdout
		godefcmd.Stderr = cmdstderr
		if err := godefcmd.Run(); err == nil {
			raw := cmdstdout.String()
			getSession().DB(dataBase).C(collection).Insert(bson.M{"_id": gFlagMD5, "raw": raw})
			success(raw)
		} else {
			fail("%v", err.Error()+"[  --  ]"+string(cmdstderr.Bytes()))
		}
	}
	success(result.Raw)
}

///////////////////////////////////////////
func fail(s string, a ...interface{}) {
	fmt.Fprint(os.Stderr, "godef: "+fmt.Sprintf(s, a...)+"\n")
	os.Exit(2)
}
func success(s string) {
	fmt.Fprint(os.Stdout, s)
	os.Exit(0)
}
