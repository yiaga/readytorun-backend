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
	"time" // Used for the delay to prevent race conditions

	"readytorun-backend/internal/database"
	"readytorun-backend/internal/models"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

// googleDriveFolderID is the ID of the Google Drive Shared Drive where files will be uploaded.
// IMPORTANT: Replace this with your actual Shared Drive ID.
// Ensure the service account (readytorun-uploader-xxxx.json) has 'Manager' access
// to this Shared Drive.
const googleDriveFolderID = "0AKMSRdTcq8gHUk9PVA"

var driveService *drive.Service

// init function initializes the Google Drive service client when the package loads.
// It reads the service account credentials and sets up the client with the necessary scope.
func init() {
	ctx := context.Background()
	// Path to your service account key file. Ensure this file is present.
	credentialFilePath := "readytorun-uploader-35abeffd60a4.json" 

	data, err := os.ReadFile(credentialFilePath)
	if err != nil {
		log.Fatalf("Unable to read client secret file %s: %v", credentialFilePath, err)
	}

	// Create a JWT config from the credentials.
	// drive.DriveScope provides full read/write access to all files and folders
	// in Google Drive that the service account has access to.
	config, err := google.JWTConfigFromJSON(data, drive.DriveScope) 
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	// Create an HTTP client using the JWT config. This client will handle authentication.
	client := config.Client(ctx)

	// Initialize the Google Drive service with the authenticated HTTP client.
	driveService, err = drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}
	log.Println("Google Drive service initialized successfully.")
}

// uploadFileToGoogleDrive handles uploading a file to the specified Google Drive folder.
// It returns the public webViewLink of the uploaded file.
func uploadFileToGoogleDrive(file multipart.File, fileName string, mimeType string, folderID string) (string, error) {
    if driveService == nil {
        return "", fmt.Errorf("Google Drive service not initialized")
    }

    tmpfile, err := os.CreateTemp("", "upload-*.tmp")
    if err != nil {
        return "", fmt.Errorf("failed to create temporary file: %w", err)
    }
    defer os.Remove(tmpfile.Name())
    defer tmpfile.Close()

    if _, err := io.Copy(tmpfile, file); err != nil {
        return "", fmt.Errorf("failed to copy file to temporary location: %w", err)
    }

    if _, err := tmpfile.Seek(0, 0); err != nil {
        return "", fmt.Errorf("failed to seek temporary file: %w", err)
    }

    fileMetadata := &drive.File{
        Name:     fileName,
        Parents:  []string{folderID},
        MimeType: mimeType,
    }

    // Upload file to Google Drive (Shared Drive safe)
    res, err := driveService.Files.Create(fileMetadata).
        Media(tmpfile).
        Fields("id, webViewLink, name").
        SupportsAllDrives(true).
        Do()
    if err != nil {
        return "", fmt.Errorf("failed to upload file to Google Drive: %w", err)
    }

    log.Printf("File '%s' uploaded to Google Drive. Link: %s\n", res.Name, res.WebViewLink)

    // Retry permission setting up to 3 times (exponential backoff)
    var permErr error
    for attempt := 1; attempt <= 3; attempt++ {
        _, permErr = driveService.Permissions.Create(res.Id, &drive.Permission{
            Type: "anyone",
            Role: "reader",
        }).SupportsAllDrives(true).Do()

        if permErr == nil {
            log.Printf("Public permissions set for file: %s", res.Name)
            break
        }

        waitTime := time.Duration(attempt*2) * time.Second
        log.Printf("Attempt %d: Failed to set permissions for file %s (ID: %s): %v. Retrying in %v...",
            attempt, res.Name, res.Id, permErr, waitTime)
        time.Sleep(waitTime)
    }

    if permErr != nil {
        log.Printf("Warning: Could not set public permissions for file %s (ID: %s): %v", res.Name, res.Id, permErr)
    }

    return res.WebViewLink, nil
}

