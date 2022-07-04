package main

import (
	"encoding/json"
	"fmt"
	"io" 
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"math/rand"
)

type Validacion struct{
	Validation_id string `json:"validation_id"`
	Ip_address string `json:"ip_address"`
 	Account_id string `json:"account_id"`
	Type string `json:"type"`
	Validation_status string `json:"validation_status"`
	Date string `json:"creation_date"`
	Detail Details `json:"details"`
	Urls Instructions `json:"instructions"`
	Attachment_status string `json:"attachment_status"`
}

type Instructions struct{
	Front_url string `json:"front_url"`
	Reverse_url string `json:"reverse_url"`
}

type Details struct{
	Document_detail document_detail `json:"document_detail"`
	Document_validations document_validations `json:"document_validations"`
}

type upload_image struct{
	Code float64 `json:"code"`
	Http_code float64 `json:"http_code"`
	Message string `json:"message"`
}

type document_detail struct{
	Birth_place string `json:"birth_place"`
	Country string `json:"country"`
	Creation_date string `json:"creation_date"`
	Date_of_birth string `json:"date_of_birth"`
	Document_number string `json:"document_number"`
	Document_type string `json:"document_type"`
	Expedition_place string `json:"expedition_place"`
	Gender string `json:"gender"`
	Height string `json:"height"`
	Issue_date string `json:"issue_date"`
	Last_name string `json:"last_name"`
	Name string `json:"name"`
	National_registrar string `json:"national_registrar"`
	Production_data string `json:"production_data"`
	Rh string `json:"rh"`
	Update_date string `json:"update_date"`
}

type document_validations struct{
	Data_consistency []document_validations_arrays `json:"data_consistency"`
	Government_database []document_validations_arrays `json:"government_database"`
	Image_analysis []document_validations_arrays `json:"image_analysis"`
	Photocopy_analysis []document_validations_arrays `json:"photocopy_analysis"`
	Photo_of_photo []document_validations_arrays `json:"photo_of_photo"`
}

type document_validations_arrays struct{
	Validation_name string `json:"validation_name"`
	Result string `json:"result"`
	Validation_type string `json:"validation_type"`
	Message string `json:"message"`
	Manually_reviewed bool `json:"manually_reviewed"`
}

// FUNCION PRINCIPAL 
func POST(w http.ResponseWriter, r *http.Request){
	r.ParseMultipartForm(2000)

	// SACAMOS LA INFORMACIÓN DE LAS IMAGENES EN LOS DATOS POST
	file_frontal, fileInfo_frontal, err_frontal := r.FormFile("Frontal")
	file_reverso, fileInfo_reverso, err_reverso := r.FormFile("Reverso")

	// CONTROLAMOS LOS ERRORES EN EL CARGE DE LA IMAGEN, NO DEBEN ESTAR VACIAS
	if err_frontal != nil || err_reverso != nil{
		log.Fatal(err_frontal)
		return
	}

	// CREAMOS UNA CARPETA Y GUARDAMOS LAS IMAGENES PARA USARLAS EN LA CARGA DEL WS
	f_frontal,err := os.OpenFile("./Imagenes/"+fileInfo_frontal.Filename,os.O_WRONLY|os.O_CREATE, 0666)

	if err != nil{
		log.Fatal(err)
		return
	}

	f_reverso,err := os.OpenFile("./Imagenes/"+fileInfo_reverso.Filename,os.O_WRONLY|os.O_CREATE, 0666)

	if err != nil{
		log.Fatal(err)
		return
	}	

	// NOS ASEGURAMOS DE CERRAR EL ARCHIVO DE LA IMAGEN
	defer f_frontal.Close()
	defer f_reverso.Close()
	
	io.Copy(f_frontal,file_frontal)
	io.Copy(f_reverso,file_reverso)	
	
	// CONSUMIMOS LA PRIMERA API PARA CREAR LA VALIDACIÓN INICIAL QUE RETORNA LAS URL PARA CARGAR LAS IMAGENES
	var arreglo Validacion
	arreglo = get_validation()	

	// OBTENEMOS DE LA ESTRUCTURA LAS URL
	url_frontal := arreglo.Urls.Front_url
	url_reverso := arreglo.Urls.Reverse_url

	// IMPRIMIMOS EL ID DE LA VALIDACIÓN
	fmt.Println(arreglo.Validation_id)
	
	// SUBIMOS LAS IMAGENES A LAS URL
	frontal := Subir_Imagen_API(url_frontal, "./Imagenes/"+fileInfo_frontal.Filename)
	reverso := Subir_Imagen_API(url_reverso, "./Imagenes/"+fileInfo_reverso.Filename)

	fmt.Println(frontal.Message)
	fmt.Println(reverso.Message)

	// ESPERAMOS 10 SEGUNDOS ANETES DE ENVIAR LA VALIDACIÓN INICIAL
	time.Sleep(10000 * time.Millisecond)

	for{
		// ENVIAMOS LA PETICION AL WS DE VALIDACIÓN ESPERANDO RESPUESTA EXITOSA O FALLIDA
		var Validacion_Final Validacion
		Validacion_Final = validacion_final(w,r,arreglo.Validation_id)		
		
		fmt.Println(Validacion_Final.Validation_status)

		if err != nil{
			log.Fatal(err)
			return
		}	

		// CUANDA YA SE HAYA TERMINADO LA VALIDACIÓN REVISAMOS QUE ES EXITOSA O NO CON LA INFORMACIÓN DE LA CEDULA
		if Validacion_Final.Validation_status == "success" || Validacion_Final.Validation_status == "failure"{
			if Validacion_Final.Attachment_status == "valid"{
				http.ServeFile(w, r, "Exitoso.php")	
				break;
			}else if Validacion_Final.Attachment_status == "invalid"{
				http.ServeFile(w, r, "Fallido.php")	
				break;
			}			
		}

		// ENTRE CADA SOLICITUD SERAN 10 SEGUNDOS PARA NO SATURAR EL SERVIDOR DE DATOS
		time.Sleep(10000 * time.Millisecond)
	}
}

