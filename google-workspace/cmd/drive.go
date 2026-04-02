package cmd

import (
	"context"
	"io"
	"os"

	"github.com/PanthroCorp-Limited/openclaw-skills/google-workspace/internal/config"
	gw "github.com/PanthroCorp-Limited/openclaw-skills/google-workspace/internal/google"
	"github.com/PanthroCorp-Limited/openclaw-skills/google-workspace/internal/oauth"
	"github.com/spf13/cobra"
)

var driveCmd = &cobra.Command{
	Use:   "drive",
	Short: "Read-only Google Drive operations",
}

func driveClient() (*gw.DriveClient, context.Context) {
	ctx := context.Background()
	key := encryptionKey()
	if key == "" {
		exitf("GOOGLE_WORKSPACE_TOKEN_KEY is not set")
	}

	cfg, err := config.Load(configDir)
	if err != nil {
		exitf("loading config: %v", err)
	}
	if !cfg.Drive {
		exitf("drive is disabled in config; run 'google-workspace config set --drive=true'")
	}

	token, err := oauth.LoadToken(configDir, key)
	if err != nil {
		exitf("%v", err)
	}

	oauthCfg := oauth.NewOAuthConfig(clientID(), clientSecret(), cfg.OAuthScopes())
	ts := oauthCfg.TokenSource(ctx, token)

	client, err := gw.NewDriveClient(ctx, ts)
	if err != nil {
		exitf("creating drive client: %v", err)
	}
	return client, ctx
}

var (
	driveListQuery      string
	driveListMaxResults int64
)

var driveListCmd = &cobra.Command{
	Use:   "list",
	Short: "List files in Google Drive",
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := driveClient()
		files, err := client.ListFiles(ctx, driveListQuery, driveListMaxResults)
		if err != nil {
			exitf("%v", err)
		}
		printJSON(files)
	},
}

var driveGetID string

var driveGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get file metadata by ID",
	Run: func(cmd *cobra.Command, args []string) {
		if driveGetID == "" {
			exitf("--id is required")
		}
		client, ctx := driveClient()
		file, err := client.GetFile(ctx, driveGetID)
		if err != nil {
			exitf("%v", err)
		}
		printJSON(file)
	},
}

var driveDownloadID string

var driveDownloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download file content to stdout",
	Long:  "Download file content to stdout. Google Docs export as plain text, Sheets as CSV, Slides as plain text. All other files download as raw bytes.",
	Run: func(cmd *cobra.Command, args []string) {
		if driveDownloadID == "" {
			exitf("--id is required")
		}
		client, ctx := driveClient()
		rc, err := client.DownloadFile(ctx, driveDownloadID)
		if err != nil {
			exitf("%v", err)
		}
		defer rc.Close()

		if _, err := io.Copy(os.Stdout, rc); err != nil {
			exitf("writing output: %v", err)
		}
	},
}

func init() {
	driveListCmd.Flags().StringVar(&driveListQuery, "query", "", "Drive search query (e.g. \"name contains 'report'\")")
	driveListCmd.Flags().Int64Var(&driveListMaxResults, "max-results", 20, "maximum number of results")

	driveGetCmd.Flags().StringVar(&driveGetID, "id", "", "file ID")

	driveDownloadCmd.Flags().StringVar(&driveDownloadID, "id", "", "file ID")

	driveCmd.AddCommand(driveListCmd, driveGetCmd, driveDownloadCmd)
	rootCmd.AddCommand(driveCmd)
}