// CreateRegistration handles the HTTP POST request for new user registrations.
// It parses form data, uploads files to Google Drive, and saves registration details to SQLite.
func CreateRegistration(w http.ResponseWriter, r *http.Request) {
    // Parse the multipart form data with a max memory of 50 MB.
    err := r.ParseMultipartForm(50 << 20) 
    if err != nil {
        http.Error(w, "Could not parse form: " + err.Error(), http.StatusBadRequest)
        log.Printf("Error parsing multipart form: %v", err)
        return
    }

    var reg models.Registration
    // Map form values to the Registration model.
    reg.Fullname = r.FormValue("fullname")
    reg.Dob = r.FormValue("dateOfBirth") // Matches frontend field name
    reg.Gender = r.FormValue("gender")
    reg.Email = r.FormValue("email")
    reg.Phone = r.FormValue("phone")
    reg.StateOfOrigin = r.FormValue("stateOfOrigin")     // Matches frontend field name
    reg.StateOfResidence = r.FormValue("stateOfResidence") // Matches frontend field name
    reg.Education = r.FormValue("education")
    reg.PreviousOffice = r.FormValue("previousOffice")     // Matches frontend field name
    reg.InterestedOffice = r.FormValue("interestedOffice") // Matches frontend field name
    reg.PreviousContest = r.FormValue("previousContest")   // Matches frontend field name
    
    // Convert string "yes"/"no" from frontend to boolean.
    reg.CardCarryingMember = (r.FormValue("partyMember") == "yes") // Matches frontend field name
    reg.Consent = (r.FormValue("consent") == "true")

    reg.Motivation = r.FormValue("motivation")
    reg.PoliticalUnderstanding = r.FormValue("politicalUnderstanding")
    reg.OtherSupport = r.FormValue("otherSupport")         // Matches frontend field name
    reg.PreferredCommunication = r.FormValue("communication") // Matches frontend field name

    // Handle array fields from frontend (assistanceNeeded[] and availability[]).
    // These are stored as JSON strings in the database.
    assistanceNeededValues := r.Form["assistanceNeeded[]"]
    if len(assistanceNeededValues) > 0 {
        jsonBytes, marshalErr := json.Marshal(assistanceNeededValues)
        if marshalErr != nil {
            log.Printf("Error marshalling assistanceNeeded: %v", marshalErr)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }
        reg.AssistanceNeeded = string(jsonBytes)
    } else {
        reg.AssistanceNeeded = "[]" // Store empty JSON array if no options selected
    }

    availabilityValues := r.Form["availability[]"]
    if len(availabilityValues) > 0 {
        jsonBytes, marshalErr := json.Marshal(availabilityValues)
        if marshalErr != nil {
            log.Printf("Error marshalling availability: %v", marshalErr)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }
        reg.Availability = string(jsonBytes)
    } else {
        reg.Availability = "[]" // Store empty JSON array if no options selected
    }

    // --- File uploads to Google Drive ---
    var partyCardLink, resumeLink string

    // Process and upload Party Card (field name "partyCard" from frontend).
    if file, handler, err := r.FormFile("partyCard"); err == nil {
        defer file.Close()
        mimeType := handler.Header.Get("Content-Type")
        if mimeType == "" {
            mimeType = getMimeTypeFromExtension(handler.Filename) // Fallback for mime type
        }
        partyCardLink, err = uploadFileToGoogleDrive(file, handler.Filename, mimeType, googleDriveFolderID)
        if err != nil {
            log.Printf("Error uploading party card: %v", err)
            http.Error(w, "Error uploading party card: " + err.Error(), http.StatusInternalServerError)
            return
        }
        reg.PartyMembershipDocLink = partyCardLink
    } else if err != http.ErrMissingFile {
        // Log error only if it's not simply a missing optional file (http.ErrMissingFile)
        log.Printf("Error getting partyCard file: %v", err)
        http.Error(w, "Error processing party card file: " + err.Error(), http.StatusInternalServerError)
        return
    }

    // Process and upload Resume (field name "resume" from frontend).
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

    // Prepare and execute the SQL INSERT statement.
    // Ensure your 'registrations' table schema matches these fields,
    // especially `party_membership_doc_link` and `cv_link` as TEXT.
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

    // Execute the statement with the registration data.
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

    // Send a success response back to the client.
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Registration successful, files uploaded to Google Drive"})
}

// getMimeTypeFromExtension is a helper function to guess the MIME type
// based on the file extension. It's a fallback if Content-Type header is missing.
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
        return "application/octet-stream" // Default for unknown types
    }
}