func get_validation() Validacion{
	
	ale := rand.Intn(101)

	url := "https://api.validations.truora.com/v1/validations"
	method := "POST"

	payload := strings.NewReader("type=document-validation&country=CO&document_type=national-id&user_authorized=true&account_id=" + string(ale))

	client := &http.Client {
	}

	// ARMAMOS LA PETICION CON LA URL Y LOS DATOS QUE SE DEBEN ENVIAR
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		log.Println(err)
		os.Exit(1)		
	}

	req.Header.Add("Truora-API-Key", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50X2lkIjoiIiwiYWRkaXRpb25hbF9kYXRhIjoie30iLCJjbGllbnRfaWQiOiJUQ0kzY2EzNDFjNGQ5Njc2MDQ2ZjI2ZDFmOGJkMDQyMDBjNyIsImV4cCI6MzIzMjgyMjA0MSwiZ3JhbnQiOiIiLCJpYXQiOjE2NTYwMjIwNDEsImlzcyI6Imh0dHBzOi8vY29nbml0by1pZHAudXMtZWFzdC0xLmFtYXpvbmF3cy5jb20vdXMtZWFzdC0xXzZRcXBPblF2NyIsImp0aSI6IjJjY2U1MWQ2LTkzMjctNDBiOS05YmIwLTcxYjQ1NDI0YTM1YSIsImtleV9uYW1lIjoiaW50ZWdyYXRpb25lc19wcm9jZXNvX2RlX3NlbGVjY2lvbiIsImtleV90eXBlIjoiYmFja2VuZCIsInVzZXJuYW1lIjoidHJ1b3JhbmFvcy1pbnRlZ3JhdGlvbmVzX3Byb2Nlc29fZGVfc2VsZWNjaW9uIn0.mMLrjtumE6zxBXc-M_PbGZsoMWTSp9NC-d9kv9jhkCg")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// HACEMOS LA PETICIÓN REQUEST A LA API
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer res.Body.Close()

	response, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}	

	// CREAMOS  EL JSON DEL RESULTADO

	var arreglo Validacion

	err = json.Unmarshal(response, &arreglo)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	return arreglo
}

