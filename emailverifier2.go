package main

import (
	"bufio"
	"fmt"
	"github.com/bobesa/go-domain-util/domainutil"
	trumail "github.com/sdwolfe32/trumail/verifier"
	"github.com/zenthangplus/goccm"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
)

var success []string
var t string
var mut sync.Mutex
var c = goccm.New(10)

func main() {
	fmt.Println("#################################################################")
	fmt.Println("Email address verifier Multi-Threaded written independently by olumide")
	fmt.Println("#################################################################")

	if (len(os.Args) < 2 || len(os.Args) > 3) {
		fmt.Println("Missing parameter, provide file names in .txt e.g program.exe sites.txt")
		os.Exit(1)
	}

	f := os.Args[1]
	//f:="applemagicmouse.txt"
	k:=strings.Split(f,".")
	t=k[0]
	file, err := os.Open(f)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	//line := 0
	scanner := bufio.NewScanner(file)
	//get number of email in the file
	for scanner.Scan() {
		//line := scanner.Text()
		//line++
		// do something with your line
		type todo struct {
			Email_Username string
			Domainname     string
		}
		var unwantedemail bool
		email := scanner.Text()
		list := []string{"policy","catalog","web","help","gov", "edu", "org", "fbi", "police", "admin", "postmaster", "webmaster", "staples",
			"press", "abuse", "security", "hostmaster", "helps", "ebay", "paypal", "help", "staples", "amazon", "microsoft",
			"donotreply", "jobs", "billing", "domain", "help", "apple", "privacy", "notes", "copyrights", "advertising","gmail","hotmail","yahoo","outlook","aol","live","mail","lycos","zoho","yandex"}
		//validate email if it is a legal email address
		emailRegexp := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
		if !emailRegexp.MatchString(email) {
			continue
		}
		//check if email consist of unwanted substring
		for _, sub := range list {
			if strings.Contains(email, sub) {
				unwantedemail = true
				break
			}
		}
		if unwantedemail {
			continue
		}
		//prepare email format
		components := strings.Split(email, "@")
		username:=components[0]
		if strings.Count(username,".")>1{continue} //count the number of dots in username,if greater than one then continue
		validdomain := domainutil.DomainSuffix(components[1])
		if validdomain != "com" && validdomain != "us" {
			continue
		}
		c.Wait()
		process(validdomain,email)

	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading input:", err)
	}
	c.WaitAllDone()

}
func process(validdomain string,email string)  {

	v := trumail.NewVerifier(validdomain, email)

	//	Address ValidFormat, Deliverable, FullInbox, HostExists, CatchAll
	// Validate a single address
	//log.Println(v.Verify(email))
	olu,_:=v.Verify(email)
	mut.Lock()
	fmt.Printf("%s ,%t\n",email,olu.Deliverable)
	if olu.Deliverable==true{
		//save file in emailist_verified.txt
		success=append(success,email)
		stringByte := strings.Join(success, "\r\n")
		// If the file doesn't exist, create it, or append to the file
		err := ioutil.WriteFile(t+"_verified.txt", []byte(stringByte), 0644)
		if err != nil {
			log.Fatal(err)
		}

	}
	mut.Unlock()
	c.Done()

}