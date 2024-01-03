package command

import (
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"bitbucket.org/dptsi/its-go/contracts"
)

const appKeyLength = 32

type GenerateAppKey struct{}

func (c *GenerateAppKey) Key() string {
	return "key:generate"
}

func (c *GenerateAppKey) Name() string {
	return "Generate app key"
}

func (c *GenerateAppKey) Description() string {
	return "Generate 32 bytes application key for encryption (Base64-encoded) and set APP_KEY value in .env"
}

func (c *GenerateAppKey) Usage() string {
	return c.Key()
}

func (c *GenerateAppKey) Handler() contracts.ScriptCommandHandler {
	return func(args []string) error {
		key := make([]byte, appKeyLength)

		if _, err := rand.Read(key); err != nil {
			return err
		}

		encodedKey := base64.StdEncoding.EncodeToString(key)
		fmt.Printf("generated app key (base64-encoded): %s\n", encodedKey)

		file, err := os.Open(".env")
		if err != nil {
			return fmt.Errorf("GenerateAppKey: handler: error opening .env file: %w", err)
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)

		lines := make([]string, 0)
		for scanner.Scan() {
			text := scanner.Text()
			if strings.HasPrefix(text, "APP_KEY") {
				text = fmt.Sprintf("APP_KEY=\"%s\"", encodedKey)
			}
			lines = append(lines, text)
		}
		if err := file.Close(); err != nil {
			return fmt.Errorf("GenerateAppKey: handler: error closing .env file: %w", err)
		}

		file, err = os.Create(".env")
		if err != nil {
			return fmt.Errorf("GenerateAppKey: handler: error creating .env file: %w", err)
		}
		if _, err := file.WriteString(strings.Join(lines, "\n")); err != nil {
			return fmt.Errorf("GenerateAppKey: handler: %w", err)
		}
		if err := file.Close(); err != nil {
			return fmt.Errorf("GenerateAppKey: handler: error closing .env file: %w", err)
		}

		return nil
	}
}
