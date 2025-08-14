package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"readytorun-backend/internal/database"
	"readytorun-backend/internal/models"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

// Define the Google Drive folder ID where files will be uploaded.
// IMPORTANT: Replace this with your actual Google Drive folder ID.
// Ensure the service account (credentials.json) has write access to this folder.
const googleDriveFolderID = "1v80Ki3tIebe3msQoNIome86RYfLtq-ni"

var driveService *drive.Service

func init() {
	// Initializes the Google Drive service client on package load.
	ctx := context.Background()
	credentialFilePath := "credentials.json"

	data, err := os.ReadFile(credentialFilePath)
	if err != nil {
		log.Fatalf("Unable to read client secret file %s: %v", credentialFilePath, err)
	}

	config, err := google.JWTConfigFromJSON(data, drive.DriveFileScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	client := config.Client(ctx)
	driveService, err = drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}
	log.Println("Google Drive service initialized successfully.")
}

// uploadFileToGoogleDrive uploads a file to a specified Google Drive folder
// and returns the public webViewLink.
func uploadFileToGoogleDrive(file multipart.File, fileName string, mimeType string, folderID string) (string, error) {
	if driveService == nil {
		return "", fmt.Errorf("Google Drive service not initialized")
	}

	// Create a temporary file to read the multipart.File content
	tmpfile, err := os.CreateTemp("", "upload-*.tmp")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer os.Remove(tmpfile.Name()) // Clean up the temporary file after function exits
	defer tmpfile.Close()

	if _, err := io.Copy(tmpfile, file); err != nil {
		return "", fmt.Errorf("failed to copy file to temporary location: %w", err)
	}

	if _, err := tmpfile.Seek(0, 0); err != nil {
		return "", fmt.Errorf("failed to seek temporary file: %w", err)
	}

	fileMetadata := &drive.File{
		Name:    fileName,
		Parents: []string{folderID},
		MimeType: mimeType,
	}

	res, err := driveService.Files.Create(fileMetadata).
		Media(tmpfile).
		Fields("id, webViewLink"). // Request ID and public URL
		SupportsAllDrives(true).
		Do()

	if err != nil {
		return "", fmt.Errorf("failed to upload file to Google Drive: %w", err)
	}

	// Set file permissions to public.
	_, err = driveService.Permissions.Create(res.Id, &drive.Permission{
		Type: "anyone",
		Role: "reader",
	}).Do()
	if err != nil {
		log.Printf("Warning: Failed to set public permissions for file %s (ID: %s): %v", fileName, res.Id, err)
	}

	log.Printf("File '%s' uploaded to Google Drive. Link: %s\n", fileName, res.WebViewLink)
	return res.WebViewLink, nil
}