func Subir_Imagen_API(URL string, RutaImagen string) upload_image{

	// COVERTIMOS LA IMAGEN A BYTES PARA ENVIARLA COMO FILE CONTENT EL LA API
	fileContent, err := ioutil.ReadFile(RutaImagen)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	text := string(fileContent)
	url := URL
	method := "PUT"

	payload := strings.NewReader(text)

	client := &http.Client {
	}

	// ARMAMOS LA PETICION CON LA URL Y LOS DATOS QUE SE DEBEN ENVIAR
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		log.Println(err)
		os.Exit(1)	
	}

	req.Header.Add("Truora-API-Key", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50X2lkIjoiIiwiYWRkaXRpb25hbF9kYXRhIjoie30iLCJjbGllbnRfaWQiOiJUQ0kzY2EzNDFjNGQ5Njc2MDQ2ZjI2ZDFmOGJkMDQyMDBjNyIsImV4cCI6MzIzMjgyMjA0MSwiZ3JhbnQiOiIiLCJpYXQiOjE2NTYwMjIwNDEsImlzcyI6Imh0dHBzOi8vY29nbml0by1pZHAudXMtZWFzdC0xLmFtYXpvbmF3cy5jb20vdXMtZWFzdC0xXzZRcXBPblF2NyIsImp0aSI6IjJjY2U1MWQ2LTkzMjctNDBiOS05YmIwLTcxYjQ1NDI0YTM1YSIsImtleV9uYW1lIjoiaW50ZWdyYXRpb25lc19wcm9jZXNvX2RlX3NlbGVjY2lvbiIsImtleV90eXBlIjoiYmFja2VuZCIsInVzZXJuYW1lIjoidHJ1b3JhbmFvcy1pbnRlZ3JhdGlvbmVzX3Byb2Nlc29fZGVfc2VsZWNjaW9uIn0.mMLrjtumE6zxBXc-M_PbGZsoMWTSp9NC-d9kv9jhkCg")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// HACEMOS LA PETICIÓN REQUEST A LA API
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		os.Exit(1)	
	}
	defer res.Body.Close()

	response, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		os.Exit(1)	
	}

	// CREAMOS  EL JSON DEL RESULTADO

	var arreglo upload_image

	err = json.Unmarshal(response, &arreglo)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	return arreglo
}

func validacion_final(w http.ResponseWriter, r *http.Request, validation_id string) Validacion{
	url := "https://api.validations.truora.com/v1/validations/" + validation_id
	method := "GET"

	client := &http.Client {
	}

	// ARMAMOS LA PETICION CON LA URL Y LOS DATOS QUE SE DEBEN ENVIAR
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	req.Header.Add("Truora-API-Key", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50X2lkIjoiIiwiYWRkaXRpb25hbF9kYXRhIjoie30iLCJjbGllbnRfaWQiOiJUQ0kzY2EzNDFjNGQ5Njc2MDQ2ZjI2ZDFmOGJkMDQyMDBjNyIsImV4cCI6MzIzMjgyMjA0MSwiZ3JhbnQiOiIiLCJpYXQiOjE2NTYwMjIwNDEsImlzcyI6Imh0dHBzOi8vY29nbml0by1pZHAudXMtZWFzdC0xLmFtYXpvbmF3cy5jb20vdXMtZWFzdC0xXzZRcXBPblF2NyIsImp0aSI6IjJjY2U1MWQ2LTkzMjctNDBiOS05YmIwLTcxYjQ1NDI0YTM1YSIsImtleV9uYW1lIjoiaW50ZWdyYXRpb25lc19wcm9jZXNvX2RlX3NlbGVjY2lvbiIsImtleV90eXBlIjoiYmFja2VuZCIsInVzZXJuYW1lIjoidHJ1b3JhbmFvcy1pbnRlZ3JhdGlvbmVzX3Byb2Nlc29fZGVfc2VsZWNjaW9uIn0.mMLrjtumE6zxBXc-M_PbGZsoMWTSp9NC-d9kv9jhkCg")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// HACEMOS LA PETICIÓN REQUEST A LA API
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer res.Body.Close()

	response, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	// CREAMOS  EL JSON DEL RESULTADO

	var arreglo Validacion

	err = json.Unmarshal(response, &arreglo)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	return arreglo	
}

func main() {
	// CUANDO SE ENVIA EL FORMULARIO CARGA LA FUNCIÓN POST
	http.HandleFunc("/POST", POST)

	// PAGINA PRINCIPAL SE CARGA EL FORMULARIO PARA LA SUBIDA DE ARCHIVOS
	http.HandleFunc("/Formulario", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "Formulario.php")
	})	

	err := http.ListenAndServe(":8080", nil)
	
	if err != nil {
		log.Fatal(err)
	}	
}