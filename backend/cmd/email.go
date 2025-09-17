package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/walterfan/lazy-rabbit-reminder/pkg/email"
	"github.com/walterfan/lazy-rabbit-reminder/pkg/log"
)

var (
	emailSubject   string
	emailBody      string
	emailTo        []string
	emailCC        []string
	testConnection bool
)

// emailCmd represents the email command
var emailCmd = &cobra.Command{
	Use:   "email",
	Short: "Send emails using the configured SMTP server",
	Long: `Send emails using the configured SMTP server.
	
This command allows you to:
- Test SMTP connection
- Send simple emails
- Send emails to multiple recipients

Configuration is loaded from environment variables:
MAIL_SERVER, MAIL_PORT, MAIL_USERNAME, MAIL_PASSWORD, etc.`,
	Run: func(cmd *cobra.Command, args []string) {
		runEmailCommand()
	},
}

// runEmailCommand runs the email command using direct email sender
func runEmailCommand() {
	// Create email sender
	sender, err := email.NewEmailSender()
	if err != nil {
		log.GetLogger().Fatalf("Failed to create email sender: %v", err)
	}

	// Test connection if requested
	if testConnection {
		fmt.Println("Testing SMTP connection...")
		if err := sender.TestConnection(); err != nil {
			log.GetLogger().Fatalf("Connection test failed: %v", err)
		}
		fmt.Println("✅ SMTP connection test successful!")

		// Show configuration (without password)
		config := sender.GetConfig()
		fmt.Printf("Configuration:\n")
		fmt.Printf("  SMTP Host: %s\n", config.SMTPHost)
		fmt.Printf("  SMTP Port: %d\n", config.SMTPPort)
		fmt.Printf("  Use SSL: %t\n", config.UseSSL)
		fmt.Printf("  Use TLS: %t\n", config.UseTLS)
		fmt.Printf("  Username: %s\n", config.Username)
		fmt.Printf("  Default From: %s\n", config.DefaultFrom)
		fmt.Printf("  Default To: %s\n", config.DefaultTo)
		return
	}

	// Send email if subject and body are provided
	if emailSubject != "" && emailBody != "" {
		fmt.Println("Sending email...")

		if len(emailTo) > 0 {
			// Send to specific recipients
			err = sender.SendEmailToRecipients(emailSubject, emailBody, emailTo, emailCC)
		} else {
			// Send to default recipient
			err = sender.SendSimpleEmail(emailSubject, emailBody)
		}

		if err != nil {
			log.GetLogger().Fatalf("Failed to send email: %v", err)
		}

		fmt.Println("✅ Email sent successfully!")
		if len(emailTo) > 0 {
			fmt.Printf("Recipients: %v\n", emailTo)
			if len(emailCC) > 0 {
				fmt.Printf("CC: %v\n", emailCC)
			}
		} else {
			fmt.Println("Sent to default recipient")
		}
	} else {
		fmt.Println("No email to send. Use --subject and --body flags, or --test to test connection.")
	}
}

func init() {
	rootCmd.AddCommand(emailCmd)

	// Email content flags
	emailCmd.Flags().StringVarP(&emailSubject, "subject", "s", "", "Email subject")
	emailCmd.Flags().StringVarP(&emailBody, "body", "b", "", "Email body")

	// Recipient flags
	emailCmd.Flags().StringSliceVarP(&emailTo, "to", "t", []string{}, "To addresses (comma-separated)")
	emailCmd.Flags().StringSliceVarP(&emailCC, "cc", "c", []string{}, "CC addresses (comma-separated)")

	// Test flag
	emailCmd.Flags().BoolVar(&testConnection, "test", false, "Test SMTP connection without sending email")

	// Examples
	emailCmd.Example = `  # Test SMTP connection
  lazy-rabbit-reminder email --test

  # Send simple email to default recipient
  lazy-rabbit-reminder email --subject "Test" --body "Hello World"

  # Send email to specific recipients
  lazy-rabbit-reminder email \
    --subject "Important Notice" \
    --body "This is important" \
    --to "user1@example.com,user2@example.com" \
    --cc "manager@example.com"

  # Send email with multiline body
  lazy-rabbit-reminder email \
    --subject "Welcome" \
    --body "Welcome to our system!

Thank you for registering.

Best regards,
Admin Team"`
}