func CreateRegistration(w http.ResponseWriter, r *http.Request) {
    err := r.ParseMultipartForm(50 << 20) // Increased to 50 MB.
    if err != nil {
        http.Error(w, "Could not parse form: " + err.Error(), http.StatusBadRequest)
        log.Printf("Error parsing multipart form: %v", err)
        return
    }

    var reg models.Registration
    reg.Fullname = r.FormValue("fullname")
    reg.Dob = r.FormValue("dateOfBirth")
    reg.Gender = r.FormValue("gender")
    reg.Email = r.FormValue("email")
    reg.Phone = r.FormValue("phone")
    reg.StateOfOrigin = r.FormValue("stateOfOrigin")
    reg.StateOfResidence = r.FormValue("stateOfResidence")
    reg.Education = r.FormValue("education")
    reg.PreviousOffice = r.FormValue("previousOffice")
    reg.InterestedOffice = r.FormValue("interestedOffice")
    reg.PreviousContest = r.FormValue("previousContest")
    
    // Handle boolean values from frontend
    reg.CardCarryingMember = (r.FormValue("partyMember") == "yes")
    reg.Consent = (r.FormValue("consent") == "true")

    reg.Motivation = r.FormValue("motivation")
    reg.PoliticalUnderstanding = r.FormValue("politicalUnderstanding")
    reg.OtherSupport = r.FormValue("otherSupport")
    reg.PreferredCommunication = r.FormValue("communication")

    // Handle array fields by marshalling them into JSON strings
    assistanceNeededValues := r.Form["assistanceNeeded[]"]
    if len(assistanceNeededValues) > 0 {
        jsonBytes, _ := json.Marshal(assistanceNeededValues)
        reg.AssistanceNeeded = string(jsonBytes)
    } else {
        reg.AssistanceNeeded = "[]"
    }

    availabilityValues := r.Form["availability[]"]
    if len(availabilityValues) > 0 {
        jsonBytes, _ := json.Marshal(availabilityValues)
        reg.Availability = string(jsonBytes)
    } else {
        reg.Availability = "[]"
    }

    // --- Handle file uploads to Google Drive ---
    var partyCardLink, resumeLink string

    // Upload Party Card
    if file, handler, err := r.FormFile("partyCard"); err == nil {
        defer file.Close()
        mimeType := handler.Header.Get("Content-Type")
        if mimeType == "" {
            mimeType = getMimeTypeFromExtension(handler.Filename)
        }
        partyCardLink, err = uploadFileToGoogleDrive(file, handler.Filename, mimeType, googleDriveFolderID)
        if err != nil {
            log.Printf("Error uploading party card: %v", err)
            http.Error(w, "Error uploading party card: " + err.Error(), http.StatusInternalServerError)
            return
        }
        reg.PartyMembershipDocLink = partyCardLink
    } else if err != http.ErrMissingFile {
        log.Printf("Error getting partyCard file: %v", err)
        http.Error(w, "Error processing party card file: " + err.Error(), http.StatusInternalServerError)
        return
    }

    // Upload Resume
    if file, handler, err := r.FormFile("resume"); err == nil {
        defer file.Close()
        mimeType := handler.Header.Get("Content-Type")
        if mimeType == "" {
            mimeType = getMimeTypeFromExtension(handler.Filename)
        }
        resumeLink, err = uploadFileToGoogleDrive(file, handler.Filename, mimeType, googleDriveFolderID)
        if err != nil {
            log.Printf("Error uploading resume: %v", err)
            http.Error(w, "Error uploading resume: " + err.Error(), http.StatusInternalServerError)
            return
        }
        reg.CVLink = resumeLink
    } else if err != http.ErrMissingFile {
        log.Printf("Error getting resume file: %v", err)
        http.Error(w, "Error processing resume file: " + err.Error(), http.StatusInternalServerError)
        return
    }

    // Prepare and execute SQL statement
    stmt, err := database.DB.Prepare(`
        INSERT INTO registrations(
            fullname, dob, gender, email, phone, state_of_origin, state_of_residence, education,
            previous_office, interested_office, previous_contest, card_carrying_member, party_membership_doc_link, cv_link, motivation,
            political_understanding, assistance_needed, other_support, availability, preferred_communication, consent
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `)
    if err != nil {
        log.Println("Prepare error:", err)
        http.Error(w, "DB error", http.StatusInternalServerError)
        return
    }
    defer stmt.Close()

    _, err = stmt.Exec(
        reg.Fullname, reg.Dob, reg.Gender, reg.Email, reg.Phone, reg.StateOfOrigin, reg.StateOfResidence,
        reg.Education, reg.PreviousOffice, reg.InterestedOffice, reg.PreviousContest, reg.CardCarryingMember,
        reg.PartyMembershipDocLink, reg.CVLink, reg.Motivation,
        reg.PoliticalUnderstanding, reg.AssistanceNeeded, reg.OtherSupport, reg.Availability,
        reg.PreferredCommunication, reg.Consent,
    )
    if err != nil {
        log.Println("Exec error:", err)
        http.Error(w, "Insert failed", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Registration successful, files uploaded to Google Drive"})
}

// getMimeTypeFromExtension is a helper to guess MIME type based on file extension.
func getMimeTypeFromExtension(filename string) string {
    ext := strings.ToLower(filepath.Ext(filename))
    switch ext {
    case ".pdf":
        return "application/pdf"
    case ".doc":
        return "application/msword"
    case ".docx":
        return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
    case ".jpg", ".jpeg":
        return "image/jpeg"
    case ".png":
        return "image/png"
    default:
        return "application/octet-stream"
    }
}