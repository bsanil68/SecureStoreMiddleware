package controllers

import (
	"SecureStore/config"
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	"github.com/hashgraph/hedera-sdk-go/v2"
	"storj.io/uplink"
)

type UploadController struct {
	StorjConfig *config.StorjConfig
}

func NewUploadController(config *config.StorjConfig) *UploadController {
	return &UploadController{StorjConfig: config}
}

func (uc *UploadController) HandleUpload(c *gin.Context) {
	// Get the file from the request
	file, header, err := c.Request.FormFile("file")

	customerID := c.PostForm("customerID")
	state := c.PostForm("state")
	version := c.PostForm("version")

	// Generate a unique document ID
	documentID := uuid.New().String()
	filename := customerID + state + version + documentID + header.Filename

	Documenthash := md5.New()
	if _, err := io.Copy(Documenthash, file); err != nil {
		fmt.Println("Error hashing file:", err)
		return
	}

	hashInBytes := Documenthash.Sum(nil)[:md5.Size] // Trim to 16 bytes
	hashString := fmt.Sprintf("%x", hashInBytes)

	fmt.Println("MD5 hash:", hashString)

	contractID, err := hedera.ContractIDFromString(os.Getenv("HEDERA_DOCTOKEN_CONTRACT_ID"))
	if err != nil {
		log.Fatalf("failed to parse contract ID: %v", err)
	}

	hederaClient := hedera.ClientForTestnet()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	// Define upload path in Storj bucket
	uploadPath := filepath.Join(customerID, state, filename)

	// Create Storj project
	ctx := context.Background()
	access, err := uc.StorjConfig.GetAccess()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not parse access grant"})
		return
	}

	project, err := uplink.OpenProject(ctx, access)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not open project"})
		return
	}
	defer project.Close()

	// Ensure the bucket exists
	_, err = project.EnsureBucket(ctx, uc.StorjConfig.BucketName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not ensure bucket"})
		return
	}

	// Start the upload
	upload, err := project.UploadObject(ctx, uc.StorjConfig.BucketName, uploadPath, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not initiate upload"})
		return
	}

	// Copy the file content to the upload
	_, err = io.Copy(upload, file)
	if err != nil {
		_ = upload.Abort()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not upload data"})
		return
	}

	// Commit the upload
	err = upload.Commit()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not commit object"})
		return
	}

	//  connect  to hedera  and  upload the hash  with  other details
	hederaClient, err = connectHedera()
	//Call the smart contract to store document details
	result, err := SaveDocumentDetailsOnBlock(hederaClient, contractID, hashString, customerID, state, version)
	if err != nil {
		log.Printf("Failed to save document details in smart contract: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save document details"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Document uploaded successfully",
		"documentHash":  Documenthash,
		"smartContract": result,
	})
}

// New method to list files based on document ID
func (uc *UploadController) ListFiles(c *gin.Context) {
	customerID := c.Param("customerID")
	if customerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "customerID is required"})
		return
	}

	ctx := context.Background()
	access, err := uc.StorjConfig.GetAccess()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not parse access grant"})
		return
	}

	project, err := uplink.OpenProject(ctx, access)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not open project"})
		return
	}
	defer project.Close()

	// List objects in the bucket
	objects := project.ListObjects(ctx, uc.StorjConfig.BucketName, &uplink.ListObjectsOptions{
		Prefix: customerID + "/", // Assuming files are stored in a folder named after the document ID
	})

	var files []string
	for objects.Next() {
		item := objects.Item()
		// Remove the document ID prefix from the file name
		fileName := strings.TrimPrefix(item.Key, customerID+"/")
		files = append(files, fileName)
	}

	if err := objects.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error listing objects"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"files": files})
}

func (uc *UploadController) GetFile(c *gin.Context) {
	fileName := c.Param("fileName")
	if fileName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "fileName is required"})
		return
	}

	ctx := context.Background()
	access, err := uc.StorjConfig.GetAccess()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not parse access grant"})
		return
	}

	project, err := uplink.OpenProject(ctx, access)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not open project"})
		return
	}
	defer project.Close()

	download, err := project.DownloadObject(ctx, uc.StorjConfig.BucketName, fileName, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not download object"})
		return
	}
	defer download.Close()

	fileContent, err := ioutil.ReadAll(download)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not read file content"})
		return
	}

	c.Data(http.StatusOK, "application/pdf", fileContent)
}
