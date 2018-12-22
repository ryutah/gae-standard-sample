package main

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	adminpb "google.golang.org/genproto/googleapis/iam/admin/v1"

	"cloud.google.com/go/iam/admin/apiv1"
)

var (
	port           = os.Getenv("PORT")
	projectID      = os.Getenv("GOOGLE_CLOUD_PROJECT")
	serviceAccount = os.Getenv("SERVICE_ACCOUNT")
)

func main() {
	e := echo.New()
	e.Use(
		middleware.Logger(),
		middleware.Recover(),
	)

	e.POST("/signbytes", signBytes)

	if port == "" {
		port = "8080"
	}
	http.ListenAndServe(":"+port, e)
}

type signBytesRequest struct {
	Payload string `json:"payload"`
}

func signBytes(c echo.Context) error {
	ctx := c.Request().Context()

	client, err := admin.NewIamClient(ctx)
	if err != nil {
		return err
	}

	payload := new(signBytesRequest)
	if err := c.Bind(payload); err != nil {
		return err
	}
	if payload.Payload == "" {
		return errors.New("invalid request")
	}

	signed := base64.RawStdEncoding.EncodeToString([]byte(payload.Payload))
	resp, err := client.SignBlob(ctx, &adminpb.SignBlobRequest{
		Name:        fmt.Sprintf("projects/%s/serviceAccounts/%s", projectID, serviceAccount),
		BytesToSign: []byte(signed),
	})
	if err != nil {
		return err
	}

	sign := base64.RawStdEncoding.EncodeToString(resp.GetSignature())
	return c.JSON(http.StatusOK, map[string]string{
		"payload":        payload.Payload,
		"signing_key_id": resp.GetKeyId(),
		"signing":        sign,
	})
}

type verifyRequest struct {
	Payload string `json:"payload"`
	Signing string `json:"signing"`
	KID     string `json:"kid"`
}

func verifySigning(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	client, err := admin.NewIamClient(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := client.GetServiceAccountKey(ctx, &adminpb.GetServiceAccountKeyRequest{
		Name:          fmt.Sprintf("projects/%s/serviceAccounts/%s/keys/%s", projectID, serviceAccount, ""),
		PublicKeyType: adminpb.ServiceAccountPublicKeyType_TYPE_X509_PEM_FILE,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	block, _ := pem.Decode(resp.GetPublicKeyData())
	if block == nil {
		http.Error(w, "failed to decode pem", http.StatusInternalServerError)
		return
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := cert.CheckSignature(x509.SHA256WithRSA, []byte(""), []byte("")); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
