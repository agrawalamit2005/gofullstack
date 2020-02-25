package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/EngineerKamesh/gofullstack/volume4/section2/gopherface/common"
	"github.com/EngineerKamesh/gofullstack/volume4/section2/gopherface/common/asyncq"
	"github.com/EngineerKamesh/gofullstack/volume4/section2/gopherface/common/authenticate"
	"github.com/EngineerKamesh/gofullstack/volume4/section2/gopherface/common/utility"
	"github.com/EngineerKamesh/gofullstack/volume4/section2/gopherface/tasks"
)

type UploadImageForm struct {
	PageTitle  string
	FieldNames []string
	Fields     map[string]string
	Errors     map[string]string
}

type BillItem struct {
	//AuditableContent        // Embedded type
	UUID         string `json:"uuid"`
	BillDate     string `json:"billdate"`
	BillAmount   int    `json:"billamount"`
	UserName     string `json:"username"`
	BillFileName string `json: "billfilename"`
	Caption      string `json:"caption"`
}

func DisplayUploadImageForm(w http.ResponseWriter, r *http.Request, u *UploadImageForm) {
	RenderGatedTemplate(w, WebAppRoot+"/templates/uploadimageform.html", u)
}

func ProcessUploadImage(w http.ResponseWriter, r *http.Request, u *UploadImageForm, e *common.Env) {

	shouldProcessThumbnailAsynchronously := true

	file, fileheader, err := r.FormFile("imagefile")

	if err != nil {
		log.Println("Encountered error when attempting to read uploaded file: ", err)
	}

	randomFileName := utility.GenerateUUID()

	if fileheader != nil {

		extension := filepath.Ext(fileheader.Filename)
		r.ParseMultipartForm(32 << 20)

		defer file.Close()

		imageFilePathWithoutExtension := "./static/uploads/images/" + randomFileName
		f, err := os.OpenFile(imageFilePathWithoutExtension+extension, os.O_WRONLY|os.O_CREATE, 0666)

		if err != nil {
			log.Println(err)
			return
		}

		defer f.Close()
		io.Copy(f, file)

		// Note: Moved the thumbnail generation logic (commented out code block below) to the
		// ImageResizeTask object in the tasks package.
		thumbnailResizeTask := tasks.NewImageResizeTask(imageFilePathWithoutExtension, extension)

		if shouldProcessThumbnailAsynchronously == true {

			asyncq.TaskQueue <- thumbnailResizeTask

		} else {

			thumbnailResizeTask.Perform()
		}

		m := make(map[string]string)
		m["thumbnailPath"] = strings.TrimPrefix(imageFilePathWithoutExtension, ".") + "_thumb.png"
		m["imagePath"] = strings.TrimPrefix(imageFilePathWithoutExtension, ".") + ".png"
		m["PageTitle"] = "Image Preview"

		if e != nil {
			//fmt.Println("reached Adding bill %s generate %s", (fileheader.Filename + extension), (imageFilePathWithoutExtension + extension))
			gfSession, err := authenticate.SessionStore.Get(r, "gopherface-session")
			if err != nil {
				log.Print(err)
				return
			}
			uuid := gfSession.Values["uuid"].(string)
			username := gfSession.Values["username"].(string)
			log.Printf("log reached Adding bill for %s original name is %s generated is %s", username, (fileheader.Filename + extension), (imageFilePathWithoutExtension + extension))
			e.DB.AddBill(uuid, (fileheader.Filename), (imageFilePathWithoutExtension + extension))
		} else {
			fmt.Println("e is NUll")
			log.Printf("log e is NUll")
		}

		RenderGatedTemplate(w, WebAppRoot+"/templates/pdfUploadConfirmation.html", m)

	} else {
		w.Write([]byte("Failed to process uploaded file!"))
	}
}

func ValidateUploadImageForm(w http.ResponseWriter, r *http.Request, u *UploadImageForm) {

	ProcessUploadImage(w, r, u, nil)

}

func ValidateUploadImageFormDB(w http.ResponseWriter, r *http.Request, u *UploadImageForm, e *common.Env) {

	ProcessUploadImage(w, r, u, e)

}
func GetPdfHandler(e *common.Env) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var uuid string = ""
		var username string = ""
		if e != nil {
			//fmt.Println("reached Adding bill %s generate %s", (fileheader.Filename + extension), (imageFilePathWithoutExtension + extension))
			gfSession, err := authenticate.SessionStore.Get(r, "gopherface-session")
			if err != nil {
				log.Print(err)
				return
			}
			uuid = gfSession.Values["uuid"].(string)
			username = gfSession.Values["username"].(string)
			//log.Printf("log reached Adding bill for %s original name is %s generated is %s", username, (fileheader.Filename + extension), (imageFilePathWithoutExtension + extension))
			//e.DB.AddBill(uuid, (fileheader.Filename), (imageFilePathWithoutExtension + extension))
		} else {
			fmt.Println("e is NUll")
			log.Printf("log e is NUll")
		}
		bill := BillItem{UUID: uuid, BillDate: "today", BillAmount: 12, UserName: username, BillFileName: "abc.pdf", Caption: "BillCaption"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(bill)
	})
}

func UploadImageHandler(w http.ResponseWriter, r *http.Request) {

	u := UploadImageForm{}
	u.Fields = make(map[string]string)
	u.Errors = make(map[string]string)
	u.PageTitle = "Upload Image"

	switch r.Method {

	case "GET":
		DisplayUploadImageForm(w, r, &u)
	case "POST":
		ValidateUploadImageForm(w, r, &u)
	default:
		DisplayUploadImageForm(w, r, &u)
	}

}

func UploadImageHandlerDB(e *common.Env) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("log reached UploadImageHandlerDB")
		fmt.Println("reached UploadImageHandlerDB")
		u := UploadImageForm{}
		u.Fields = make(map[string]string)
		u.Errors = make(map[string]string)
		u.PageTitle = "Upload Image DB"

		switch r.Method {

		case "GET":
			DisplayUploadImageForm(w, r, &u)
		case "POST":
			ValidateUploadImageFormDB(w, r, &u, e)
		default:
			DisplayUploadImageForm(w, r, &u)
		}

	})
}
