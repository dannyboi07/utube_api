package controller

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"utube/redis"
	"utube/schema"
	"utube/utils"

	"github.com/go-chi/chi/v5"
)

var maxVideoFileSize int = 1000 * 1000 * 100 // 2MB

func CreateVideo(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxVideoFileSize))

	jDec := json.NewDecoder(r.Body)
	jDec.DisallowUnknownFields()

	var (
		createVideoReq schema.CreateVideoRequest
		statusCode     int
		err            error
	)
	statusCode, err = utils.JsonParseErr(jDec.Decode(&createVideoReq))
	if err != nil {
		utils.WriteApiErrMessage(w, statusCode, err.Error())
		utils.Log.Println("Failed to parse request, err:", err)
		return
	}

	var (
		notValid bool
		reason   string
	)
	notValid, reason = createVideoReq.Validate()
	if notValid {
		utils.WriteApiErrMessage(w, http.StatusBadRequest, reason)
		utils.Log.Println("Invalid create video request: ", reason)
		return
	}

	// Connect with Redis, create further logic
	err = redis.SetCreateVideoReq(redis.CreateVideoCached{
		Name:      *createVideoReq.Name,
		Size:      *createVideoReq.Size,
		Uploaded:  0,
		UserId:    0,
		UserEmail: "danielchettiar@gmail.com",
	})
	if err != nil {
		utils.WriteApiErrMessage(w, 0, "Failed to register request to create video")
		utils.Log.Println("Failed to register request to create video in cache, err:", err)
		return
	}

	utils.WriteDataToResponse(w, schema.CreateVideoResponse{
		Name: *createVideoReq.Name,
		Size: *createVideoReq.Size,
		Type: *createVideoReq.Type,
	})
}

func UploadStatus(w http.ResponseWriter, r *http.Request) {
	var userDetails map[string]interface{} = r.Context().Value("userDetails").(map[string]interface{})
	var userId uint64 = userDetails["userId"].(uint64)
	var videoName string = chi.URLParam(r, "video_name")

	var (
		result string
		err    error
	)

	result, err = redis.GetVideoUploadStatus(userId, videoName)
	utils.Log.Println("result", result, err)
}

func UploadChunk(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxVideoFileSize))

	jDec := json.NewDecoder(r.Body)
	jDec.DisallowUnknownFields()

	var (
		uploadVideoChunkRequest schema.UploadVideoChunkRequest
		statusCode              int
		err                     error
	)
	statusCode, err = utils.JsonParseErr(jDec.Decode(&uploadVideoChunkRequest))
	if err != nil {
		utils.WriteApiErrMessage(w, statusCode, err.Error())
		utils.Log.Println("Failed to parse request, err:", err)
		return
	}

	utils.Log.Println("request", *uploadVideoChunkRequest.Name, *uploadVideoChunkRequest.Data)
	data, _ := base64.StdEncoding.DecodeString(*uploadVideoChunkRequest.Data)
	utils.Log.Println("decoded", data, string(data))
	utils.Log.Println(len(string(data)), len(data))
}

func UploadVideo(w http.ResponseWriter, r *http.Request) {
	// r.Body = http.MaxBytesReader(w, r.Body, int64(maxVideoFileSize))

	var (
		requestReader *multipart.Reader
		err           error
	)
	requestReader, err = r.MultipartReader()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		utils.Log.Println("Controller error parsing as MultipartReader, err: ", err)
		return
	}

	var (
		logString  string
		statusCode int
	)
	var lastModifiedField *multipart.Part
	lastModifiedField, err = requestReader.NextPart()
	logString, statusCode, err = utils.ValidateFormField("last_modified", lastModifiedField, err, false)
	if err != nil {
		http.Error(w, err.Error(), statusCode)
		utils.Log.Println(logString)
		return
	}

	var videoField *multipart.Part
	videoField, err = requestReader.NextPart()
	logString, statusCode, err = utils.ValidateFormField("video_file", videoField, err, true)
	if err != nil {
		http.Error(w, err.Error(), statusCode)
		utils.Log.Println(logString)
		return
	}

	valid, videoFile, mimeType := utils.ValidVideoFileType(videoField)
	if !valid {
		utils.Log.Println("video file invalid")
		return
	}

	bytes, _ := ioutil.ReadAll(videoFile)
	os.WriteFile(fmt.Sprintf("video%s", mimeType.Extension()), bytes, 0666)
}
