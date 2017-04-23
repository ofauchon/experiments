package main

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
	"net/http"

	"gopkg.in/gomail.v2"
	"github.com/benbjohnson/phantomjs"
	"github.com/spf13/viper"
)


//
// Download page content
//
func download(pUrl string) []byte {
	log.Println("Downloading", pUrl)
	resp, err := http.Get(pUrl)
	if err != nil {
		log.Fatal(err)
	} else {
		defer resp.Body.Close()
		bytes, _ := ioutil.ReadAll(resp.Body)
		return bytes
	}
	return nil
}

//
// Use PhantomJs to render the page
//
func shot(pUrl string) error {

	// Start the process once.
	if err := phantomjs.DefaultProcess.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer phantomjs.DefaultProcess.Close()

	p := phantomjs.DefaultProcess

	page, err := p.CreateWebPage()

	if err != nil {
		return err
	}
	defer page.Close()

	// Open a URL.
	if err := page.Open(pUrl); err != nil {
		return err
	}

	// Setup the viewport and render the results view.
	if err := page.SetViewportSize(1024, 800); err != nil {
		return err
	}

//    page.SetClipRect(rectJson{Top:0, Left:0, Width: 1024, Height: 1024})

	if err := page.Render(".temp.png", "png", 25); err != nil {
		return err
	}

	return nil
}

//
// sendMail
//
func sendMail(pSub, pMsg, pTo, pFrom , pUrl string) {


	m := gomail.NewMessage()

    m.SetHeader("From", pFrom);
    m.SetHeader("To", pTo);
    m.SetHeader("Subject", pSub);
    m.Embed(".temp.png");
    m.SetBody("text/html",`<a href="` + pUrl + `">Lien Original</a><img src="cid:.temp.png" alt="My image" />`) 
    d := gomail.NewDialer(viper.GetString("smtp.host"), viper.GetInt("smtp.port") , viper.GetString("smtp.user"), viper.GetString("smtp.pass")  )
    if err := d.DialAndSend(m); err != nil {
        panic(err)
    }
}

func main() {


	log.Println("Start Alert System")
	viper.SetConfigName("lbc-alert")
	viper.AddConfigPath("$HOME/.config")
	viper.AddConfigPath("./config")

     err := viper.ReadInConfig()
  	if err != nil {
    	fmt.Println("Config file not found...")
        os.Exit(1)
  	} 
	fmt.Println("Config file found ...")

	cUrl:=viper.GetString("alert1.url")
	cFrom:=viper.GetString("alert1.from")
	cTo:=viper.GetString("alert1.to")
	cSubject:=viper.GetString("alert1.subject")

    fmt.Println("url:",cUrl, " subject: ",cSubject," from:",cFrom," to:",cTo)


	cur_md5 := md5.Sum([]byte{0})
	old_md5 := md5.Sum([]byte{0})
	pause := 60

	for {
		html := download(cUrl)
		cur_md5 = md5.Sum(html)

		if old_md5 != cur_md5 {
			log.Println("MD5 Changed ", old_md5, cur_md5)
			log.Println("Creating screenshot for page: ", cUrl)
			shot(cUrl);
			log.Println("Sending mail from:",cFrom, " to:", cTo)
			sendMail(cSubject, string(html), cTo, cFrom, cUrl)
			old_md5=cur_md5

		} else {
			log.Println("MD5 did not change")
		}
		log.Println("Going to sleep for   ", pause, " seconds")
		time.Sleep(time.Second * 300)
	}

}
